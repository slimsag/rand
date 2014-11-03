// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
	"azul3d.org/v1/gfx"
	"runtime"
	"sync"
	"time"
)

var (
	workers struct {
		sync.RWMutex
		count int
	}
	workChan = make(chan func())
	timeout  = 5 * time.Second
)

// Worker executes functions sent over the work channel. If no work is received
// after the timeout then the goroutine exits and a new worker will be spawned
// later.
func worker(f func()) {
	workers.Lock()
	workers.count++
	workers.Unlock()
	f()
	for {
		select {
		case fn := <-workChan:
			fn()
		case <-time.After(timeout):
			workers.Lock()
			workers.count--
			workers.Unlock()
			return
		}
	}
}

// work submits work to a worker goroutine. If no worker goroutine exists then
// a new one is spawned, but only if max == 0 or the number of workers is less
// than max, otherwise f() is executed in this goroutine.
func work(f func(), max int) {
	select {
	case workChan <- f:

	default:
		workers.RLock()
		spawn := max == 0 || workers.count < max
		workers.RUnlock()
		if spawn {
			go worker(f)
		} else {
			f()
		}
	}
}

var maxTraversalWorkers = -1

func (n *Node) traverse(callback Traverser, wg *sync.WaitGroup) {
	if maxTraversalWorkers == -1 {
		maxTraversalWorkers = runtime.GOMAXPROCS(-1)
	}
	wg.Add(1)
	work(func() {
		if callback(n) == nil {
			// Stop traversing this node.
			goto end
		}

		n.access.RLock()
		for c := 0; c < 8; c++ {
			child := n.children[c]
			if child == nil {
				continue
			}
			child.traverse(callback, wg)
		}
		n.access.RUnlock()

	end:
		wg.Done()
	}, maxTraversalWorkers)
}

// Traverser is used to traverse the octree. It is used by the Tree.Traverse
// method.
type Traverser func(n *Node) Traverser

// Traverse begins a parralel traversal of the octree. It starts at the root of
// the tree and traverses downward recursively into each octant. At each octant
// node found it invokes the callback. If the callback returns non-nil then the
// returned function is used to handle the traversal of that node's children.
//
// This function does not return until the traversal is finished completely.
func (t *Tree) Traverse(callback Traverser) {
	t.RLock()
	wg := new(sync.WaitGroup)
	t.root.traverse(callback, wg)
	wg.Wait()
	t.RUnlock()
}

// In performs a search on the octree for objects within the search area
// defined by s. The search is executed in parralel and this function returns
// immedietly.
//
// The searcher s must be a Container -- any other searcher will cause a panic.
//
// The results channel has the valid search results sent over it, and when the
// search finishes the channel is closed.
//
// If non-nil, the stop channel can be used to halt the search permanently.
func (t *Tree) In(s Searcher, results chan gfx.Boundable, stop chan struct{}) {
	container, ok := s.(Container)
	if !ok {
		panic("In(): invalid searcher: need Container")
	}
	work(func() {
		t.containerSearch(container, results, stop)
	}, 0)
}

// Intersect performs a search on the octree for objects intersecting the
// search area defined by s. The search is executed in parralel and this
// function returns immedietly.
//
// Intersection searching can be more efficient if a Container is used -- but
// both Container and Intersector searches are accepted (any other searcher
// will cause a panic).
//
// The results channel has the valid search results sent over it, and when the
// search finishes the channel is closed.
//
// If non-nil, the stop channel can be used to halt the search permanently.
func (t *Tree) Intersect(s Searcher, results chan gfx.Boundable, stop chan struct{}) {
	switch s.(type) {
	case Container, Intersector:
	default:
		panic("Intersect(): invalid searcher: need Container or Intersector")
	}
	work(func() {
		t.intersectorSearch(s, results, stop)
	}, 0)
}

// Closest performs a search on the octree for objects closest to the search
// area defined by s. The search is executed in parralel and this function
// returns immedietly.
//
// The searcher s must be a Distancer -- any other searcher will cause a panic.
//
// The results channel has the valid search results sent over it, and when the
// search finishes the channel is closed.
//
// If non-nil, the stop channel can be used to halt the search permanently.
func (t *Tree) Closest(s Searcher, results chan gfx.Boundable, stop chan struct{}) {
	distancer, ok := s.(Distancer)
	if !ok {
		// FIXME:
		//panic("Closest(): invalid searcher: need Distancer")
	}
	work(func() {
		t.distancerSearch(distancer, results, stop)
	}, 0)
}

// Searcher defines a single searcher. It's value type must be Intersector or Container or
// Distancer.
type Searcher interface{}

type notContainer struct {
	Intersector
	actual Container
}

func (c notContainer) Contains(b gfx.Boundable) bool {
	return !c.actual.Contains(b)
}

type notIntersector struct {
	actual Intersector
}

func (i notIntersector) Intersects(b gfx.Boundable) bool {
	return !i.actual.Intersects(b)
}

// Not returns a negated form of s. It panics if s is not a valid searcher. It
// is useful for doing opposite-searches like so:
//  // Not within the rectangle r:
//  tree.In(Not(Rect3(r)), results, stop)
//
//  // Not intersecting the rectangle r:
//  tree.Intersect(Not(Rect3(r)), results, stop)
func Not(s Searcher) Searcher {
	switch t := s.(type) {
	case Intersector:
		return notIntersector{
			actual: t,
		}

	case Container:
		return notContainer{
			Intersector: notIntersector{
				actual: Intersector(t),
			},
			actual: t,
		}

	default:
		panic("Not(): Invalid searcher")
	}
}

type andContainer struct {
	Intersector
	a, b Container
}

func (c andContainer) Contains(b gfx.Boundable) bool {
	return c.a.Contains(b) && c.b.Contains(b)
}

type andIntersector struct {
	a, b Intersector
}

func (i andIntersector) Intersects(b gfx.Boundable) bool {
	return i.a.Intersects(b) && i.b.Intersects(b)
}

// And returns a searcher which performs both a and b. It panics if a or b is
// not a valid searcher, or if they are not the same type of searcher. It is
// useful for combining searchers:
//  // All objects within the sphere s and not in the rectangle r.
//  tree.In(And(Sphere(s), Not(Rect3(r))), results, stop)
//
//  // All objects intersecting the sphere s and not intersecting the rectangle r.
//  tree.In(And(Sphere(s), Not(Rect3(r))), results, stop)
func And(a, b Searcher) Searcher {
	switch t := a.(type) {
	case Intersector:
		bt, ok := b.(Intersector)
		if !ok {
			panic("And(): searchers are of different types")
		}
		return andIntersector{
			a: t,
			b: bt,
		}

	case Container:
		bt, ok := b.(Container)
		if !ok {
			panic("And(): searchers are of different types")
		}
		return andContainer{
			Intersector: andIntersector{
				a: Intersector(t),
				b: Intersector(bt),
			},
			a: t,
			b: bt,
		}

	default:
		panic("Not(): Invalid searcher")
	}
}
