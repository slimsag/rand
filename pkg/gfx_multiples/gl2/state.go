// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gl2

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/native/gl"
	"image"
)

func convertStencilOp(o gfx.StencilOp) int32 {
	switch o {
	case gfx.SKeep:
		return gl.KEEP
	case gfx.SZero:
		return gl.ZERO
	case gfx.SReplace:
		return gl.REPLACE
	case gfx.SIncr:
		return gl.INCR
	case gfx.SIncrWrap:
		return gl.INCR_WRAP
	case gfx.SDecr:
		return gl.DECR
	case gfx.SDecrWrap:
		return gl.DECR_WRAP
	case gfx.SInvert:
		return gl.INVERT
	}
	panic("failed to convert")
}

func convertCmp(c gfx.Cmp) int32 {
	switch c {
	case gfx.Always:
		return gl.ALWAYS
	case gfx.Never:
		return gl.NEVER
	case gfx.Less:
		return gl.LESS
	case gfx.LessOrEqual:
		return gl.LEQUAL
	case gfx.Greater:
		return gl.GREATER
	case gfx.GreaterOrEqual:
		return gl.GEQUAL
	case gfx.Equal:
		return gl.EQUAL
	case gfx.NotEqual:
		return gl.NOTEQUAL
	}
	panic("failed to convert")
}

func convertBlendOp(o gfx.BlendOp) int32 {
	switch o {
	case gfx.BZero:
		return gl.ZERO
	case gfx.BOne:
		return gl.ONE
	case gfx.BSrcColor:
		return gl.SRC_COLOR
	case gfx.BOneMinusSrcColor:
		return gl.ONE_MINUS_SRC_COLOR
	case gfx.BDstColor:
		return gl.DST_COLOR
	case gfx.BOneMinusDstColor:
		return gl.ONE_MINUS_DST_COLOR
	case gfx.BSrcAlpha:
		return gl.SRC_ALPHA
	case gfx.BOneMinusSrcAlpha:
		return gl.ONE_MINUS_SRC_ALPHA
	case gfx.BDstAlpha:
		return gl.DST_ALPHA
	case gfx.BOneMinusDstAlpha:
		return gl.ONE_MINUS_DST_ALPHA
	case gfx.BConstantColor:
		return gl.CONSTANT_COLOR
	case gfx.BOneMinusConstantColor:
		return gl.ONE_MINUS_CONSTANT_COLOR
	case gfx.BConstantAlpha:
		return gl.CONSTANT_ALPHA
	case gfx.BOneMinusConstantAlpha:
		return gl.ONE_MINUS_CONSTANT_ALPHA
	case gfx.BSrcAlphaSaturate:
		return gl.SRC_ALPHA_SATURATE
	}
	panic("failed to convert")
}

func convertBlendEq(eq gfx.BlendEq) int32 {
	switch eq {
	case gfx.BAdd:
		return gl.FUNC_ADD
	case gfx.BSub:
		return gl.FUNC_SUBTRACT
	case gfx.BReverseSub:
		return gl.FUNC_REVERSE_SUBTRACT
	}
	panic("failed to convert")
}

func convertRect(rect, bounds image.Rectangle) (x, y int32, width, height uint32) {
	// We must flip the Y axis because image.Rectangle uses top-left as
	// the origin but OpenGL uses bottom-left as the origin.
	y = int32(bounds.Dy() - (rect.Min.Y + rect.Dy())) // bottom
	height = uint32(rect.Dy())                        // top

	x = int32(rect.Min.X)
	width = uint32(rect.Dx())
	return
}

var glDefaultStencil = gfx.StencilState{
	WriteMask: 0xFFFF,
	Fail:      gfx.SKeep,
	DepthFail: gfx.SKeep,
	DepthPass: gfx.SKeep,
	Cmp:       gfx.Always,
}

var glDefaultBlend = gfx.BlendState{
	Color:    gfx.Color{0, 0, 0, 0},
	SrcRGB:   gfx.BOne,
	DstRGB:   gfx.BZero,
	SrcAlpha: gfx.BOne,
	DstAlpha: gfx.BZero,
	RGBEq:    gfx.BAdd,
	AlphaEq:  gfx.BAdd,
}

