// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"testing"

	"azul3d.org/gfx.v1"
	math "azul3d.org/lmath.v1"
)

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
