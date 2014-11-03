package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"

	"azul3d.org/chippy.v1"
	"azul3d.org/gfx.v1"
	"azul3d.org/gfx/window.v2"
	"azul3d.org/keyboard.v1"
	"azul3d.org/lmath.v1"
	"azul3d.org/octree.v1"
)

func randVec3() lmath.Vec3 {
	return lmath.Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func random(size, posScale float64) (r gfx.Bounds) {
	f := func() float64 {
		return (rand.Float64() * 2.0) - 1.0
	}

	r.Min = lmath.Vec3{
		f() * size,
		f() * size,
		f() * size,
	}
	r.Max = r.Min.Add(lmath.Vec3{
		rand.Float64() * size,
		rand.Float64() * size,
		rand.Float64() * size,
	})

	// Random position
	pos := lmath.Vec3{f(), f(), f()}
	pos = pos.MulScalar(posScale)
	r.Max = r.Max.Add(pos)
	r.Min = r.Min.Add(pos)
	return r
}

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w window.Window, r gfx.Renderer) {
	// Setup a camera to use a perspective projection.
	camera := gfx.NewCamera()
	camFOV := 60.0
	camNear := 0.0001
	camFar := 1000.0
	camera.SetPersp(r.Bounds(), camFOV, camNear, camFar)

	// Move the camera -2 on the Y axis (back two units away from the triangle
	// object).
	camera.SetPos(lmath.Vec3{0, -2, 0})

	// Create a triangle object.
	triangle := gfx.NewObject()
	triangle.Shader = Wireframe
	triangle.Meshes = []*gfx.Mesh{gfx.NewMesh()}
	triangle.Meshes[0].GenerateBary()
	triangle.FaceCulling = gfx.NoFaceCulling
	triangle.AlphaMode = gfx.AlphaToCoverage

	tree := octree.New()
	level := 0

	go func() {
		event := w.Events()
		for e := range event {
			switch ev := e.(type) {
			case chippy.ResizedEvent:
				// Update the camera's projection matrix for the new width and
				// height.
				camera.Lock()
				camera.SetPersp(r.Bounds(), camFOV, camNear, camFar)
				camera.Unlock()

			case keyboard.TypedEvent:
				if ev.Rune == ' ' || ev.Rune == 'b' {
					for i := 0; i < 100; i++ {
						s := random(0.01, 0.15)
						if ev.Rune == 'b' {
							s = random(0.01, 0.5)
						}
						tree.Add(s)
					}

					// Create new mesh and ask the renderer to load it.
					newMesh := OctreeMesh(tree, level)
					onLoad := make(chan *gfx.Mesh, 1)
					r.LoadMesh(newMesh, onLoad)
					<-onLoad

					// Swap the mesh.
					triangle.Lock()
					triangle.Meshes[0] = newMesh
					triangle.Unlock()
				}

				if ev.Rune == '1' || ev.Rune == '2' {
					if ev.Rune == '1' {
						level--
					} else {
						level++
					}
					fmt.Println("level", level)

					// Create new mesh and ask the renderer to load it.
					newMesh := OctreeMesh(tree, level)
					onLoad := make(chan *gfx.Mesh, 1)
					r.LoadMesh(newMesh, onLoad)
					<-onLoad

					// Swap the mesh.
					triangle.Lock()
					triangle.Meshes[0] = newMesh
					triangle.Unlock()
				}

				if ev.Rune == 's' || ev.Rune == 'S' {
					fmt.Println("Writing screenshot to file...")
					// Download the image from the graphics hardware and save
					// it to disk.
					complete := make(chan image.Image, 1)
					r.Download(image.Rect(0, 0, 0, 0), complete)
					img := <-complete // Wait for download to complete.

					// Save to png.
					f, err := os.Create("screenshot.png")
					if err != nil {
						log.Fatal(err)
					}
					err = png.Encode(f, img)
					if err != nil {
						log.Fatal(err)
					}
					fmt.Println("Wrote texture to screenshot.png")
				}
			}
		}
	}()

	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// Update the rotation.
		dt := r.Clock().Dt()
		triangle.RLock()
		rot := triangle.Transform.Rot()
		if w.Keyboard.Down(keyboard.ArrowLeft) {
			rot.Z += 90 * dt
		}
		if w.Keyboard.Down(keyboard.ArrowRight) {
			rot.Z -= 90 * dt
		}
		if w.Keyboard.Down(keyboard.ArrowUp) {
			rot.X += 20 * dt
		}
		if w.Keyboard.Down(keyboard.ArrowDown) {
			rot.X -= 20 * dt
		}
		triangle.Transform.SetRot(rot)
		triangle.RUnlock()

		// Draw the triangle to the screen.
		r.Draw(image.Rect(0, 0, 0, 0), triangle, camera)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop, nil)
}
