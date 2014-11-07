// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import "azul3d.org/lmath.v1"

// LineSegment represents a 2D line segment composed of a start and end point.
type lineSegment struct {
	Start, End lmath.Vec2
}

func (l lineSegment) sideOfPoint(p lmath.Vec2) float64 {
	ax := l.Start.X
	bx := l.End.X
	ay := l.Start.Y
	by := l.End.Y
	return (bx-ax)*(p.Y-ay) - (by-ay)*(p.X-ax)
}

// Intersect tests the two line segments for intersection. It returns hit==true
// if the segments intersection, or returns hit==false if there was no
// intersection / the segments are collinear.
//
// The returned point is the intersection point, or a zero vector.
//
// The algorithm used is the one described at:
//  http://processingjs.org/learning/custom/intersect/
func (a lineSegment) intersect(b lineSegment) (p lmath.Vec2, hit bool) {
	x1 := a.Start.X
	y1 := a.Start.Y
	x2 := a.End.X
	y2 := a.End.Y

	x3 := b.Start.X
	y3 := b.Start.Y
	x4 := b.End.X
	y4 := b.End.Y

	// Compute a1, b1, c1, where line joining points 1 and 2
	// is "a1 x + b1 y + c1 = 0".
	a1 := y2 - y1
	b1 := x1 - x2
	c1 := (x2 * y1) - (x1 * y2)

	// Compute r3 and r4.
	r3 := ((a1 * x3) + (b1 * y3) + c1)
	r4 := ((a1 * x4) + (b1 * y4) + c1)

	// Check signs of r3 and r4. If both point 3 and point 4 lie on
	// same side of line 1, the line segments do not intersect.
	if (r3 != 0) && (r4 != 0) && sameSign(r3, r4) {
		return lmath.Vec2Zero, false
	}

	// Compute a2, b2, c2
	a2 := y4 - y3
	b2 := x3 - x4
	c2 := (x4 * y3) - (x3 * y4)

	// Compute r1 and r2
	r1 := (a2 * x1) + (b2 * y1) + c2
	r2 := (a2 * x2) + (b2 * y2) + c2

	// Check signs of r1 and r2. If both point 1 and point 2 lie
	// on same side of second line segment, the line segments do
	// not intersect.
	if (r1 != 0) && (r2 != 0) && sameSign(r1, r2) {
		return lmath.Vec2Zero, false
	}

	// Line segments intersect: compute intersection point.
	denom := (a1 * b2) - (a2 * b1)
	if denom == 0 {
		// Collinear
		return lmath.Vec2Zero, false
	}

	var offset float64
	if denom < 0 {
		offset = -denom / 2
	} else {
		offset = denom / 2
	}

	// Intersection point.
	var x, y float64

	// The denom/2 is to get rounding instead of truncating. It
	// is added or subtracted to the numerator, depending upon the
	// sign of the numerator.
	num := (b1 * c2) - (b2 * c1)
	if num < 0 {
		x = (num - offset) / denom
	} else {
		x = (num + offset) / denom
	}

	num = (a2 * c1) - (a1 * c2)
	if num < 0 {
		y = (num - offset) / denom
	} else {
		y = (num + offset) / denom
	}
	return lmath.Vec2{x, y}, true
}

func sameSign(a, b float64) bool {
	return a*b >= 0
}
