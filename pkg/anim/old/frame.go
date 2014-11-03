// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anim

import (
	"azul3d.org/v1/gfx"
	"sync"
)

var framePool sync.Pool

// Frame describes a single animation frame. A frame effectively describes the
// transformation of each bone in a skeleton at this frame.
type Frame []*gfx.Transform

// Release releases this frame so that it may be reused by other callers to
// NewFrame() in an effort to relieve GC pressure.
func (f Frame) Release() {
	framePool.Put(f)
}

// NewFrame returns a new frame whose cap is n and whose length is zero.
//
// NewFrame() may re-use frames that were previously released in an effort to
// relieve GC pressure.
func NewFrame(n int) Frame {
	i := framePool.Get()
	if i == nil {
		return make(Frame, n)
	}
	f := i.(Frame)
	if cap(f) < n {
		f = append(f, make(Frame, n-cap(f))...)
	}
	return f[:0:n]
}
