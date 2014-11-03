package octree

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"testing"
)

func less(min, max math.Vec3) bool {
	if min.X > max.X || min.Y > max.Y || min.Z > max.Z {
		return false
	}
	return true
}

func TestChildFits(t *testing.T) {
	o := &Octant{
		Bounds: gfx.AABB{
			Min: math.Vec3{0, 0, 0},
			Max: math.Vec3{1, 1, 1},
		},
	}

	try := func(tgt int, name string, want gfx.AABB) {
		if !less(want.Min, want.Max) {
			t.Log("want: Invalid AABB")
			t.Fail()
		}
		idx := o.childFits(want)
		if idx != tgt {
			t.Log("Invalid index for", name, idx, "want", tgt)
			t.Fail()
		}
	}

	try(0, "Top, Front, Left", gfx.AABB{
		Min: math.Vec3{0, .5, .5},
		Max: math.Vec3{.5, 1, 1},
	})

	try(1, "Top, Front, Right", gfx.AABB{
		Min: math.Vec3{.5, .5, .5},
		Max: math.Vec3{1, 1, 1},
	})

	try(2, "Top, Back, Left", gfx.AABB{
		Min: math.Vec3{0, 0, .5},
		Max: math.Vec3{.5, .5, 1},
	})

	try(3, "Top, Back, Right", gfx.AABB{
		Min: math.Vec3{.5, 0, .5},
		Max: math.Vec3{1, .5, 1},
	})

	try(4, "Bottom, Front, Left", gfx.AABB{
		Min: math.Vec3{0, .5, 0},
		Max: math.Vec3{.5, 1, .5},
	})

	try(5, "Bottom, Front, Right", gfx.AABB{
		Min: math.Vec3{.5, .5, 0},
		Max: math.Vec3{1, 1, .5},
	})

	try(6, "Bottom, Back, Left", gfx.AABB{
		Min: math.Vec3{0, 0, 0},
		Max: math.Vec3{.5, .5, .5},
	})

	try(7, "Bottom, Back, Right", gfx.AABB{
		Min: math.Vec3{.5, 0, 0},
		Max: math.Vec3{1, .5, .5},
	})
}

func TestExpansion(t *testing.T) {
	try := func(into int, name string, want gfx.AABB) {
		if !less(want.Min, want.Max) {
			t.Log("want: Invalid AABB")
			t.Fail()
		}

		o := &Octant{
			Bounds: gfx.AABB{
				Min: math.Vec3{0, 0, 0},
				Max: math.Vec3{1, 1, 1},
			},
		}

		o.expand(into)

		for _, c := range o.Children {
			if c != nil && c.Depth != o.Depth+1 {
				t.Log("Invalid depth found.")
				t.Log("got", o.Depth)
				t.Log("want", o.Depth+1)
				t.Fail()
			}
		}

		if !less(o.Bounds.Min, o.Bounds.Max) {
			t.Log("expanded bounds: Invalid AABB")
			t.Fail()
		}
		if !o.Bounds.Equals(want) {
			t.Log(name)
			t.Log("got", o.Bounds)
			t.Log("want", want)
			t.Fail()
		}
	}

	// Determine the expansion amounts.
	leftMin, leftMax := 0.0, 2.0
	rightMin, rightMax := -1.0, 1.0
	frontMin, frontMax := -1.0, 1.0
	backMin, backMax := 0.0, 2.0
	topMin, topMax := -1.0, 1.0
	bottomMin, bottomMax := 0.0, 2.0

	try(0, "0 - Top, Front, Left", gfx.AABB{
		Min: math.Vec3{leftMin, frontMin, topMin},
		Max: math.Vec3{leftMax, frontMax, topMax},
	})

	try(1, "1 - Top, Front, Right", gfx.AABB{
		Min: math.Vec3{rightMin, frontMin, topMin},
		Max: math.Vec3{rightMax, frontMax, topMax},
	})

	try(2, "2 - Top, Back, Left", gfx.AABB{
		Min: math.Vec3{leftMin, backMin, topMin},
		Max: math.Vec3{leftMax, backMax, topMax},
	})

	try(3, "3 - Top, Back, Right", gfx.AABB{
		Min: math.Vec3{rightMin, backMin, topMin},
		Max: math.Vec3{rightMax, backMax, topMax},
	})

	try(4, "4 - Bottom, Front, Left", gfx.AABB{
		Min: math.Vec3{leftMin, frontMin, bottomMin},
		Max: math.Vec3{leftMax, frontMax, bottomMax},
	})

	try(5, "5 - Bottom, Front, Right", gfx.AABB{
		Min: math.Vec3{rightMin, frontMin, bottomMin},
		Max: math.Vec3{rightMax, frontMax, bottomMax},
	})

	try(6, "6 - Bottom, Back, Left", gfx.AABB{
		Min: math.Vec3{leftMin, backMin, bottomMin},
		Max: math.Vec3{leftMax, backMax, bottomMax},
	})

	try(7, "7 - Bottom, Back, Right", gfx.AABB{
		Min: math.Vec3{rightMin, backMin, bottomMin},
		Max: math.Vec3{rightMax, backMax, bottomMax},
	})
}

