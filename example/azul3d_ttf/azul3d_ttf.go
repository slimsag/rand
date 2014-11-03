// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Example - Demonstrates texture coordinates.
package main

import (
	"go/build"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"azul3d.org/chippy.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/lmath.v1"
	"azul3d.org/text.v1"
)

// This helper function is not an important example concept, please ignore it.
//
// absPath the absolute path to an file given one relative to the examples
// directory:
//  $GOPATH/src/azul3d.org/examples.v1
var examplesDir string

func absPath(relPath string) string {
	if len(examplesDir) == 0 {
		// Find assets directory.
		for _, path := range filepath.SplitList(build.Default.GOPATH) {
			path = filepath.Join(path, "src/azul3d.org/examples.v1")
			if _, err := os.Stat(path); err == nil {
				examplesDir = path
				break
			}
		}
	}
	return filepath.Join(examplesDir, relPath)
}

func loadShaderSources(s *gfx.Shader, vert, frag string) {
	vertFile, err := os.Open(vert)
	if err != nil {
		panic(err)
	}

	s.GLSLVert, err = ioutil.ReadAll(vertFile)
	if err != nil {
		panic(err)
	}

	fragFile, err := os.Open(frag)
	if err != nil {
		panic(err)
	}

	s.GLSLFrag, err = ioutil.ReadAll(fragFile)
	if err != nil {
		panic(err)
	}
}

func mustLoadShader(name, vert, frag string) *gfx.Shader {
	shader := gfx.NewShader(name)
	loadShaderSources(shader, vert, frag)
	return shader
}

// gfxLoop is responsible for drawing things to the window.
func gfxLoop(w window.Window, r gfx.Renderer) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camNear := 0.01
	camFar := 1000.0
	camera.SetOrtho(r.Bounds(), camNear, camFar)

	// Move the camera back two units away from the card.
	camera.SetPos(lmath.Vec3{0, -2, 0})

	// Create a simple shader.
	shader := mustLoadShader(
		"SimpleShader",
		absPath("azul3d_ttf/shader.vert"),
		absPath("azul3d_ttf/shader.frag"),
	)

	// Create a card mesh.
	var font text.Font
	var err error
	font, err = text.LoadFontFile(absPath("assets/fonts/vera/Vera.ttf"))
	if err != nil {
		log.Fatal(err)
	}

	txt, err := text.New(font, "`Hello World!")
	if err != nil {
		log.Fatal(err)
	}

	txt.Object.Shader = shader

	/*
		buf := truetype.NewGlyphBuf()
		err = buf.Load(font, font.FUnitsPerEm(), font.Index('T'), truetype.NoHinting)
		if err != nil {
			log.Fatal(err)
		}

		textMesh := gfx.NewMesh()
		tess := text.NewTess()
		tess.AppendGlyph(textMesh, font, 'T')

		// Create a card object.
		card := gfx.NewObject()
		card.AlphaMode = gfx.AlphaToCoverage
		card.FaceCulling = gfx.NoFaceCulling
		card.Shader = shader
		card.Meshes = []*gfx.Mesh{textMesh}


						err = buf.Load(font, font.FUnitsPerEm(), font.Index(ev.Rune), truetype.NoHinting)
						if err != nil {
							log.Fatal(err)
						}

						textMesh := gfx.NewMesh()
						text.AppendMesh(buf, textMesh)

						textMesh := gfx.NewMesh()
						tess := text.NewTess()
						tess.Append(textMesh, font, 'T')

						card.Lock()
						card.Meshes = []*gfx.Mesh{textMesh}
						card.Unlock()
	*/

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
				if ev.Rune == 'r' {
					txt.Lock()
					txt.Shader.Lock()
					txt.Shader.Reset()
					loadShaderSources(
						txt.Shader,
						absPath("azul3d_ttf/shader.vert"),
						absPath("azul3d_ttf/shader.frag"),
					)
					txt.Shader.Unlock()
					txt.Unlock()

				} else if ev.Rune == 'b' {
					txt.Shader.Lock()

					var enabled bool
					v, ok := txt.Shader.Inputs["Enabled"]
					if ok {
						enabled = v.(bool)
					}
					txt.Shader.Inputs["Enabled"] = !enabled
					txt.Shader.Unlock()
				} else {
					txt.Set(string(ev.Rune))
				}
			}
		}
	}()

	for {
		// Center the card in the window.
		b := r.Bounds()
		txt.SetPos(lmath.Vec3{float64(b.Dx()) / 2.0, 0, float64(b.Dy()) / 2.0})

		// Scale the card to fit the window.
		s := float64(b.Dy()) / 2.0 // Card is two units wide, so divide by two.
		txt.SetScale(lmath.Vec3{s, s, s})

		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{.7, .7, .7, .7})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		//h := b.Dy() / 2.0
		//r.Clear(image.Rect(b.Min.X, h-5, b.Max.X, h), gfx.Color{0, 0, 1, 1})

		// Draw the textured card.
		r.Draw(image.Rect(0, 0, 0, 0), txt.Object, camera)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
