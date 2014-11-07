// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgfx

import (
	"math"

	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
)

// Polygon represents a single 2D polygon composed of all the points in the
// slice. A set of polygons can be tesselated into triangles for display using
// the Tesselator type exposed by this package.
type Polygon []lmath.Vec2

// FIXME: remove
// Bounds calculates and returns the minimum and maximum bounding box points of
// the polygon, p.
func (p Polygon) Bounds() (min, max lmath.Vec2) {
	if len(p) == 0 {
		return
	}
	min.X = math.MaxFloat64
	min.Y = math.MaxFloat64
	max.X = -math.MaxFloat64
	max.Y = -math.MaxFloat64
	for _, v := range p {
		if v.X < min.X {
			min.X = v.X
		}
		if v.X > max.X {
			max.X = v.X
		}

		if v.Y < min.Y {
			min.Y = v.Y
		}
		if v.Y > max.Y {
			max.Y = v.Y
		}
	}
	return
}

// Tess defines a tesselator capable of tesselating 2D polygons into a set of
// triangles suitable for rendering on a GPU.
type Tess struct {
	// The winding rule to use during tesselation.
	WindingRule

	hSweeper
	vSweeper
}

// Tesselate performs tesselation of the polygons, p, appending the resulting
// triangles to the given vertices slice (which may be nil), which is then
// returned (it therefor works similarly to the builtin append function).
//
// Although the polygons are only 2D, the returned set of vertices is 3D with
// the Y component being unused (that is, X is horizontal and Z is vertical),
// so that data can be submitted to the GPU.
func (t *Tess) Tesselate(p []Polygon, vertices []gfx.Vec3) []gfx.Vec3 {
	var prev, hits []lmath.Vec2
	for _, poly := range p {
		// Copy the polygon into the sweeper, expanding if needed.
		copy(t.hSweeper, poly)
		if len(t.hSweeper) < len(poly) {
			need := len(poly) - len(t.hSweeper)
			t.hSweeper = append(t.hSweeper, poly[len(poly)-need:]...)
		}

		// Sweep the points from left-to-right horizontally.
		t.hSweeper.sweep()

		// For every sorted point left-to-right.
		for _, lr := range t.hSweeper {
			// The sweep line goes from the bottom to the top. We test it for
			// intersections against the polygon.
			sweepLine := lineSegment{
				Start: lmath.Vec2{X: lr.X, Y: -math.MaxFloat64},
				End:   lmath.Vec2{X: lr.X, Y: math.MaxFloat64},
			}

			// Find all of the intersection points between the vertical sweep
			// line and the polygon's connected vertices (forming lines).
			hits = t.collect(sweepLine, poly, hits[:0])

			// Triangulate the hit (intersection) buffer.
			vertices = t.triangulate(prev, hits, vertices)

			// Swap the hit buffers to avoid a copy.
			prev, hits = hits, prev
		}
	}
	return vertices
}

// collect compiles a list of all intersection points between the sweep line
// and each line segment defined by every two vertices of the given polygon.
//
// The results are appended to the given buffer and returned.
func (t *Tess) collect(sweepLine lineSegment, p Polygon, buf []lmath.Vec2) []lmath.Vec2 {
	for i := 1; i < len(p); i += 2 {
		hitPoint, status := sweepLine.intersect(lineSegment{
			Start: p[i],
			End:   p[i-1],
		})
		if status == hit {
			buf = append(buf, hitPoint)
		}
	}
	return buf
}

// triangulate performs triangulation of the given points, where prev is the
// previous set of intersection points, and buf is the current set of
// intersection points with the sweep line.
func (t *Tess) triangulate(prev, buf []lmath.Vec2, vertices []gfx.Vec3) []gfx.Vec3 {
	// Copy this hit points into the vertical sweeper.
	t.vSweeper = append(t.vSweeper, prev...)
	t.vSweeper = append(t.vSweeper, buf...)

	// Sweep the points from top-to-bottom vertically.
	t.vSweeper.sweep()

	for _, p := range t.vSweeper {
		vertices = append(vertices, gfx.Vec3{
			X: float32(p.X),
			Z: float32(p.Y),
		})
	}
	return vertices
}

// NewTess returns a new initialized tesselator. The returned tesselator has
// a winding rule of non-zero.
func NewTess() *Tess {
	return &Tess{
		WindingRule: NonZero,
	}
}
