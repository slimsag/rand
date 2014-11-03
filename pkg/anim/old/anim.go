// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anim

import (
	"sync"
)

var animPool sync.Pool

// Anim describes a single skeletal animation, it's name, frame rate, etc. It
// is safe for concurrent access.
type Anim struct {
	access   sync.RWMutex
	name     string
	playRate float64
	frames   []Frame
}

// Name returns the name of this animation.
func (a *Anim) Name() string {
	a.access.RLock()
	name := a.name
	a.access.RUnlock()
	return name
}

// SetPlayRate sets the rate at which this animation plays at, e.g. 24.3-FPS
// would be written as:
//  a.SetPlayRate(24.3)
func (a *Anim) SetPlayRate(rate float64) {
	a.access.Lock()
	a.playRate = rate
	a.access.Unlock()
}

// PlayRate returns the rate at which this animation plays at (e.g. 24.3 FPS).
func (a *Anim) PlayRate() float64 {
	a.access.RLock()
	rate := a.playRate
	a.access.RUnlock()
	return rate
}

// NumFrames returns the number of frames that this animation contains.
func (a *Anim) NumFrames() int {
	a.access.RLock()
	numFrames := len(a.frames)
	a.access.RUnlock()
	return numFrames
}

// Frame returns the nth frame of this animation. To iterate over all the
// frames of this animation one could write:
//  for n := 0; n < a.NumFrames(); n++ {
//      a.Frame(n)
//  }
func (a *Anim) Frame(n int) Frame {
	a.access.RLock()
	frame := a.frames[n]
	a.access.RUnlock()
	return frame
}

// Release releases this animation and all of it's frames so that it may be
// reused by other callers to NewAnim() in an effort to relieve GC pressure.
func (a *Anim) Release() {
	a.access.RLock()
	for _, f := range a.frames {
		f.Release()
	}
	a.access.RUnlock()
	animPool.Put(a)
}

// NewAnim returns a new animation with the given name and set of frames. If
// possible NewAnim will reuse a previously released animation in an effort to
// relieve GC pressure.
func NewAnim(name string, frames []Frame) *Anim {
	i := animPool.Get()
	if i == nil {
		return &Anim{
			name: name,
		}
	}
	a := i.(*Anim)
	a.name = name
	a.playRate = 0
	a.frames = frames
	return a
}
