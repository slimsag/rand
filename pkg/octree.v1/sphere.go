// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
)

type sphere lmath.Sphere

func (s sphere) Intersects(b gfx.Boundable) bool {
	return lmath.Sphere(s).OverlapsRect3(b.Bounds())
}

func (s sphere) Contains(b gfx.Boundable) bool {
	return b.Bounds().InSphere(lmath.Sphere(s))
}

// Sphere returns a Container usable for searching for objects inside or
// intersecting with the given 3D sphere, for instance:
//  tree.In(Sphere(1.0, math.Vec3{0, 0, 0}), results, stop)
//  tree.Intersect(Sphere(1.0, math.Vec3{0, 0, 0}), results, stop)
func Sphere(s lmath.Sphere) Container {
	return sphere(s)
}
