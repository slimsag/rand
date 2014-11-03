// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"fmt"
)

func largestAxis(a AABB) float64 {
	sz := a.Size()
	if sz.X > sz.Y && sz.X > sz.Z {
		return sz.X
	}
	if sz.Y > sz.X && sz.Y > sz.Z {
		return sz.Y
	}
	return sz.Z
}

type Octant struct {
	// The axis aligned bounding box that encapsulates this octant.
	AABB AABB

	// The depth at which this octant lies.
	Depth int

	// All 8 child octants, any of which may be nil. In specific order like so:
	//    0-----1
	//   /|    /|
	//  / |   / |
	// 2-----3  |
	// |  4--|--5
	// | /   | /
	// |/    |/
	// 6-----7
	Octants [8]*Octant

	// A map of spatial objects to empty structs. A map is used for faster
	// Has(spatial) lookups, the empty structs server no purpose.
	Objects map[Boundable]struct{}
}

func (o *Octant) Query(point math.Vec3) []Boundable {
	return nil
}

func (o *Octant) QueryAABB(b AABB) []Boundable {
	return nil
}

func (o *Octant) Has(s Boundable) bool {
	return true
}

// findOctant finds and returns the child octant index that sb could fit into.
// Returns -1 if sb cannot fit into this octant or any child octant.
func (o *Octant) findOctant(sb AABB) int {
	min := o.AABB.Min
	max := o.AABB.Max
	center := o.AABB.Center()

	inLeft := sb.Min.X > min.X && sb.Max.X < center.X
	inRight := sb.Min.X > center.X && sb.Max.X < max.X

	inBack := sb.Min.Y > min.Y && sb.Max.Y < center.Y
	inFront := sb.Min.Y > center.Y && sb.Max.Y < max.Y

	inBottom := sb.Min.Z > min.Z && sb.Max.Z < center.Z
	inTop := sb.Min.Z > center.Z && sb.Max.Z < max.Z
	//fmt.Println("findOctant", o.AABB)
	//fmt.Println("sb", sb)
	//fmt.Println("inLeft", inLeft, "inRight", inRight, "inBack", inBack, "inFront", inFront, "inBottom", inBottom, "inTop", inTop)

	if (inLeft && inRight) || (inBack && inFront) || (inBottom && inTop) {
		return -1
	}

	//    0-----1
	//   /|    /|
	//  / |   / |
	// 2-----3  |
	// |  4--|--5
	// | /   | /
	// |/    |/
	// 6-----7
	if inLeft {
		// 0, 2, 4, 6
		if inTop {
			// 0, 2
			if inFront {
				return 0
			} else if inBack {
				return 2
			}
		} else if inBottom {
			// 4, 6
			if inFront {
				return 4
			} else if inBack {
				return 6
			}
		}
	} else if inRight {
		// 1, 3, 5, 7
		if inTop {
			// 1, 3
			if inFront {
				return 1
			} else if inBack {
				return 3
			}
		} else if inBottom {
			// 5, 7
			if inFront {
				return 5
			} else if inBack {
				return 7
			}
		}
	}

	return -1
}

func (o *Octant) octantBounds(i int) (b AABB) {
	min := o.AABB.Min
	max := o.AABB.Max
	center := o.AABB.Center()

	//    0-----1
	//   /|    /|
	//  / |   / |
	// 2-----3  |
	// |  4--|--5
	// | /   | /
	// |/    |/
	// 6-----7
	isBottom := i == 4 || i == 5 || i == 6 || i == 7
	isLeft := i == 0 || i == 2 || i == 4 || i == 6
	isBack := i == 2 || i == 3 || i == 6 || i == 7

	if isLeft {
		b.Min.X = min.X
		b.Max.X = center.X
	} else {
		b.Min.X = center.X
		b.Max.X = max.X
	}

	if isBack {
		b.Min.Y = min.Y
		b.Max.Y = center.Y
	} else {
		b.Min.Y = center.Y
		b.Max.Y = max.Y
	}

	if isBottom {
		b.Min.Z = min.Z
		b.Max.Z = center.Z
	} else {
		b.Min.Z = center.Z
		b.Max.Z = max.Z
	}
	return
}

// findDistantOctant finds the smallest distant child octant that sb can fit
// into or returns nil.
// If nil is returned sb may still fit into the octant 'o', as only distant
// child octants are checked.
func (o *Octant) findDistantOctant(sb AABB) *Octant {
	dstIndex := o.findOctant(sb)
	if dstIndex != -1 {
		// Create dst if non-existant.
		if o.Octants[dstIndex] == nil {
			o.Octants[dstIndex] = &Octant{
				Depth:   o.Depth + 1,
				AABB:    o.octantBounds(dstIndex),
				Objects: make(map[Boundable]struct{}),
			}
		}
		dst := o.Octants[dstIndex]
		FIXMEMAXDEPTH := 2
		if dst.Depth+1 < FIXMEMAXDEPTH {
			r := dst.findDistantOctant(sb)
			if r != nil {
				return r
			}
		}
		return dst
	}
	return nil
}

