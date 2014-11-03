// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ice

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
)

func alphaMode(s string) gfx.AlphaMode {
	switch s {
	case "AlphaBlend":
		return gfx.AlphaBlend
	case "BinaryAlpha":
		return gfx.BinaryAlpha
	case "AlphaToCoverage":
		return gfx.AlphaToCoverage
	}
	return gfx.AlphaMode(0)
}

func cmp(s string) gfx.Cmp {
	switch s {
	case "Always":
		return gfx.Always
	case "Never":
		return gfx.Never
	case "Less":
		return gfx.Less
	case "LessOrEqual":
		return gfx.LessOrEqual
	case "Greater":
		return gfx.Greater
	case "GreaterOrEqual":
		return gfx.GreaterOrEqual
	case "Equal":
		return gfx.Equal
	case "NotEqual":
		return gfx.NotEqual
	}
	return gfx.Cmp(0)
}

func faceCulling(s string) gfx.FaceCullMode {
	switch s {
	case "BackFaceCulling":
		return gfx.BackFaceCulling
	case "FrontFaceCulling":
		return gfx.FrontFaceCulling
	}
	return gfx.FaceCullMode(0)
}

func blendOp(s string) gfx.BlendOp {
	switch s {
	case "Zero":
		return gfx.BZero
	case "One":
		return gfx.BOne
	case "SrcColor":
		return gfx.BSrcColor
	case "OneMinusSrcColor":
		return gfx.BOneMinusSrcColor
	case "DstColor":
		return gfx.BDstColor
	case "OneMinusDstColor":
		return gfx.BOneMinusDstColor
	case "SrcAlpha":
		return gfx.BSrcAlpha
	case "OneMinusSrcAlpha":
		return gfx.BOneMinusSrcAlpha
	case "DstAlpha":
		return gfx.BDstAlpha
	case "OneMinusDstAlpha":
		return gfx.BOneMinusDstAlpha
	case "ConstantColor":
		return gfx.BConstantColor
	case "OneMinusConstantColor":
		return gfx.BOneMinusConstantColor
	case "ConstantAlpha":
		return gfx.BConstantAlpha
	case "OneMinusConstantAlpha":
		return gfx.BOneMinusConstantAlpha
	case "SrcAlphaSaturate":
		return gfx.BSrcAlphaSaturate
	}
	return gfx.BlendOp(0)
}

func blendEq(s string) gfx.BlendEq {
	switch s {
	case "Add":
		return gfx.BAdd
	case "Sub":
		return gfx.BSub
	case "ReverseSub":
		return gfx.BReverseSub
	}
	return gfx.BlendEq(0)
}

func stencilOp(s string) gfx.StencilOp {
	switch s {
	case "Keep":
		return gfx.SKeep
	case "Zero":
		return gfx.SZero
	case "Replace":
		return gfx.SReplace
	case "Incr":
		return gfx.SIncr
	case "IncrWrap":
		return gfx.SIncrWrap
	case "Decr":
		return gfx.SDecr
	case "DecrWrap":
		return gfx.SDecrWrap
	case "Invert":
		return gfx.SInvert
	}
	return gfx.StencilOp(0)
}

func texFormat(s string) gfx.TexFormat {
	switch s {
	case "RGBA":
		return gfx.RGBA
	case "RGB":
		return gfx.RGB
	case "DXT1":
		return gfx.DXT1
	case "DXT1RGBA":
		return gfx.DXT1RGBA
	case "DXT3":
		return gfx.DXT3
	case "DXT5":
		return gfx.DXT5
	}
	return gfx.TexFormat(0)
}

func texWrap(s string) gfx.TexWrap {
	switch s {
	case "Repeat":
		return gfx.Repeat
	case "Clamp":
		return gfx.Clamp
	case "BorderColor":
		return gfx.BorderColor
	case "Mirror":
		return gfx.Mirror
	}
	return gfx.TexWrap(0)
}

func texFilter(s string) gfx.TexFilter {
	switch s {
	case "Nearest":
		return gfx.Nearest
	case "Linear":
		return gfx.Linear
	case "NearestMipmapNearest":
		return gfx.NearestMipmapNearest
	case "LinearMipmapNearest":
		return gfx.LinearMipmapNearest
	case "NearestMipmapLinear":
		return gfx.NearestMipmapLinear
	case "LinearMipmapLinear":
		return gfx.LinearMipmapLinear
	}
	return gfx.TexFilter(0)
}

func color(v [4]float32) gfx.Color {
	return gfx.Color{v[0], v[1], v[2], v[3]}
}

func texCoord(v [2]float32) gfx.TexCoord {
	return gfx.TexCoord{v[0], v[1]}
}

func gfxVec3(v [3]float32) gfx.Vec3 {
	return gfx.Vec3{v[0], v[1], v[2]}
}

func vec3(v [3]float64) math.Vec3 {
	return math.Vec3{v[0], v[1], v[2]}
}

func quat(v [4]float64) math.Quat {
	return math.Quat{v[0], v[1], v[2], v[3]}
}

func mat4(v [16]float64) math.Mat4 {
	return math.Mat4{
		[4]float64{v[0], v[1], v[2], v[3]},
		[4]float64{v[4], v[5], v[6], v[7]},
		[4]float64{v[8], v[9], v[10], v[11]},
		[4]float64{v[12], v[13], v[14], v[15]},
	}
}
