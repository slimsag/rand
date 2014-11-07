// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/gfx.v1"
)

// Container describes any object that can answer if a given boundable (octree
// object or node) is contained within the arbitrary search area defined by
// the interface.
type Container interface {
	Intersector

	// Contains should test if the given boundable object is completely within
	// the search area and return the result.
	//
	// This method must be safe to call from multiple goroutines concurrently.
	Contains(b gfx.Boundable) bool
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
}*/
