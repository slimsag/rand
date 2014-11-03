// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"sync"
)

// Transformable represents a generic interface to any object that can return a
// transformation matrix.
type Transformable interface {
	Mat4() math.Mat4
}

// CoordSpace represents a single coordinate space conversion.
//
// World space is the top-most 'world' or 'global' space. A transform whose
// parent is nil explicitly means the parent is the 'world'.
//
// Local space is the local space that this transform defines. Conceptually
// you may think of a transform as positioning, scaling, shearing, etc it's
// (local) space relative to it's parent.
//
// Parent space is the local space of any given transform's parent. If the
// transform does not have a parent then parent space is identical to world
// space.
type CoordConv uint8

const (
	// LocalToWorld converts from local space to world space.
	LocalToWorld CoordConv = iota

	// WorldToLocal converts from world space to local space.
	WorldToLocal

	// ParentToWorld converts from parent space to world space.
	ParentToWorld

	// WorldToParent converts from world space to parent space.
	WorldToParent
)

// Transform represents a 3D transformation to a coordinate space. A transform
// effectively defines the position, scale, shear, etc of the local space,
// therefore it is sometimes said that a Transform is a coordinate space.
//
// It can be safely used from multiple goroutines concurrently. It is built
// from various components such as position, scale, and shear values and may
// use euler or quaternion rotation. It supports a hierarchial tree system of
// transforms to create complex transformations.
//
// When in doubt about coordinate spaces it can be helpful to think about the
// fact that each vertex of an object is considered to be in it's local space
// and is then converted to world space for display.
//
// Since world space serves as the common factor among all transforms (i.e.
// any value in any transform's local space can be converted to world space and
// back) converting between world and local/parent space can be extremely
// useful for e.g. relative movement/rotation to another object's transform.
type Transform struct {
	access sync.RWMutex

	// The parent transform, or nil if there is none.
	parent     *Transform
	lastParent *Transform

	// A pointer to the built (i.e. cached) transformation matrix or nil if a
	// rebuild is required.
	built *math.Mat4

	// Pointers to the matrices describing local-to-world and world-to-local
	// space conversions.
	localToWorld, worldToLocal *math.Mat4

	// A pointer to a quaternion rotation, or nil if euler rotation is in use.
	quat *math.Quat

	// The position, rotation, scaling, and shearing components.
	pos, rot, scale, shear math.Vec3
}

// Equals tells if the two transforms are equal.
func (t *Transform) Equals(other *Transform) bool {
	t.access.RLock()
	other.access.RLock()

	// Compare parent pointers.
	if t.parent != other.parent {
		goto fail
	}

	// Two-step quaternion comparison.
	if (t.quat != nil) != (other.quat != nil) {
		goto fail
	}
	if t.quat != nil && !(*t.quat).Equals(*other.quat) {
		goto fail
	}

	// Compare position, rotation, scale, and shear.
	if !t.pos.Equals(other.pos) {
		goto fail
	}
	if !t.rot.Equals(other.rot) {
		goto fail
	}
	if !t.scale.Equals(other.scale) {
		goto fail
	}
	if !t.shear.Equals(other.shear) {
		goto fail
	}

	t.access.RUnlock()
	other.access.RUnlock()
	return true

fail:
	t.access.RUnlock()
	other.access.RUnlock()
	return false
}

// build builds and stores the transformation matrix from the components of
// this transform.
func (t *Transform) build() {
	if t.built != nil && (t.lastParent != nil && t.parent != nil && t.lastParent.Equals(t.parent)) {
		// No update is required.
		return
	}
	if t.parent != nil {
		t.lastParent = t.parent.Copy()
	}

	// Apply rotation
	var hpr math.Vec3
	if t.quat != nil {
		// Use quaternion rotation.
		hpr = (*t.quat).Hpr(math.CoordSysZUpRight)
	} else {
		// Use euler rotation.
		hpr = t.rot.XyzToHpr().Radians()
	}

	// Compose upper 3x3 matrics using scale, shear, and HPR components.
	scaleShearHpr := math.Mat3Compose(t.scale, t.shear, hpr, math.CoordSysZUpRight)

	// Build this space's transformation matrix.
	built := math.Mat4Identity.SetUpperMat3(scaleShearHpr)
	built = built.SetTranslation(t.pos)
	t.built = &built

	// Build the local-to-world transformation matrix.
	ltw := built
	if t.parent != nil {
		ltw = ltw.Mul(t.parent.Convert(LocalToWorld))
	}
	t.localToWorld = &ltw

	// Build the world-to-local transformation matrix.
	wtl, _ := built.Inverse()
	if t.parent != nil {
		parent := t.parent.Convert(WorldToLocal)
		wtl = wtl.Mul(parent)
	}
	t.worldToLocal = &wtl
}