func (r *Renderer) clearLastState() {
	// Ensure that these values match up with the default OpenGL state values.
	r.stateColorWrite(true, true, true, true)
	r.stateDithering(true)
	r.stateStencilTest(false)
	r.stateStencilOp(glDefaultStencil, glDefaultStencil)
	r.stateStencilFunc(glDefaultStencil, glDefaultStencil)
	r.stateStencilMask(0xFFFF, 0xFFFF)
	r.stateDepthFunc(gfx.Less)
	r.stateDepthTest(false)
	r.stateDepthWrite(true)
	r.stateFaceCulling(gfx.NoFaceCulling)
	r.stateProgram(0)
	r.stateBlend(false)
	r.stateBlendColor(glDefaultBlend.Color)
	r.stateBlendFuncSeparate(glDefaultBlend)
	r.stateBlendEquationSeparate(glDefaultBlend)
	r.stateAlphaToCoverage(false)
	r.stateClearColor(gfx.Color{0.0, 0.0, 0.0, 0.0})
	r.stateClearDepth(1.0)
	r.stateClearStencil(0)
}

func (r *Renderer) stateScissor(rect image.Rectangle) {
	// Only if the (final) scissor rectangle has changed do we need to make the
	// OpenGL call.
	bounds := r.Bounds()

	// If the rectangle is empty use the entire area.
	if rect.Empty() {
		rect = bounds
	} else {
		// Intersect the rectangle with the renderer's bounds.
		rect = bounds.Intersect(rect)
	}

	if r.last.scissor != rect {
		// Store the new scissor rectangle.
		r.last.scissor = rect
		x, y, width, height := convertRect(rect, bounds)
		r.render.Scissor(x, y, width, height)
	}
}

func (r *Renderer) stateColorWrite(cr, g, b, a bool) {
	cw := [4]bool{cr, g, b, a}
	if r.last.colorWrite != cw {
		r.last.colorWrite = cw
		r.render.ColorMask(
			gl.GLBool(cr),
			gl.GLBool(g),
			gl.GLBool(b),
			gl.GLBool(a),
		)
	}
}

func (r *Renderer) stateDithering(enabled bool) {
	if r.last.dithering != enabled {
		r.last.dithering = enabled
		if enabled {
			r.render.Enable(gl.DITHER)
		} else {
			r.render.Disable(gl.DITHER)
		}
	}
}

func (r *Renderer) stateStencilTest(stencilTest bool) {
	if r.last.stencilTest != stencilTest {
		r.last.stencilTest = stencilTest
		if stencilTest {
			r.render.Enable(gl.STENCIL_TEST)
		} else {
			r.render.Disable(gl.STENCIL_TEST)
		}
	}
}

func (r *Renderer) stateStencilOp(front, back gfx.StencilState) {
	if r.last.stencilOpFront != front || r.last.stencilOpBack != back {
		r.last.stencilOpFront = front
		r.last.stencilOpBack = back
		if front == back {
			// We can save a few calls.
			r.render.StencilOpSeparate(
				gl.FRONT_AND_BACK,
				convertStencilOp(front.Fail),
				convertStencilOp(front.DepthFail),
				convertStencilOp(front.DepthPass),
			)
		} else {
			r.render.StencilOpSeparate(
				gl.FRONT,
				convertStencilOp(front.Fail),
				convertStencilOp(front.DepthFail),
				convertStencilOp(front.DepthPass),
			)
			r.render.StencilOpSeparate(
				gl.BACK,
				convertStencilOp(back.Fail),
				convertStencilOp(back.DepthFail),
				convertStencilOp(back.DepthPass),
			)
		}
	}
}

func (r *Renderer) stateStencilFunc(front, back gfx.StencilState) {
	if r.last.stencilFuncFront != front || r.last.stencilFuncBack != back {
		r.last.stencilFuncFront = front
		r.last.stencilFuncBack = back
		if front == back {
			// We can save a few calls.
			r.render.StencilFuncSeparate(
				gl.FRONT_AND_BACK,
				convertCmp(front.Cmp),
				int32(front.Reference),
				uint32(front.ReadMask),
			)
		} else {
			r.render.StencilFuncSeparate(
				gl.FRONT,
				convertCmp(front.Cmp),
				int32(front.Reference),
				uint32(front.ReadMask),
			)
			r.render.StencilFuncSeparate(
				gl.BACK,
				convertCmp(back.Cmp),
				int32(back.Reference),
				uint32(back.ReadMask),
			)
		}
	}
}

