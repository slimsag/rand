package ntree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"math/rand"
	"testing"
)

func randVec3() math.Vec3 {
	return math.Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
}

func random() gfx.Spatial {
	min := randVec3().MulScalar(.2)
	max := min.Add(randVec3().MulScalar(.2))
	o := randVec3()
	min = min.Add(o)
	max = max.Add(o)
	return gfx.Bounds{min, max}
}

func TestInTree(t *testing.T) {
	tree := New()

	n := 100000
	for i := 0; i < n; i++ {
		tree.Add(random())
	}

	r := math.Rect3{
		Min: math.Vec3{-.2, -.2, -.2},
		Max: math.Vec3{.2, .2, .2},
	}
	results := make(chan gfx.Spatial)
	stop := make(chan struct{}, 1)
	tree.In(r, results, stop)
	nResult := 0
	for {
		_, ok := <-results
		if !ok {
			break
		}
		nResult++
		if nResult == n {
			stop <- struct{}{}
			break
		}
	}
	if nResult != n {
		t.Log("nResult", nResult, "want", n)
		t.Fail()
	}
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
