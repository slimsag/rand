// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"log"
	"sort"
)

// Distancer is any object that can answer how far an object is away from an
// arbitrary search point defined by this interface.
type Distancer interface {
	// Distance should calculate and return how far the given boundable object
	// is away from the search point and return the result. Distance is not
	// measured in any particular units, instead it is measured in arbitrary
	// units.
	//
	// This method must be safe to call from multiple goroutines concurrently.
	Distance(b gfx.Boundable) float64
}

func (t *Tree) distancerSearch(search Distancer, results chan gfx.Boundable, stop chan struct{}) {
	sendResult := func(r gfx.Boundable) (stopSearch bool) {
		select {
		case results <- r:
		case <-stop:
			return false
		}
		return true
	}

	lastRadius := 0.0
	radius := 0.1

	p := math.Vec3{-.2, -.2, -.1}
	var s chan struct{}

	buf := byDist{
		Objects: make([]gfx.Boundable, 0, 128),
		Target:  p,
	}
	sphere := math.Sphere{
		Radius: radius,
		Center: p,
	}
	//lastSphere := math.Sphere{
	//	Radius: lastRadius,
	//	Center: p,
	//}
	for {
		r := make(chan gfx.Boundable, 128)
		s = make(chan struct{}, 1)
		//t.Intersect(And(Sphere(sphere), Not(Sphere(lastSphere))), r, s)
		t.Intersect(Not(Sphere(sphere)), r, s)

		// Load into the buffer.
		for result := range r {
			buf.Objects = append(buf.Objects, result)
		}
		log.Printf("objects=%d last=%f now=%f", len(buf.Objects), lastRadius, radius)

		// Sort the buffer.
		sort.Sort(buf)

		for _, result := range buf.Objects {
			if !sendResult(result) {
				goto end
			}
		}
		buf.Objects = buf.Objects[:0]

		// Expand radius and search again.
		//lastRadius = radius
		radius += 1.1
		if radius > 10.0 {
			goto end
		}
	}

end:
	s <- struct{}{}
	close(results)
}

/*
func (t *Tree) containerSearch(search Container, results chan gfx.Boundable, stop chan struct{}) {
	sendResult := func(r gfx.Boundable) (stopSearch bool) {
		select {
		case results <- r:
		case <-stop:
			return false
		}
		return true
	}

	var running struct {
		sync.Mutex
		stopped bool
	}

	// Tells if the search is stopped.
	stopped := func() bool {
		running.Lock()
		select {
		case <-stop:
			running.stopped = true
		default:
		}
		stopped := running.stopped
		running.Unlock()
		return stopped
	}

	// sends all results without testing as if they where all valid.
	var sendAllResults func(n *Node) Traverser
	sendAllResults = func(n *Node) Traverser {
		for oct := 0; oct < 9; oct++ {
			nObjects := n.NumObjects(oct)
			for o := 0; o < nObjects; o++ {
				obj := n.Object(oct, o)
				if !sendResult(obj) {
					return nil
				}
			}
		}
		return sendAllResults
	}

	var trav Traverser
	trav = func(n *Node) (t Traverser) {
		if stopped() {
			return nil
		}

		// If the node is not at all intersecting, then there is no need to
		// continue traversing this node.
		if !search.Intersects(n) {
			return nil
		}

		// If the node is completely contained, then all of it's children are valid
		// results.
		if search.Contains(n) {
			return sendAllResults(n)
		}

		// Test each one of this node's objects to see if it is a valid result.
		for oct := 0; oct < 9; oct++ {
			nObjects := n.NumObjects(oct)
			for o := 0; o < nObjects; o++ {
				obj := n.Object(oct, o)

				// If the object is not completely contained, then it is not a
				// valid search result.
				if !search.Contains(obj) {
					continue
				}

				if !sendResult(obj) {
					return nil
				}
			}
		}

		// Continue searching child octants.
		return trav
	}
	t.Traverse(trav)
	close(results)
}
*/