func (r *Renderer) stateStencilMask(front, back uint) {
	if r.last.stencilMaskFront != front || r.last.stencilMaskBack != back {
		r.last.stencilMaskFront = front
		r.last.stencilMaskBack = back
		if front == back {
			// We can save a call.
			r.render.StencilMaskSeparate(gl.FRONT_AND_BACK, uint32(front))
		} else {
			r.render.StencilMaskSeparate(gl.FRONT, uint32(front))
			r.render.StencilMaskSeparate(gl.BACK, uint32(back))
		}
	}
}

func (r *Renderer) stateDepthFunc(df gfx.Cmp) {
	if r.last.depthFunc != df {
		r.last.depthFunc = df
		r.render.DepthFunc(convertCmp(df))
	}
}

func (r *Renderer) stateDepthTest(enabled bool) {
	if r.last.depthTest != enabled {
		r.last.depthTest = enabled
		if enabled {
			r.render.Enable(gl.DEPTH_TEST)
		} else {
			r.render.Disable(gl.DEPTH_TEST)
		}
	}
}

func (r *Renderer) stateDepthWrite(enabled bool) {
	if r.last.depthWrite != enabled {
		r.last.depthWrite = enabled
		if enabled {
			r.render.DepthMask(gl.GLBool(true))
		} else {
			r.render.DepthMask(gl.GLBool(false))
		}
	}
}

func (r *Renderer) stateFaceCulling(m gfx.FaceCullMode) {
	if r.last.faceCulling != m {
		r.last.faceCulling = m
		switch m {
		case gfx.BackFaceCulling:
			r.render.Enable(gl.CULL_FACE)
			r.render.CullFace(gl.BACK)
		case gfx.FrontFaceCulling:
			r.render.Enable(gl.CULL_FACE)
			r.render.CullFace(gl.FRONT)
		default:
			r.render.Disable(gl.CULL_FACE)
		}
	}
}

func (r *Renderer) stateProgram(p uint32) {
	if r.last.program != p {
		r.last.program = p
		r.render.UseProgram(p)
	}
}

func (r *Renderer) stateBlend(blend bool) {
	if r.last.blend != blend {
		r.last.blend = blend
		if blend {
			r.render.Enable(gl.BLEND)
		} else {
			r.render.Disable(gl.BLEND)
		}
	}
}

func (r *Renderer) stateBlendColor(c gfx.Color) {
	if r.last.blendColor != c {
		r.last.blendColor = c
		r.render.BlendColor(c.R, c.G, c.B, c.A)
	}
}

func (r *Renderer) stateBlendFuncSeparate(s gfx.BlendState) {
	if r.last.blendFuncSeparate != s {
		r.last.blendFuncSeparate = s
		r.render.BlendFuncSeparate(
			convertBlendOp(s.SrcRGB),
			convertBlendOp(s.DstRGB),
			convertBlendOp(s.SrcAlpha),
			convertBlendOp(s.SrcAlpha),
		)
	}
}

func (r *Renderer) stateBlendEquationSeparate(s gfx.BlendState) {
	if r.last.blendEquationSeparate != s {
		r.last.blendEquationSeparate = s
		r.render.BlendEquationSeparate(
			convertBlendEq(s.RGBEq),
			convertBlendEq(s.AlphaEq),
		)
	}
}

func (r *Renderer) stateAlphaToCoverage(alphaToCoverage bool) {
	if r.last.alphaToCoverage != alphaToCoverage {
		r.last.alphaToCoverage = alphaToCoverage
		if r.gpuInfo.AlphaToCoverage {
			if alphaToCoverage {
				r.render.Enable(gl.SAMPLE_ALPHA_TO_COVERAGE)
			} else {
				r.render.Disable(gl.SAMPLE_ALPHA_TO_COVERAGE)
			}
		}
	}
}

func (r *Renderer) stateClearColor(color gfx.Color) {
	if r.last.clearColor != color {
		r.last.clearColor = color
		r.render.ClearColor(color.R, color.G, color.B, color.A)
	}
}

func (r *Renderer) stateClearDepth(depth float64) {
	if r.last.clearDepth != depth {
		r.last.clearDepth = depth
		r.render.ClearDepth(depth)
	}
}

func (r *Renderer) stateClearStencil(stencil int) {
	if r.last.clearStencil != stencil {
		r.last.clearStencil = stencil
		r.render.ClearStencil(int32(stencil))
	}
}
