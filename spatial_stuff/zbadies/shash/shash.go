package shash

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"fmt"
	gmath "math"
)

type Index struct {
	AvgLengthSq float64
	S           gfx.Spatial
}

type Shash struct {
	AvgLengthSq float64
	Table       [][][]Index
}

func (s *Shash) DistIndex(dist math.Vec3) int {
	dx := gmath.Floor(dist.X*93563.0) / s.AvgLengthSq
	dy := gmath.Floor(dist.Y*89041.0) / s.AvgLengthSq
	dz := gmath.Floor(dist.Z*84631.0) / s.AvgLengthSq
	return int(dx + dy + dz)
}

func (s *Shash) AngleIndex(angle math.Vec3) int {
	ax := gmath.Floor(angle.X * 93563.0)
	ay := gmath.Floor(angle.Y * 89041.0)
	az := gmath.Floor(angle.Z * 84631.0)
	return int(ax + ay + az)
}

func (s *Shash) Add(sp gfx.Spatial) {
	// Find the distance from the center of the spatial's bounding box to the
	// center of the world.
	b := sp.Bounds()
	dist := b.Center().Sub(math.Vec3Zero)
	s.AvgLengthSq += dist.LengthSq()
	s.AvgLengthSq /= 2.0

	// Find the angle between the spatial's position in space relative to the
	// center of the world.
	angle, _ := dist.Normalized()

	// Computer the start and end indices.
	start := s.AngleIndex(angle)
	end := s.AngleIndex(angle)
	distStart := s.DistIndex(b.Min)
	distEnd := s.DistIndex(b.Max)

	if start > end || distStart > distEnd {
		panic("ups")
	}

	//fmt.Println(start, end, "|", distStart, distEnd)
	_ = fmt.Println

	// Insert into table.
	for a := start; a <= end; a *= len(s.Table) {
		ak := a % len(s.Table)
		for d := distStart; d <= distEnd; d *= len(s.Table) {
			dk := d % len(s.Table)
			s.Table[ak][dk] = append(s.Table[ak][dk], Index{
				AvgLengthSq: s.AvgLengthSq,
				S:           sp,
			})
		}
	}
}

func (s *Shash) Remove(sp gfx.Spatial) bool {
	return false
}

func NewSize(size int) *Shash {
	s := &Shash{
		Table: make([][][]Index, size),
	}
	for j := 0; j < size; j++ {
		s.Table[j] = make([][]Index, size)
	}
	return s
}

func New() *Shash {
	return NewSize(64)
}