// Implements Transformable interface by simply returning the local-to-world
// matrix.
func (t *Transform) Mat4() math.Mat4 {
	return t.Convert(LocalToWorld)
}

// LocalMat4 returns a matrix describing the space that this transform defines.
// It is the matrix that is built out of the components of this transform, it
// does not include any parent transformation, etc.
func (t *Transform) LocalMat4() math.Mat4 {
	t.access.Lock()
	t.build()
	l := *t.built
	t.access.Unlock()
	return l
}

// SetParent sets a parent transform for this transform to effectively inherit
// from. This allows creating complex hierarchies of transformations.
//
// e.g. setting the parent of a camera's transform to the player's transform
// makes it such that the camera follows the player.
func (t *Transform) SetParent(p *Transform) {
	t.access.Lock()
	if t.parent != p {
		t.built = nil
		t.parent = p
	}
	t.access.Unlock()
}

// Parent returns the parent of this transform, as previously set.
func (t *Transform) Parent() *Transform {
	t.access.RLock()
	p := t.parent
	t.access.RUnlock()
	return p
}

// SetQuat sets the quaternion rotation of this transform.
//
// The last call to either SetQuat or SetRot is what effictively determines
// whether quaternion or euler rotation will be used by this transform.
func (t *Transform) SetQuat(q math.Quat) {
	t.access.Lock()
	if (*t.quat) != q {
		t.built = nil
		t.quat = &q
	}
	t.access.Unlock()
}

// Quat returns the quaternion rotation of this transform. If this transform is
// instead using euler rotation (see IsQuat) then a quaternion is created from
// the euler rotation of this transform and returned.
//
// The last call to either SetQuat or SetRot is what effictively determines
// whether quaternion or euler rotation will be used by this transform.
func (t *Transform) Quat() math.Quat {
	var q math.Quat
	t.access.RLock()
	if t.quat != nil {
		q = *t.quat
	} else {
		// Convert euler rotation to quaternion.
		q = math.QuatFromHpr(t.rot.XyzToHpr().Radians(), math.CoordSysZUpRight)
	}
	t.access.RUnlock()
	return q
}

// IsQuat tells if this transform is currently utilizing quaternion or euler
// rotation.
//
// The last call to either SetQuat or SetRot is what effictively determines
// whether quaternion or euler rotation will be used by this transform.
func (t *Transform) IsQuat() bool {
	t.access.RLock()
	isQuat := t.quat != nil
	t.access.RUnlock()
	return isQuat
}

// SetRot sets the euler rotation of this transform in degrees about their
// respective axis (e.g. if r.X == 45 then it is 45 degrees around the X
// axis).
//
// The last call to either SetQuat or SetRot is what effictively determines
// whether quaternion or euler rotation will be used by this transform.
func (t *Transform) SetRot(r math.Vec3) {
	t.access.Lock()
	if t.rot != r {
		t.built = nil
		t.quat = nil
		t.rot = r
	}
	t.access.Unlock()
}

// Rot returns the euler rotation of this transform. If this transform is
// instead using quaternion (see IsQuat) rotation then it is converted to euler
// rotation and returned.
//
// The last call to either SetQuat or SetRot is what effictively determines
// whether quaternion or euler rotation will be used by this transform.
func (t *Transform) Rot() math.Vec3 {
	var r math.Vec3
	t.access.RLock()
	if t.quat == nil {
		r = t.rot
	} else {
		// Convert quaternion rotation to euler rotation.
		r = (*t.quat).Hpr(math.CoordSysZUpRight).HprToXyz().Degrees()
	}
	t.access.RUnlock()
	return r
}

// SetPos sets the local position of this transform.
func (t *Transform) SetPos(p math.Vec3) {
	t.access.Lock()
	if t.pos != p {
		t.built = nil
		t.pos = p
	}
	t.access.Unlock()
}

// Pos returns the local position of this transform.
func (t *Transform) Pos() math.Vec3 {
	t.access.RLock()
	p := t.pos
	t.access.RUnlock()
	return p
}

