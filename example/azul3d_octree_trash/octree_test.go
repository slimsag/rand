// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"math/rand"
	"testing"
)

func randomAABB(sz float64) (r AABB) {
	f := func() float64 {
		return rand.Float64() * sz
	}

	r.Max = math.Vec3{f(), f(), f()}
	r.Min = math.Vec3{
		r.Max.X - f(),
		r.Max.Y - f(),
		r.Max.Z - f(),
	}

	return r
}

type mySpatial struct {
	aabb AABB
}

func (m mySpatial) AABB() AABB { return m.aabb }

func TestOctree(t *testing.T) {
	octree := NewOctree()
	a := mySpatial{
		aabb: AABB{
			Min: math.Vec3{-1, -1, -1},
			Max: math.Vec3{1, 1, 1},
		},
	}

	b := mySpatial{
		aabb: AABB{
			Min: math.Vec3{-1, -1, -1},
			Max: math.Vec3{1, 1, 1},
		},
		center: math.Vec3{5, 5, 5},
	}

	// x/2
	c := mySpatial{
		aabb: AABB{
			Min: math.Vec3{-.5, -.5, -.5},
			Max: math.Vec3{1, 1, 1},
		},
		center: math.Vec3{0, 0, 0},
	}

	octree.Add(a)
	octree.Add(b)
	octree.Add(c)

	if !octree.Has(a) {
		t.Fatal("Missing octant a")
	}
	if !octree.Has(b) {
		t.Fatal("Missing octant b")
	}
	if !octree.Has(c) {
		t.Fatal("Missing octant c")
	}

	//for i := 0; i < 1000; i++ {
	//	octree.Add(randomAABB(100))
	//}

	x := AABB{
		Min: math.Vec3{-1, -1, -1},
		Max: math.Vec3{6, 6, 6},
	}
	if !octree.AABB.Equals(x) {
		t.Fail()
	}
}
