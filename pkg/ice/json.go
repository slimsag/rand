// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ice

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"fmt"
	"image"
)

type jsonScene struct {
	Props      map[string]interface{}
	Cameras    map[string]*jsonCamera
	Meshes     map[string]*jsonMesh
	Textures   map[string]*jsonTexture
	Objects    map[string]*jsonObject
	Transforms map[string]*jsonTransform
	Shaders    map[string]*jsonShader
}

func (j *jsonScene) scene(resolver Resolver) *Scene {
	s := &Scene{
		Props:   j.Props,
		Cameras: make(map[string]*gfx.Camera, len(j.Cameras)),
		Objects: make(map[string]*gfx.Object, len(j.Objects)),
	}

	// Declare a function responsible for finding a *gfx.Transform from it's
	// JSON counterpart.
	gfxTransforms := make(map[string]*gfx.Transform, len(j.Transforms))
	var findTransform func(name string) *gfx.Transform
	findTransform = func(name string) *gfx.Transform {
		// It's possible this can hit infinite recursion with bad parent
		// listings.. not sure how that should be handled.
		t, ok := gfxTransforms[name]
		if !ok {
			t = j.Transforms[name].transform(findTransform)
			gfxTransforms[name] = t
		}
		return t
	}

	// Declare a function responsible for finding a *gfx.Shader from it's JSON
	// counterpart.
	gfxShaders := make(map[string]*gfx.Shader, len(j.Shaders))
	var findShader func(name string) *gfx.Shader
	findShader = func(name string) *gfx.Shader {
		s, ok := gfxShaders[name]
		if !ok {
			s = j.Shaders[name].shader(name)
			gfxShaders[name] = s
		}
		return s
	}

	// Declare a function responsible for finding a *gfx.Texture from it's JSON
	// counterpart.
	gfxTextures := make(map[string]*gfx.Texture, len(j.Textures))
	var findTexture func(name string) *gfx.Texture
	findTexture = func(name string) *gfx.Texture {
		tex, ok := gfxTextures[name]
		if !ok {
			jTex, ok := j.Textures[name]
			if !ok {
				// FIXME: bubble up
				panic(fmt.Errorf("texture not listed: %q", name))
			}
			tex = jTex.texture()
			gfxTextures[name] = tex

			// Load the texture resource.
			rc, err := resolver.Resolve(*jTex.Source)
			if err != nil {
				// FIXME: bubble up
				panic(err)
			}
			defer rc.Close()

			// Decode the image.
			img, _, err := image.Decode(rc)
			if err != nil {
				// FIXME: bubble up
				panic(err)
			}
			tex.Source = img
		}
		return tex
	}

	for name, cam := range j.Cameras {
		s.Cameras[name] = cam.camera(j.Meshes, findTexture, findTransform, findShader)
	}
	for name, obj := range j.Objects {
		s.Objects[name] = obj.object(j.Meshes, findTexture, findTransform, findShader)
	}
	return s
}

type jsonCamera struct {
	Object *jsonObject
}

func (j *jsonCamera) camera(meshes map[string]*jsonMesh, findTexture func(name string) *gfx.Texture, findTransform func(name string) *gfx.Transform, findShader func(name string) *gfx.Shader) *gfx.Camera {
	c := gfx.NewCamera()
	c.Object = j.Object.object(meshes, findTexture, findTransform, findShader)
	return c
}

type jsonObject struct {
	Props         map[string]interface{}
	OcclusionTest *bool
	State         jsonState
	Shader        *string
	Transform     *string
	Meshes        []string
	Textures      []string
}

func (j *jsonObject) object(meshes map[string]*jsonMesh, findTexture func(name string) *gfx.Texture, findTransform func(name string) *gfx.Transform, findShader func(name string) *gfx.Shader) *gfx.Object {
	o := gfx.NewObject()
	if j.OcclusionTest != nil {
		o.OcclusionTest = *j.OcclusionTest
	}
	//o.Props = j.Props
	o.Meshes = make([]*gfx.Mesh, len(j.Meshes))
	for i, m := range j.Meshes {
		o.Meshes[i] = meshes[m].mesh()
	}
	o.Textures = make([]*gfx.Texture, len(j.Textures))
	for i, t := range j.Textures {
		o.Textures[i] = findTexture(t)
	}
	if j.Shader != nil {
		o.Shader = findShader(*j.Shader)
	}
	if j.Transform != nil {
		o.Transform = findTransform(*j.Transform)
	}
	return o
}

