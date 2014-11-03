// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"sort"
	"testing"
)

// ByDist sorts a list of objects based on their distance away from a target
// point.
type byDist struct {
	// The list of objects to sort.
	Objects []gfx.Boundable

	// The target position to compare against. The list is sorted based off
	// each object's distance away from this position (typically this is the
	// camera's position).
	Target math.Vec3
}

// Implements sort.Interface.
func (b byDist) Len() int {
	return len(b.Objects)
}

// Implements sort.Interface.
func (b byDist) Swap(i, j int) {
	b.Objects[i], b.Objects[j] = b.Objects[j], b.Objects[i]
}

// Implements sort.Interface.
func (b byDist) Less(ii, jj int) bool {
	i := b.Objects[ii]
	j := b.Objects[jj]

	// Calculate the distance from each object to the target position.
	iDist := i.Bounds().SqDistToPoint(b.Target)
	jDist := j.Bounds().SqDistToPoint(b.Target)

	// If i is further away from j (greater value) then it should sort first.
	return iDist > jDist
}

func BenchmarkClosestLinear1k(b *testing.B) {
	tree := New()

	n := 1000
	p := math.Vec3{-.2, -.2, -.1}
	bd := byDist{
		Objects: make([]gfx.Boundable, 0, n),
		Target:  p,
	}
	for i := 0; i < n; i++ {
		o := random()
		bd.Objects = append(bd.Objects, o)
		tree.Add(o)
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		sort.Sort(bd)
	}
}

func BenchmarkClosest1k(b *testing.B) {
	tree := New()

	n := 1000
	for i := 0; i < n; i++ {
		o := random()
		tree.Add(o)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		results := make(chan gfx.Boundable)
		tree.Closest(nil, results, nil)
		nResults := 0
		for {
			_, ok := <-results
			if !ok {
				break
			}
			nResults++
		}
		if len(results) != n {
			b.Log("results", nResults)
			b.Log("want", n)
			b.Fail()
		}
	}
}

/*
func TestIn5k(t *testing.T) {
	tree := New()

	n := 5000

	r := math.Rect3{
		Min: math.Vec3{-.2, -.2, -.1},
		Max: math.Vec3{.2, .2, .1},
	}

	lookup := make(map[gfx.Boundable]struct{}, n)
	for i := 0; i < n; i++ {
		o := random()
		if o.Bounds().In(r) {
			lookup[o] = struct{}{}
		}
		tree.Add(o)
	}

	results := make(chan gfx.Boundable)
	tree.In(Rect3(r), results, nil)
	nResult := 0
	for {
		result, ok := <-results
		if !ok {
			break
		}
		_, ok = lookup[result]
		if !ok {
			t.Log("Got invalid result", result)
			t.Fail()
		}
		nResult++
	}
	if nResult != len(lookup) {
		t.Log("nResult", nResult, "want", len(lookup))
		t.Fail()
	}

	// Test that the bounds of each node are at least sane.
	var validate func(n *Node, b math.Rect3)
	validate = func(n *Node, b math.Rect3) {
		if !n.bounds.In(b) {
			t.Log("invalid bounds:", n.bounds)
			t.Log("does not fit in:", b)
			t.Fail()
		}
		for c := 0; c < 8; c++ {
			child := n.Child(ChildIndex(c))
			if child == nil {
				continue
			}
			validate(child, n.bounds)
		}
	}
	root := tree.Root()
	validate(root, root.bounds)
}

func benchInTree(n int, b *testing.B) {
	tree := New()

	r := math.Rect3{
		Min: math.Vec3{-.1, -.1, -.1},
		Max: math.Vec3{.1, .1, .1},
	}

	for i := 0; i < n; i++ {
		tree.Add(random())
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		results := make(chan gfx.Boundable, 32)
		tree.In(Rect3(r), results, nil)
		for {
			_, ok := <-results
			if !ok {
				break
			}
		}
	}
}

func BenchmarkInTree1k(b *testing.B) {
	benchInTree(1000, b)
}

func BenchmarkInTree10k(b *testing.B) {
	benchInTree(10000, b)
}

func BenchmarkInTree25k(b *testing.B) {
	benchInTree(25000, b)
}

func BenchmarkInTree50k(b *testing.B) {
	benchInTree(50000, b)
}

func BenchmarkInTree300k(b *testing.B) {
	benchInTree(300000, b)
}
*/
