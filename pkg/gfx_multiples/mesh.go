// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"sync"
)

// Spatial represents any object that can return it's axis-aligned bounding
// box.
type Spatial interface {
	// Bounds returns the axis-aligned bounding box of this spatial object.
	Bounds() math.Rect3
}

// Bounds is a simple datatype which implements the Spatial interface.
type Bounds math.Rect3

// Bounds implements the Spatial interface.
func (b Bounds) Bounds() math.Rect3 {
	return math.Rect3(b)
}

// TexCoordSet represents a single texture coordinate set for a mesh.
type TexCoordSet struct {
	// The slice of texture coordinates for the set.
	Slice []TexCoord

	// Weather or not the texture coordinates of this set have changed since
	// the last time the mesh was loaded. If set to true the renderer should
	// take note and re-upload the data slice to the graphics hardware.
	Changed bool
}

// NativeMesh represents the native object of a mesh, typically only renderers
// create these.
type NativeMesh Destroyable

// Mesh represents a single mesh made up of several components. A mesh may or
// may not be made up of indexed vertices, etc, depending on whether or not
// len(m.Indices) == 0 holds true.
// In the event that a mesh is indexed, m.Indices holds the indices and it can
// be expected that each other slice (Vertices for instance) will hold at least
// enough elements (or be nil) such that the each index will not be out of
// bounds.
//
// Clients are responsible for utilizing the RWMutex of the mesh when using it
// or invoking methods.
type Mesh struct {
	sync.RWMutex

	// The native object of this mesh. Once loaded the renderer using this mesh
	// must assign a value to this field. Typically clients should not assign
	// values to this field at all.
	NativeMesh

	// Weather or not this mesh is currently loaded or not.
	Loaded bool

	// If true then when this mesh is loaded the sources of it will be kept
	// instead of being set to nil (which allows them to be garbage collected).
	KeepDataOnLoad bool

	// Weather or not the mesh will be dynamically updated. Only used as a hint
	// to increase performence of dynamically updated meshes, does not actually
	// control whether or not a mesh may be dynamically updated.
	Dynamic bool

	// AABB is the axis aligned bounding box of this mesh. There may not be one
	// if AABB.Empty() == true, but one can be calculate using the
	// CalculateBounds() method.
	AABB math.Rect3

	// A slice of indices, if non-nil then this slice contains indices into
	// each other slice (such as Vertices) and this is a indexed mesh.
	// The indices are uint32 (instead of int) for compatability with graphics
	// hardware.
	Indices []uint32

	// Weather or not the indices have changed since the last time the mesh
	// was loaded. If set to true the renderer should take note and
	// re-upload the data slice to the graphics hardware.
	IndicesChanged bool

	// The slice of vertices for the mesh.
	Vertices []Vec3

	// Weather or not the vertices have changed since the last time the
	// mesh was loaded. If set to true the renderer should take note and
	// re-upload the data slice to the graphics hardware.
	VerticesChanged bool

	// The slice of vertex colors for the mesh.
	Colors []Color

	// Weather or not the vertex colors have changed since the last time
	// the mesh was loaded. If set to true the renderer should take note
	// and re-upload the data slice to the graphics hardware.
	ColorsChanged bool

	// A slice of barycentric coordinates for the mesh.
	Bary []Vec3

	// Whether or not the barycentric coordinates have changed since the last
	// time the mesh was loaded. If set to true the renderer should take note
	// and re-upload the data slice to the graphics hardware.
	BaryChanged bool

	// A slice of texture coordinate sets for the mesh, there may be
	// multiple sets which directly relate to multiple textures on a
	// object.
	TexCoords []TexCoordSet
}

