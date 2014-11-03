// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"sync"
	"image"

	"azul3d.org/lmath.v1"
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

// Drawable represents any object that can be drawn to a canvas.
type Drawable interface {
	// DrawTo should draw this object to the given canvas, using the given
	// bounding box and camera. For instance one could write:
	//  c.Draw(bounds, o, cam)
	//
	// It is permitted to ignore the bounds and camera parameters and instead
	// substitute them with your own, in this way an object can explicitly
	// control the area it will be drawn to / it's camera.
	DrawTo(c Canvas, bounds image.Rectangle, cam *Camera)
}

// DrawFunc implements the Drawable interface by simply calling this function
// whenever DrawTo is called.
type DrawFunc func(c Canvas, bounds image.Rectangle, cam *Camera)

// DrawTo implements the Drawable interface by simply invoking 
func (d DrawFunc) DrawTo(c Canvas, bounds image.Rectangle, cam *Camera) {
	d(c, bounds, cam)
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

	// The shader program to be used during rendering the object.
	*Shader

	// A slice of meshes which make up the object. The order in which the
	// meshes appear in this slice also affects the order in which they are
	// sent to the graphics card.
	//
	// This is a slice specifically to allow renderer implementations to
	// optimize the number of draw calls that must occur to render
	// consecutively listed meshes here (this allows for 'hardware' instancing
	// support).
	Meshes []*Mesh

	// A slice of textures which are used to texture the meshes of this object.
	// The order in which the textures appear in this slice is also the order
	// in which they are sent to the graphics card.
	Textures []*Texture

	// CachedBounds represents the pre-calculated cached bounding box of this
	// object. Note that the bounds are only calculated once Object.Bounds() is
	// invoked.
	//
	// If you make changes to the vertices of any mesh associated with this
	// object, or if you add / remove meshes from this object, the bounds will
	// not reflect this automatically. Instead, you must clear the cached
	// bounds explicitly:
	//  o.Lock()
	//  o.CachedBounds = nil
	//  o.Unlock()
	//
	// And then simply invoke o.Bounds() again to calculate the bounds again.
	CachedBounds *lmath.Rect3
}

// DrawTo implements the Drawable interface by simply calling:
//  c.Draw(bounds, o, cam)
func (o *Object) DrawTo(c Canvas, bounds image.Rectangle, cam *Camera) {
	c.Draw(bounds, o, cam)
}

// Bounds implements the Boundable interface. The returned bounding box takes
// into account all of the mesh's bounding boxes, transformed into world space.
//
// The bounding box is cached (see o.CachedBounds) so that multiple calls to
// this method are fast. If you make changes to the vertices, or add/remove
// meshes from this object you need to explicitly clear the cached bounds so
// that the next call to Bounds() will calculate the bounding box again:
//  o.Lock()
//  o.CachedBounds = nil
//  o.Unlock()
//
// You do not need to clear the cached bounds if the transform of the object
// has changed (as it is applied after calculation of the bounding box).
//
// This method properly read-locks the object.
func (o *Object) Bounds() lmath.Rect3 {
	var b lmath.Rect3
	o.Lock()
	// Do we have a cached bounding box? If so, use it.
	if o.CachedBounds != nil {
		b = *o.CachedBounds
	} else {
		// Calculate the bounding box then.
		for i, m := range o.Meshes {
			if i == 0 {
				b = m.Bounds()
			} else {
				b = b.Union(m.Bounds())
			}
		}

		// Cache it for later.
		o.CachedBounds = &b
	}
	if o.Transform != nil {
		b.Min = o.Transform.ConvertPos(b.Min, LocalToWorld)
		b.Max = o.Transform.ConvertPos(b.Max, LocalToWorld)
		b = b.Union(b)
	}
	o.Unlock()
	return b
}

// Compare compares this object's state (including shader and textures) against
// the other one and determines if it should sort before the other one for
// state sorting purposes.
//
// The object's read lock must be held for this method to operate safely.
func (o *Object) Compare(other *Object) bool {
	if o == other {
		return true
	}

	// Compare shaders.
	if o.Shader != other.Shader {
		return false
	}

	// Compare textures.
	for i, tex := range o.Textures {
		if other.Textures[i] != tex {
			return false
		}
	}

	// Compare state then.
	return o.State.Compare(other.State)
}

// Copy returns a new copy of this Object. Explicitily not copied is the native
// object. The transform is copied via it's Copy() method. The shader is only
// copied by pointer.
//
// The object's read lock must be held for this method to operate safely.
func (o *Object) Copy() *Object {
	cpy := &Object{
		OcclusionTest: o.OcclusionTest,
		State:         o.State,
		Transform:     o.Transform.Copy(),
		Shader:        o.Shader,
		Meshes:        make([]*Mesh, len(o.Meshes)),
		Textures:      make([]*Texture, len(o.Textures)),
	}
	copy(cpy.Meshes, o.Meshes)
	copy(cpy.Textures, o.Textures)
	return cpy
}

// Reset resets this object to it's default (NewObject) state.
//
// The object's write lock must be held for this method to operate safely.
func (o *Object) Reset() {
	o.NativeObject = nil
	o.OcclusionTest = false
	o.State = DefaultState
	o.Transform = NewTransform()
	o.Shader = nil

	// Nil out each mesh pointer.
	for i := 0; i < len(o.Meshes); i++ {
		o.Meshes[i] = nil
	}
	o.Meshes = o.Meshes[:0]

	// Nil out each texture pointer.
	for i := 0; i < len(o.Textures); i++ {
		o.Textures[i] = nil
	}
	o.Textures = o.Textures[:0]
}

// Destroy destroys this object for use by other callees to NewObject. You must
// not use it after calling this method. This makes an implicit call to
// o.NativeObject.Destroy.
//
// The object's write lock must be held for this method to operate safely.
func (o *Object) Destroy() {
	if o.NativeObject != nil {
		o.NativeObject.Destroy()
	}
	o.Reset()
	objPool.Put(o)
}

var objPool = sync.Pool{
	New: func() interface{} {
		return &Object{
			State:     DefaultState,
			Transform: NewTransform(),
		}
	},
}

// NewObject creates and returns a new object with:
//  o.State == DefaultState
//  o.Transform == DefaultTransform
func NewObject() *Object {
	return objPool.Get().(*Object)
}
