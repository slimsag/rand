// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"azul3d.org/v1/gfx"
	gmath "azul3d.org/v1/math"
	"fmt"
	"math"
	"sort"
)

// entry represents a single entry stored in a node within the tree.
type entry struct {
	// The bounding box of all the children of this entry.
	bounds gmath.Rect3

	// The child node of this entry.
	child *node

	// The spatial object this entry represents.
	obj gfx.Spatial
}

// Tree represents a 3D R* Tree, which can be used for efficient searching of
// 3D spatial data.
type Tree struct {
	// When splitting occurs on a leaf node it will be split such that each
	// split portion has no less than this many children.
	MinChildren int

	// The amount of children a leaf node can have before it will be split.
	MaxChildren int

	// The number of objects inserted into the tree.
	size int

	// The root node of the tree.
	root *node

	// The current depth of the tree.
	depth int
}

// Size returns the number of objects that have been inserted into the tree.
func (t *Tree) Size() int {
	return t.size
}

// Depth returns the current depth of the tree.
func (t *Tree) Depth() int {
	return t.depth
}

// Insert inserts the given spatial object into the tree.
func (t *Tree) Insert(s gfx.Spatial) {
	t.size++
	e := entry{
		bounds: s.Bounds(),
		obj:    s,
	}
	t.insert(e, 1)
}

// New returns a new, initialized, 3D R* Tree.
//
// minChildren specifies the minimum number a leaf node can have after it is
// split (i.e., splitting will occur such that each leaf node will have min
// children).
//
// maxChildren is the amount of children a leaf node can have before it will be
// split.
//
// For example:
//  tree := New(32, 64)
func New(minChildren, maxChildren int) *Tree {
	return &Tree{
		MinChildren: minChildren,
		MaxChildren: maxChildren,
		root: &node{
			leaf:  true,
			level: 1,
		},
	}
}

// insert adds the entry to the tree node at the given level.
func (t *Tree) insert(e entry, level int) {
	leaf := t.chooseNode(t.root, e, level)
	leaf.entries = append(leaf.entries, e)
	fmt.Printf("insert %p into %p\n", e.obj, leaf)

	// update parent pointer if necessary
	if e.child != nil {
		e.child.parent = leaf
	}

	// split the leaf node if it overflows.
	var split *node
	if len(leaf.entries) > t.MaxChildren {
		leaf, split = leaf.split(t.MinChildren)
	}
	root, splitRoot := t.adjustTree(leaf, split)
	if splitRoot != nil {
		oldRoot := root
		t.depth++
		t.root = &node{
			parent: nil,
			level:  t.depth,
			entries: []entry{
				entry{
					bounds: oldRoot.calcBounds(),
					child:  oldRoot,
				},
				entry{
					bounds: splitRoot.calcBounds(),
					child:  splitRoot,
				},
			},
		}
		oldRoot.parent = t.root
		splitRoot.parent = t.root
	}
}

// chooseNode finds the node at the specified level to which e should be added.
func (t *Tree) chooseNode(n *node, e entry, level int) *node {
	if n.leaf || n.level == level {
		return n
	}

	// Find the entry whose bounds need the least enlargement to include the
	// spatial.
	diff := math.MaxFloat64
	var chosen entry
	for _, en := range n.entries {
		entrySize := en.bounds.Size().LengthSq()
		u := en.bounds.Union(e.bounds)
		d := u.Size().LengthSq() - entrySize
		if d < diff || (d == diff && entrySize < chosen.bounds.Size().LengthSq()) {
			diff = d
			chosen = en
		}
	}

	return t.chooseNode(chosen.child, e, level)
}

// adjustTree splits overflowing nodes and propagates the changes upwards.
func (t *Tree) adjustTree(n, nn *node) (*node, *node) {
	// Let the caller handle root adjustments.
	if n == t.root {
		return n, nn
	}

	// Re-size the bounding box of n to account for lower-level changes.
	en := n.getEntry()
	en.bounds = n.calcBounds()

	// If nn is nil, then we're just propagating changes upwards.
	if nn == nil {
		return t.adjustTree(n.parent, nil)
	}

	// Otherwise, these are two nodes resulting from a split.
	// n was reused as the "left" node, but we need to add nn to n.parent.
	enn := entry{nn.calcBounds(), nn, nil}
	n.parent.entries = append(n.parent.entries, enn)

	// If the new entry overflows the parent, split the parent and propagate.
	if len(n.parent.entries) > t.MaxChildren {
		n0, n1 := n.parent.split(t.MinChildren)
		return t.adjustTree(n0, n1)
	}

	// Otherwise keep propagating changes upwards.
	return t.adjustTree(n.parent, nil)
}

func assign(e entry, group *node) {
	if e.child != nil {
		e.child.parent = group
	}
	group.entries = append(group.entries, e)
}