// Copy returns a new copy of this Mesh. Depending on how large the mesh is
// this may be an expensive operation. Explicitly not copied over is the native
// mesh, the OnLoad slice, and the loaded and changed statuses (Loaded,
// IndicesChanged, VerticesChanged, etc).
//
// The mesh's read lock must be held for this method to operate safely.
func (m *Mesh) Copy() *Mesh {
	cpy := &Mesh{
		sync.RWMutex{},
		nil,   // Native mesh -- not copied.
		false, // Loaded status -- not copied.
		m.KeepDataOnLoad,
		m.Dynamic,
		m.AABB,
		make([]uint32, len(m.Indices)),
		false, // IndicesChanged -- not copied.
		make([]Vec3, len(m.Vertices)),
		false, // VerticesChanged -- not copied.
		make([]Color, len(m.Colors)),
		false, // ColorsChanged -- not copied.
		make([]Vec3, len(m.Bary)),
		false, // BaryChanged -- not copied.
		make([]TexCoordSet, len(m.TexCoords)),
	}

	copy(cpy.Indices, m.Indices)
	copy(cpy.Vertices, m.Vertices)
	copy(cpy.Colors, m.Colors)
	copy(cpy.Bary, m.Bary)
	for index, set := range m.TexCoords {
		setCpy := TexCoordSet{
			Slice: make([]TexCoord, len(set.Slice)),
		}
		copy(setCpy.Slice, set.Slice)
		cpy.TexCoords[index] = setCpy
	}
	return cpy
}

// CanDraw reports if this mesh is valid for drawing. Cases where a mesh is not
// valid for drawing are as follows:
//  len(m.Vertices) == 0
//  len(m.Colors) != len(m.Vertices)
//  len(m.Bary) != len(m.Vertices)
//  Any m.TexCoord whose len(texCoordSet.Slice) != len(m.Vertices)
func (m *Mesh) CanDraw() bool {
	if len(m.Vertices) == 0 {
		return false
	}
	if len(m.Colors) != len(m.Vertices) {
		return false
	}
	if len(m.Bary) != len(m.Vertices) {
		return false
	}
	for _, tcs := range m.TexCoords {
		if len(tcs.Slice) != len(m.Vertices) {
			return false
		}
	}
	return true
}

// Bounds implements the Spatial interface. It is thread-safe and performs
// locking automatically. If the AABB of this mesh is empty then the bounds are
// calculated.
func (m *Mesh) Bounds() math.Rect3 {
	m.Lock()
	if m.AABB.Empty() {
		m.CalculateBounds()
	}
	bounds := m.AABB
	m.Unlock()
	return bounds
}

// GenerateBary generates the barycentric coordinates for this mesh.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) GenerateBary() {
	var (
		bci = -1
		v   Vec3
	)
	for _ = range m.Vertices {
		// Add barycentric coordinates.
		bci++
		switch bci % 3 {
		case 0:
			v = Vec3{1, 0, 0}
		case 1:
			v = Vec3{0, 1, 0}
		case 2:
			v = Vec3{0, 0, 1}
		}
		m.Bary = append(m.Bary, v)
	}
}

// CalculateBounds calculates a new axis aligned bounding box for this mesh.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) CalculateBounds() {
	var bb math.Rect3
	if len(m.Vertices) > 0 {
		for _, v32 := range m.Vertices {
			v := v32.Vec3()
			bb.Min = bb.Min.Min(v)
			bb.Max = bb.Max.Max(v)
		}
	}
	m.AABB = bb
}

// HasChanged tells if any of the data slices of the mesh are marked as having
// changed.
//
// The mesh's read lock must be held for this method to operate safely.
func (m *Mesh) HasChanged() bool {
	if m.IndicesChanged || m.VerticesChanged || m.ColorsChanged || m.BaryChanged {
		return true
	}
	for _, texCoordSet := range m.TexCoords {
		if texCoordSet.Changed {
			return true
		}
	}
	return false
}

// ClearData sets the data slices of this mesh to nil if m.KeepDataOnLoad is
// set to false.
//
// The mesh's write lock must be held for this method to operate safely.
func (m *Mesh) ClearData() {
	if !m.KeepDataOnLoad {
		m.Indices = nil
		m.Vertices = nil
		m.Colors = nil
		m.Bary = nil
		m.TexCoords = nil
	}
}