// Add adds the given spatial to this octant.
//
// If the spatial is too large to fit inside of this octant then in-place
// expansion will occur, where this octant will become the first octant of the
// new one (i.e. 'o' will be moved to o.Octants[0]).
//
// The spatial will be fit into subdivided octants recursively up to the given
// limit of o.MaxDepth.
func (o *Octant) Add(s Boundable) {
	sb := s.AABB()

	// If the spatial is too large to fit into this octant then we expand.
	if !o.AABB.Contains(sb) {
		// The spatial is too large to fit into this octant. We will expand our
		// octant to be large enough to contain both the spatial and our old
		// octant.
		old := o.AABB

		largest := largestAxis(o.AABB.Fit(sb))
		o.AABB = AABB{
			Min: math.Vec3{-largest, -largest, -largest},
			Max: math.Vec3{largest, largest, largest},
		}

		/*
			largest := largestAxis(o.AABB.Fit(sb))
			if !o.AABB.Empty() {
				// If this octant was not the empty-starting octant then we must
				// move this octant to become the [-x, +y, +z] child octant.
				firstChild := &Octant{
					Depth: o.Depth + 1,
					AABB: o.AABB,
					Octants: o.Octants,
					Objects: o.Objects,
				}
				o.Octants = [8]*Octant{firstChild}
			}
			o.AABB = AABB{
				Min: math.Vec3{-largest, -largest, -largest},
				Max: math.Vec3{largest, largest, largest},
			}
		*/

		fmt.Println("Expand to", o.AABB)
		//fmt.Println("    New:", o.AABB)
		//fmt.Println("    Size:", o.AABB.Size())
		//fmt.Println("    :")
		fmt.Println("    Small:", sb)
		fmt.Println("    Small-Center:", sb.Center())
		fmt.Println("    Small-Size:", sb.Size())
		//fmt.Println("    Size:", sb.Size())
		fmt.Println("    :")
		fmt.Println("    Big:", old)
		fmt.Println("    Big-Center:", old.Center())
		fmt.Println("    Big-Size:", old.Size())
		//fmt.Println("    Size:", old.Size())
		fmt.Println("    :")
		o.Objects = make(map[Boundable]struct{})
	}

	if !o.AABB.Contains(sb) {
		fmt.Println("Failure to contain:", sb)
		fmt.Println("    Size:", sb.Size())
		fmt.Println("Expansion failed!")
		//panic("Expansion failed!")
	}

	// Find an octant (this one or a distant child one) that the spatial can be
	// fit into.
	dst := o.findDistantOctant(sb)
	if dst == nil {
		dst = o
	}
	fmt.Println("fit in", dst.Depth, len(dst.Objects))
	dst.Objects[s] = struct{}{}
}

func (o *Octant) Remove(s Boundable) {
}

func (o *Octant) Update(s Boundable) {
	o.Remove(s)
	o.Add(s)
}

/*
func (o *Octant) Has(s Spatial) bool {
	// Check if it's inside this octant.
	_, ok := o.Objects[s]
	if ok {
		return true
	}

	// Recursive check to see if children octants contain the spatial.
	for _, octant := range o.Octants {
		if octant != nil && octant.Has(s) {
			return true
		}
	}
	return false
}

func (o *Octant) expand(s Spatial) {
	// Create a bounding box that fits both o.AABB and sb; half that is the new
	// octant size.
	sb := s.AABB()
	f := o.AABB.Fit(sb)
	octantSize := AABB{
		Min: math.Vec3{
			X: f.Min.X / 2.0,
			Y: f.Min.Y / 2.0,
			Z: f.Min.Z / 2.0,
		},
		Max: math.Vec3{
			X: f.Max.X / 2.0,
			Y: f.Max.Y / 2.0,
			Z: f.Max.Z / 2.0,
		},
	}

	// Octant 'o' will become the top-left (Octants[0]) octant, and we will
	// then have seven new empty octants.
	fmt.Println("Octant", "Expanded->", f)
	o.AABB = f
	topLeft := 	&Octant{
		AABB: octantSize,
		Octants: o.Octants,
		Objects: o.Objects,
	}
	o.Octants = [8]*Octant{topLeft}
	o.Objects = make(map[Spatial]struct{})
}

func (o *Octant) place(s Spatial, sb AABB) bool {
}

func (o *Octant) Add(s Spatial) {
	// Remove the spatial if it already exists inside this octree.
	o.Remove(s)

	sb := s.AABB()
	if !o.AABB.Contains(sb) {
		// This octant is too small to fit the spatial. We must expand.
		o.expand(s)
	}

	// Try to place the spatial into a child octant.
	if o.place(s, sb) {
		return
	}

	// We couldn't place the spatial into a child octant, but we do know that
	// it at least fits inside

	if !o.place(s, sb) {
		// We couldn't place the spatial into a child octant, but we DO
	}

	o.
	octantIndex := o.findIndex(sb)
	if octantIndex

	for _, octant := range o.Octants {
		if o.findIndex(sb) !=
		if o.couldFit(octantIndex, sb) {

		}

		if octant != nil && octant.AABB.Contains(sb) {
			fmt.Println("Octant", "Fit->", sb)
			octant.Objects[s] = struct{}{}
			return
		}
	}

	// The spatial couldn't fit into any of our child octants. It's placed
	// inside 'o' then.
	o.Objects[s] = struct{}{}
	fmt.Println("Octant", "Store->", sb)
}

func (o *Octant) Remove(s Spatial) bool {
	// Check if it's inside this octant.
	_, ok := o.Objects[s]
	if ok {
		// It's inside this octant so delete it.
		// TODO: consider de-expanding?
		delete(o.Objects, s)
		return true
	}

	// Recursive check to see if children octants contain the spatial.
	for _, octant := range o.Octants {
		if octant != nil && octant.Remove(s) {
			return true
		}
	}
	return false
}
*/

func NewOctree() *Octant {
	return &Octant{}
}
