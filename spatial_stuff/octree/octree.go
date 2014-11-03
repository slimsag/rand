package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"fmt"
)

type Octant struct {
	Depth  int
	Bounds math.Rect3

	// In order of:
	// 0 - Top, Front, Left
	// 1 - Top, Front, Right
	// 2 - Top, Back, Left
	// 3 - Top, Back, Right
	// 4 - Bottom, Front, Left
	// 5 - Bottom, Front, Right
	// 6 - Bottom, Back, Left
	// 7 - Bottom, Back, Right
	Children [8]*Octant

	Objects []gfx.Spatial
}

// depthAdd adds the given value to each octant's depth value, recursively.
func (o *Octant) depthAdd(d int) {
	o.Depth += d
	for _, c := range o.Children {
		if c != nil {
			c.depthAdd(d)
		}
	}
}

// childBounds returns a bounding box for the given child index. Normally this
// is used for creating the child octants.
func (o *Octant) childBounds(i int) math.Rect3 {
	// Determine where the child octant is located.
	isTop := i == 0 || i == 1 || i == 2 || i == 3
	isFront := i == 0 || i == 1 || i == 4 || i == 5
	isRight := i == 1 || i == 3 || i == 5 || i == 7

	// The child bounds is this octant's bounds shrunk in the proper direction
	// by half this octant's size.
	c := o.Bounds
	halfSize := o.Bounds.Size().DivScalar(2)
	if isTop {
		c.Min.Z += halfSize.Z
	} else {
		c.Max.Z -= halfSize.Z
	}

	if isFront {
		c.Min.Y += halfSize.Y
	} else {
		c.Max.Y -= halfSize.Y
	}

	if isRight {
		c.Min.X += halfSize.X
	} else {
		c.Max.X -= halfSize.X
	}
	return c
}

// childFits returns a child octant index where the given bounding box can fit.
// It returns -1 if there is no child octant that can fit the bounding box.
func (o *Octant) childFits(b math.Rect3) int {
	ob := o.Bounds
	center := ob.Center()

	if b.Max.X <= center.X {
		// Left
		if b.Max.Y <= center.Y {
			// Back, Left
			if b.Max.Z <= center.Z {
				// Bottom, Back, Left
				return 6
			} else if b.Min.Z >= center.Z {
				// Top, Back, Left
				return 2
			}
		} else if b.Min.Y >= center.Y {
			// Front, Left
			if b.Max.Z <= center.Z {
				// Bottom, Front, Left
				return 4
			} else if b.Min.Z >= center.Z {
				// Top, Front, Left
				return 0
			}
		}
	} else if b.Min.X >= center.X {
		// Right
		if b.Max.Y <= center.Y {
			// Back, Right
			if b.Max.Z <= center.Z {
				// Bottom, Back, Right
				return 7
			} else if b.Min.Z >= center.Z {
				// Top, Back, Right
				return 3
			}
		} else if b.Min.Y >= center.Y {
			// Front, Right
			if b.Max.Z <= center.Z {
				// Bottom, Front, Right
				return 5
			} else if b.Min.Z >= center.Z {
				// Top, Front, Right
				return 1
			}
		}
	}
	return -1
}

// expand performs octree expansion, it should only be performed on the root of
// the octree. It operates such that the octant, o, becomes the child octant i.
func (o *Octant) expand(i int) {
	newChild := &Octant{
		Bounds:   o.Bounds,
		Children: o.Children,
		Objects:  o.Objects,
	}
	newChild.depthAdd(1)

	// Nil the children.
	for ci := range o.Children {
		o.Children[ci] = nil
	}

	// Store this octant as a child.
	o.Children[i] = newChild

	// Nil the object list.
	o.Objects = nil

	// Determine the direction in which we should expand:
	goDown := i == 0 || i == 1 || i == 2 || i == 3
	goBack := i == 0 || i == 1 || i == 4 || i == 5
	goLeft := i == 1 || i == 3 || i == 5 || i == 7

	// Create a proper bounding box (twice the old size).
	size := o.Bounds.Size()
	if goDown {
		o.Bounds.Min.Z -= size.Z
	} else {
		o.Bounds.Max.Z += size.Z
	}

	if goBack {
		o.Bounds.Min.Y -= size.Y
	} else {
		o.Bounds.Max.Y += size.Y
	}

	if goLeft {
		o.Bounds.Min.X -= size.X
	} else {
		o.Bounds.Max.X += size.X
	}
}

func (o *Octant) add(b gfx.Spatial, maxDepth int) {
	if o.Depth+1 <= maxDepth {
		bounds := b.Bounds()
		childIndex := o.childFits(bounds)
		fmt.Println("\n\nChild fits")
		fmt.Println("Boundable:", b)
		fmt.Println("Child:", childIndex, o.childBounds(childIndex))
		if childIndex != -1 {
			child := &Octant{
				Bounds: o.childBounds(childIndex),
				Depth:  o.Depth + 1,
			}
			o.Children[childIndex] = child
			child.add(b, maxDepth)
			return
		}
	}

	o.Objects = append(o.Objects, b)
	fmt.Printf("add %v %p\n", b, o)
	return
}

type Root struct {
	MaxDepth     int
	MaxExpansion int
	*Octant
}

func (r *Root) Add(b gfx.Spatial) {
	//bounds := b.AABB()
	if r.Octant.Bounds.Empty() {
		//r.Octant.Bounds = bounds
		r.Octant.Bounds = math.Rect3{
			Min: math.Vec3{-1, -1, -1},
			Max: math.Vec3{1, 1, 1},
		}
	}
	/*
		if !r.Octant.Bounds.Contains(bounds) {
			fmt.Println("")
			fmt.Println("CANNOT FIT", bounds)
			fmt.Println("HAVE", r.Octant.Bounds)
			r.Octant.expand(1)
			fmt.Println("EXPAND TO", r.Octant.Bounds)
		}
	*/
	r.Octant.add(b, 10) //0xFFFFFFFFFFFF)
}

func (r *Root) Remove(b gfx.Spatial) {

}

func New() *Root {
	return &Root{
		Octant: new(Octant),
	}
}
