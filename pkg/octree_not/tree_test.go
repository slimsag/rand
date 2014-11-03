// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"math/rand"
	"testing"
)

func random() (r gfx.Bounds) {
	f := func() float64 {
		return (rand.Float64() * 2.0) - 1.0
	}

	size := .01
	posScale := .5

	r.Min = math.Vec3{
		f() * size,
		f() * size,
		f() * size,
	}
	r.Max = r.Min.Add(math.Vec3{
		rand.Float64() * size,
		rand.Float64() * size,
		rand.Float64() * size,
	})

	// Random position
	pos := math.Vec3{f(), f(), f()}
	pos = pos.MulScalar(posScale)
	r.Max = r.Max.Add(pos)
	r.Min = r.Min.Add(pos)
	return r
}

func offset(r gfx.Bounds) gfx.Bounds {
	f := func() float64 {
		return (rand.Float64() * 2.0) - 1.0
	}

	scale := 0.000001
	v := math.Vec3{f(), f(), f()}
	v = v.MulScalar(scale)

	r.Min = r.Min.Add(v)
	r.Max = r.Max.Add(v)
	return r
}

func TestNumObjects(t *testing.T) {
	tree := New()
	for i := 0; i < 500; i++ {
		tree.Add(random())
	}
	if tree.NumObjects() != 500 {
		t.Log("NumObjects", tree.NumObjects)
		t.Log("want 500")
		t.Fail()
	}
	for i := 0; i < 250; i++ {
		tree.Add(random())
	}
	if tree.NumObjects() != 750 {
		t.Log("NumObjects", tree.NumObjects)
		t.Log("want 750")
		t.Fail()
	}

	// Just for sanity:
	if tree.NumNodes() <= 1 {
		t.Log("NumNodes incorrect..", tree.NumNodes())
		t.Log("want > 2")
		t.Fail()
	}
}

func TestRemoval(t *testing.T) {
	tree := New()
	objs := make([]gfx.Boundable, 512)
	for i := 0; i < len(objs); i++ {
		o := random()
		objs[i] = o
		tree.Add(o)
	}
	fail := 0
	for _, o := range objs {
		if !tree.Remove(o) {
			fail++
		}
	}
	if fail > 0 {
		t.Logf("Failed to remove %d/%d objects from the tree.", fail, len(objs))
		t.Fail()
	}
}

// Benchmarks the cost of updating a single boundable in the octree by adding,
// removing, then adding the new version of it.
// This is to see how much faster the Update method is in comparison.
var benchObjs []gfx.Boundable

const nUpdateObjects = 300000

func BenchmarkUpdateDumb(b *testing.B) {
	objs := benchObjs
	for len(objs) < nUpdateObjects {
		objs = append(objs, random())
	}
	benchObjs = objs

	tree := New()
	for _, o := range objs {
		tree.Add(o)
	}
	b.ResetTimer()

	fail := 0
	success := 0
	for n := 0; n < b.N; n++ {
		oldIndex := n % len(objs)
		oldObj := objs[oldIndex]
		newObj := random()
		if !tree.Remove(oldObj) {
			fail++
		} else {
			success++
		}
		tree.Add(newObj)
		objs[oldIndex] = newObj
	}
	if fail > 0 {
		b.Logf("Removal failure ratio: fail=%d success=%d.", fail, success)
		b.Fail()
	}
}

// Ideally should perform just as well generally as BenchmarkUpdateDumb. This
// is where the new position for each boundable cannot benifit from temporal
// coherence because it was moved to a new (random) location likely outside
// it's octant.
func BenchmarkUpdateWorst(b *testing.B) {
	objs := benchObjs
	for len(objs) < nUpdateObjects {
		objs = append(objs, random())
	}
	benchObjs = objs

	tree := New()
	for _, o := range objs {
		tree.Add(o)
	}

	b.ResetTimer()

	fail := 0
	success := 0
	for n := 0; n < b.N; n++ {
		oldIndex := n % len(objs)
		oldObj := objs[oldIndex]
		newObj := random()
		if !tree.Update(oldObj, newObj) {
			fail++
		} else {
			success++
		}
		objs[oldIndex] = newObj
	}
	if fail > 0 {
		b.Logf("Update failure ratio: fail=%d success=%d.", fail, success)
		b.Fail()
	}
}

// Ideally should perform much better than both BenchmarkUpdateDumb and
// BenchmarkUpdateWorst.
func BenchmarkUpdateFast(b *testing.B) {
	objs := benchObjs
	for len(objs) < nUpdateObjects {
		objs = append(objs, random())
	}
	benchObjs = objs

	tree := New()
	for _, o := range objs {
		tree.Add(o)
	}

	b.ResetTimer()

	fail := 0
	success := 0
	for n := 0; n < b.N; n++ {
		oldIndex := n % len(objs)
		oldObj := objs[oldIndex]
		newObj := offset(oldObj.(gfx.Bounds))
		if !tree.Update(oldObj, newObj) {
			fail++
		} else {
			success++
		}
		objs[oldIndex] = newObj
	}
	if fail > 0 {
		b.Logf("Update failure ratio: fail=%d success=%d.", fail, success)
		b.Fail()
	}
}
