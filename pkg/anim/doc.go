// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package anim provides a generic system for vertex animation.
//
// This package implements two of the most common forms of vertex animation,
// skeletal and morph-target animation. Although breifly covered below, they
// are covered in significantly more depth on e.g. Wikipedia:
//
// http://en.wikipedia.org/wiki/Skeletal_animation
//
// http://en.wikipedia.org/wiki/Morph_target_animation
//
// Skeletal Animation
//
// Skeletal animation is by far the most common means of animating vertex data.
// At the basic level, a given object has a list of bones which are used to
// manipulate the vertex data of the object's meshes. You can have any number
// of bones which affect an object.
//
// Multiple bones can affect the same exact vertex, when this occurs the result
// is that the two positions are blended together to produce the final result.
//
// Morph-Target Animation
//
// Morph target animation (sometimes called Per-Vertex Animation, Shape
// Interpolation, or Blend Shapes) is another popular means of animating vertex
// data. The basic idea is to take multiple shapes and blend between them over
// some period of time (e.g. making a star shape slowly morph into a monkey
// shape). Multiple shapes can be blended together to produce the final result,
// which makes it very useful for things like facial animation.
//
// One of the disadvantages of morph target animation is that it consumes much
// more memory than skeletal animation does. Each morph target is essentially
// a copy of the vertex positions of a mesh, just 'morphed' into place.
package anim

type Manager struct {
}

func (m *Manager) Tick() {
}

// Limitations
//
// To gain the best performance, animation is performed almost solely in a
// a vertex shader which does all of the heavy-lifting. In fact, this is how
// all modern games perform vertex animation. Although it provides the best
// performance (compared to, e.g. CPU-based animation) it does come with a few
// limitations.
//
// 1. A single vertex may only be affected by four bones at the same time. You
// may have any number of affecting a mesh, but only four can affect the exact
// same vertex at any given time. All others are just ignored.
//
// 2. A single mesh may only be unfluenced by N (configurable) morph targets
// at the same time. This is because each morph target consumes a single
// varying input to the graphics shader.

// 2. A single shape may only be affected by N Morph targets at the same time.
// Since a morph target is essentially a complete copy of the mesh's vertex
// positions, just 'morphed', it means that much more data must be streamed to
// the graphics hardware. In OpenGL for instance, this means taking up a single
// varying input for every morph target.
//
// The N morph targets limit is configurable, but by default is set to 16.
//
// Since it is common to have multiple morph targets active at the same time,
// but there is the aforementioned limit, instead of simply limiting the number
// of total morph targets an object can use to say 16, we use a different
// tactic.
//
// Consider that the N morph targets is not a limit but instead
//Every few frames the CPU determines the importance level of each
// morph target,
//
// Every few frames the CPU determines the importance level of each morphtarget.
// For every morph target affecting the mesh a level

// Since it is often desired to have multiple morph targets active at the same
// time, a

//  Because it is an entire copy
// positions,
//
// The flaw with shape animation is how much memory it consumes. For every
// shape animation, a entire set of vertex data must be stored. Most of the
// time this is not a huge problem, but if your models use an abnormally large
// number of vertices or many shape animations: you will use more graphics
// memory.
