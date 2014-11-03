package main

import (
	"azul3d.org/v1/chippy"
	"azul3d.org/v1/chippy/keyboard"
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/gfx/gl2"
	"azul3d.org/v1/gfx/window"
	"azul3d.org/v1/math"
	"fmt"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
)

var wireVert = []byte(`
#version 120

attribute vec3 Vertex;
attribute vec4 Color;

varying vec3 vBC;

uniform float rx;

mat4 rotationMatrix(vec3 axis, float angle)
{
    axis = normalize(axis);
    float s = sin(angle);
    float c = cos(angle);
    float oc = 1.0 - c;
    
    return mat4(oc * axis.x * axis.x + c,           oc * axis.x * axis.y - axis.z * s,  oc * axis.z * axis.x + axis.y * s,  0.0,
                oc * axis.x * axis.y + axis.z * s,  oc * axis.y * axis.y + c,           oc * axis.y * axis.z - axis.x * s,  0.0,
                oc * axis.z * axis.x - axis.y * s,  oc * axis.y * axis.z + axis.x * s,  oc * axis.z * axis.z + c,           0.0,
                0.0,                                0.0,                                0.0,                                1.0);
}

void main() {
	vBC = Color.xyz;
	gl_Position = rotationMatrix(vec3(0, 1, 0), rx) * vec4(Vertex, 1.0);
}
`)

var wireFrag = []byte(`
#version 120
#extension GL_OES_standard_derivatives: enable

varying vec3 vBC;

float edgeFactor() {
	vec3 d = fwidth(vBC);
	vec3 a3 = smoothstep(vec3(0.0), d*1.5, vBC);
	return min(min(a3.x, a3.y), a3.z);
}

void main() {
	if(gl_FrontFacing){
		gl_FragColor = vec4(0.0, 0.0, 0.0, (1.0-edgeFactor())*0.95);
	} else{
		gl_FragColor = vec4(0.0, 0.0, 0.0, (1.0-edgeFactor())*0.6);
	}
	if(gl_FragColor.a < 0.5) {
		discard;
	}
}
`)

func appendCube(m *gfx.Mesh, scale float32) {
	addV := func(x, y, z float32) {
		m.Vertices = append(m.Vertices, gfx.Vec3{x, y, z})
	}

	s := scale

	addV(-s, -s, -s)
	addV(-s, -s, s)
	addV(-s, s, s)

	addV(s, s, -s)
	addV(-s, -s, -s)
	addV(-s, s, -s)

	addV(s, -s, s)
	addV(-s, -s, -s)
	addV(s, -s, -s)

	addV(s, s, -s)
	addV(s, -s, -s)
	addV(-s, -s, -s)

	addV(-s, -s, -s)
	addV(-s, s, s)
	addV(-s, s, -s)

	addV(s, -s, s)
	addV(-s, -s, s)
	addV(-s, -s, -s)

	addV(-s, s, s)
	addV(-s, -s, s)
	addV(s, -s, s)

	addV(s, s, s)
	addV(s, -s, -s)
	addV(s, s, -s)

	addV(s, -s, -s)
	addV(s, s, s)
	addV(s, -s, s)

	addV(s, s, s)
	addV(s, s, -s)
	addV(-s, s, -s)

	addV(s, s, s)
	addV(-s, s, -s)
	addV(-s, s, s)

	addV(s, s, s)
	addV(-s, s, s)
	addV(s, -s, s)
}

func octreeMesh(octree *gfx.Octant) *gfx.Mesh {
	m := new(gfx.Mesh)

	// Slice to zero.
	m.Vertices = m.Vertices[:0]
	m.Colors = m.Colors[:0]

	bci := -1
	nextBC := func() gfx.Color {
		bci++
		switch bci % 3 {
		case 0:
			return gfx.Color{1, 0, 0, 0}
		case 1:
			return gfx.Color{0, 1, 0, 0}
		case 2:
			return gfx.Color{0, 0, 1, 0}
		}
		panic("never here.")
	}

	var add func(o *gfx.Octant)
	add = func(o *gfx.Octant) {
		// Add vertices.
		before := len(m.Vertices)
		appendCube(m, float32(o.AABB.Size().X/2.0))
		center := o.AABB.Center()
		for i, v := range m.Vertices[before:] {
			vert := v.Vec3()
			vert = vert.Add(center)
			m.Vertices[before+i] = gfx.ConvertVec3(vert)
		}
		fmt.Println("Octant", o.Depth, "-", len(o.Objects), "objects")

		for s := range o.Objects {
			before := len(m.Vertices)
			appendCube(m, float32(s.AABB().Size().X/2.0))
			center := s.AABB().Center()
			//fmt.Println(center, s.AABB())
			for i, v := range m.Vertices[before:] {
				vert := v.Vec3()
				vert = vert.Add(center)
				m.Vertices[before+i] = gfx.ConvertVec3(vert)
			}
		}

		for _, octant := range o.Octants {
			if octant != nil {
				add(octant)
			}
		}
	}
	add(octree)

	for _ = range m.Vertices {
		// Add barycentric coordinates.
		m.Colors = append(m.Colors, nextBC())
	}
	return m
}

