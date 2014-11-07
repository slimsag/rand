// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
)

func appendSq(vertices []gfx.Vec3, a gfx.Vec3, s float32) []gfx.Vec3 {
	v := func(x, y float32) {
		vertices = append(vertices, gfx.Vec3{x, 0, y})
	}

	left := a.X - s
	right := a.X + s
	bottom := a.Z - s
	top := a.Z + s

	v(left, bottom)
	v(left, top)
	v(right, top)

	v(left, bottom)
	v(right, top)
	v(right, bottom)
	return vertices
}

func appendLine(vertices []gfx.Vec3, a, b gfx.Vec3, w float32) []gfx.Vec3 {
	for t := 0.0; t < 1.0; t += 1.0 / 90.0 {
		av := a.Vec3()
		bv := b.Vec3()
		ip := lmath.Vec3{
			X: lmath.Lerp(av.X, bv.X, t),
			Z: lmath.Lerp(av.Z, bv.Z, t),
		}
		vertices = appendSq(vertices, gfx.ConvertVec3(ip), w)
	}
	return vertices

	/*
		v(a.X, a.Z)
		v(b.X, b.Z)
		v(a.X+w, a.Z+w)

		v(b.X, b.Z)
		v(a.X, a.Z)
		v(b.X+w, b.Z+w)
	*/
}

func (m *GlyphMesher) debugPoints(vertices []gfx.Vec3, g QuadGlyph, scale, pointSize float32, onCurve bool) []gfx.Vec3 {
	last := 0
	for _, index := range g.Indices {
		contour := g.Points[last:index]
		for p := 0; p < len(contour); p++ {
			pt := contour[p]
			if (!onCurve && p%2 != 0) || (onCurve && p%2 == 0) {
				v := gfx.Vec3{float32(pt.X) * scale, 0, float32(pt.Y) * scale}
				vertices = appendSq(vertices, v, pointSize)
			}
		}
	}
	return vertices
}
