// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

import "sync"

// NativeShader represents the native object of a shader. Typically only
// renderers will create these.
type NativeShader Destroyable

// Shader represents a single shader program.
//
// Clients are responsible for utilizing the RWMutex of the shader when using
// it or invoking methods.
type Shader struct {
	sync.RWMutex

	// The native object of this shader. Once loaded (if no compiler error
	// occured) then the renderer using this shader must assign a value to this
	// field. Typically clients should not assign values to this field at all.
	NativeShader

	// Weather or not this shader is currently loaded or not.
	Loaded bool

	// If true then when this shader is loaded the data sources of it will be
	// kept instead of being set to nil (which allows them to be garbage
	// collected).
	KeepDataOnLoad bool

	// The name of the shader, optional (used in the shader compilation error
	// log).
	Name string

	// The GLSL vertex shader source.
	GLSLVert []byte

	// The GLSL fragment shader.
	GLSLFrag []byte

	// A map of names and values to use as inputs for the shader program while
	// rendering. Values must be of the following data types or else they will
	// be ignored:
	//  bool
	//  float32
	//  []float32
	//  gfx.Vec3
	//  []gfx.Vec3
	//  gfx.Mat4
	//  []gfx.Mat4
	Inputs map[string]interface{}

	// The error log from compiling the shader program, if any. Only set once
	// the shader is loaded.
	Error []byte
}

// CanDraw reports if this shader is valid for drawing. Cases where a shader is
// not valid for drawing are as follows:
//  len(s.GLSLVert) == 0
//  len(s.GLSLFrag) == 0
//  len(s.Error) > 0
func (s *Shader) CanDraw() bool {
	if len(s.GLSLVert) == 0 {
		return false
	}
	if len(s.GLSLFrag) == 0 {
		return false
	}
	if len(s.Error) > 0 {
		return false
	}
	return true
}

// Copy returns a new copy of this Shader. Explicitly not copied over is the
// native shader, the OnLoad slice, the Loaded status, and error log slice.
func (s *Shader) Copy() *Shader {
	cpy := &Shader{
		sync.RWMutex{},
		nil,   // Native shader -- not copied.
		false, // Loaded status -- not copied.
		s.KeepDataOnLoad,
		s.Name,
		make([]byte, len(s.GLSLVert)),
		make([]byte, len(s.GLSLFrag)),
		make(map[string]interface{}, len(s.Inputs)),
		nil, // Error slice -- not copied.
	}
	copy(cpy.GLSLVert, s.GLSLVert)
	copy(cpy.GLSLFrag, s.GLSLFrag)
	for name := range s.Inputs {
		cpy.Inputs[name] = s.Inputs[name]
	}
	return cpy
}

// ClearData sets the data slices (s.GLSLVert, s.Error, etc) of this shader to
// nil if s.KeepDataOnLoad is set to false.
func (s *Shader) ClearData() {
	if !s.KeepDataOnLoad {
		s.GLSLVert = nil
		s.GLSLFrag = nil
		s.Error = nil
	}
}

// NewShader returns a new, initialized *Shader object with the given name.
func NewShader(name string) *Shader {
	return &Shader{
		Name:   name,
		Inputs: make(map[string]interface{}),
	}
}