// SetScale sets the local scale of this transform (e.g. a scale of
// math.Vec3{2, 1.5, 1} would make an object appear twice as large on the local
// X axis, one and a half times larger on the local Y axis, and would not scale
// on the local Z axis at all).
func (t *Transform) SetScale(s math.Vec3) {
	t.access.Lock()
	if t.scale != s {
		t.built = nil
		t.scale = s
	}
	t.access.Unlock()
}

// Scale returns the local scacle of this transform.
func (t *Transform) Scale() math.Vec3 {
	t.access.RLock()
	s := t.scale
	t.access.RUnlock()
	return s
}

// SetShear sets the local shear of this transform.
func (t *Transform) SetShear(s math.Vec3) {
	t.access.Lock()
	if t.shear != s {
		t.built = nil
		t.shear = s
	}
	t.access.Unlock()
}

// Shear returns the local shear of this transform.
func (t *Transform) Shear() math.Vec3 {
	t.access.RLock()
	s := t.shear
	t.access.RUnlock()
	return s
}

// Reset sets all of the values of this transform to the default ones.
func (t *Transform) Reset() {
	t.access.Lock()
	t.parent = nil
	t.built = nil
	t.localToWorld = nil
	t.worldToLocal = nil
	t.quat = nil
	t.pos = math.Vec3Zero
	t.rot = math.Vec3Zero
	t.scale = math.Vec3One
	t.shear = math.Vec3Zero
	t.access.Unlock()
}

// Copy returns a new transform with all of it's values set equal to t (i.e. a
// copy of this transform).
func (t *Transform) Copy() *Transform {
	t.access.RLock()
	cpy := &Transform{
		parent: t.parent,
		pos:    t.pos,
		rot:    t.rot,
		scale:  t.scale,
		shear:  t.shear,
	}
	if t.built != nil {
		builtCpy := *t.built
		cpy.built = &builtCpy
	}
	if t.localToWorld != nil {
		ltwCpy := *t.localToWorld
		cpy.localToWorld = &ltwCpy
	}
	if t.worldToLocal != nil {
		wtlCpy := *t.worldToLocal
		cpy.worldToLocal = &wtlCpy
	}
	if t.quat != nil {
		quatCpy := *t.quat
		cpy.quat = &quatCpy
	}
	t.access.RUnlock()
	return cpy
}

// NewTransform returns a new *Transform with the default values (a uniform
// scale of one).
func NewTransform() *Transform {
	return &Transform{
		scale: math.Vec3One,
	}
}

// Convert returns a matrix which performs the given coordinate space
// conversion.
func (t *Transform) Convert(c CoordConv) math.Mat4 {
	switch c {
	case LocalToWorld:
		t.access.Lock()
		t.build()
		ltw := *t.localToWorld
		t.access.Unlock()
		return ltw

	case WorldToLocal:
		t.access.Lock()
		t.build()
		wtl := *t.worldToLocal
		t.access.Unlock()
		return wtl

	case ParentToWorld:
		t.access.Lock()
		t.build()
		ltw := *t.localToWorld
		local := *t.built
		t.access.Unlock()

		// Reverse the local transform:
		localInv, _ := local.Inverse()
		return localInv.Mul(ltw)

	case WorldToParent:
		t.access.Lock()
		t.build()
		wtl := *t.worldToLocal
		local := *t.built
		t.access.Unlock()
		return local.Mul(wtl)
	}
	panic("Convert(): invalid conversion")
}

// ConvertPos converts the given point, p, using the given coordinate space
// conversion. For instance to convert a point in local space into world space:
//  t.ConvertPos(p, LocalToWorld)
func (t *Transform) ConvertPos(p math.Vec3, c CoordConv) math.Vec3 {
	return p.TransformMat4(t.Convert(c))
}

// ConvertRot converts the given rotation, r, using the given coordinate space
// conversion. For instance to convert a rotation in local space into world
// space:
//  t.ConvertRot(p, LocalToWorld)
func (t *Transform) ConvertRot(r math.Vec3, c CoordConv) math.Vec3 {
	m := t.Convert(c)
	q := math.QuatFromHpr(r.XyzToHpr().Radians(), math.CoordSysZUpRight)
	m = q.ExtractToMat4().Mul(m)
	q = math.QuatFromMat3(m.UpperMat3())
	return q.Hpr(math.CoordSysZUpRight).HprToXyz().Degrees()
}
