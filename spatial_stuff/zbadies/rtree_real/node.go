package rtree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
)

type Node struct {
	// The parent of this node.
	Parent *Node

	// Bounds represents the bounding rectangle of this node and all entries
	// below it.
	Bounds math.Rect3

	// The height of this node in the tree.
	Height int

	// A slice of child node entries.
	Entries []*Node

	// An object, if any.
	Object gfx.Spatial
}

func (n *Node) IsLeaf() bool {
	return len(n.Entries) == 0
}

func (n *Node) enlargmentTo(r math.Rect3) float64 {
	// Find the rectangle that can fit both n.Bounds and r.
	u := n.Bounds.Union(r)

	// Grab it's squared length from it's min and max points.
	// XXX: Right order?
	return u.Max.Sub(u.Min).LengthSq()
}

// calcBounds calculates a bounding box needed to encapsulate all the entries
// of this node and returns it.
func (n *Node) calcBounds() math.Rect3 {
	b := n.Bounds
	for _, e := range n.Entries {
		b = b.Union(e.Bounds)
	}
	return b
}

// ChooseLeaf selects a leaf node in which to place a new index entry, e.
func (n *Node) ChooseLeaf(e math.Rect3) *Node {
	if n.IsLeaf() {
		// If N is a leaf node, return it.
		return n
	}
	// Choose the subtree. If N is not a leaf node, let F be the entry in N
	// whose rectangle needs the least enlargment to include E. Resolve ties by
	// choosing the entry with the rectangle of smallest area.
	var (
		f             *Node
		enlargmentToF float64
	)
choosing:
	for _, entry := range n.Entries {
		if f == nil {
			f = entry
			enlargmentToF = f.enlargmentTo(e)
			continue choosing
		}
		enlargmentToEntry := entry.enlargmentTo(e)
		if enlargmentToF == enlargmentToEntry {
			// Tie. Resolve by choosing the entry with the rectangle of the
			// smallest area.
			if f.Bounds.Area() < entry.Bounds.Area() {
				continue choosing
			}
			f = entry
			enlargmentToF = enlargmentToEntry
			continue choosing
		}
		if enlargmentToEntry < enlargmentToF {
			// The enlargment needed from entry-to-e is smaller than the
			// enlargment needed from f-to-e.
			f = entry
			enlargmentToF = enlargmentToEntry
			continue choosing
		}
	}

	// Descend until a leaf is reached.
	return f.ChooseLeaf(e)
}
