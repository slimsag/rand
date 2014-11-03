package ntree

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

	size := .1
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

func TestInTree500k(t *testing.T) {
	tree := New()

	n := 500000

	r := math.Rect3{
		Min: math.Vec3{-.2, -.2, -.1},
		Max: math.Vec3{.2, .2, .1},
	}

	lookup := make(map[gfx.Spatial]struct{}, n)
	for i := 0; i < n; i++ {
		o := random()
		if o.Bounds().In(r) {
			lookup[o] = struct{}{}
		}
		tree.Add(o)
	}

	results := make(chan gfx.Spatial)
	tree.In(r, results, nil)
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
		for _, c := range n.Children {
			validate(c, n.bounds)
		}
	}
	validate(tree.Root, tree.Root.bounds)
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
		results := make(chan gfx.Spatial)
		tree.In(r, results, nil)
		nResult := 0
		for {
			_, ok := <-results
			if !ok {
				break
			}
			nResult++
		}
	}
}

func BenchmarkInTree10k(b *testing.B) {
	benchInTree(10000, b)
}

func BenchmarkInTree20k(b *testing.B) {
	benchInTree(20000, b)
}

func BenchmarkInTree30k(b *testing.B) {
	benchInTree(30000, b)
}

func BenchmarkInTree40k(b *testing.B) {
	benchInTree(40000, b)
}

func BenchmarkInTree50k(b *testing.B) {
	benchInTree(50000, b)
}

/*
func TestNearestTree(t *testing.T) {
	tree := New()

	n := 1000000
	for i := 0; i < n; i++ {
		tree.Add(random())
	}
	if tree.OutsideCount() > 0 {
		t.Log("outside count > 0")
		t.Fail()
	}

	start := time.Now()

	p := randVec3()
	results := make(chan gfx.Spatial, 32)
	stop := make(chan struct{})
	tree.Nearest(p, results, stop)
	nResult := 0
	lastResult := 0.0
	for {
		r, ok := <-results
		if !ok {
			break
		}
		distToR := r.Bounds().Closest(p).LengthSq()
		t.Log(distToR < lastResult)
		//if lastResult != 0 && distToR < lastResult {
		//	t.Log("last result", lastResult)
		//	t.Log("this result", distToR)
		//	t.Fail()
		//}
		lastResult = distToR
		nResult++
		if nResult == 1000000 {
			stop <- struct{}{}
			break
		}
	}
	if nResult != n/2 {
		//t.Log("nResult", nResult, "want", n)
		//t.Fail()
	}

	t.Log("search time", time.Since(start))
}
*/
