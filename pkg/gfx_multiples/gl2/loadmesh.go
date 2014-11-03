// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl2

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/native/gl"
	"runtime"
	"unsafe"
)

// nativeMesh is stored inside the *Mesh.Native interface and stores vertex
// buffer object ID's.
type nativeMesh struct {
	indices                     uint32
	vertices                    uint32
	colors                      uint32
	bary                        uint32
	texCoords                   []uint32
	verticesCount, indicesCount uint32
	r                           *Renderer
}

func finalizeMesh(n *nativeMesh) {
	n.r.meshesToFree.Lock()
	n.r.meshesToFree.slice = append(n.r.meshesToFree.slice, n)
	n.r.meshesToFree.Unlock()
}

// Implements gfx.Destroyable interface.
func (n *nativeMesh) Destroy() {
	finalizeMesh(n)
}

func (r *Renderer) createVBO() (vboId uint32) {
	// Generate new VBO.
	r.loader.GenBuffers(1, &vboId)
	r.loader.Execute()
	return
}

func (r *Renderer) updateVBO(usageHint int32, dataSize uintptr, dataLength int, data unsafe.Pointer, vboId uint32) {
	// Bind the VBO now.
	r.loader.BindBuffer(gl.ARRAY_BUFFER, vboId)

	// Fill the VBO with the data.
	r.loader.BufferData(
		gl.ARRAY_BUFFER,
		dataSize*uintptr(dataLength),
		data,
		usageHint,
	)
	r.loader.Execute()
}

func (r *Renderer) deleteVBO(vboId *uint32) {
	// Delete the VBO.
	if *vboId == 0 {
		return
	}
	r.loader.DeleteBuffers(1, vboId)
	r.loader.Execute()
	*vboId = 0 // Just for safety.
}

func (r *Renderer) freeMeshes() {
	// Lock the list.
	r.meshesToFree.Lock()

	// Free the meshes.
	for _, native := range r.meshesToFree.slice {
		// Delete single VBO's.
		r.loader.DeleteBuffers(1, &native.indices)
		r.loader.DeleteBuffers(1, &native.vertices)
		r.loader.DeleteBuffers(1, &native.colors)
		r.loader.DeleteBuffers(1, &native.bary)

		// Delete texture coords buffers.
		if len(native.texCoords) > 0 {
			r.loader.DeleteBuffers(uint32(len(native.texCoords)), &native.texCoords[0])
		}

		// Flush and execute OpenGL commands.
		r.loader.Flush()
		r.loader.Execute()
	}

	// Slice to zero, and unlock.
	r.meshesToFree.slice = r.meshesToFree.slice[:0]
	r.meshesToFree.Unlock()
}

