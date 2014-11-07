// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
)

type rect lmath.Rect3

func (c rect) Intersects(b gfx.Boundable) bool {
	return lmath.Rect3(c).Overlaps(b.Bounds())
}
func (c rect) Contains(b gfx.Boundable) bool {
	return b.Bounds().In(lmath.Rect3(c))
}

// Rect3 returns a Container usable for searching for objects inside or
// intersecting with the given 3D rectangle, for instance:
//  tree.In(Rect3(r), results, stop)
//  tree.Intersect(Rect3(r), results, stop)
func Rect3(r lmath.Rect3) Container {
	return rect(r)
}
