// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import (
	"azul3d.org/v1/math"
	"testing"
)

func TestTransform(t *testing.T) {
	tf := NewTransform()
	tf.SetScale(math.Vec3{2, 4, 6})
	pos := math.Vec3{1, 2, 3}
	tf.SetPos(pos)
	m := tf.Mat4()
	if m[0][0] != 2 || m[1][1] != 4 || m[2][2] != 6 {
		t.Log("scale invalid")
		t.Log(m)
		t.Fail()
	}
	if !m.Translation().Equals(pos) {
		t.Log("pos invalid")
		t.Log(m)
		t.Fail()
	}
}

func TestTransformRel(t *testing.T) {
	a := NewTransform()
	a.SetPos(math.Vec3{10, 0, 2})

	b := NewTransform()
	b.SetPos(math.Vec3{10, 5, 0})
	b.SetParent(a)

	c := NewTransform()
	c.SetPos(math.Vec3{10, 0, 5})
	c.SetParent(b)

	ltw := c.Convert(LocalToWorld)
	want := math.Vec3{30, 5, 7}
	if !ltw.Translation().Equals(want) {
		t.Log("local-to-world invalid")
		t.Log("want (world)", want)
		t.Log("got (world)", ltw.Translation())
		t.Log(ltw)
		t.Fail()
	}

	wtl := c.Convert(WorldToLocal)
	want = math.Vec3{-30, -5, -7}
	if !wtl.Translation().Equals(want) {
		t.Log("world-to-local invalid")
		t.Log("want (world)", want)
		t.Log("got (world)", wtl.Translation())
		t.Log(wtl)
		t.Fail()
	}

	wtp := c.Convert(WorldToParent)
	want = math.Vec3{-20, -5, -2}
	if !wtp.Translation().Equals(want) {
		t.Log("world-to-parent invalid")
		t.Log("want (world)", want)
		t.Log("got (world)", wtp.Translation())
		t.Log(wtp)
		t.Fail()
	}
}

func TestTransformPointToWorld(t *testing.T) {
	a := NewTransform()
	a.SetPos(math.Vec3{0, 0, -50})

	b := NewTransform()
	b.SetPos(math.Vec3{-25, -35, -50})
	b.SetParent(a)

	p := b.ConvertPos(math.Vec3{50, 0, 0}, LocalToWorld)
	want := math.Vec3{25, -35, -100}
	if !p.Equals(want) {
		t.Log("got (world)", p)
		t.Log("want (world)", want)
		t.Fail()
	}
}

func TestTransformPointToLocal(t *testing.T) {
	a := NewTransform()
	a.SetPos(math.Vec3{0, 0, -50})

	b := NewTransform()
	b.SetPos(math.Vec3{0, 0, -50})
	b.SetParent(a)

	p := b.ConvertPos(math.Vec3{50, 0, 0}, LocalToWorld)
	p = b.ConvertPos(p, WorldToLocal)
	want := math.Vec3{50, 0, 0}
	if !p.Equals(want) {
		t.Log("got (local)", p)
		t.Log("want (local)", want)
		t.Fail()
	}
}

func TestTransformRotToWorld(t *testing.T) {
	a := NewTransform()
	a.SetRot(math.Vec3{0, 0, 45})

	b := NewTransform()
	b.SetRot(math.Vec3{45, 0, 45})
	b.SetParent(a)

	p := b.ConvertRot(math.Vec3{-45, 0, 0}, LocalToWorld)
	want := math.Vec3{0, 0, 90}
	if !p.Equals(want) {
		t.Log("got (world)", p)
		t.Log("want (world)", want)
		t.Fail()
	}
}
