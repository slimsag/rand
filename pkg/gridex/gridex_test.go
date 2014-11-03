package gridex

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

	size := .013123
	posScale := .5123213

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

func TestLocality(t *testing.T) {
	// Size of the 3D gridex table (sz*sz*sz):
	sz := 2

	// Number of randoms to add for locality testing.
	n := 1000000

	// Margin of error for locality hash:
	margin := 750

	g := New(sz)
	counts := make([]int, sz*sz*sz)
	for i := 0; i < n; i++ {
		idx := g.Index(random().Bounds().Center())
		counts[idx]++
	}
	bestLocality := n / (sz * sz * sz)
	for _, c := range counts {
		if c > (bestLocality+margin) || c < (bestLocality-margin) {
			t.Log(c, "not in margin", bestLocality, "+-", margin)
			t.Fail()
		}
	}
}

func TestIter(t *testing.T) {
	sz := 16
	g := New(sz)
	visited := make([]bool, sz*sz*sz)
	fsz := float64(sz)
	for x := 0.0; x < fsz; x += 1 {
		for y := 0.0; y < fsz; y += 1 {
			for z := 0.0; z < fsz; z += 1 {
				p := math.Vec3{x, y, z}
				idx := g.Index(p)
				if visited[idx] {
					t.Log("Got invalid iter index:", idx, p)
					t.Fail()
				}
				visited[idx] = true
			}
		}
	}

	for i, v := range visited {
		if !v {
			t.Log("Never visited:", i)
			t.Fail()
		}
	}
}

func TestIn(t *testing.T) {
	sz := 4
	g := New(sz)
	n := 10
	r := math.Rect3{
		Min: math.Vec3{-1, -1, -1},
		Max: math.Vec3{1, 1, 1},
	}
	visited := make(map[gfx.Boundable]bool)
	for i := 0; i < n; i++ {
		s := random()
		g.Add(s)
		if s.Bounds().In(r) {
			visited[s] = false
		}
	}

	results := make(chan gfx.Boundable, 32)
	g.In(r, results, nil)

	for {
		result, ok := <-results
		if !ok {
			break
		}
		visited[result] = true
		if !result.Bounds().In(r) {
			t.Log("Bad result", result.Bounds())
			t.Fail()
		}
	}

	for sp, v := range visited {
		if !v {
			t.Log("Never visited:", sp)
			for i, d := range g.Data {
				for _, s := range d {
					if s == sp {
						t.Log("here>", i)
					}
				}
			}
			t.Fail()
		}
	}
}

var addRemoveList []gfx.Boundable

func benchAddRemove(amount int, b *testing.B) {
	g := New(16)
	if len(addRemoveList) > amount {
		addRemoveList = addRemoveList[:amount]
	}
	for n := 0; n < amount; n++ {
		if len(addRemoveList) == amount {
			break
		}
		addRemoveList = append(addRemoveList, random())
	}
	if len(addRemoveList) != amount {
		b.Log("len() reports", len(addRemoveList), "want", amount)
		b.Fail()
		return
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for _, s := range addRemoveList {
			g.Add(s)
			if !g.Remove(s) {
				b.Log("failed to remove")
				b.Fail()
			}
		}
	}
}

func BenchmarkAddRemove1K(b *testing.B) {
	benchAddRemove(1000, b)
}

func BenchmarkAddRemove5K(b *testing.B) {
	benchAddRemove(5000, b)
}

func BenchmarkAddRemove10K(b *testing.B) {
	benchAddRemove(10000, b)
}

var (
	inList []gfx.Boundable
)

func benchIn(n int, b *testing.B) {
	sz := 8
	g := New(sz)
	r := math.Rect3{
		Min: math.Vec3{-1, -1, -1},
		Max: math.Vec3{1, 1, 1},
	}
	if len(inList) > n {
		inList = inList[:n]
	}
	for i := 0; i < n; i++ {
		if len(inList) == n {
			break
		}
		inList = append(inList, random())
	}
	if len(inList) != n {
		b.Log("len() reports", len(inList), "want", n)
		b.Fail()
		return
	}

	visited := make(map[gfx.Boundable]bool)
	for _, s := range inList {
		g.Add(s)
		if s.Bounds().In(r) {
			visited[s] = false
		}
	}
	if len(inList) != n {
		panic("test is invalid")
	}
	if len(visited) == 0 {
		panic("test is invalid")
	}
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		results := make(chan gfx.Boundable, 32)
		g.In(r, results, nil)

		for {
			result, ok := <-results
			if !ok {
				break
			}
			visited[result] = true
			if !result.Bounds().In(r) {
				b.Log("Bad result", result.Bounds())
				b.Fail()
			}
		}

		for i, v := range visited {
			if !v {
				b.Log("Never visited:", i)
				b.Fail()
				return
			}
		}
	}
}

func BenchmarkIn1K(b *testing.B) {
	benchIn(1000, b)
}

func BenchmarkIn5K(b *testing.B) {
	benchIn(5000, b)
}

func BenchmarkIn10K(b *testing.B) {
	benchIn(10000, b)
}

func BenchmarkIn100K(b *testing.B) {
	benchIn(100000, b)
}

func BenchmarkIn250K(b *testing.B) {
	benchIn(250000, b)
}

/*
func TestNearest(t *testing.T) {
	// Test searching against 1,000,000 random spatials in a non-unique way,
	// (i.e. there may often be some (mostly 2x) duplicate searches for every
	// spatial object, but it costs less memory than a unique search).
	g := New(2)
	n := 1000
	for i := 0; i < n; i++ {
		g.Add(random())
	}

	last := 0.0
	p := math.Vec3Zero
	k := g.Nearest(p, func(s gfx.Boundable) bool {
		c := s.Bounds().Center()
		dist := c.Sub(p).LengthSq()
		t.Log(dist, last)
		if dist < last {
			t.Log(dist, last)
			t.Fail()
			return false
		}
		last = dist
		return true
	})
	if k != n {
		t.Log("k, n", k, n)
		t.Fail()
	}
}
*/
