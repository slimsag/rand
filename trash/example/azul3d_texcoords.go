// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates texture coordinates.
package main

import (
	"azul3d.org/v1/chippy"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/gfx/window"
	"azul3d.org/v1/keyboard"
	"azul3d.org/v1/math"
	"image"
	"io/ioutil"
	"log"
	"os"
)

func mustLoadShader(name, vert, frag string) *gfx.Shader {
	shader := gfx.NewShader(name)

	vertFile, err := os.Open(vert)
	if err != nil {
		panic(err)
	}

	shader.GLSLVert, err = ioutil.ReadAll(vertFile)
	if err != nil {
		panic(err)
	}

	fragFile, err := os.Open(frag)
	if err != nil {
		panic(err)
	}

	shader.GLSLFrag, err = ioutil.ReadAll(fragFile)
	if err != nil {
		panic(err)
	}

	return shader
}

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camNear := 0.01
	camFar := 1000.0
	camera.SetOrtho(r.Bounds(), camNear, camFar)

	// Move the camera back two units away from the card.
	camera.SetPos(math.Vec3{0, -2, 0})

	// Create a simple shader.
	shader := mustLoadShader(
		"SimpleShader",
		"src/azul3d.org/v1/examples/azul3d_texcoords/shader.vert",
		"src/azul3d.org/v1/examples/azul3d_texcoords/shader.frag",
	)

	// Load the picture.
	f, err := os.Open("src/azul3d.org/v1/assets/textures/texture_coords_1024x1024.png")
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// Create new texture.
	tex := &gfx.Texture{
		Source:    img,
		MinFilter: gfx.Linear,
		MagFilter: gfx.Linear,
		Format:    gfx.DXT1RGBA,
	}

	// Create a card mesh.
	cardMesh := &gfx.Mesh{
		Vertices: []gfx.Vec3{
			// Bottom-left triangle.
			{-1, 0, -1},
			{1, 0, -1},
			{-1, 0, 1},

			// Top-right triangle.
			{-1, 0, 1},
			{1, 0, -1},
			{1, 0, 1},
		},
		TexCoords: []gfx.TexCoordSet{
			{
				Slice: []gfx.TexCoord{
					{0, 1},
					{1, 1},
					{0, 0},

					{0, 0},
					{1, 1},
					{1, 0},
				},
			},
		},
	}

	// Create a card object.
	card := gfx.NewObject()
	card.AlphaMode = gfx.AlphaToCoverage
	card.Shader = shader
	card.Textures = []*gfx.Texture{tex}
	card.Meshes = []*gfx.Mesh{cardMesh}

	go func() {
		for e := range w.Events() {
			switch ev := e.(type) {
			case chippy.ResizedEvent:
				// Update the camera's projection matrix for the new width and
				// height.
				camera.Lock()
				camera.SetOrtho(r.Bounds(), camNear, camFar)
				camera.Unlock()
			case keyboard.TypedEvent:
				if ev.Rune == 'b' {
					card.Shader.Lock()

					var enabled bool
					v, ok := card.Shader.Inputs["Enabled"]
					if ok {
						enabled = v.(bool)
					}
					card.Shader.Inputs["Enabled"] = !enabled
					card.Shader.Unlock()
				}
			}
		}
	}()

	for {
		// Center the card in the window.
		b := r.Bounds()
		card.SetPos(math.Vec3{float64(b.Dx()) / 2.0, 0, float64(b.Dy()) / 2.0})

		// Scale the card to fit the window.
		s := float64(b.Dy()) / 2.0 // Card is two units wide, so divide by two.
		card.SetScale(math.Vec3{s, s, s})

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Draw the textured card.
		r.Draw(image.Rect(0, 0, 0, 0), card, camera)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop)
}