type jsonTransform struct {
	Pos, Rot, Scale, Shear *[3]float64
	Quat                   *[4]float64
	Mat4                   *[16]float64
	Parent                 *string
}

func (j *jsonTransform) transform(findTransform func(name string) *gfx.Transform) *gfx.Transform {
	t := gfx.NewTransform()
	if j.Pos != nil {
		t.SetPos(vec3(*j.Pos))
	}
	if j.Rot != nil {
		t.SetRot(vec3(*j.Rot))
	}
	if j.Scale != nil {
		t.SetScale(vec3(*j.Scale))
	}
	if j.Shear != nil {
		t.SetShear(vec3(*j.Shear))
	}
	if j.Quat != nil {
		t.SetQuat(quat(*j.Quat))
	}
	// FIXME:
	//if j.Mat4 != nil {
	//	t.SetMat4(mat4(*j.Mat4))
	//}
	if j.Parent != nil {
		t.SetParent(findTransform(*j.Parent))
	}
	return t
}

type jsonShader struct {
	KeepDataOnLoad     *bool
	GLSLSources        []string
	GLSLVert, GLSLFrag []string
	Inputs             map[string]interface{}
}

func (j *jsonShader) shader(name string) *gfx.Shader {
	s := gfx.NewShader(name)
	if j.KeepDataOnLoad != nil {
		s.KeepDataOnLoad = *j.KeepDataOnLoad
	}
	// FIXME: j.GLSLSources filepaths
	for _, line := range j.GLSLVert {
		s.GLSLVert = append(s.GLSLVert, []byte(line)...)
	}
	for _, line := range j.GLSLFrag {
		s.GLSLFrag = append(s.GLSLFrag, []byte(line)...)
	}
	s.Inputs = j.Inputs
	return s
}

type jsonState struct {
	AlphaMode                                   *string
	Blend                                       *jsonBlendState
	WriteRed, WriteGreen, WriteBlue, WriteAlpha *bool
	Dithering                                   *bool
	DepthTest, DepthWrite                       *bool
	DepthCmp                                    *string
	StencilTest                                 *bool
	FaceCulling                                 *string
	StencilFront, StencilBack                   *jsonStencilState
}

func (j *jsonState) state() gfx.State {
	s := gfx.DefaultState
	if j.AlphaMode != nil {
		s.AlphaMode = alphaMode(*j.AlphaMode)
	}
	if j.Blend != nil {
		s.Blend = j.Blend.blendState()
	}
	if j.WriteRed != nil {
		s.WriteRed = *j.WriteRed
	}
	if j.WriteGreen != nil {
		s.WriteGreen = *j.WriteGreen
	}
	if j.WriteBlue != nil {
		s.WriteBlue = *j.WriteBlue
	}
	if j.WriteAlpha != nil {
		s.WriteAlpha = *j.WriteAlpha
	}
	if j.Dithering != nil {
		s.Dithering = *j.Dithering
	}
	if j.DepthTest != nil {
		s.DepthTest = *j.DepthTest
	}
	if j.DepthWrite != nil {
		s.DepthWrite = *j.DepthWrite
	}
	if j.DepthCmp != nil {
		s.DepthCmp = cmp(*j.DepthCmp)
	}
	if j.StencilTest != nil {
		s.StencilTest = *j.StencilTest
	}
	if j.FaceCulling != nil {
		s.FaceCulling = faceCulling(*j.FaceCulling)
	}
	if j.StencilFront != nil {
		s.StencilFront = j.StencilFront.stencilState()
	}
	if j.StencilBack != nil {
		s.StencilBack = j.StencilBack.stencilState()
	}
	return s
}

type jsonBlendState struct {
	Color              *[4]float32
	SrcRGB             *string
	DstRGB             *string
	SrcAlpha, DstAlpha *string
	RGBEq, AlphaEq     *string
}

