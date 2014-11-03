// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anim

import (
	"azul3d.org/v1/gfx"
	"sync"
	"time"
)

// Object represents a single animation object. It maintains a graphics
// object which it solely represents, a list of key frames, a current-frame
// index, and various other settings.
//
// The bind-pose for the animation object is *always* frame zero, therefore an
// animation object will always have a key frame at index zero: o.KeyFrame(0).
type Object struct {
	access                     sync.RWMutex
	obj                        *gfx.Object
	dirty                      bool
	keyFrames                  map[int]*KeyFrame
	maxFrame, frame, frameRate int
	playing                    bool
	lastFrameTime              time.Duration
	activeFrame                *KeyFrame
}

// interpFrame finds an interpolated frame for the given index and returns it.
func (o *Object) interpFrame(index int) *KeyFrame {
	var (
		leftFrame, rightFrame   *KeyFrame
		leftWeight, rightWeight float64
		ok                      bool
	)

	// Move left on the timeline until we find a key frame.
	for f := index; f > 0; f-- {
		leftFrame, ok = o.keyFrames[f]
		if ok {
			leftWeight = float64(f) / float64(index)
			break
		}
	}

	// Move right on the timeline until we find a key frame.
	for f := index; f < o.maxFrame; f++ {
		rightFrame, ok = o.keyFrames[f]
		if ok {
			rightWeight = float64(index) / float64(f)
			break
		}
	}

	// Blend the two frames together.
	if !math.Equals(leftWeight+rightWeight, 1.0) {
		panic("interpolated weights are not equal to 1.0")
	}

	/*
		// A list of per-bone poses for this key frame.
		Bones []gfx.Mat4

		// A list of per-shape morph weights, each in the range of zero to one.
		MorphWeights []float32
	*/
	return nil
}

// InterpFrame returns a interpolated frame, such that if you request a frame
// index which does not contain a key frame, one is instead generated from the
// two nearest key frames.
func (o *Object) InterpFrame(index int) *KeyFrame {
	o.access.RLock()
	k := o.interpFrame(index)
	o.access.RUnlock()
	return k
}

// Update moves this animation object foward by the given duration, d. Most
// clients will not use this method directly, but will instead use a Manager
// which performs this operation automatically.
func (o *Object) Update(d time.Duration) {
	o.access.Lock()

	fps := time.Second / time.Duration(o.frameRate)
	sinceLastFrame := o.lastFrameTime + d
	framesPassed := int(sinceLastFrame / fps)
	if framesPassed > 0 {
		o.frame += framesPassed
		o.lastFrameTime = sinceLastFrame % fps
		o.dirty = true
	}

	if o.dirty {
		o.activeFrame = o.interpFrame(o.frame)
	}

	o.access.Unlock()
}

// ActiveFrame returns a copy of the active key frame of this animation object.
func (o *Object) ActiveFrame() *KeyFrame {
	o.access.RLock()
	cpy := o.activeFrame.Copy()
	o.access.RUnlock()
	return cpy
}

// SetKeyFrame sets (or replaces) the key frame at the given index. It does not
// make a copy of the key frame, so you may wish to pass in a copy yourself:
//  o.SetKeyFrame(index, k.Copy())
func (o *Object) SetKeyFrame(index int, k *KeyFrame) {
	// FIXME: expansion and de-expansion
	o.access.Lock()
	o.keyFrames[index] = k
	o.access.Unlock()
}

// KeyFrame returns a copy of the key frame at the given index, or nil if there
// is no key frame at the given index.
func (o *Object) KeyFrame(index int) *KeyFrame {
	o.access.RLock()
	cpy := o.keyFrames[index]
	if cpy != nil {
		cpy = cpy.Copy()
	}
	o.access.RUnlock()
	return cpy
}

// NumFrames returns the number of frames this animation contains. Zero is the
// first frame, and the last frame (o.NumFrames()-1) is the last keyframe of
// the animation.
func (o *Object) NumFrames() int {
	o.access.RLock()
	v := o.maxFrame
	o.access.RUnlock()
	return v
}

// SetFrame sets the current frame of this animation, any number in the range
// [0, o.NumFrames()] is valid, where zero is the first frame. Numbers outside
// that range are clamped.
func (o *Object) SetFrame(i int) {
	o.access.Lock()
	if i < 0 {
		i = 0
	} else if i > o.maxFrame {
		i = o.maxFrame
	}
	if o.frame != i {
		o.frame = i
		o.dirty = true
	}
	o.access.Unlock()
}

// Frame tells what frame this animation is currently on.
func (o *Object) Frame() int {
	o.access.RLock()
	v := o.frame
	o.access.RUnlock()
	return v
}

// SetFrameRate sets the rate at which the animation plays, measured in frames
// per second. Negative values are accepted and cause the animation to play in
// reverse.
func (o *Object) SetFrameRate(fps int) {
	o.access.Lock()
	o.frameRate = fps
	o.access.Unlock()
}

// FrameRate returns the rate at which the animation plays, measured in frames
// per second. Negative values can be returned and effectively cause the
// animation to play in reverse.
func (o *Object) FrameRate() int {
	o.access.RLock()
	fps := o.frameRate
	o.access.RUnlock()
	return fps
}

// SetPlaying sets whether this animation should be playing or not. It will
// continue off from where it was last paused.
func (o *Object) SetPlaying(playing bool) {
	o.access.Lock()
	if o.playing != playing {
		o.playing = playing
		o.dirty = true
	}
	o.access.Unlock()
}

// Playing tells if this animation is currently playing.
func (o *Object) Playing() bool {
	o.access.RLock()
	v := o.playing
	o.access.RUnlock()
	return v
}

// Play is short-handed for:
//  o.SetPlaying(true)
func (o *Object) Play() {
	o.SetPlaying(true)
}

// Pause is short-handed for:
//  o.SetPlaying(false)
func (o *Object) Pause() {
	o.SetPlaying(false)
}

// Stop is short-handed for:
//  o.SetPlaying(false)
//  o.SetFrame(0)
func (o *Object) Stop() {
	o.SetPlaying(false)
	o.SetFrame(0)
}

// New returns a new animation object given a graphics object. The returned
// animation object solely represents the given graphics object.
//
// The key frame map is not copied, instead it is directly used as the internal
// map of key frames.
//
// The returned animation object has the defaults of:
//  o.SetFrame(0)
//  o.SetPlaying(true)
//  o.SetFrameRate(24)
func New(obj *gfx.Object, keyFrames map[int]*KeyFrame) *Object {
	// Scan the map of key frames for the maximum frame number.
	maxFrame := 0
	for frameNumber := range keyFrames {
		if frameNumber > maxFrame {
			maxFrame = frameNumber
		}
	}

	o := &Object{
		obj:       obj,
		dirty:     true,
		playing:   true,
		maxFrame:  maxFrame,
		frameRate: 24,
		keyFrames: keyFrames,
	}
	return o
}
