// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
)

type frustum math.Mat4

func (ff frustum) Intersects(bb gfx.Boundable) bool {
	f := math.Mat4(ff)
	b := bb.Bounds()

	// If any single corner of the 3D rectangle is in the frustum, then it is
	// intersecting.
	var inside bool
	for _, corner := range b.Corners() {
		if _, inside = f.Project(corner); inside {
			return true
		}
	}
	return false
}

func (ff frustum) Contains(bb gfx.Boundable) bool {
	f := math.Mat4(ff)
	b := bb.Bounds()

	// Every corner of the 3D rectangle must be in the frustum -- thus we can
	// say that if no corner is in the frustum, the rectangle is not contained.
	var inside bool
	for _, corner := range b.Corners() {
		if _, inside = f.Project(corner); !inside {
			return false
		}
	}
	return true
}

// Frustum returns a Container usable for searching for objects inside or
// intersecting with the given viewing frustum (projection) matrix, for
// instance:
//  tree.In(Frustum(f), results, stop)
//  tree.Intersect(Frustum(f), results, stop)
//
// The matrix may be composed (e.g. view * projection).
func Frustum(f math.Mat4) Container {
	return frustum(f)
}
