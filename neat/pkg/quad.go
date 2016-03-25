// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgfx

import (
	"azul3d.org/lmath.v1"
)

// QuadCurve represents a single 2D quadratic bezier curve with start, end, and
// a single control point.
type QuadCurve struct {
	Start, Control, End lmath.Vec2
}

// QuadPath represents a path composed of 2D quadratic bezier curves.
type QuadPath []QuadCurve
