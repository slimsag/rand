// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gfx

// BlendState represents the blend state to use when rendering an object whose
// AlphaMode == BlendedAlpha.
type BlendState struct {
	// The constant blend color to be used (e.g. with BConstantColor).
	Color Color

	// Specifies the blend operand to use for the source RGB components.
	// All predefined BlendOp constants may be used.
	SrcRGB BlendOp

	// Specifies the blend operand to use for the destination RGB components.
	// All predefined BlendOp constants may be used except BSrcAlphaSaturate.
	DstRGB BlendOp

	// Specifies the blending operation to use between source and destination
	// alpha components.
	SrcAlpha, DstAlpha BlendOp

	// Specifies the blending equation to use for RGB and Alpha components,
	// respectively.
	RGBEq, AlphaEq BlendEq
}

// Compare compares this state against the other one using DefaultBlendState as
// a reference when inequality occurs and returns whether or not this state
// should sort before the other one for purposes of state sorting.
func (b BlendState) Compare(other BlendState) bool {
	if b == other {
		return true
	}
	if b.Color != other.Color {
		return b.Color == DefaultBlendState.Color
	}
	if b.SrcRGB != other.SrcRGB {
		return b.SrcRGB == DefaultBlendState.SrcRGB
	}
	if b.DstRGB != other.DstRGB {
		return b.DstRGB == DefaultBlendState.DstRGB
	}
	if b.SrcRGB != other.SrcRGB {
		return b.SrcRGB == DefaultBlendState.SrcRGB
	}
	if b.SrcAlpha != other.SrcAlpha {
		return b.SrcAlpha == DefaultBlendState.SrcAlpha
	}
	if b.DstAlpha != other.DstAlpha {
		return b.DstAlpha == DefaultBlendState.DstAlpha
	}
	if b.RGBEq != other.RGBEq {
		return b.RGBEq == DefaultBlendState.RGBEq
	}
	if b.AlphaEq != other.AlphaEq {
		return b.AlphaEq == DefaultBlendState.AlphaEq
	}
	return true
}

// The default blend state to use for graphics objects (by default it works
// well for premultiplied alpha blending).
var DefaultBlendState = BlendState{
	Color:    Color{0, 0, 0, 0},
	SrcRGB:   BOne,
	SrcAlpha: BOne,
	DstRGB:   BOneMinusSrcAlpha,
	DstAlpha: BOneMinusSrcAlpha,
	RGBEq:    BAdd,
	AlphaEq:  BAdd,
}

// BlendOp represents a single blend operand, e.g. BOne, BOneMinusSrcAlpha.
type BlendOp uint8

const (
	BZero BlendOp = iota
	BOne
	BSrcColor
	BOneMinusSrcColor
	BDstColor
	BOneMinusDstColor
	BSrcAlpha
	BOneMinusSrcAlpha
	BDstAlpha
	BOneMinusDstAlpha
	BConstantColor
	BOneMinusConstantColor
	BConstantAlpha
	BOneMinusConstantAlpha

	// Not applicable for use in BlendState.SrcRGB.
	BSrcAlphaSaturate
)

// BlendEq represents a single blend equation to use when blending RGB or Alpha
// components in the color buffer, it must be one of BlendAdd, BlendSubtract, BlendReverseSubtract.
type BlendEq uint8

const (
	// BAdd represents a blending equation where the src and dst colors are
	// added to eachother to produce the result.
	BAdd BlendEq = iota

	// BSub represents a blending equation where the src and dst colors are
	// subtracted from eachother to produce the result.
	BSub

	// BReverseSub represents a blending equation where the src and dst colors
	// are reverse-subtracted from eachother to produce the result.
	BReverseSub
)
