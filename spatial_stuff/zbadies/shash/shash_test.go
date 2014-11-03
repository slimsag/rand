package shash

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

func TestZombies(t *testing.T) {
	h := New()
	n := 10
	for i := 0; i < n; i++ {
		h.Add(random())
	}
	t.Log("Spawned", n, "zombies.")

	for i, byAngle := range h.Table {
		for j, byDist := range byAngle {
			t.Log(i, j, len(byDist))
		}
	}
}
