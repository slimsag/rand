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

// Table represents a three-dimensional gridex table.
type Table struct {
	Data [][]gfx.Boundable
	Size int
}

// Index returns the data index for the given point in space.
func (t *Table) Index(p gmath.Vec3) int {
	p.X *= primeX
	p.Y *= primeY
	p.Z *= primeZ
	x := int(math.Abs(p.X)) % t.Size
	y := int(math.Abs(p.Y)) % t.Size
	z := int(math.Abs(p.Z)) % t.Size
	return x + t.Size*(y+t.Size*z)
}

func (t *Table) eachIndex(r gmath.Rect3, cb func(i int)) {
	for x := r.Min.X; x <= r.Max.X; x += 1 {
		for y := r.Min.Y; y <= r.Max.Y; y += 1 {
			for z := r.Min.Z; z <= r.Max.Z; z += 1 {
				cb(t.Index(gmath.Vec3{x, y, z}))
			}
		}
	}
}

func (t *Table) Add(s gfx.Boundable) {
	sb := s.Bounds()
	t.eachIndex(sb, func(idx int) {
		fmt.Println("add", idx)
		t.Data[idx] = append(t.Data[idx], s)
	})
}

func (t *Table) Remove(s gfx.Boundable) (ok bool) {
	sb := s.Bounds()
	t.eachIndex(sb, func(idx int) {
		found := -1
		for i, v := range t.Data[idx] {
			if v == s {
				found = i
				break
			}
		}
		if found == -1 {
			ok = false
			return
		}
		t.Data[idx] = append(t.Data[idx][:found], t.Data[idx][found+1:]...)
	})
	ok = true
	return
}

func (t *Table) rectSearch(r gmath.Rect3, results chan gfx.Boundable, stop chan struct{}, search func(s gfx.Boundable) bool) {
	t.eachIndex(r, func(idx int) {
		fmt.Println("search", idx)
		for _, s := range t.Data[idx] {
			if idx == 0 {
				fmt.Println("send>", s)
			}
			if !search(s) {
				return
			}
		}
	})
	close(results)
}

// In performs a search for all spatials within the table that are completely
// inside the given rectangle. This function returns immedietly.
//
// The search will be executed in a seperate goroutine, results will be sent
// over the given results channel (e.g. with a buffer size of 32) until the
// search completes or is halted.
//
// If non-nil, the stop channel can be used to halt the search. Whenever a
// struct{}{} is sent over the channel the search will be indefinitely halted.
//
// The results channel will be closed when the search complets or is halted.
func (t *Table) In(r gmath.Rect3, results chan gfx.Boundable, stop chan struct{}) {
	go t.rectSearch(r, results, stop, func(s gfx.Boundable) bool {
		if s.Bounds().In(r) {
			select {
			case results <- s:
			case <-stop:
				return false
			}
		}
		return true
	})
}

// Intersect performs a search for all spatials within the table that are
// intersecting with the given rectangle. This function returns immedietly.
//
// The search will be executed in a seperate goroutine, results will be sent
// over the given results channel (e.g. with a buffer size of 32) until the
// search completes or is halted.
//
// If non-nil, the stop channel can be used to halt the search. Whenever a
// struct{}{} is sent over the channel the search will be indefinitely halted.
//
// The results channel will be closed when the search complets or is halted.
func (t *Table) Intersect(r gmath.Rect3, results chan gfx.Boundable, stop chan struct{}) {
	go t.rectSearch(r, results, stop, func(s gfx.Boundable) bool {
		if _, ok := s.Bounds().Intersect(r); ok {
			select {
			case results <- s:
			case <-stop:
				return false
			}
		}
		return true
	})
}

/*
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

func (t *Table) Add(s gfx.Boundable) {
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

func (t *Table) Remove(s gfx.Boundable) {
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

func (t *Table) Nearest(p gmath.Vec3, callback func(s gfx.Boundable) bool) int {
	var (
		kp      []gfx.Boundable
		ki      int
		visited = make(map[gfx.Boundable]struct{})
	)
	t.NearestChunks(p, func(i int) bool {
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
*/

func New(size int) *Table {
	return &Table{
		Data: make([][]gfx.Boundable, size*size*size),
		Size: size,
	}
}