func (j *jsonBlendState) blendState() gfx.BlendState {
	s := gfx.DefaultBlendState
	if j.Color != nil {
		s.Color = color(*j.Color)
	}
	if j.SrcRGB != nil {
		s.SrcRGB = blendOp(*j.SrcRGB)
	}
	if j.DstRGB != nil {
		s.DstRGB = blendOp(*j.DstRGB)
	}
	if j.SrcAlpha != nil {
		s.SrcAlpha = blendOp(*j.SrcAlpha)
	}
	if j.DstAlpha != nil {
		s.DstAlpha = blendOp(*j.DstAlpha)
	}
	if j.RGBEq != nil {
		s.RGBEq = blendEq(*j.RGBEq)
	}
	if j.AlphaEq != nil {
		s.AlphaEq = blendEq(*j.AlphaEq)
	}
	return s
}

type jsonStencilState struct {
	WriteMask *uint
	ReadMask  *uint
	Reference *uint
	Fail      *string
	DepthFail *string
	DepthPass *string
	Cmp       *string
}

func (j *jsonStencilState) stencilState() gfx.StencilState {
	s := gfx.DefaultStencilState
	if j.WriteMask != nil {
		s.WriteMask = *j.WriteMask
	}
	if j.ReadMask != nil {
		s.ReadMask = *j.ReadMask
	}
	if j.Reference != nil {
		s.Reference = *j.Reference
	}
	if j.Fail != nil {
		s.Fail = stencilOp(*j.Fail)
	}
	if j.DepthFail != nil {
		s.DepthFail = stencilOp(*j.DepthFail)
	}
	if j.DepthPass != nil {
		s.DepthPass = stencilOp(*j.DepthPass)
	}
	if j.Cmp != nil {
		s.Cmp = cmp(*j.Cmp)
	}
	return s
}

type jsonAABB struct {
	Min [3]float64
	Max [3]float64
}

func (j jsonAABB) rect3() math.Rect3 {
	return math.Rect3{
		Min: math.Vec3{j.Min[0], j.Min[1], j.Min[2]},
		Max: math.Vec3{j.Max[0], j.Max[1], j.Max[2]},
	}
}

type jsonMesh struct {
	KeepDataOnLoad bool
	Dynamic        bool
	AABB           jsonAABB
	Indices        []float64
	Vertices       [][3]float32
	Colors         [][4]float32
	TexCoords      [][][2]float32
}

func (j *jsonMesh) mesh() *gfx.Mesh {
	m := &gfx.Mesh{
		KeepDataOnLoad: j.KeepDataOnLoad,
		Dynamic:        j.Dynamic,
		Indices:        make([]uint32, len(j.Indices)),
		Vertices:       make([]gfx.Vec3, len(j.Vertices)),
		Colors:         make([]gfx.Color, len(j.Colors)),
		TexCoords:      make([]gfx.TexCoordSet, len(j.TexCoords)),
	}
	m.AABB = j.AABB.rect3()
	for i, v := range j.Indices {
		m.Indices[i] = uint32(v)
	}
	for i, v := range j.Vertices {
		m.Vertices[i] = gfxVec3(v)
	}
	for i, v := range j.Colors {
		m.Colors[i] = color(v)
	}
	for tcs, _ := range j.TexCoords {
		m.TexCoords[tcs].Slice = make([]gfx.TexCoord, len(j.TexCoords[tcs]))
		for i, v := range j.TexCoords[tcs] {
			m.TexCoords[tcs].Slice[i] = texCoord(v)
		}
	}
	return m
}

type jsonTexture struct {
	KeepDataOnLoad       *bool
	Source               *string
	Format               *string
	WrapU, WrapV         *string
	BorderColor          *[4]float32
	MinFilter, MagFilter *string
}

func (j *jsonTexture) texture() *gfx.Texture {
	t := new(gfx.Texture)
	if j.KeepDataOnLoad != nil {
		t.KeepDataOnLoad = *j.KeepDataOnLoad
	}
	// FIXME: load from file
	//if j.Source != nil {
	//	t.Source = *j.Source
	//}
	if j.Format != nil {
		t.Format = texFormat(*j.Format)
	}
	if j.WrapU != nil {
		t.WrapU = texWrap(*j.WrapU)
	}
	if j.WrapV != nil {
		t.WrapV = texWrap(*j.WrapV)
	}
	if j.BorderColor != nil {
		t.BorderColor = color(*j.BorderColor)
	}
	if j.MinFilter != nil {
		t.MinFilter = texFilter(*j.MinFilter)
	}
	if j.MagFilter != nil {
		t.MagFilter = texFilter(*j.MagFilter)
	}
	return t
}
