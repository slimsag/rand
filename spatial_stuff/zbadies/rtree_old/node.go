// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	gmath "azul3d.org/v1/math"
)

// node represents a single node within the tree.
type node struct {
	// The parent of this node.
	parent *node

	// Whether or not this node is considered a leaf node.
	leaf bool

	// The level at which this node resides within the tree.
	level int

	// The list of entries
	entries []entry
}

// calcBounds calculates the bounds of all entries in the node.
func (n *node) calcBounds() (bounds gmath.Rect3) {
	for _, e := range n.entries {
		bounds = bounds.Union(e.bounds)
	}
	return
}

// getEntry returns a pointer to the entry for the node n from n's parent.
func (n *node) getEntry() *entry {
	var e *entry
	for i := range n.parent.entries {
		if n.parent.entries[i].child == n {
			e = &n.parent.entries[i]
			break
		}
	}
	return e
}

// split splits a node into two groups while attempting to minimize the
// bounding-box area of the resulting groups.
func (n *node) split(minGroupSize int) (left, right *node) {
	// find the initial split
	l, r := n.pickSeeds()
	leftSeed, rightSeed := n.entries[l], n.entries[r]

	// get the entries to be divided between left and right
	remaining := append(n.entries[:l], n.entries[l+1:r]...)
	remaining = append(remaining, n.entries[r+1:]...)

	// setup the new split nodes, but re-use n as the left node
	left = n
	left.entries = []entry{leftSeed}
	right = &node{
		parent:  n.parent,
		leaf:    n.leaf,
		level:   n.level,
		entries: []entry{rightSeed},
	}

	if rightSeed.child != nil {
		rightSeed.child.parent = right
	}
	if leftSeed.child != nil {
		leftSeed.child.parent = left
	}

	// distribute all of n's old entries into left and right.
	for len(remaining) > 0 {
		next := pickNext(left, right, remaining)
		e := remaining[next]

		if len(remaining)+len(left.entries) <= minGroupSize {
			assign(e, left)
		} else if len(remaining)+len(right.entries) <= minGroupSize {
			assign(e, right)
		} else {
			assignGroup(e, left, right)
		}

		remaining = append(remaining[:next], remaining[next+1:]...)
	}

	return
}