func TestChildBounds(t *testing.T) {
	try := func(childIndex int, name string, want gfx.AABB) {
		if !less(want.Min, want.Max) {
			t.Log("want: Invalid AABB")
			t.Fail()
		}
		o := &Octant{
			Bounds: gfx.AABB{
				Min: math.Vec3{5, 5, 5},
				Max: math.Vec3{15, 15, 15},
			},
		}

		cb := o.childBounds(childIndex)
		if !less(cb.Min, cb.Max) {
			t.Log("child bounds: Invalid AABB")
			t.Fail()
		}
		if !cb.Equals(want) {
			t.Log(name)
			t.Log("got", cb)
			t.Log("want", want)
			t.Fail()
		}
	}

	// Determine the child bounds amounts.
	leftMin, leftMax := 5.0, 10.0
	rightMin, rightMax := 10.0, 15.0
	frontMin, frontMax := 10.0, 15.0
	backMin, backMax := 5.0, 10.0
	topMin, topMax := 10.0, 15.0
	bottomMin, bottomMax := 5.0, 10.0

	try(0, "0 - Top, Front, Left", gfx.AABB{
		Min: math.Vec3{leftMin, frontMin, topMin},
		Max: math.Vec3{leftMax, frontMax, topMax},
	})

	try(1, "1 - Top, Front, Right", gfx.AABB{
		Min: math.Vec3{rightMin, frontMin, topMin},
		Max: math.Vec3{rightMax, frontMax, topMax},
	})

	try(2, "2 - Top, Back, Left", gfx.AABB{
		Min: math.Vec3{leftMin, backMin, topMin},
		Max: math.Vec3{leftMax, backMax, topMax},
	})

	try(3, "3 - Top, Back, Right", gfx.AABB{
		Min: math.Vec3{rightMin, backMin, topMin},
		Max: math.Vec3{rightMax, backMax, topMax},
	})

	try(4, "4 - Bottom, Front, Left", gfx.AABB{
		Min: math.Vec3{leftMin, frontMin, bottomMin},
		Max: math.Vec3{leftMax, frontMax, bottomMax},
	})

	try(5, "5 - Bottom, Front, Right", gfx.AABB{
		Min: math.Vec3{rightMin, frontMin, bottomMin},
		Max: math.Vec3{rightMax, frontMax, bottomMax},
	})

	try(6, "6 - Bottom, Back, Left", gfx.AABB{
		Min: math.Vec3{leftMin, backMin, bottomMin},
		Max: math.Vec3{leftMax, backMax, bottomMax},
	})

	try(7, "7 - Bottom, Back, Right", gfx.AABB{
		Min: math.Vec3{rightMin, backMin, bottomMin},
		Max: math.Vec3{rightMax, backMax, bottomMax},
	})
}
