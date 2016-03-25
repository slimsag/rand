// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgfx

import (
	"azul3d.org/lmath.v1"
)

// WindingRule is a function that can effectively determine if any arbitrary
// point in space is considered within a solid or empty region of a polygonal
// surface.
//
// The predefined winding rules NonZero and EvenOdd are among the most common,
// and their results are shown in this picture on Wikipedia:
//  http://en.wikipedia.org/wiki/Nonzero-rule#mediaviewer/File:Even-odd_and_non-zero_winding_fill_rules.png
type WindingRule func(p lmath.Vec2) bool

// NonZero implements the non-zero winding rule as outlined at:
//  http://en.wikipedia.org/wiki/Nonzero-rule
func NonZero(p lmath.Vec2) bool {
	panic("not implemented")
}

// EvenOdd implements the even-odd winding rule as outlined at:
//  http://en.wikipedia.org/wiki/Even%E2%80%93odd_rule
func EvenOdd(p lmath.Vec2) bool {
	panic("not implemented")
}
