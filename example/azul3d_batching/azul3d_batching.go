// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Draws too many triangles.
package main

import (
	"azul3d.org/v1/chippy"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/gfx/window"
	"fmt"
	"image"
)

import _ "net/http/pprof"
import "net/http"

func init() {
	go func() {
		fmt.Println(http.ListenAndServe(":6060", nil))
	}()
}

var glslVert = []byte(`
#version 120

attribute vec4 Vertex;
attribute vec4 Color;

varying vec4 frontColor;

void main()
{
	frontColor = Color;
	gl_Position = Vertex;
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
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	// Create a simple shader.
	shader := &gfx.Shader{
		Name:     "SimpleShader",
		GLSLVert: glslVert,
		GLSLFrag: glslFrag,
	}

	// Preload the shader (useful for seeing shader errors, if any).
	onLoad := make(chan *gfx.Shader, 1)
	r.LoadShader(shader, onLoad)
	go func() {
		<-onLoad
		shader.RLock()
		if shader.Loaded {
			fmt.Println("Shader loaded")
		} else {
			fmt.Println(string(shader.Error))
		}
		shader.RUnlock()
	}()

	n := 25000
	// Create a new batch.
	triangles := make([]*gfx.Object, 0, n)
	for i := 0; i < n; i++ {
		// Create a triangle object.
		triangle := gfx.NewObject()
		triangle.Shader = shader
		triangle.Meshes = []*gfx.Mesh{
			&gfx.Mesh{
				Vertices: []gfx.Vec3{
					// Top
					{-.5, 0, 0},
					{.5, 0, 0},
					{0, 1, 0},
				},
				Colors: []gfx.Color{
					// Top
					{1, 0, 0, 1},
					{0, 1, 0, 1},
					{0, 0, 1, 1},
				},
			},
		}
		triangles = append(triangles, triangle)
	}

	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Clear a few rectangles on the window using different background
		// colors.
		r.Clear(image.Rect(0, 100, 640, 380), gfx.Color{0, 1, 0, 1})
		r.Clear(image.Rect(100, 100, 540, 380), gfx.Color{1, 0, 0, 1})
		r.Clear(image.Rect(100, 200, 540, 280), gfx.Color{0, 0.5, 0.5, 1})
		r.Clear(image.Rect(200, 200, 440, 280), gfx.Color{1, 1, 0, 1})

		// Draw triangles using the batcher.
		for _, triangle := range triangles {
			r.Draw(image.Rect(50, 50, 640-50, 480-50), triangle, nil)
		}

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop)
}