// assignGroup chooses one of two groups to which a node should be added.
func assignGroup(e entry, left, right *node) {
	leftBounds := left.calcBounds()
	rightBounds := right.calcBounds()
	leftEnlarged := leftBounds.Union(e.bounds)
	rightEnlarged := rightBounds.Union(e.bounds)

	leftBoundsSize := leftBounds.Size().LengthSq()
	rightBoundsSize := rightBounds.Size().LengthSq()

	// first, choose the group that needs the least enlargement
	leftDiff := leftEnlarged.Size().LengthSq() - leftBoundsSize
	rightDiff := rightEnlarged.Size().LengthSq() - rightBoundsSize
	if diff := leftDiff - rightDiff; diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group that has smaller area
	if diff := leftBoundsSize - rightBoundsSize; diff < 0 {
		assign(e, left)
		return
	} else if diff > 0 {
		assign(e, right)
		return
	}

	// next, choose the group with fewer entries
	if diff := len(left.entries) - len(right.entries); diff <= 0 {
		assign(e, left)
		return
	}
	assign(e, right)
}

// pickSeeds chooses two child entries of n to start a split.
func (n *node) pickSeeds() (int, int) {
	left, right := 0, 1
	maxWastedSpace := -1.0
	for i, e1 := range n.entries {
		for j, e2 := range n.entries[i+1:] {
			e1Size := e1.bounds.Size().LengthSq()
			e2Size := e2.bounds.Size().LengthSq()
			d := e1.bounds.Union(e2.bounds).Size().LengthSq() - e1Size - e2Size
			if d > maxWastedSpace {
				maxWastedSpace = d
				left, right = i, j+i+1
			}
		}
	}
	return left, right
}

// pickNext chooses an entry to be added to an entry group.
func pickNext(left, right *node, entries []entry) (next int) {
	maxDiff := -1.0
	leftBounds := left.calcBounds()
	rightBounds := right.calcBounds()
	leftBoundsSize := leftBounds.Size().LengthSq()
	rightBoundsSize := rightBounds.Size().LengthSq()
	for i, e := range entries {
		d1 := leftBounds.Union(e.bounds).Size().LengthSq() - leftBoundsSize
		d2 := rightBounds.Union(e.bounds).Size().LengthSq() - rightBoundsSize
		d := math.Abs(d1 - d2)
		if d > maxDiff {
			maxDiff = d
			next = i
		}
	}
	return
}

// Deletion

// Delete removes an object from the tree.  If the object is not found, returns
// false, otherwise returns true.
//
// Implemented per Section 3.3 of "R-trees: A Dynamic Index Structure for
// Spatial Searching" by A. Guttman, Proceedings of ACM SIGMOD, p. 47-57, 1984.
func (t *Tree) Delete(obj gfx.Spatial) bool {
	n := t.findLeaf(t.root, obj)
	if n == nil {
		panic("findLeaf failed!")
		return false
	}

	ind := -1
	for i, e := range n.entries {
		if e.obj == obj {
			ind = i
		}
	}
	if ind < 0 {
		return false
	}

	n.entries = append(n.entries[:ind], n.entries[ind+1:]...)

	t.condenseTree(n)
	t.size--

	if !t.root.leaf && len(t.root.entries) == 1 {
		t.root = t.root.entries[0].child
	}
	return true
}

// findLeaf finds the leaf node containing obj.
func (t *Tree) findLeaf(n *node, obj gfx.Spatial) *node {
	if n.leaf {
		return n
	}
	// if not leaf, search all candidate subtrees
	objBounds := obj.Bounds()
	for _, e := range n.entries {
		fmt.Printf("\nin %p\n", n, objBounds.In(e.bounds))
		fmt.Println("objBounds", objBounds)
		fmt.Println("e.bounds", e.bounds)
		if objBounds.In(e.bounds) {
			leaf := t.findLeaf(e.child, obj)
			if leaf == nil {
				continue
			}
			// check if the leaf actually contains the object
			for _, leafEntry := range leaf.entries {
				if leafEntry.obj == obj {
					return leaf
				}
			}
		}
	}
	return nil
}

// condenseTree deletes underflowing nodes and propagates the changes upwards.
func (t *Tree) condenseTree(n *node) {
	deleted := []*node{}

	for n != t.root {
		if len(n.entries) < t.MinChildren {
			// remove n from parent entries
			entries := []entry{}
			for _, e := range n.parent.entries {
				if e.child != n {
					entries = append(entries, e)
				}
			}
			if len(n.parent.entries) == len(entries) {
				panic("Failed to remove entry from parent")
			}
			n.parent.entries = entries

			// only add n to deleted if it still has children
			if len(n.entries) > 0 {
				deleted = append(deleted, n)
			}
		} else {
			// just a child entry deletion, no underflow
			n.getEntry().bounds = n.calcBounds()
		}
		n = n.parent
	}

	for _, n := range deleted {
		// reinsert entry so that it will remain at the same level as before
		e := entry{
			bounds: n.calcBounds(),
			child:  n,
		}
		t.insert(e, n.level+1)
	}
}

// Searching

