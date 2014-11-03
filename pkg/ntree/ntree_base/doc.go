// Package ntree implements a generic 3D N tree.
//
// A N tree is a spatial indexing data structure similar to a Quadtree or
// Octree, in fact an N tree can represent many different types of spatial
// tree's that resolve around N cubes in each dimension by changing the divisor
// of the tree (it can effectively represent a Quadtree and also an Octree, but
// by default it represents a viginti septum tree which has three cubes in each
// axis direction).
//
// Unlike Quadtree's and Octree's the N tree is fully capable of expansion
// allowing it to consume areas that are largely outside of the tree.
//
// Additionally a N tree is more memory efficient than a Quadtree or Octree due
// to the fact that subdivided nodes are created and deleted as needed to
// encapsulate spatial objects, whereas most Quadtree's and Octree's allocate
// all child nodes whenever a subdivision occurs.
//
// Limits can be imposed on both the depth (how far the tree can be subdivided)
// and expansion (how far the tree can expand outward to encapsulate spatial
// objects residing outside the tree). After expansion occurs if spatial
// objects still reside outside the tree then they are placed in a linear list
// (a slice) such that they still remain functional within the tree.
package ntree
