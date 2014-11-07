// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"unicode/utf8"

	"azul3d.org/gfx.v1"
)

// Object represents a single 2D text object.
type Object struct {
	*gfx.Object
	font Font
	str  string
}

// SetFont sets the font that this text object is using.
func (o *Object) SetFont(f Font) {
	o.Lock()
	o.font = f
	o.Unlock()
}

// Font returns the font that this text object is using.
func (o *Object) Font() Font {
	o.RLock()
	font := o.font
	o.RUnlock()
	return font
}

// Set sets the string that this text object represents.
func (o *Object) Set(s string) error {
	o.Lock()
	o.str = s
	err := o.build()
	o.Unlock()
	return err
}

// String returns the string that this text object represents.
func (o *Object) String() string {
	o.RLock()
	s := o.str
	o.RUnlock()
	return s
}

func (o *Object) build() error {
	mesh := o.Meshes[0]
	mesh.VerticesChanged = true

	// Grab a glyph mesher to use for the generation of the new string.
	m := FindGlyphMesher(o.font)
	r, size := utf8.DecodeRuneInString(o.str)
	_ = size
	return m.Append(mesh, o.font.Index(r))
}

// New creates and returns a new text object given a string.
func New(font Font, s string) (*Object, error) {
	o := &Object{
		gfx.NewObject(),
		font,
		s,
	}
	o.Object.AlphaMode = gfx.AlphaToCoverage
	o.Object.FaceCulling = gfx.NoFaceCulling // FIXME: remove
	o.Object.Shader = nil                    // FIXME: add default
	o.Object.Meshes = []*gfx.Mesh{gfx.NewMesh()}
	err := o.build()
	return o, err
}
