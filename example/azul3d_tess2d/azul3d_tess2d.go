// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates 2D polygon tesselation.
package main

import (
	"image"

	"log"

	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/lmath.v1"
	"azul3d.org/vgfx.v1"
)

var glslVert = []byte(`
#version 120

attribute vec3 Vertex;
attribute vec4 Color;

uniform mat4 MVP;

varying vec4 frontColor;

void main()
{
	frontColor = Color;
	gl_Position = MVP * vec4(Vertex, 1.0);
}
`)

var glslFrag = []byte(`
#version 120

varying vec4 frontColor;

void main()
{
	gl_FragColor = frontColor;
}
`)

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w window.Window, r gfx.Renderer) {
	// Setup a camera to use a orthographic projection.
	camera := gfx.NewCamera()
	camNear := 0.01
	camFar := 1000.0
	camera.SetOrtho(r.Bounds(), camNear, camFar)

	// Move the camera back two units away from the scene.
	camera.SetPos(lmath.Vec3{0, -2, 0})

	// Create a simple shader.
	shader := gfx.NewShader("SimpleShader")
	shader.GLSLVert = glslVert
	shader.GLSLFrag = glslFrag

	// Create a mesh.
	mesh := gfx.NewMesh()

	// Tesselate the polygons into triangles.
	tess := vgfx.NewTess()
	mesh.Vertices = tess.Tesselate(polygons, mesh.Vertices)
	mesh.Vertices = mesh.Vertices[:len(mesh.Vertices)-(len(mesh.Vertices)%3)]
	log.Println(len(mesh.Vertices), len(mesh.Vertices)%3)
	/*
		mesh.Vertices = []gfx.Vec3{
			// Top
			{0, 0, 1},
			{-.5, 0, 0},
			{.5, 0, 0},

			// Bottom-Left
			{-.5, 0, 0},
			{-1, 0, -1},
			{0, 0, -1},

			// Bottom-Right
			{.5, 0, 0},
			{0, 0, -1},
			{1, 0, -1},
		}
	*/

	// Create a graphics object.
	obj := gfx.NewObject()
	obj.Shader = shader
	obj.Meshes = []*gfx.Mesh{mesh}
	obj.State.FaceCulling = gfx.NoFaceCulling

	for {
		b := r.Bounds()
		aspectRatio := float64(b.Dx()) / float64(b.Dy())
		scale := float64(b.Dy()) / 2.0
		obj.Transform.SetScale(lmath.Vec3{
			X: scale * (1.0 / aspectRatio),
			Y: scale,
			Z: scale,
		})
		obj.Transform.SetPos(lmath.Vec3{
			X: float64(b.Dy()) / 2.0,
			Z: float64(b.Dx()) / 2.0,
		})

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Draw the object.
		r.Draw(image.Rect(0, 0, 0, 0), obj, camera)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
