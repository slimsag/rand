package gridex

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"math/rand"
	"testing"
)

func random() gfx.Spatial {
	o := math.Vec3{
		rand.Float64() * float64(rand.Int()),
		rand.Float64() * float64(rand.Int()),
		rand.Float64() * float64(rand.Int()),
	}
	min := math.Vec3{rand.Float64(), rand.Float64(), rand.Float64()}
	max := min.Add(math.Vec3{rand.Float64(), rand.Float64(), rand.Float64()})
	min = min.Add(o)
	max = max.Add(o)
	return gfx.Bounds{min, max}
}

func TestNearestChunks(t *testing.T) {
	// Test searching against 1,000,000 random spatials in a non-unique way,
	// (i.e. there may often be some (mostly 2x) duplicate searches for every
	// spatial object, but it costs less memory than a unique search).
	g := New(16)
	n := 1000000
	for i := 0; i < n; i++ {
		g.Add(random())
	}

	c := 0
	g.NearestChunks(math.Vec3Zero, func(i int) bool {
		c++
		return true
	})
	nChunks := 16 * 16 * 16
	if c != nChunks {
		t.Log("chunk count", c, "want", nChunks)
		t.Fail()
	}
}

func TestNearest(t *testing.T) {
	// Test searching against 1,000,000 random spatials in a non-unique way,
	// (i.e. there may often be some (mostly 2x) duplicate searches for every
	// spatial object, but it costs less memory than a unique search).
	g := New(2)
	n := 10000
	for i := 0; i < n; i++ {
		g.Add(random())
	}

	last := 0.0
	p := math.Vec3Zero
	k := g.Nearest(p, func(s gfx.Spatial) bool {
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
