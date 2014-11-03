// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collision

import (
	"azul3d.org/v1/math"
	"fmt"
)

// Plane3 represents a single 3D plane specified by a distance from the origin
// and a surface normal.
//
// Conceptually a plane can be thought of as a flat surface extending
// infinitely into space in all directions.
type Plane3 struct {
	// The distance from the origin.
	Dist float64

	// The normal vector of the plane.
	Normal math.Vec3
}

// String returns a string representation of the plane p.
func (p Plane3) String() string {
	return fmt.Sprintf("Plane3(Pos=%v, Normal=%v)", p.Pos, p.Normal)
}

// Dist calculates the signed distance of q to the plane p.
//
// Implemented as described in:
//  Real-Time Collision Detection, 5.1.1 "Closest Point on Plane to Point".
func (p Plane3) Dist(q math.Vec3) float64 {
	return q.Dot(p.Normal) - p.Pos
}

// Closest returns the closest point on the plane p nearest to the point q.
//
// Implemented as described in:
//  Real-Time Collision Detection, 5.1.1 "Closest Point on Plane to Point".
func (p Plane3) Closest(q math.Vec3) math.Vec3 {
	t := p.Normal.Dot(q) - p.Pos
	return q.Sub(t).Mul(p.Normal)
}
