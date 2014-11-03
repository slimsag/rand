package main

import (
	"azul3d.org/v1/chippy"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/gfx/gl2"
	"azul3d.org/v1/gfx/window"
	"azul3d.org/v1/keyboard"
	"azul3d.org/v1/math"
	"azul3d.org/v1/ntree"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
)

func randVec3() math.Vec3 {
	return math.Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func random(size, posScale float64) (r gfx.Bounds) {
	f := func() float64 {
		return (rand.Float64() * 2.0) - 1.0
	}

	r.Min = math.Vec3{
		f() * size,
		f() * size,
		f() * size,
	}
	r.Max = r.Min.Add(math.Vec3{
		rand.Float64() * size,
		rand.Float64() * size,
		rand.Float64() * size,
	})

	// Random position
	pos := math.Vec3{f(), f(), f()}
	pos = pos.MulScalar(posScale)
	r.Max = r.Max.Add(pos)
	r.Min = r.Min.Add(pos)
	return r
}

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	w.SetSize(640, 640)
	w.SetPositionCenter(chippy.DefaultScreen())
	glr := r.(*gl2.Renderer)
	glr.UpdateBounds(image.Rect(0, 0, 640, 640))

	// Create a camera.

	// Wait for the shader to load (not strictly required).
	onLoad := make(chan *gfx.Shader, 1)
	r.LoadShader(Wireframe, onLoad)
	go func() {
		<-onLoad
		Wireframe.RLock()
		if Wireframe.Loaded {
			fmt.Println("Shader loaded")
		} else {
			fmt.Println(string(Wireframe.Error))
		}
		Wireframe.RUnlock()
	}()

	// Create a triangle object.
	triangle := gfx.NewObject()
	triangle.Shader = Wireframe
	triangle.Meshes = []*gfx.Mesh{
		&gfx.Mesh{},
	}
	triangle.Meshes[0].GenerateBary()
	triangle.FaceCulling = gfx.NoFaceCulling
	triangle.AlphaMode = gfx.AlphaToCoverage

	tree := ntree.New()
	level := 0

	go func() {
		events := w.Events()
		for {
			e := <-events
			kev, ok := e.(keyboard.TypedEvent)
			if ok {
				if kev.Rune == ' ' || kev.Rune == 'b' {
					for i := 0; i < 1; i++ {
						s := random(0.1, 0.15)
						if kev.Rune == 'b' {
							s = random(0.1, 0.5)
						}
						tree.Add(s)
					}

					// Create new mesh and ask the renderer to load it.
					newMesh := NTreeMesh(tree, level)
					onLoad := make(chan *gfx.Mesh, 1)
					r.LoadMesh(newMesh, onLoad)
					<-onLoad

					// Swap the mesh.
					triangle.Lock()
					triangle.Meshes[0] = newMesh
					triangle.Unlock()
				}

				if kev.Rune == '1' || kev.Rune == '2' {
					if kev.Rune == '1' {
						level--
					} else {
						level++
					}
					fmt.Println("level", level)

					// Create new mesh and ask the renderer to load it.
					newMesh := NTreeMesh(tree, level)
					onLoad := make(chan *gfx.Mesh, 1)
					r.LoadMesh(newMesh, onLoad)
					<-onLoad

					// Swap the mesh.
					triangle.Lock()
					triangle.Meshes[0] = newMesh
					triangle.Unlock()
				}

				if kev.Rune == 's' || kev.Rune == 'S' {
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
		r.Draw(image.Rect(0, 0, 0, 0), triangle, nil)

		// Render the whole frame.
		r.Render()
	}
}

func main() {
	window.Run(gfxLoop)
}
