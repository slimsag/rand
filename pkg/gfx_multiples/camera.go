// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"image"
)

var (
	// Get an matrix which will translate our matrix from ZUpRight to YUpRight
	zUpRightToYUpRight = math.CoordSysZUpRight.ConvertMat4(math.CoordSysYUpRight)
)

// Camera represents a camera object, it may be moved in 3D space using the
// objects transform and the viewing frustum controls how the camera views
// things. Since a camera is in itself also an object it may also have visible
// meshes attatched to it, etc.
type Camera struct {
	*Object

	// The projection matrix of the camera, which is responsible for projecting
	// world coordinates into device coordinates.
	Projection Mat4
}

// SetOrtho sets this camera's Projection matrix to an orthographic one.
//
// The view parameter is the viewing rectangle for the orthographic
// projection in window coordinates.
//
// The near and far parameters describe the minimum closest and maximum
// furthest clipping points of the view frustum.
//
// Clients who need advanced control over how the orthographic viewing frustum
// is set up may use this method's source as a reference (e.g. to change the
// center point, which this method sets at the bottom-left).
//
// Write access is required for this method to operate safely.
func (c *Camera) SetOrtho(view image.Rectangle, near, far float64) {
	w := float64(view.Dx())
	w = float64(int((w / 2.0)) * 2)
	h := float64(view.Dy())
	h = float64(int((h / 2.0)) * 2)
	m := math.Mat4Ortho(0, w, 0, h, near, far)
	c.Projection = ConvertMat4(m)
}

// SetPersp sets this camera's Projection matrix to an perspective one.
//
// The view parameter is the viewing rectangle for the orthographic
// projection in window coordinates.
//
// The fov parameter is the Y axis field of view (e.g. some games use 75) to
// use.
//
// The near and far parameters describe the minimum closest and maximum
// furthest clipping points of the view frustum.
//
// Clients who need advanced control over how the perspective viewing frustum
// is set up may use this method's source as a reference (e.g. to change the
// center point, which this method sets at the center).
//
// Write access is required for this method to operate safely.
func (c *Camera) SetPersp(view image.Rectangle, fov, near, far float64) {
	aspectRatio := float64(view.Dx()) / float64(view.Dy())
	m := math.Mat4Perspective(fov, aspectRatio, near, far)
	c.Projection = ConvertMat4(m)
}

// Project returns a 2D point in normalized device space coordinates given a 3D
// point in the world.
//
// If ok=false is returned then the point is outside of the camera's view and
// the returned point may not be meaningful.
func (c *Camera) Project(p3 math.Vec3) (p2 math.Vec2, ok bool) {
	cameraInv, _ := c.Object.Transform.Mat4().Inverse()
	cameraInv = cameraInv.Mul(zUpRightToYUpRight)

	projection := c.Projection.Mat4()
	vp := cameraInv.Mul(projection)

	p4 := math.Vec4{p3.X, p3.Y, p3.Z, 1.0}
	p4 = p4.Transform(vp)
	if p4.W == 0 {
		p2 = math.Vec2Zero
		ok = false
		return
	}

	recipW := 1.0 / p4.W
	p2 = math.Vec2{p4.X * recipW, p4.Y * recipW}

	xValid := (p2.X >= -1) && (p2.X <= 1)
	yValid := (p2.Y >= -1) && (p2.Y <= 1)
	ok = (p4.W > 0) && xValid && yValid
	return
}

// NewCamera returns a new *Camera with the default values.
func NewCamera() *Camera {
	return &Camera{
		NewObject(),
		ConvertMat4(math.Mat4Identity),
	}
}
