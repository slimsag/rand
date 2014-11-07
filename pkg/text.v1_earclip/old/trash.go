// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
)

func appendSq(m *gfx.Mesh, a gfx.Vec3, s float32) {
	v := func(x, y float32) {
		m.Vertices = append(m.Vertices, gfx.Vec3{x, 0, y})
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
}

func appendLine(m *gfx.Mesh, a, b gfx.Vec3, w float32) {
	v := func(x, y float32) {
		m.Vertices = append(m.Vertices, gfx.Vec3{x, 0, y})
	}

	for t := 0.0; t < 1.0; t += 1.0 / 90.0 {
		av := a.Vec3()
		bv := b.Vec3()
		ip := lmath.Vec3{
			X: lmath.Lerp(av.X, bv.X, t),
			Z: lmath.Lerp(av.Z, bv.Z, t),
		}
		appendSq(m, gfx.ConvertVec3(ip), w)
	}
	return

	v(a.X, a.Z)
	v(b.X, b.Z)
	v(a.X+w, a.Z+w)

	v(b.X, b.Z)
	v(a.X, a.Z)
	v(b.X+w, b.Z+w)

	/*
		left := a.X
		bottom := a.Z
		right := b.X
		top := b.Z

		w := float32(0.007)
		bottom -= w
		top += w
		left -= w
		right += w

		v(left, bottom)
		v(left, top)
		v(right, top)

		v(left, bottom)
		v(right, top)
		v(right, bottom)
	*/
}
