// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package octree implements a truly dynamic octree.
//
// Octrees normally use a fixed size for the tree specified at creation time,
// this octree implementation does not require that and instead employs a
// method of dynamically expanding the octree by replacing the root node with
// a new one twice the size of the previous one, such that the old root becomes
// a child octant of this new root.
//
// The octree implementation allows traversing the octree as well as high-level
// searches for objects intersecting or completely contained within some
// defined space (a 3D rectangle, sphere, or viewing frustum).
//
// TODO: add Contains
// TODO: add Intersects
// TODO: nearest to point, nearest to ... Distancer?
package octree
