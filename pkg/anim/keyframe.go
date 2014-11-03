// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package anim

import (
	"azul3d.org/v1/gfx"
)

// A key frame is composed of a list of bone poses (i.e. their position,
// rotation, etc) stored as 4x4 matrices, and a list of per-shape morph weights
// in the range of zero to one.
type KeyFrame struct {
	// A list of per-bone poses for this key frame.
	Bones []gfx.Mat4

	// A list of per-shape morph weights, each in the range of zero to one.
	MorphWeights []float32
}

// Copy returns a new 1:1 copy of this key frame.
func (k *KeyFrame) Copy() *KeyFrame {
	cpy := &KeyFrame{
		Bones:        make([]gfx.Mat4, len(k.Bones)),
		MorphWeights: make([]float32, len(k.MorphWeights)),
	}
	copy(cpy.Bones, k.Bones)
	copy(cpy.MorphWeights, k.MorphWeights)
	return cpy
}
