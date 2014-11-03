// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates loading and rendering an Ice file.
package main

import (
	"azul3d.org/v1/chippy"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/gfx/window"
	"azul3d.org/v1/ice"
	"azul3d.org/v1/keyboard"
	"azul3d.org/v1/math"
	"azul3d.org/v1/mouse"
	"image"
	_ "image/jpeg"
	"log"
	"os"
)

var glslVert = []byte(`
#version 120

attribute vec3 Vertex;
attribute vec2 TexCoord0;

uniform mat4 MVP;

varying vec2 tc0;

void main()
{
	tc0 = TexCoord0;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
`)

var glslFrag = []byte(`
#version 120

varying vec2 tc0;

uniform sampler2D Texture0;
uniform bool BinaryAlpha;

void main()
{
	gl_FragColor = texture2D(Texture0, tc0);
	//gl_FragColor = vec4(1.0, 0.0, 0.0, 1.0);//texture2D(Texture0, tc0);
	if(BinaryAlpha && gl_FragColor.a < 0.5) {
		discard;
	}
}
`)

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	// Load the Ice file.
	scene, err := ice.LoadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camFOV := 75.0
	camNear := 0.1
	camFar := 100.0
	camera.SetPersp(r.Bounds(), camFOV, camNear, camFar)

	// Move the camera -2 on the Y axis (back two units away from the triangle
	// object).
	//camera.SetPos(math.Vec3{0, -50, 10})
	//camera.SetPos(math.Vec3{0, -5, 2})
	camera.SetPos(math.Vec3{0, -7, 3})

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSLVert = glslVert
	shader.GLSLFrag = glslFrag

	// Preload the shader (useful for seeing shader errors, if any).
	onLoad := make(chan *gfx.Shader, 1)
	r.LoadShader(shader, onLoad)
	go func() {
		<-onLoad
		shader.RLock()
		if !shader.Loaded {
			log.Println(string(shader.Error))
		}
		shader.RUnlock()
	}()

	// Assign the shader to each object in the scene.
	for _, o := range scene.Objects {
		o.Shader = shader
		o.State.FaceCulling = gfx.NoFaceCulling
		//var verts = make([]gfx.Vec3, 0, len(o.Meshes[0].Indices))
		//for _, v := range o.Meshes[0].Indices {
		//	verts = append(verts, o.Meshes[0].Vertices[v])
		//}
		//o.Meshes[0].Vertices = verts
		//o.Meshes[0].Indices = nil
		//if len(o.Meshes[0].Indices) > 5 {
		//	log.Println(name, len(o.Meshes[0].Indices))
		//	bad := o.Meshes[0].Indices[743]
		//	log.Println(o.Meshes[0].Vertices[bad])
		//	o.Meshes[0].Indices = o.Meshes[0].Indices[744-3:744]
		//}
	}

	// Start a goroutine to handle window events and move the camera around.
	go func() {
		event := w.Events()
		for {
			select {
			case e := <-event:
				switch ev := e.(type) {
				case keyboard.TypedEvent:
					if ev.Rune == 'm' {
						// Toggle MSAA now.
						msaa := !r.MSAA()
						r.SetMSAA(msaa)
						log.Println("MSAA Enabled?", msaa)
					}

				case mouse.Event:
					if ev.Button == mouse.Left && ev.State == mouse.Down {
						w.SetCursorGrabbed(!w.CursorGrabbed())
					}
				}
			}
		}
	}()

	event := w.Events()
	for {
	camEvents:
		for {
			select {
			case e := <-event:
				switch ev := e.(type) {
				case chippy.ResizedEvent:
					// Update the camera's projection matrix for the new width and
					// height.
					camera.Lock()
					camera.SetPersp(r.Bounds(), camFOV, camNear, camFar)
					camera.Unlock()

				case chippy.CursorPositionEvent:
					if w.CursorGrabbed() {
						dt := r.Clock().Dt()
						camera.Lock()
						camRot := camera.Rot()
						camRot.Z -= 4 * ev.X * dt
						camRot.X -= 4 * ev.Y * dt
						camera.SetRot(camRot)
						camera.Unlock()
					}
				}
			default:
				break camEvents
			}
		}

		// Move the camera now.
		if w.CursorGrabbed() {
			dt := r.Clock().Dt()
			var local, parent math.Vec3
			speed := 16.0
			if w.Keyboard.Down(keyboard.A) {
				local.X -= speed * dt
			}
			if w.Keyboard.Down(keyboard.D) {
				local.X += speed * dt
			}
			if w.Keyboard.Down(keyboard.W) {
				local.Y += speed * dt
			}
			if w.Keyboard.Down(keyboard.S) {
				local.Y -= speed * dt
			}
			if w.Keyboard.Down(keyboard.LeftCtrl) {
				parent.Z -= speed * dt
			}
			if w.Keyboard.Down(keyboard.LeftShift) {
				parent.Z += speed * dt
			}
			camera.Lock()
			worldSpace := camera.ConvertPos(local, gfx.LocalToWorld)
			parentSpace := camera.ConvertPos(worldSpace, gfx.WorldToParent)
			camera.SetPos(parentSpace.Add(parent))
			camera.Unlock()
		}

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Draw each model in the scene.
		for _, o := range scene.Objects {
			r.Draw(image.Rect(0, 0, 0, 0), o, camera)
		}

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Usage: azul3d_ice [file.ice][file.json]")
	}
	window.Run(gfxLoop)
}
