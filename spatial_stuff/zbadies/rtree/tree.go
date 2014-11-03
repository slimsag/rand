package rtree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
)

type Node struct {
	parent *Node

	// The bounding box of this node and all child nodes below it.
	bounds math.Rect3

	// A list of children nodes.
	Children []*Node

	// A list of objects stored in this node.
	Objects []gfx.Spatial
}

// Implements gfx.Spatial interface.
func (n *Node) Bounds() math.Rect3 {
	return n.bounds
}

func (n *Node) calcBounds() math.Rect3 {
	b := n.bounds
	for i, c := range n.Children {
		if i == 0 {
			b = c.bounds
			continue
		}
		b = b.Union(c.bounds)
	}
	for _, o := range n.Objects {
		b = b.Union(o.Bounds())
	}
	return b
}

func (n *Node) fits(sb math.Rect3) bool {
	if sb.In(n.bounds) {
		return true
	}
	return false
}

func (n *Node) findNode(sb math.Rect3) *Node {
	// First try to place it into this node.
	if n.fits(sb) {
		// Try our children, recursively.
		for _, c := range n.Children {
			childFound := c.findNode(sb)
			if childFound != nil {
				return childFound
			}
		}
		return n
	}
	return nil
}

func (n *Node) overlap(b math.Rect3) (o float64) {
	if n.parent != nil {
		i, ok := n.parent.bounds.Intersect(b)
		if ok {
			o += i.Size().X
		}
		//o += n.parent.overlap(b)
	}
	return
}

func (n *Node) split(newRoot *Node) (a, b *Node) {
	a = &Node{parent: n}
	b = &Node{parent: n}

	// Count the number of objects that can fit in each split axis.
	var (
		nb                                                   = n.bounds
		axCount, bxCount, ayCount, byCount, azCount, bzCount int
		ax, bx, ay, by, az, bz                               = nb, nb, nb, nb, nb, nb
	)
	half := nb.Size().DivScalar(2)
	ax.Min.X = ax.Min.X + half.X
	bx.Max.X = bx.Max.X - half.X
	ay.Min.Y = ay.Min.Y + half.Y
	by.Max.Y = by.Max.Y - half.Y
	az.Min.Z = ay.Min.Z + half.Z
	bz.Max.Z = by.Max.Z - half.Z

	for _, c := range n.Objects {
		cb := c.Bounds()
		if cb.In(ax) {
			axCount++
		}
		if cb.In(bx) {
			bxCount++
		}
		if cb.In(ay) {
			ayCount++
		}
		if cb.In(by) {
			byCount++
		}
		if cb.In(az) {
			azCount++
		}
		if cb.In(bz) {
			bzCount++
		}
	}

	// Choose the best split axis.
	var (
		aBest, bBest           = ax, bx
		aBestCount, bBestCount = axCount, bxCount
	)
	if ayCount > aBestCount {
		aBest = ay
		aBestCount = ayCount
	}
	if azCount > aBestCount {
		aBest = az
		aBestCount = azCount
	}
	if byCount > bBestCount {
		bBest = by
		bBestCount = byCount
	}
	if bzCount > bBestCount {
		bBest = bz
		bBestCount = bzCount
	}

	// Place the objects in the best axis bounds.
	a.bounds = aBest
	b.bounds = bBest
	for _, c := range n.Objects {
		cb := c.Bounds()
		if cb.In(a.bounds) {
			a.Objects = append(a.Objects, c)
		} else if cb.In(b.bounds) {
			b.Objects = append(b.Objects, c)
		} else {
			//newRoot.Objects = append(newRoot.Objects, c)
			//newRoot.place(cb, c)
			//n.root.place(c)
		}
	}
	return
}

/*
func (n *Node) split(newRoot *Node) (a, b *Node) {
	a = &Node{parent: n}
	b = &Node{parent: n}

	// Count the number of objects that can fit in each split axis.
	var(
		nb = n.bounds
		axCount, bxCount, ayCount, byCount, azCount, bzCount int
		ax, bx, ay, by, az, bz = nb, nb, nb, nb, nb, nb
	)
	half := nb.Size().DivScalar(2)
	ax.Min.X = ax.Min.X + half.X
	bx.Max.X = bx.Max.X - half.X
	ay.Min.Y = ay.Min.Y + half.Y
	by.Max.Y = by.Max.Y - half.Y
	az.Min.Z = ay.Min.Z + half.Z
	bz.Max.Z = by.Max.Z - half.Z

	for _, c := range n.Objects {
		cb := c.Bounds()
		if cb.In(ax) {
			axCount++
		}
		if cb.In(bx) {
			bxCount++
		}
		if cb.In(ay) {
			ayCount++
		}
		if cb.In(by) {
			byCount++
		}
		if cb.In(az) {
			azCount++
		}
		if cb.In(bz) {
			bzCount++
		}
	}

	// Choose the best split axis.
	var(
		aBest, bBest = ax, bx
		aBestCount, bBestCount = axCount, bxCount
	)
	if ayCount > aBestCount {
		aBest = ay
		aBestCount = ayCount
	}
	if azCount > aBestCount {
		aBest = az
		aBestCount = azCount
	}
	if byCount > bBestCount {
		bBest = by
		bBestCount = byCount
	}
	if bzCount > bBestCount {
		bBest = bz
		bBestCount = bzCount
	}

	// Place the objects in the best axis bounds.
	a.bounds = aBest
	b.bounds = bBest
	for _, c := range n.Objects {
		cb := c.Bounds()
		if cb.In(a.bounds) {
			a.Objects = append(a.Objects, c)
		} else if cb.In(b.bounds) {
			b.Objects = append(b.Objects, c)
		} else {
			newRoot.place(cb, c)
			//n.root.place(c)
		}
	}
	return
}*/

func (n *Node) place(sb math.Rect3, s gfx.Spatial) {
	// Expand existing bounds.
	n.bounds = n.bounds.Union(sb)

	// Append the object.
	n.Objects = append(n.Objects, s)

	if len(n.Objects) > 256 {
		newRoot := Node{
			parent:   n.parent,
			Children: n.Children,
		}
		a, b := n.split(&newRoot)
		newRoot.Children = append(newRoot.Children, a)
		newRoot.Children = append(newRoot.Children, b)
		*n = newRoot
		n.bounds = n.calcBounds() // FIXME
	}
}

type Tree struct {
	// The root node of the tree.
	root *Node

	// The number of objects in the tree currently.
	count int
}

// Root returns the root node of the tree. It returns nil if there are no
// objects in the tree yet.
func (t *Tree) Root() *Node {
	return t.root
}

func (t *Tree) Count() int {
	return t.count
}

// Add adds the spatial to the tree.
func (t *Tree) Add(s gfx.Spatial) {
	t.count++
	sb := s.Bounds()
	if sb.Empty() {
		panic("Tree.Add(): cannot add empty spatial")
	}

	if t.root == nil {
		t.root = &Node{
			parent:  t.root,
			bounds:  sb,
			Objects: []gfx.Spatial{s},
		}
		return
	}

	n := t.root.findNode(sb)
	if n != nil {
		n.place(sb, s)
		return
	}

	t.root.place(sb, s)
	t.root.sanitizeBounds()
	return
}

func (n *Node) sanitizeBounds() {
	for _, c := range n.Children {
		c.sanitizeBounds()
	}
	if !n.bounds.Equals(n.calcBounds()) {
		panic("insane bounds")
	}
}

func New() *Tree {
	return &Tree{}
}
