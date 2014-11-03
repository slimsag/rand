// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"sync"
)

// ChildIndex represents a single child octant number in the range of 0-7, most
// clients will utilize the constants which define specific child octants, e.g.
// TopFrontLeft.
type ChildIndex int

// Left tells if this child is located on the left.
func (c ChildIndex) Left() bool {
	return c == TopFrontLeft || c == TopBackLeft || c == BottomFrontLeft || c == BottomBackLeft
}

// Back tells if this child is located on the back.
func (c ChildIndex) Back() bool {
	return c == TopBackLeft || c == TopBackRight || c == BottomBackLeft || c == BottomBackRight
}

// Bottom tells if this child is located on the bottom.
func (c ChildIndex) Bottom() bool {
	return c == BottomFrontLeft || c == BottomFrontRight || c == BottomBackLeft || c == BottomBackRight
}

const (
	TopFrontLeft ChildIndex = iota
	TopFrontRight
	TopBackLeft
	TopBackRight
	BottomFrontLeft
	BottomFrontRight
	BottomBackLeft
	BottomBackRight
)

type entry struct {
	b      gfx.Boundable
	bounds *math.Rect3
}

// Node represents a single node within the octree. It maintains a list of
// all eight child octant's and a slice of objects stored in the node residing
// in each octant (or not). Like *Tree, nodes are safe for access from multiple
// goroutines.
type Node struct {
	access   *sync.RWMutex
	level    int
	bounds   math.Rect3
	parent   *Node
	children [8]*Node

	// objects per-octant (last index is objects not in any octant).
	objects [][]*entry
}

// Level returns the level in the tree (where the root node is zero and each
// child level increments by one).
func (n *Node) Level() int {
	n.access.RLock()
	l := n.level
	n.access.RUnlock()
	return l
}

// Bounds implements the gfx.Boundable interface by returning this node's bounds.
func (n *Node) Bounds() math.Rect3 {
	n.access.RLock()
	b := n.bounds
	n.access.RUnlock()
	return b
}

// Child returns the child octant of this node at the given index. In addition
// to the predefined constants (e.g. TopFrontLeft) you may wish to create your
// own child index in the range of 0-7.
//
// If nil is returned then there is currently no child node for the given
// index.
func (n *Node) Child(i ChildIndex) *Node {
	n.access.RLock()
	c := n.children[i]
	n.access.RUnlock()
	return c
}

// NumObjects returns the number of objects in this node for the given octant
// index. Indices 0-7 refer to each octant, and the special index 8 refers to
// the objects not contained within any octant for this node.
func (n *Node) NumObjects(oct int) int {
	n.access.RLock()
	l := len(n.objects[oct])
	n.access.RUnlock()
	return l
}

// Object returns the object in this node at the given index inside the given
// octant. To iterate over all objects within this node, for instance:
//  for oct := 0; i < 9; oct++ {
//      for i := 0; i < n.NumObjects(oct); i++ {
//          obj := n.Object(oct, i)
//      }
//  }
func (n *Node) Object(oct, i int) gfx.Boundable {
	n.access.RLock()
	o := n.objects[oct][i].b
	n.access.RUnlock()
	return o
}

// childBounds returns a bounding box for the given child index. The child
// bounds is this node's bounds shrunk by half in the proper direction.
func (n *Node) childBounds(i ChildIndex) math.Rect3 {
	c := n.bounds
	halfSize := c.Size().DivScalar(2)
	if i.Bottom() {
		c.Max.Z -= halfSize.Z
	} else {
		c.Min.Z += halfSize.Z
	}

	if i.Back() {
		c.Max.Y -= halfSize.Y
	} else {
		c.Min.Y += halfSize.Y
	}

	if i.Left() {
		c.Max.X -= halfSize.X
	} else {
		c.Min.X += halfSize.X
	}
	return c
}

// childFits returns a child octant index where the given bounding box can fit.
// It returns -1 if there is no child octant that can fit the bounding box.
func (n *Node) childFits(b math.Rect3) ChildIndex {
	nb := n.bounds
	center := nb.Center()

	if b.Max.X <= center.X {
		// Left
		if b.Max.Y <= center.Y {
			// Back, Left
			if b.Max.Z <= center.Z {
				return BottomBackLeft
			} else if b.Min.Z >= center.Z {
				return TopBackLeft
			}
		} else if b.Min.Y >= center.Y {
			// Front, Left
			if b.Max.Z <= center.Z {
				return BottomFrontLeft
			} else if b.Min.Z >= center.Z {
				return TopFrontLeft
			}
		}
	} else if b.Min.X >= center.X {
		// Right
		if b.Max.Y <= center.Y {
			// Back, Right
			if b.Max.Z <= center.Z {
				return BottomBackRight
			} else if b.Min.Z >= center.Z {
				return TopBackRight
			}
		} else if b.Min.Y >= center.Y {
			// Front, Right
			if b.Max.Z <= center.Z {
				return BottomFrontRight
			} else if b.Min.Z >= center.Z {
				return TopFrontRight
			}
		}
	}
	return -1
}

