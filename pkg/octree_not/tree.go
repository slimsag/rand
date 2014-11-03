// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"fmt"
	"sync"
)

const paranoid = false

// Tree represents a single octree. It is safe for access from multiple
// goroutines.
type Tree struct {
	sync.RWMutex
	splitFactor          int
	numObjects, numNodes int
	root                 *Node

	nodeByObject map[gfx.Boundable]*Node
}

// Root returns the current root node of the octree. The root node returned by
// this method will be different as expansion of the octree occurs.
func (t *Tree) Root() *Node {
	t.RLock()
	r := t.root
	t.RUnlock()
	return r
}

func (t *Tree) Objects() map[gfx.Boundable]*Node {
	return t.nodeByObject
}

// NumObjects returns the number of objects currently added to the tree.
func (t *Tree) NumObjects() int {
	t.RLock()
	n := t.numObjects
	t.RUnlock()
	return n
}

// NumNodes returns the number of nodes the tree is currently composed of.
func (t *Tree) NumNodes() int {
	t.RLock()
	n := t.numNodes
	t.RUnlock()
	return n
}

// Add adds the given boundable to the tree.
func (t *Tree) Add(b gfx.Boundable) {
	bb := b.Bounds()
	t.Lock()
	t.add(b, bb, t.root)
	t.Unlock()
}

// Is literally Add() but without the lock held.
func (t *Tree) add(b gfx.Boundable, bb math.Rect3, target *Node) {
	// Increment the object count.
	t.numObjects++

	// Find a place (a distant child node, possibly) where bb can fit.
find:
	place, childIndex := target.findPlace(bb)
	if paranoid {
		if place != nil && !bb.In(place.bounds) {
			fmt.Println("bb is ", bb)
			fmt.Println("not in", place.bounds)
			panic("findPlace() has failed")
		}
	}
	if place == nil {
		if target != t.root {
			panic("add(): invalid target node")
		}
		// See if expansion is required to accomodate bb.
		for !bb.In(t.root.bounds) {
			// Expand in-place until we accomodate bb.
			newRoot := t.root.expand(bb)
			if newRoot != nil {
				t.numNodes++
				t.root = newRoot
			}
		}

		// No place for bb exists, we'll add it to the root node then.
		place = t.root

		if paranoid {
			if !bb.In(place.bounds) {
				fmt.Println("bb is ", bb)
				fmt.Println("not in", place.bounds)
				panic("expansion has failed.")
			}
		}
	}

	// Split the node if required.
	if count := place.split(t.splitFactor); count > 0 {
		// Increment the node count by the number of nodes created in the
		// split.
		t.numNodes += count

		// Continue the search from the beginning.
		goto find
	}

	// Add the object to the octant's index.
	if childIndex == -1 {
		childIndex = 8
	}
	place.objects[childIndex] = append(place.objects[childIndex], &entry{
		b:      b,
		bounds: &bb,
	})
	if paranoid {
		if !bb.In(place.bounds) {
			fmt.Println("bb is ", bb)
			fmt.Println("not in", place.bounds)
			panic("Add() has failed")
		}
	}

	// Add the object to the map of nodes by object.
	t.nodeByObject[b] = place
}

// Has tells if the tree has the boundable object b inside it. Internally the
// tree keeps a map of all objects which makes this a rather quick operation.
func (t *Tree) Has(b gfx.Boundable) bool {
	t.RLock()
	_, ok := t.nodeByObject[b]
	t.RUnlock()
	return ok
}

// decimate removes the node n from the tree if it does not have any objects or
// child nodes.
func (t *Tree) decimate(n *Node) {
	if n.parent != nil {
		// No parent node, it must be the root node. We can't decimate it.
		return
	}

	for _, octObjs := range n.objects {
		if len(octObjs) > 0 {
			// The node has objects, don't decimate it.
			return
		}
	}

	for _, child := range n.children {
		if child != nil {
			// The node has children, don't decimate it.
			return
		}
	}

	// Remove from the parent's children list.
	n.parent = nil
	for nodeIndex, node := range n.parent.children {
		if n == node {
			n.parent.children[nodeIndex] = nil
			t.numNodes--
		}
	}
}

