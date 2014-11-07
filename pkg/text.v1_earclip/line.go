// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

// LineSegment represents a 2D line segment composed of a start and end point.
type lineSegment struct {
	Start, End Point
}

func (l lineSegment) sideOfPoint(p Point) int {
	ax := l.Start.X
	bx := l.End.X
	ay := l.Start.Y
	by := l.End.Y
	return (bx-ax)*(p.Y-ay) - (by-ay)*(p.X-ax)
}

// isBetween tells if the point, p is between the two points composing this
// line segment.
func (l lineSegment) isBetween(p Point) bool {
	ax := l.Start.X
	ay := l.Start.Y
	bx := l.End.X
	by := l.End.Y
	cx := p.X
	cy := p.Y
    crossProduct := (cy - ay) * (bx - ax) - (cx - ax) * (by - ay)
    if absInt(crossProduct) != 0 {
		return false
	}

    dotProduct := (cx - ax) * (bx - ax) + (cy - ay) * (by - ay)
	if dotProduct < 0 {
		return false
	}

    squaredLengthba := (bx - ax) * (bx - ax) + (by - ay) * (by - ay)
    if dotProduct > squaredLengthba {
		return false
	}
	return true
}

const (
	noHit int = iota
	collinear
	hit
)

// Intersect tests the two line segments for intersection.
//
// The returned point is the intersection point, or a zero vector.
//
// The algorithm used is the one described at:
//  http://processingjs.org/learning/custom/intersect/
func (a lineSegment) intersect(b lineSegment) (p Point, status int) {
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
		return Point{}, noHit
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
		return Point{}, noHit
	}

	// Line segments intersect: compute intersection point.
	denom := (a1 * b2) - (a2 * b1)
	if denom == 0 {
		// Collinear
		return Point{}, collinear
	}

	var offset int
	if denom < 0 {
		offset = -denom / 2
	} else {
		offset = denom / 2
	}

	// Intersection point.
	var x, y int

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
	return Point{X: x, Y: y}, hit
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func sameSign(a, b int) bool {
	return a*b >= 0
}