// expand performs expansion of the root node, n, towards r's center. If a new
// root node is created then it is returned and n.parent is set to the new root
// node.
func (n *Node) expand(r math.Rect3) *Node {
	nb := n.bounds
	if nb.Empty() {
		// For starting bounds we will (squarely) encapsulate the rectangle.
		rsz := r.Size()
		s := rsz.X
		if rsz.Y > s {
			s = rsz.Y
		}
		if rsz.Z > s {
			s = rsz.Z
		}
		rcenter := r.Center()
		s /= 2
		s *= 32
		startSize := math.Vec3{s, s, s}
		n.bounds = math.Rect3{
			Min: rcenter.Sub(startSize),
			Max: rcenter.Add(startSize),
		}
		return nil
	}

	// Expansion occurs by growing the octree such that the root node becomes
	// a new node whose child is the old root node. Thus we can simply
	// determine in which direction the new root should be extended (by twice
	// the old root's size) by comparing the centres of the old root and the
	// rectangle in question.
	c := nb.Center()
	rc := r.Center()
	sz := nb.Size()

	// Expand by becoming the opposite octant of r's closest point to nb.
	var ci ChildIndex
	if rc.Z > c.Z {
		// Top
		if rc.Y > c.Y {
			// Top, Front
			if rc.X > c.X {
				ci = TopFrontRight
			} else {
				ci = TopFrontLeft
			}
		} else {
			// Top, Back
			if rc.X > c.X {
				ci = TopBackRight
			} else {
				ci = TopBackLeft
			}
		}
	} else {
		// Bottom
		if rc.Y > c.Y {
			// Bottom, Front
			if rc.X > c.X {
				ci = BottomFrontRight
			} else {
				ci = BottomFrontLeft
			}
		} else {
			// Bottom, Back
			if rc.X > c.X {
				ci = BottomBackRight
			} else {
				ci = BottomBackLeft
			}
		}
	}

	expDown := ci.Bottom()
	expBack := ci.Back()
	expLeft := ci.Left()

	fb := n.bounds
	if expDown {
		fb.Min.Z -= sz.Z
	} else {
		fb.Max.Z += sz.Z
	}

	if expBack {
		fb.Min.Y -= sz.Y
	} else {
		fb.Max.Y += sz.Y
	}

	if expLeft {
		fb.Min.X -= sz.X
	} else {
		fb.Max.X += sz.X
	}

	newRoot := &Node{
		access:  n.access,
		bounds:  fb,
		level:   n.level + 1,
		objects: make([][]*entry, 9),
	}
	newRoot.children[ci] = n
	n.parent = newRoot
	return newRoot
}

// findPlace finds a place in this node or any node below it in the tree where
// r can be placed.
func (n *Node) findPlace(r math.Rect3) (*Node, ChildIndex) {
	if !r.In(n.bounds) {
		return nil, -1
	}
	childIndex := n.childFits(r)
	if childIndex != -1 && n.children[childIndex] != nil {
		cn, ci := n.children[childIndex].findPlace(r)
		if cn != nil {
			return cn, ci
		}
	}
	return n, -1
}

// split considers splitting this node for the addition of a new object. It
// returns the number of nodes created in the split, if any. The split factor
// defines the number of objects that can reside in a node before a split must
// occur.
func (n *Node) split(splitFactor int) int {
	nObjects := 0
	if n.objects != nil {
		for i := 0; i < 9; i++ {
			nObjects += len(n.objects[i])
		}
	}
	if nObjects+1 < splitFactor {
		return 0
	}

	splitCount := 0
	for i := 0; i < 8; i++ {
		if n.children[i] != nil {
			continue
		}
		splitCount++

		child := &Node{
			access:  n.access,
			level:   n.level + 1,
			bounds:  n.childBounds(ChildIndex(i)),
			objects: make([][]*entry, 9),
			parent:  n,
		}
		child.objects[i] = n.objects[i]
		n.objects[i] = nil
		n.children[i] = child
	}
	return splitCount
}
