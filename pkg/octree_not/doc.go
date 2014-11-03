// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package octree implements a truly dynamic octree.
//
// Most octrees use a fixed size specified at creation time, but this octree
// implementation allows for dynamic expansion of the octree by replacing the
// root node with a new one twice the size of the previous one.
//
// The octree implementation allows traversing the octree as well as high-level
// searches for objects intersecting or completely in some defined space (a 3D
// rectangle, sphere, or viewing frustum (projection) matrices).
package octree
