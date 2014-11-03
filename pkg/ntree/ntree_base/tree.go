package ntree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
)

type Node struct {
	Level    int
	bounds   math.Rect3
	Children []*Node
	Objects  []gfx.Spatial
}

// Bounds implements the gfx.Spatial interface.
func (n *Node) Bounds() math.Rect3 {
	return n.bounds
}

// traverse calls the callback for this node, and then for it's children,
// recursively until the callback returns false or there are no more children.
//
// The return value does not serve a purpose for callers.
func (n *Node) traverse(cb func(n *Node) bool) bool {
	if !cb(n) {
		return false
	}
	for _, c := range n.Children {
		if !c.traverse(cb) {
			return false
		}
	}
	return true
}

// find finds this or a distance node where the given rectangle, r, would have
// been placed.
func (n *Node) find(r math.Rect3) *Node {
	if !r.In(n.bounds) {
		// Not inside this node at all.
		return nil
	}

	// Check children...
	for _, child := range n.Children {
		ccn := child.find(r)
		if ccn != nil {
			return ccn
		}
	}
	return n
}

// childSize returns the size of each child node given the divisor.
func (n *Node) childSize(divisor math.Vec3) math.Vec3 {
	return n.bounds.Size().Div(divisor)
}

// Takes a child node position and returns it's bounds.
func (n *Node) childBounds(pos, divisor math.Vec3) math.Rect3 {
	size := n.childSize(divisor)
	return math.Rect3{
		Min: pos,
		Max: pos.Add(size),
	}
}

// childFits finds a direct child of this node that can fit the rectangle, r,
// and returns it's bounds.
func (n *Node) childFits(r math.Rect3, divisor math.Vec3) (b math.Rect3, ok bool) {
	// Find a child path to create.
	nb := n.bounds
	sz := n.childSize(divisor)
	for x := nb.Min.X; x <= nb.Max.X; x += sz.X {
		for y := nb.Min.Y; y <= nb.Max.Y; y += sz.Y {
			for z := nb.Min.Z; z <= nb.Max.Z; z += sz.Z {
				cb := n.childBounds(math.Vec3{x, y, z}, divisor)
				if r.In(cb) {
					return cb, true
				}
			}
		}
	}
	return math.Rect3Zero, false
}

// createPath creates nodes along the needed path up to maxDepth where the
// given spatial's bounds, sb, can be placed.
func (n *Node) createPath(divisor math.Vec3, maxDepth int, sb math.Rect3) *Node {
	nb := n.bounds
	if !sb.In(nb) {
		// Not inside this node at all.
		return nil
	}
	if n.Level+1 > maxDepth {
		// Any child would exceed the maximum depth.
		return n
	}

	// Check existing children...
	for _, child := range n.Children {
		ccn := child.find(sb)
		if ccn != nil {
			// An existing child node can fit the spatials bounds, sb, so ask
			// it to create the path instead.
			return child.createPath(divisor, maxDepth, sb)
		}
	}

	// Find a child path to create.
	db, ok := n.childFits(sb, divisor)
	if ok {
		// Create the child.
		child := &Node{
			Level:  n.Level + 1,
			bounds: db,
		}
		n.Children = append(n.Children, child)
		ccn := child.createPath(divisor, maxDepth, sb)
		if ccn != nil {
			return ccn
		}
		return child
	}
	return n
}

// Tree represents a single N tree.
type Tree struct {
	divisor             math.Vec3
	startScale          float64
	maxDepth, maxExpand int
	count               int
	Root                *Node
	outside             []gfx.Spatial
}

// SetDivisor sets the divisor of this N tree to the given one. The divisor is
// what effectively determines how many nodes make up each level of the N tree.
//
// For instance the following divisor would effectively create a quadtree using
// the X and Y axis:
//  SetDivisor(math.Vec3{2, 2, 1})
//
// And this would create a octree using the X, Y, and Z axis:
//  SetDivisor(math.Vec3{2, 2, 2})
//
// By default the divisor is used to create a viginti septem tree:
//  SetDivisor(math.Vec3{3, 3, 3})
//
// Returned is the old divisor that was in use before the call.
//
// It is not advised to change this value after the tree has had spatials added
// to it.
func (t *Tree) SetDivisor(divisor math.Vec3) (old math.Vec3) {
	old = t.divisor
	t.divisor = divisor
	return
}

// Divisor returns the divisor of this N tree. For information about what the
// divisor is, see the SetDivisor() method.
//
// The default value is math.Vec3{3, 3, 3}
func (t *Tree) Divisor() math.Vec3 {
	return t.divisor
}

// SetStartScale sets the scaling amount for the intial bounds of the root node
// of the tree. The first item added to the tree determines the root node size,
// and it is scaled according to this value.
//
// The default value is 2.0 (twice the size of the first item added).
//
// Returned is the old start scale that was in use before the call.
//
// It is not advised to change this value after the tree has had spatials added
// to it.
func (t *Tree) SetStartScale(startScale float64) (old float64) {
	old = t.startScale
	t.startScale = startScale
	return
}