// Remove tries to remove the given boundable object from the tree. The bounds
// that the boundable returns are useless to this method (i.e. they could have
// changed since insertion into the tree) as internally a map is used for all
// objects in the tree.
//
// If the object is not in the tree, false is returned. Otherwise true is
// returned.
func (t *Tree) Remove(b gfx.Boundable) bool {
	t.Lock()
	defer t.Unlock()

	// Find in the node map and delete it.
	n, ok := t.nodeByObject[b]
	if !ok {
		return false
	}
	delete(t.nodeByObject, b)

	var oct int
	for oct = 0; oct < 9; oct++ {
		for index, o := range n.objects[oct] {
			if o.b == b {
				// This is the object.
				n.objects[oct][index] = nil
				n.objects[oct] = append(n.objects[oct][:index], n.objects[oct][index+1:]...)
				t.numObjects--
				t.decimate(n)
				return true
			}
		}
	}

	panic("Failed to remove object.")
}

// Update updates the given boundable object old by moving it to the new spot
// described by b. It is functionally equivilent to:
//  result := t.Remove(old)
//  t.Add(b)
// It is faster than the above code because it can leverage temporal coherence
// when the object has not moved very far and is still inside it's octant.
func (t *Tree) Update(old, b gfx.Boundable) bool {
	var (
		bb            = b.Bounds()
		existed       bool
		oct, objIndex int
		o             *entry
	)

	t.Lock()
	defer t.Unlock()

	// Find in the node map and delete it.
	n, ok := t.nodeByObject[old]
	if !ok {
		return existed
	}
	addTarget := t.root

	for oct = 0; oct < 9; oct++ {
		objs := n.objects[oct]
		for objIndex, o = range objs {
			if o.b != old {
				// This object is not it.
				continue
			}

			// This is the object we want.
			existed = true
			fittingChild := int(n.childFits(bb))
			fitsNode := bb.In(n.bounds)
			if oct != 9 && fittingChild == oct {
				// It fits inside a child octant still: just replace it.
				goto replacal
			} else if fitsNode {
				if oct == 9 && fittingChild == -1 {
					// It fits inside this octant still: just replace it.
					goto replacal
				}

				// It still fits in this octant node though: add it to that node.
				addTarget = n
			}
			goto removal
		}
	}
	panic("Failed to remove object.")

removal:
	delete(t.nodeByObject, old)
	n.objects[oct][objIndex] = nil
	n.objects[oct] = append(n.objects[oct][:objIndex], n.objects[oct][objIndex+1:]...)
	t.numObjects--
	t.add(b, bb, addTarget)
	t.decimate(n)
	return existed

replacal:
	delete(t.nodeByObject, old)
	t.nodeByObject[b] = n
	n.objects[oct][objIndex] = &entry{
		b:      b,
		bounds: &bb,
	}
	return existed
}

// NewTree returns a new octree with a split factor of k and initial root
// bounds b.
//
// The split factor k specifies how many objects can live in a single node
// before the node must be split into eight new octants.
//
// The bounds b specify the intial size of the octree root, if the bounds are
// empty then the first object added to the tree is encapsulated. Remember that
// for expansion to occur the root size must be doubled, so using overly large
// intial root sizes may prove counterintuitive.
func NewTree(k int, b math.Rect3) *Tree {
	t := &Tree{
		numNodes:     1,
		splitFactor:  k,
		nodeByObject: make(map[gfx.Boundable]*Node, 128),
	}
	t.root = &Node{
		access:  &t.RWMutex,
		bounds:  b,
		objects: make([][]*entry, 9),
	}
	return t
}

// New is short-hand for:
//  NewTree(100, math.Rect3Zero)
func New() *Tree {
	return NewTree(100, math.Rect3Zero)
}
