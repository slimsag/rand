// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"testing"
)

func TestRTree(t *testing.T) {
	tree := New(2, 4)

	mk := func(min, max float64) gfx.Spatial {
		return gfx.Bounds{
			Min: math.Vec3{min, min, min},
			Max: math.Vec3{max, max, max},
		}
	}

	a := mk(4, 5)
	b := mk(2, 3)
	c := mk(1, 2)
	d := mk(-1, 1)
	e := mk(-2, -1)
	f := mk(-3, -2)

	tree.Insert(a)
	tree.Insert(b)
	tree.Insert(c)
	tree.Insert(d)
	tree.Insert(e)
	tree.Insert(f)

	if tree.Size() != 6 {
		t.Log("tree.Size() reports incorrect size(1)!")
		t.Fail()
	}

	if !tree.Delete(a) {
		t.Log("tree.Delete() has failed(1)!")
	}
	if !tree.Delete(b) {
		t.Log("tree.Delete() has failed(2)!")
	}

	if tree.Size() != 4 {
		t.Log("tree.Size() reports incorrect size(2)!")
		t.Fail()
	}

	t.Log(tree.Depth())
	t.Fail()

	if tree.NearestNeighbor(math.Vec3{-2, -2, -2}) != e {
		t.Log("NearestNeighbor(1) failed!")
		t.Fail()
	}

	if tree.NearestNeighbor(math.Vec3{-3, -3, -2}) != f {
		t.Log("NearestNeighbor(2) failed!")
		t.Fail()
	}

	nb := tree.NearestNeighbors(3, math.Vec3{1, 1, 1})
	if len(nb) != 3 {
		t.Log("NearestNeighbors returned incorrect length.")
		t.Fail()
	}
	if nb[0] != c {
		t.Log("NearestNeighbors(1) failed!")
		t.Fail()
	}
	if nb[1] != d {
		t.Log("NearestNeighbors(2) failed!")
		t.Fail()
	}
	if nb[2] != e {
		t.Log("NearestNeighbors(3) failed!")
		t.Fail()
	}

	r := math.Rect3{
		Min: f.Bounds().Min,
		Max: a.Bounds().Max,
	}
	is := tree.Intersect(r, -1, nil)
	if len(is) != 4 {
		t.Log("Intersect returned incorrect length.")
		t.Fail()
	}
}
