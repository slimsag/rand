package rtree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
)

// Tree represents an rectangle tree.
type Tree struct {
	Min, Max int
	Root     *Node
}

// Search searches the rectangle tree for any entries that have intersecting
// rectangles with the given search rectangle, r.
//
// Each time an entry is found to intersect with r the callback is executed, if
// the callback returns true then the search is continued otherwise the search
// immedietly returns.
func (t *Tree) Search(r math.Rect3, callback func(s gfx.Spatial) bool) {
	if !t.Root.IsLeaf() {
		for _, e := range t.Root.Entries {
			_, intersects := e.Bounds.Intersect(r)
			if intersects {
				t.Search(r, callback)
			}
		}
	}
	if !callback(t.Root.Object) {
		return
	}
}

// Insert inserts the given entry into the tree, splitting the tree as needed.
func (t *Tree) Insert(e gfx.Spatial) {
	var l, ll *Node

	// Select a leaf node, l, in which to place the entry, e.
	eBounds := e.Bounds()
	l = t.Root.ChooseLeaf(eBounds)

	// Add the entry to the leaf node. If L has room for another entry, install
	// it. Otherwise invoke SplitNode to obtain L and LL containing E and all
	// of the old entries of L.
	if len(l.Entries) < t.Max {
		// We have room to store the entry.
		l.Entries = append(l.Entries, &Node{
			Parent: l,
			Bounds: eBounds,
			Height: t.Root.Height,
			Object: e,
		})
	} else {
		// We must split the node.
		l, ll = l.SplitNode()
	}

	// Propagate the changes to the tree upward. Invoke AdjustTree on L, also
	// passing LL if a split was performed.
	// XXX: AdjustTree must check if ll != nil to see if a split was performed.
	root, splitRoot := t.AdjustTree(l, ll)
	if splitRoot != nil {
		// Grow the tree taller. If node split propagation caused the root to
		// split then create a new root whose children are the two resulting
		// nodes.
		oldRoot := root
		t.Root.Height++
		t.Root = &Node{
			Parent: nil, // XXX: It's the root but is nil okay?
			Height: t.Root.Height,
		}
		t.Root.Entries = []*Node{
			&Node{
				Parent: t.Root,
				Bounds: splitRoot.calcBounds(),
				Height: t.Root.Height,
			},
			&Node{
				Parent: t.Root,
				Bounds: oldRoot.calcBounds(),
				Height: t.Root.Height,
			},
		}
		t.Root.Bounds = t.Root.calcBounds()
	}
}

// Ascends from a leaf node, l, to the root adjusting covering rectangles and
// propagating node splits as necessary.
func (t *Tree) AdjustTree(n, nn *Node) (root, splitRoot *Node) {
	// Set N=L, if L was split set NN to the resulting second node.

	if n == t.Root {
		if nn != nil {
			// Build a new root and add children.
		}
		/*
		   if (nn != null)
		   {
		     // build new root and add children.
		     root = buildRoot(false);
		     root.children.add(n);
		     n.parent = root;
		     root.children.add(nn);
		     nn.parent = root;
		   }
		   tighten(root);
		   return;
		*/
		// If N is the root, stop.
		return n, nn
	}

	// Adjust the covering rectangle in the parent entry. Let P be the parent
	// node of N, and let En be N's entry in P. Adjust En-I so that it tightly
	// encloses all entry rectangles in N.
	t.Tighten(n)
	if nn != nil {
		t.Tighten(nn)
		if len(n.Parent.Children) > t.MaxEntries {
			s0, s1 := splitNode(n.Parent)
			t.AdjustTree(s0, s1)
		}
	}
	if n.Parent != nil {
		t.AdjustTree(n.Parent, nil)
	}
}

// New returns a new initialized *Tree.
func New(min, max int) *Tree {
	if min <= 0 || min > max {
		panic("min <= 0 || min > max")
	}
	if max <= 0 {
		panic("max <= 0")
	}
	return &Tree{
		Root: new(Node), // XXX: What about bounds?
		Min:  min,
		Max:  max,
	}
}