// Intersect determines which spatial objects in the tree intersect with the
// given rectangle, appends them to the given slice and then returns it.
//
// Since the results slice is just appended to, you may reuse and pre-allocate
// it's memory for additional performance benifits.
//
// If the limit parameter is -1 then an unlimited number of objects are
// appended to the slice, otherwise only the given limit of objects are.
func (t *Tree) Intersect(r gmath.Rect3, limit int, results []gfx.Spatial) []gfx.Spatial {
	return t.intersect(limit, t.root, r, results)
}

func (t *Tree) intersect(k int, n *node, r gmath.Rect3, results []gfx.Spatial) []gfx.Spatial {
	for _, e := range n.entries {
		if k >= 0 && len(results) >= k {
			break
		}

		if !e.bounds.Intersect(r).Equals(gmath.Rect3Zero) {
			if n.leaf {
				results = append(results, e.obj)
			} else {
				margin := k - len(results)
				results = append(results, t.intersect(margin, e.child, r, results)...)
			}
		}
	}
	return results
}

// NearestNeighbor returns the closest object to the specified point.
func (t *Tree) NearestNeighbor(p gmath.Vec3) gfx.Spatial {
	obj, _ := t.nearestNeighbor(p, t.root, math.MaxFloat64, nil)
	return obj
}

// utilities for sorting slices of entries

type entrySlice struct {
	entries []entry
	dists   []float64
	pt      gmath.Vec3
}

func (s entrySlice) Len() int { return len(s.entries) }

func (s entrySlice) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
	s.dists[i], s.dists[j] = s.dists[j], s.dists[i]
}

func (s entrySlice) Less(i, j int) bool {
	return s.dists[i] < s.dists[j]
}

func sortEntries(p gmath.Vec3, entries []entry) ([]entry, []float64) {
	sorted := make([]entry, len(entries))
	dists := make([]float64, len(entries))
	for i := 0; i < len(entries); i++ {
		sorted[i] = entries[i]
		bounds := entries[i].bounds
		dists[i] = p.Sub(bounds.Min).LengthSq()
	}
	sort.Sort(entrySlice{sorted, dists, p})
	return sorted, dists
}

func pruneEntries(p gmath.Vec3, entries []entry, minDists []float64) []entry {
	minMinMaxDist := math.MaxFloat64
	for i := range entries {
		minMaxDist := p.Sub(entries[i].bounds.Max).LengthSq()
		if minMaxDist < minMinMaxDist {
			minMinMaxDist = minMaxDist
		}
	}
	// remove all entries with minDist > minMinMaxDist
	pruned := []entry{}
	for i := range entries {
		if minDists[i] <= minMinMaxDist {
			pruned = append(pruned, entries[i])
		}
	}
	return pruned
}

func (t *Tree) nearestNeighbor(p gmath.Vec3, n *node, d float64, nearest gfx.Spatial) (gfx.Spatial, float64) {
	if n.leaf {
		for _, e := range n.entries {
			dist := p.Sub(e.bounds.Min).Length()
			if dist < d {
				d = dist
				nearest = e.obj
			}
		}
	} else {
		branches, dists := sortEntries(p, n.entries)
		branches = pruneEntries(p, branches, dists)
		for _, e := range branches {
			subNearest, dist := t.nearestNeighbor(p, e.child, d, nearest)
			if dist < d {
				d = dist
				nearest = subNearest
			}
		}
	}

	return nearest, d
}

func (t *Tree) NearestNeighbors(k int, p gmath.Vec3) []gfx.Spatial {
	dists := make([]float64, k)
	objs := make([]gfx.Spatial, k)
	for i := 0; i < k; i++ {
		dists[i] = math.MaxFloat64
		objs[i] = nil
	}
	objs, _ = t.nearestNeighbors(k, p, t.root, dists, objs)
	return objs
}

// insert obj into nearest and return the first k elements in increasing order.
func insertNearest(k int, dists []float64, nearest []gfx.Spatial, dist float64, obj gfx.Spatial) ([]float64, []gfx.Spatial) {
	i := 0
	for i < k && dist >= dists[i] {
		i++
	}
	if i >= k {
		return dists, nearest
	}

	left, right := dists[:i], dists[i:k-1]
	updatedDists := make([]float64, k)
	copy(updatedDists, left)
	updatedDists[i] = dist
	copy(updatedDists[i+1:], right)

	leftObjs, rightObjs := nearest[:i], nearest[i:k-1]
	updatedNearest := make([]gfx.Spatial, k)
	copy(updatedNearest, leftObjs)
	updatedNearest[i] = obj
	copy(updatedNearest[i+1:], rightObjs)

	return updatedDists, updatedNearest
}

func (t *Tree) nearestNeighbors(k int, p gmath.Vec3, n *node, dists []float64, nearest []gfx.Spatial) ([]gfx.Spatial, []float64) {
	if n.leaf {
		for _, e := range n.entries {
			dist := p.Sub(e.bounds.Min).Length()
			dists, nearest = insertNearest(k, dists, nearest, dist, e.obj)
		}
	} else {
		branches, branchDists := sortEntries(p, n.entries)
		branches = pruneEntries(p, branches, branchDists)
		for _, e := range branches {
			nearest, dists = t.nearestNeighbors(k, p, e.child, dists, nearest)
		}
	}
	return nearest, dists
}
