// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anim

import (
	"azul3d.org/v1/clock"
	"azul3d.org/v1/gfx"
	"sync"
)

var (
	animClock = clock.New()
	skelPool  = new(sync.Pool)
)

// Skeleton represents a single skeleton used to animate a graphics object. It
// is safe for concurrent access.
type Skeleton struct {
	access sync.RWMutex
	name   string
	object *gfx.Object
	anims  []*Anim
}

// Play causes the skeleton to play
func (s *Skeleton) Play(start, end int)

// Release releases this skeleton and all of it's animations so that it may be
// reused by other callers to New() in an effort to relieve GC pressure.
func (s *Skeleton) Release() {
	s.access.RLock()
	for _, a := range s.anims {
		a.Release()
	}
	s.access.RUnlock()
	skelPool.Put(s)
}

// New returns a new skeleton which will apply animation to the given graphics
// object by locking the object immedietly and setting it's shader and shader
// inputs appropriately, as needed.
//
// The skeleton will be composed of the given animations.
//
// If possible New() will reuse a previously release skeleton in an effort to
// relieve GC pressure.
func New(name string, o *gfx.Object, anims []*Anim) *Skeleton {
	var s *Skeleton
	i := skelPool.Get()
	if i != nil {
		s = i.(*Skeleton)
	} else {
		s = new(Skeleton)
	}
	s.name = name
	s.object = o
	s.anims = anims
	return s
}

/*

// Skeleton represents a single skeleton used for animating a graphics object's
// meshes.
type Skeleton struct {
	// The graphics object to be updated by this skeleton.
	*gfx.Object

	access sync.RWMutex

	// A list of bones that make up this skeleton. These must not be changed
	// post initialization.
	bones []*gfx.Transform

	// A list of animations for this skeleton.
	anims []*Anim

	// Whether or not the animation is currently playing or not.
	playing bool

	// The current frame number of the active animation.
	frame int

	// The active animation, nil if there is none.
	anim *Anim

	// The active transformation of each bone.
	activeBones []*gfx.Transform

	// The time (measured as time passed since the program started) at which
	// the last frame was played. If enough time has passed then the animation
	// frame will increase.
	lastFrameTime time.Duration

	// The last frame used. If it is different from Frame then the
	// skeleton is updated.
	lastFrame int

	// The last animation used. If it is different from Anim then the
	// skeleton is updated.
	lastAnim *Anim
}

// FindAnim finds and returns the named animation from the list of animations
// held by this skeleton. Returns nil if none is found.
func (s *Skeleton) FindAnim(animName string) *Anim {
	for _, a := range s.anims {
		if a.name == animName {
			return a
		}
	}
	return nil
}

// Update updates the currently active transformation of each bone. It updates
// the shader inputs, etc. This must be called each frame before the skeleton
// is drawn.
func (s *Skeleton) Update() {
	if !s.playing {
		return
	}

	now := animClock.Time()
	update := s.frame != s.lastFrame || s.anim != s.lastAnim
	if !update {
		// Compare frame times then.
		diffSeconds := float64(now - s.lastFrameTime) / float64(time.Second)
		if diffSeconds > s.anim.playRate {
			update = true
		}
	}

	if update {
		// Update the last-known time, animation, and frame number.
		s.lastFrameTime = now
		s.lastAnim = s.anim
		s.lastFrame = s.frame
	}
}

// New returns a new skeleton with the given name.
func New(name string) *Skeleton {
	return &Skeleton{}
}
*/