func (r *Renderer) LoadMesh(m *gfx.Mesh, done chan *gfx.Mesh) {
	// Lock the mesh until we are done loading it.
	m.Lock()
	if m.Loaded && !m.HasChanged() {
		// Mesh is already loaded and has not changed, signal completion and
		// return after unlocking.
		m.Unlock()
		select {
		case done <- m:
		default:
		}
		return
	}

	f := func() {
		// Find the native mesh, creating a new one if none exists.
		var native *nativeMesh
		if !m.Loaded {
			native = new(nativeMesh)
			native.r = r
		} else {
			native = m.NativeMesh.(*nativeMesh)
		}

		// Determine usage hint.
		usageHint := gl.STATIC_DRAW
		if m.Dynamic {
			usageHint = gl.DYNAMIC_DRAW
		}

		// Update Indices VBO.
		if !m.Loaded || m.IndicesChanged {
			if len(m.Indices) == 0 {
				// Delete indices VBO.
				r.deleteVBO(&native.indices)
			} else {
				if native.indices == 0 {
					// Create indices VBO.
					native.indices = r.createVBO()
				}
				// Update indices VBO.
				r.updateVBO(
					usageHint,
					unsafe.Sizeof(m.Indices[0]),
					len(m.Indices),
					unsafe.Pointer(&m.Indices[0]),
					native.indices,
				)
				native.indicesCount = uint32(len(m.Indices))
			}
			m.IndicesChanged = false
		}

		// Update Vertices VBO.
		if !m.Loaded || m.VerticesChanged {
			if len(m.Vertices) == 0 {
				// Delete vertices VBO.
				r.deleteVBO(&native.vertices)
				native.verticesCount = 0
			} else {
				if native.vertices == 0 {
					// Create vertices VBO.
					native.vertices = r.createVBO()
				}
				// Update vertices VBO.
				r.updateVBO(
					usageHint,
					unsafe.Sizeof(m.Vertices[0]),
					len(m.Vertices),
					unsafe.Pointer(&m.Vertices[0]),
					native.vertices,
				)
				native.verticesCount = uint32(len(m.Vertices))
			}
			m.VerticesChanged = false
		}

		// Update Colors VBO.
		if !m.Loaded || m.ColorsChanged {
			if len(m.Colors) == 0 {
				// Delete colors VBO.
				r.deleteVBO(&native.colors)
			} else {
				if native.colors == 0 {
					// Create colors VBO.
					native.colors = r.createVBO()
				}
				// Update colors VBO.
				r.updateVBO(
					usageHint,
					unsafe.Sizeof(m.Colors[0]),
					len(m.Colors),
					unsafe.Pointer(&m.Colors[0]),
					native.colors,
				)
			}
			m.ColorsChanged = false
		}

		// Update Bary VBO.
		if !m.Loaded || m.BaryChanged {
			if len(m.Bary) == 0 {
				// Delete bary VBO.
				r.deleteVBO(&native.bary)
			} else {
				if native.bary == 0 {
					// Create bary VBO.
					native.bary = r.createVBO()
				}
				// Update bary VBO.
				r.updateVBO(
					usageHint,
					unsafe.Sizeof(m.Bary[0]),
					len(m.Bary),
					unsafe.Pointer(&m.Bary[0]),
					native.bary,
				)
			}
			m.BaryChanged = false
		}

		// Any texture coordinate sets that were removed should have their
		// VBO's deleted.
		deletedMax := len(m.TexCoords)
		if deletedMax > len(native.texCoords) {
			deletedMax = len(native.texCoords)
		}
		deleted := native.texCoords[:deletedMax]
		native.texCoords = native.texCoords[:deletedMax]
		for _, vbo := range deleted {
			r.deleteVBO(&vbo)
		}

		// Any texture coordinate sets that were added should have VBO's
		// created.
		added := m.TexCoords[len(native.texCoords):]
		toUpdate := m.TexCoords
		for _, set := range added {
			vbo := r.createVBO()
			native.texCoords = append(native.texCoords, vbo)

			// Update the VBO.
			r.updateVBO(
				usageHint,
				unsafe.Sizeof(set.Slice[0]),
				len(set.Slice),
				unsafe.Pointer(&set.Slice[0]),
				vbo,
			)
		}

		// And finally, any texture coordinate sets that were changed need to
		// have their VBO's updated.
		for index, set := range toUpdate {
			if set.Changed {
				// Update the VBO.
				r.updateVBO(
					usageHint,
					unsafe.Sizeof(set.Slice[0]),
					len(set.Slice),
					unsafe.Pointer(&set.Slice[0]),
					native.texCoords[index],
				)
			}
		}

		// Ensure no buffer is active when we leave (so that OpenGL state is untouched).
		r.loader.BindBuffer(gl.ARRAY_BUFFER, 0)

		// Flush and execute OpenGL commands.
		r.loader.Flush()
		r.loader.Execute()

		// If the mesh is not loaded, then we need to assign the native mesh
		// and create a finalizer to free the native mesh later.
		if !m.Loaded {
			// Assign the native mesh.
			m.NativeMesh = native

			// Attach a finalizer to the mesh that will later free it.
			runtime.SetFinalizer(native, finalizeMesh)
		}

		// Set the mesh to loaded, clear any data slices if they are not wanted.
		m.Loaded = true
		m.ClearData()

		// Unlock, signal completion, and return.
		m.Unlock()
		select {
		case done <- m:
		default:
		}
		return
	}

	select {
	case r.LoaderExec <- f:
	default:
		go func() {
			r.LoaderExec <- f
		}()
	}
}
