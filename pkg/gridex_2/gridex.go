package gridex

import (
	"azul3d.org/v1/gfx"
	gmath "azul3d.org/v1/math"
	"fmt"
	"math"
)

const (
	primeX = 104009.0
	primeY = 194101.0
	primeZ = 115561.0
)

type byDist struct {
	p gmath.Vec3
	s []gfx.Spatial
}

func (b byDist) Len() int      { return len(b.s) }
func (b byDist) Swap(i, j int) { b.s[i], b.s[j] = b.s[j], b.s[i] }
func (b byDist) Less(i, j int) bool {
	iDist := b.s[i].Bounds().Center().Sub(b.p).LengthSq()
	jDist := b.s[j].Bounds().Center().Sub(b.p).LengthSq()
	return iDist < jDist
}

// A gridex table in Table[xyz][k] format.
type Table struct {
	Data [][]gfx.Spatial
	Size int
}

// Index returns the data index for the given point in space.
func (t *Table) Index(p gmath.Vec3) int {
	x := int(math.Floor(p.X*primeX)) % t.Size
	y := int(math.Floor(p.Y*primeY)) % t.Size
	z := int(math.Floor(p.Z*primeZ)) % t.Size
	return x + t.Size*(y+t.Size*z)
}

func (t *Table) Add(s gfx.Spatial) {
	sb := s.Bounds()

	// Add minimum.
	i := t.Index(sb.Min)
	t.Data[i] = append(t.Data[i], s)

	// Add center.
	i = t.Index(sb.Center())
	t.Data[i] = append(t.Data[i], s)

	// Add maximum.
	i = t.Index(sb.Max)
	t.Data[i] = append(t.Data[i], s)
}

func (t *Table) Remove(s gfx.Spatial) {
	sb := s.Bounds()

	// Remove minimum.
	i := t.Index(sb.Min)
	for ki, k := range t.Data[i] {
		if k == s {
			t.Data[i] = append(t.Data[i][ki:], t.Data[i][ki+1:]...)
		}
	}

	// Remove center.
	i = t.Index(sb.Center())
	for ki, k := range t.Data[i] {
		if k == s {
			t.Data[i] = append(t.Data[i][ki:], t.Data[i][ki+1:]...)
		}
	}

	// Remove maximum.
	i = t.Index(sb.Max)
	for ki, k := range t.Data[i] {
		if k == s {
			t.Data[i] = append(t.Data[i][ki:], t.Data[i][ki+1:]...)
		}
	}
	return
}

func (t *Table) NearestChunks(p gmath.Vec3, callback func(i int) bool) {
	search := func(p gmath.Vec3) bool {
		index := t.Index(p)
		if !callback(index) {
			return false
		}
		return true
	}

	s := float64(t.Size)
	for x := 0.0; x < s; x++ {
		for y := 0.0; y < s; y++ {
			for z := 0.0; z < s; z++ {
				if !search(gmath.Vec3{p.X + x, p.Y + y, p.Z + z}) {
					return
				}
			}
		}
	}
}

func (t *Table) Nearest(p gmath.Vec3, callback func(s gfx.Spatial) bool) int {
	var (
		kp      []gfx.Spatial
		ki      int
		visited = make(map[gfx.Spatial]struct{})
	)
	t.NearestChunks(p, func(i int) bool {
		/*
			chunk := t.Data[i]
			fmt.Println(i, len(t.Data[i]))
			if len(chunk) == 0 {
				return true
			}
			if len(kp) < len(chunk) {
				kp = append(kp, make([]gfx.Spatial, len(chunk)-len(kp))...)
			}
			kp = kp[:len(chunk)]
			copy(kp, chunk)
			fmt.Println(len(kp))
		*/

		kp = t.Data[i]
		fmt.Println(len(kp))
		gfx.InsertionSort(byDist{
			p: p,
			s: kp,
		})
		for _, s := range kp {
			_, v := visited[s]
			if !v {
				visited[s] = struct{}{}
				ki++
				if !callback(s) {
					return false
				}
			}
		}
		return true
	})
	return ki
}

func New(size int) *Table {
	return &Table{
		Data: make([][]gfx.Spatial, size*size*size),
		Size: size,
	}
}
