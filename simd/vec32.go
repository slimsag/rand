package simd

import (
	"math"
)

// Vec32 is a four-component 64-bit floating point vector.
type Vec32 [4]float32

// Vec32Add returns the result of element-wise a + b.
var Vec32Add func(a, b Vec32) Vec32

// Implemented in Vec32.s
func sse2Vec32Add(a, b Vec32) Vec32

// Implemented in Vec32.s
func avxVec32Add(a, b Vec32) Vec32

func goVec32Add(a, b Vec32) Vec32 {
	return Vec32{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
		a[3] + b[3],
	}
}

// Vec32Sub returns the result of element-wise a - b.
var Vec32Sub func(a, b Vec32) Vec32

// Implemented in Vec32.s
func sse2Vec32Sub(a, b Vec32) Vec32

// Implemented in Vec32.s
func avxVec32Sub(a, b Vec32) Vec32

func goVec32Sub(a, b Vec32) Vec32 {
	return Vec32{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
		a[3] - b[3],
	}
}

// Vec32Mul returns the result of element-wise a * b.
var Vec32Mul func(a, b Vec32) Vec32

// Implemented in Vec32.s
func sse2Vec32Mul(a, b Vec32) Vec32

// Implemented in Vec32.s
func avxVec32Mul(a, b Vec32) Vec32

func goVec32Mul(a, b Vec32) Vec32 {
	return Vec32{
		a[0] * b[0],
		a[1] * b[1],
		a[2] * b[2],
		a[3] * b[3],
	}
}

// Vec32Div returns the result of element-wise a / b.
var Vec32Div func(a, b Vec32) Vec32

// Implemented in Vec32.s
func sse2Vec32Div(a, b Vec32) Vec32

// Implemented in Vec32.s
func avxVec32Div(a, b Vec32) Vec32

func goVec32Div(a, b Vec32) Vec32 {
	return Vec32{
		a[0] / b[0],
		a[1] / b[1],
		a[2] / b[2],
		a[3] / b[3],
	}
}

// Vec32Eq returns the result of element-wise a == b.
var Vec32Eq func(a, b Vec32) bool

// Implemented in Vec32.s
func sse2Vec32Eq(a, b Vec32) bool

// Implemented in Vec32.s
func avxVec32Eq(a, b Vec32) bool

func goVec32Eq(a, b Vec32) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}

// Vec32Min returns the result of element-wise math.Min(a, b).
var Vec32Min func(a, b Vec32) Vec32

// Implemented in Vec32.s
func sse2Vec32Min(a, b Vec32) Vec32

// Implemented in Vec32.s
func avxVec32Min(a, b Vec32) Vec32

func goVec32Min(a, b Vec32) Vec32 {
	return Vec32{
		float32(math.Min(float64(a[0]), float64(b[0]))),
		float32(math.Min(float64(a[1]), float64(b[1]))),
		float32(math.Min(float64(a[2]), float64(b[2]))),
		float32(math.Min(float64(a[3]), float64(b[3]))),
	}
}

// Vec32Max returns the result of element-wise math.Max(a, b).
var Vec32Max func(a, b Vec32) Vec32

// Implemented in Vec32.s
func sse2Vec32Max(a, b Vec32) Vec32

// Implemented in Vec32.s
func avxVec32Max(a, b Vec32) Vec32

func goVec32Max(a, b Vec32) Vec32 {
	return Vec32{
		float32(math.Max(float64(a[0]), float64(b[0]))),
		float32(math.Max(float64(a[1]), float64(b[1]))),
		float32(math.Max(float64(a[2]), float64(b[2]))),
		float32(math.Max(float64(a[3]), float64(b[3]))),
	}
}
