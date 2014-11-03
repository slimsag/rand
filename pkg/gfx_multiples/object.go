// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"sync"
)

// Destroyable defines a destroyable object. Once an object is destroyed it may
// still be used, but typically doing so is not good and would e.g. involve
// reloading the entire object and cause performance issues.
//
// Clients should invoke the Destroy() method when they are done utilizing the
// object or else doing so will be left up to a runtime Finalizer.
type Destroyable interface {
	// Destroy destroys this object. Once destroyed the object can still be
	// used but doing so is not advised for performance reasons (e.g. requires
	// reloading the entire object).
	//
	// This method is safe to invoke from multiple goroutines concurrently.
	Destroy()
}

// NativeObject represents a native graphics object, they are normally only
// created by renderers.
type NativeObject interface {
	Destroyable

	// If the GPU supports occlusion queries (see GPUInfo.OcclusionQuery) and
	// OcclusionTest is set to true on the graphics object, then this method
	// will return the number of samples that passed the depth and stencil
	// testing phases the last time the object was drawn. If occlusion queries
	// are not supported then -1 will be returned.
	//
	// This method is safe to invoke from multiple goroutines concurrently.
	SampleCount() int
}

// Object represents a single graphics object for rendering, it has a
// transformation matrix which is applied to each vertex of each mesh, it
// has a shader program, meshes, and textures used for rendering the object.
//
// Clients are responsible for utilizing the RWMutex of the object when using
// it or invoking methods.
type Object struct {
	sync.RWMutex

	// The native object of this graphics object. The renderer using this
	// graphics object must assign a value to this field after a call to
	// Draw() has finished before unlocking the object.
	NativeObject

	// Whether or not this object should be occlusion tested. See also the
	// SampleCount() method of NativeObject.
	OcclusionTest bool

	// The render state of this object.
	State

	// The transformation of the object.
	*Transform

	// The shader programs to be used when rendering each mesh of this object.
	Shaders []*Shader

	// A slice of meshes which make up the object. The order in which the
	// meshes appear in this slice also affects the order in which they are
	// sent to the graphics card.
	//
	// Each mesh will be drawn with each corresponding shader of this Object.
	Meshes []*Mesh

	// The texture slices which are used to texture each of the meshes of this
	// object. The order in which the textures appear in this slice is also the
	// order in which they are sent to the graphics card.
	Textures [][]*Texture
}

// CanDraw tells if this object can be drawn. Cases where it cannot be drawn
// are where:
//  len(o.Meshes) == 0
//  len(o.Shaders) != len(o.Meshes)
//  len(o.Textures) < len(o.Meshes)
//  Any *Shader who reports !s.CanDraw().
//  Any *Texture who reports !t.CanDraw().
//  Any *Mesh who reports !m.CanDraw().
func (o *Object) CanDraw() bool {
	if len(o.Meshes) == 0 {
		return false
	}
	if len(o.Shaders) != len(o.Meshes) {
		return false
	}
	if len(o.Textures) != len(o.Meshes) {
		return false
	}
	for _, s := range o.Shaders {
		if !s.CanDraw() {
			return false
		}
	}
	for _, texSet := range o.Textures {
		for _, tex := range texSet {
			if !tex.CanDraw() {
				return false
			}
		}
	}
	for _, m := range o.Meshes {
		if !m.CanDraw() {
			return false
		}
	}
	return true
}

// Bounds implements the Spatial interface. The returned bounding box takes
// into account all of the mesh's bounding boxes, transformed into world space.
//
// This method properly read-locks the object.
func (o *Object) Bounds() math.Rect3 {
	var b math.Rect3
	o.RLock()
	for i, m := range o.Meshes {
		if i == 0 {
			b = m.Bounds()
		} else {
			b = b.Union(m.Bounds())
		}
	}
	if o.Transform != nil {
		b.Min = o.Transform.ConvertPos(b.Min, LocalToWorld)
		b.Max = o.Transform.ConvertPos(b.Max, LocalToWorld)
		b = b.Union(b)
	}
	o.RUnlock()
	return b
}

// Compare compares this object's state (including shader and textures) against
// the other one and determines if it should sort before the other one for
// state sorting purposes.
func (o *Object) Compare(other *Object) bool {
	if o == other {
		return true
	}

	// Compare shaders.
	for i, shader := range o.Shaders {
		if shader != other.Shaders[i] {
			return false
		}
	}

	// Compare textures.
	for ts, texSet := range o.Textures {
		for t, tex := range texSet {
			if other.Textures[ts][t] != tex {
				return false
			}
		}
	}

	// Compare state then.
	return o.State.Compare(other.State)
}

// NewObject creates and returns a new object with:
//  o.State == DefaultState
//  o.Transform == DefaultTransform
func NewObject() *Object {
	return &Object{
		State:     DefaultState,
		Transform: NewTransform(),
	}
}