func randomAABB(size, posScale float64) (r gfx.AABB) {
	f := func() float64 {
		return (rand.Float64() * 2.0) - 1.0
	}

	max := f() * size
	min := f() * size
	r.Max = math.Vec3{max, max, max}
	r.Min = math.Vec3{
		max - min,
		max - min,
		max - min,
	}

	// Center
	center := r.Center()
	r.Min = r.Min.Sub(center)
	r.Max = r.Max.Sub(center)

	// Random position
	pos := math.Vec3{f(), f(), f()}
	pos = pos.MulScalar(posScale)
	r.Max = r.Max.Add(pos)
	r.Min = r.Min.Add(pos)
	return r
}

type mySpatial struct {
	aabb gfx.AABB
}

func (m mySpatial) AABB() gfx.AABB { return m.aabb }

// gfxLoop is responsible for drawing things to the window. This loop must be
// independent of the Chippy main loop.
func gfxLoop(w *chippy.Window, r gfx.Renderer) {
	w.SetSize(640, 640)
	w.SetPositionCenter(chippy.DefaultScreen())
	glr := r.(*gl2.Renderer)
	glr.UpdateBounds(image.Rect(0, 0, 640, 640))

	// Create a perspective viewing frustum matrix.
	width, height := 640, 640
	aspectRatio := float64(width) / float64(height)
	viewMat := gfx.ConvertMat4(math.Mat4Perspective(75.0, aspectRatio, 0.001, 1000.0))

	// Create a camera.
	camera := &gfx.Camera{
		Object:  new(gfx.Object),
		Frustum: viewMat,
	}
	_ = camera

	// Create a wireframe shader.
	shader := &gfx.Shader{
		Name:     "wireframe shader",
		GLSLVert: wireVert,
		GLSLFrag: wireFrag,
		Inputs:   make(map[string]interface{}),
	}

	// Wait for the shader to load (not strictly required).
	onLoad := make(chan *gfx.Shader, 1)
	r.LoadShader(shader, onLoad)
	<-onLoad
	if shader.Loaded {
		fmt.Println("Shader loaded")
	} else {
		fmt.Println(string(shader.Error))
	}

	// Create a triangle object.
	triangle := &gfx.Object{
		Shader: shader,
		State:  gfx.DefaultState,
		Meshes: []*gfx.Mesh{
			&gfx.Mesh{
				Vertices: []gfx.Vec3{
					// Top
					{0, .9, 0},
					{-.9, -.9, 0},
					{.9, -.9, 0},
				},
				Colors: []gfx.Color{
					// Top
					{1, 0, 0, 0},
					{0, 1, 0, 0},
					{0, 0, 1, 0},
				},
			},
		},
	}
	triangle.State.FaceCulling = gfx.NoFaceCulling

	octree := gfx.NewOctree()

	update := make(chan *gfx.Object)
	go func() {
		events := w.Events()
		for {
			e := <-events
			kev, ok := e.(keyboard.TypedEvent)
			if ok {
				if kev.Rune == 'f' || kev.Rune == 'b' {
					for i := 0; i < 100; i++ {
						s := mySpatial{
							aabb: randomAABB(0.1, 0.25),
						}
						if kev.Rune == 'b' {
							s = mySpatial{
								aabb: randomAABB(0.1, 0.5),
							}
						}
						//s.aabb.Min.Z = -.1
						//s.aabb.Max.Z = .1
						//octree = gfx.NewOctree()
						octree.Add(s)
						//fmt.Println(s.AABB().Center(), s.AABB())
					}

					// Create new mesh and ask the renderer to load it.
					newMesh := octreeMesh(octree)
					onLoad := make(chan *gfx.Mesh, 1)
					r.LoadMesh(newMesh, onLoad)
					<-onLoad

					// Take ownership of the triangle.
					<-update

					// Swap the mesh.
					triangle.Meshes[0] = newMesh

					// Give back ownership.
					update <- triangle
				}

				if kev.Rune == 'q' {
					// Take ownership of the triangle.
					<-update

					// Update rotation.
					v, ok := triangle.Shader.Inputs["rx"]
					rx := float32(0.0)
					if ok {
						rx = v.(float32)
					}
					rx += 0.1
					triangle.Shader.Inputs["rx"] = rx

					// Give back ownership.
					update <- triangle
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

	triangleDrawn := make(chan *gfx.Object, 1)
	for {
		// Clear the entire area (empty rectangle means "the whole area").
		r.Clear(image.Rect(0, 0, 0, 0), gfx.Color{1, 1, 1, 1})
		r.ClearDepth(image.Rect(0, 0, 0, 0), 1.0)

		// See if someone else needs ownership of the triangle before we draw.
		select {
		case update <- triangle:
			// Wait for them to give ownership back.
			<-update
		default:
		}

		// Draw the triangle to the screen.
		r.Draw(image.Rect(0, 0, 0, 0), triangle, triangleDrawn)

		// Render the whole frame.
		r.Render()

		select {
		case <-triangleDrawn:
			// Allow updates to the triangle if needed.
			select {
			case update <- triangle:
				<-update
			default:
			}
		}
	}
}

func main() {
	window.Run(gfxLoop)
}
