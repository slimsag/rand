package ntree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"fmt"
	"runtime"
	"sort"
)

var (
	work chan func()
)

// worker simply executes tasks in parralel.
func worker() {
	for {
		f, ok := <-work
		if !ok {
			return
		}
		f()
	}
}

func init() {
	n := runtime.GOMAXPROCS(-1)
	work = make(chan func(), n)
	for w := 0; w < n; w++ {
		go worker()
	}
}

var drop int

func rectRecurse(n *Node, r math.Rect3, results chan gfx.Spatial, stop chan struct{}, within bool) bool {
	// If the node's bounds do not even intersect with the search rectangle
	// then there is no point traversing the node further.
	//_, ok := n.bounds.Intersect(r)
	//fmt.Println("drop node", ok, n.bounds)
	//fmt.Println("drop rect", ok, r)
	if _, ok := n.bounds.Intersect(r); !ok {
		drop++
		fmt.Println("drop", drop)
		return true
	}

	// If the node's bounds are completely within the search rectangle we
	// do not need to search each individual object at all.
	if n.bounds.In(r) {
		for _, o := range n.Objects {
			// Wait to send until successful or stopped.
			select {
			case results <- o:
			case <-stop:
				return false
			}
		}
		goto children
	}

	// We must check each individual object's bounds to determine if it is
	// within or intersecting the search rectangle.
	for _, o := range n.Objects {
		ob := o.Bounds()
		if within && !ob.In(r) {
			// We only want objects within the rectangle, this is not one.
			continue
		} else if _, ok := ob.Intersect(r); !within && !ok {
			// We only want objects intersecting with the rectangle, this
			// is not one.
			continue
		}

		// Wait to send until successful or stopped.
		select {
		case results <- o:
			continue
		case <-stop:
			return false
		}
	}

children:
	// Recursively go deeper into the tree.
	for _, c := range n.Children {
		if !rectRecurse(c, r, results, stop, within) {
			return false
		}
	}
	return true
}

// rectSearch is the backend for both within and intersecting searches of the
// N tree. If within is true only objects within the rectangle are sent over
// the channel, otherwise any that intersect are.
func (t *Tree) rectSearch(r math.Rect3, results chan gfx.Spatial, stop chan struct{}, within bool) {
	// Perform a linear search across all of the objects outside of the N tree.
	for _, o := range t.outside {
		ob := o.Bounds()
		if within && !ob.In(r) {
			// We only want objects within the rectangle, this is not one.
			continue
		} else if _, ok := ob.Intersect(r); !within && !ok {
			// We only want objects intersecting with the rectangle, this
			// is not one.
			continue
		}

		// Wait to send until successful or stopped.
		select {
		case results <- o:
			continue
		case <-stop:
			close(results)
			return
		}
	}

	// Begin searching at the root node.
	rectRecurse(t.Root, r, results, stop, within)

	select {
	case <-stop:
	default:
	}
	close(results)
}

// In performs a search of the N tree to find all spatial objects that are
// completely contained within the given rectangle.
//
// This function will return immedietly and the search will be executed in
// parralel. Results will be sent over the given results channel (e.g. with a
// buffer size of 32) until the search is completed or halted. The channel will
// be closed when the last result is sent.
//
// If the stop channel is not nil, then when a struct{}{} is sent over the stop
// channel the search will be indefinitely halted.
func (t *Tree) In(r math.Rect3, results chan gfx.Spatial, stop chan struct{}) {
	work <- func() {
		t.rectSearch(r, results, stop, true)
	}
}

// Intersect performs a search of the N tree to find all spatial objects that
// are intersecting with the given rectangle.
//
// This function will return immedietly and the search will be executed in
// parralel. Results will be sent over the given results channel (e.g. with a
// buffer size of 32) until the search is completed or halted. The channel will
// be closed when the last result is sent.
//
// If the stop channel is not nil, then when a struct{}{} is sent over the stop
// channel the search will be indefinitely halted.
func (t *Tree) Intersect(r math.Rect3, results chan gfx.Spatial, stop chan struct{}) {
	work <- func() {
		t.rectSearch(r, results, stop, false)
	}
}

func closestObject(target math.Vec3, objs []gfx.Spatial) (closest gfx.Spatial, index int) {
	var (
		closestDist float64
	)
	for i, s := range objs {
		sDist := s.Bounds().Closest(target).Sub(target).LengthSq()
		if i == 0 || sDist < closestDist {
			closestDist = sDist
			closest = s
			index = i
		}
	}
	return
}

func collectNodes(n *Node, nodes []gfx.Spatial) []gfx.Spatial {
	nodes = append(nodes, n)
	for _, c := range n.Children {
		nodes = collectNodes(c, nodes)
	}
	return nodes
}

// Nearest performs a search of the N tree to find all spatial objects that are
// closest to the given point.
//
// This function will return immedietly and the search will be executed in
// parralel. Results will be sent over the given results channel (e.g. with a
// buffer size of 32) in order of closest to furthest away until the search is
// completed or halted. The channel will be closed when the last result is
// sent.
//
// If the stop channel is not nil, then when a struct{}{} is sent over the stop
// channel the search will be indefinitely halted.
func (t *Tree) Nearest(p math.Vec3, results chan gfx.Spatial, stop chan struct{}) {
	work <- func() {
		// Collect all of the nodes and sort them by nearest to p.
		nByDist := gfx.ByDist{
			Objects: collectNodes(t.Root, nil),
			Target:  p,
		}
		sort.Sort(nByDist)

		sortResults := make(chan []gfx.Spatial)
		// Sort each node's spatial objects by nearest to p.
		for _, s := range nByDist.Objects {
			go func() {
				n := s.(*Node)
				sByDist := gfx.ByDist{
					Objects: n.Objects,
					Target:  p,
				}
				sort.Sort(sByDist)
				sortResults <- sByDist.Objects
			}()
		}

		for rn := 0; rn < len(nByDist.Objects); rn++ {
			r := <-sortResults
			// Wait to send until successful or stopped.
			for _, o := range r {
				select {
				case results <- o:
				case <-stop:
					goto end
				}
			}
		}

	end:
		select {
		case <-stop:
		default:
		}
		close(results)
		return
		// FIXME: t.outside...
	}
}
