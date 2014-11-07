// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package octree

import (
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

// Traverse performs a parralel traversal of the octree in multiple goroutines.
// It begins at the root octree node and then traverses recursively down into
// all child nodes. At each node traversed it invokes the traverser callback,
// and if that callback returns a non-nil traverser function then it is invoked
// to handle that node's children.
//
// This function does not return until the traversal is finished completely.
func (t *Tree) Traverse(callback Traverser) {
	t.RLock()
	wg := new(sync.WaitGroup)
	t.root.traverse(callback, wg)
	wg.Wait()
	t.RUnlock()
}