// StartScale returns the starting scale for the intial bounds of the root node
// of the tree. For more information about what the start scale is, see the
// SetStartScale() method.
func (t *Tree) StartScale() float64 {
	return t.startScale
}

// SetMaxDepth sets the maximum depth for the N tree. This controls how many
// levels of nodes may be added below the root node (level 0). This value does
// not have an effect on expansion of the tree.
//
// It is not advised to change this value after the tree has had spatials added
// to it.
func (t *Tree) SetMaxDepth(maxDepth int) (old int) {
	old = t.maxDepth
	t.maxDepth = maxDepth
	return
}

// MaxDepth returns the maximum depth of the N tree. For more information about
// what this value is see the SetMaxDepth() method.
func (t *Tree) MaxDepth() int {
	return t.maxDepth
}

// SetMaxExpand sets the maximum expansion for the N tree. This controls how
// many levels of nodes may be added above the root node (level 0). This value
// is independant of the maximum depth of the tree.
//
// It is not advised to change this value after the tree has had spatials added
// to it.
func (t *Tree) SetMaxExpand(maxExpand int) (old int) {
	old = t.maxExpand
	t.maxExpand = maxExpand
	return
}

// MaxExpand returns the maximum expansion of the N tree. For more information
// about what this value is see the SetMaxExpand() method.
func (t *Tree) MaxExpand() int {
	return t.maxExpand
}

// Count returns the number of spatials currently added to the N tree. The
// returned number includes those spatials that cannot fit into the tree due to
// expansion limits.
func (t *Tree) Count() int {
	return t.count
}

// OutsideCount returns the number of spatials that are stored in the N tree
// and are outside of the root node's bounds due to expansion limits.
func (t *Tree) OutsideCount() int {
	return len(t.outside)
}

// Add adds the given spatial to the N tree.
func (t *Tree) Add(s gfx.Spatial) {
	t.count++
	sb := s.Bounds()
	if t.Root == nil {
		size := sb.Size()
		if size.Y > size.X {
			size.X = size.X
			size.Y = size.X
			size.Z = size.X
		}
		if size.Z > size.X {
			size.X = size.X
			size.Y = size.X
			size.Z = size.X
		}
		s := t.startScale
		t.Root = &Node{
			bounds: math.Rect3{
				Min: sb.Min.MulScalar(s),
				Max: sb.Min.Add(size).MulScalar(s),
			},
		}
	}
	for !sb.In(t.Root.bounds) {
		rootBounds := t.Root.bounds
		rootSize := rootBounds.Size()
		if t.Root.Level < -t.maxExpand {
			break
		}
		expanded := &Node{
			Level: t.Root.Level - 1,
			bounds: math.Rect3{
				Min: rootBounds.Min.Sub(rootSize),
				Max: rootBounds.Max.Add(rootSize),
			},
			Children: []*Node{t.Root},
		}
		t.Root = expanded
	}

	// Create a path of nodes, subdividings as needed to insert the object into
	// the tree.
	p := t.Root.createPath(t.divisor, t.maxDepth, sb)
	if p == nil {
		// Doesn't fit in the tree.
		t.outside = append(t.outside, s)
		return
	}
	p.Objects = append(p.Objects, s)
}

// Remove tries to remove the given spatial from the N tree.
//
// It should be explicitly noted that it is only possible to remove a spatial
// if Bounds() returns the same identical bounds as when it was added, this is
// for efficiency reasons. Clients wishing to not have to keep track can just
// use a map of spatials to their gfx.Bounds of when they where added to the N
// tree.
//
// This method returns true if the spatial was removed or false if it could not
// be located in the tree due to:
//  1. The spatial's bounds having changed since the last time it was added.
//  2. The spatial having already been removed.
func (t *Tree) Remove(s gfx.Spatial) bool {
	if t.Root == nil {
		return false
	}
	sb := s.Bounds()
	found := t.Root.find(sb)
	if found == nil {
		goto outside
	}
	for i, o := range found.Objects {
		if o == s {
			found.Objects = append(found.Objects[:i], found.Objects[i+1:]...)
			t.count--
			return true
		}
	}

outside:
	for i, o := range t.outside {
		if o == s {
			t.outside = append(t.outside[:i], t.outside[i+1:]...)
			t.count--
			return true
		}
	}
	return false
}

// New returns a new N-tree with the default options:
//  Divisor: math.Vec3{3, 3, 3}
//  MaxDepth: 8
//  MaxExpand: 8
//  StartScale: 1.0
func New() *Tree {
	t := &Tree{
		divisor:    math.Vec3{3, 3, 3},
		maxDepth:   8,
		maxExpand:  8,
		startScale: 2.0,
	}
	return t
}
