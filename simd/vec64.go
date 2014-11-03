package simd

import (
	"math"
)

// Vec64 is a four-component 64-bit floating point vector.
type Vec64 [4]float64

// Vec64Add returns the result of element-wise a + b.
var Vec64Add func(a, b Vec64) Vec64

// Implemented in vec64.s
func sse2Vec64Add(a, b Vec64) Vec64

// Implemented in vec64.s
func avxVec64Add(a, b Vec64) Vec64

func goVec64Add(a, b Vec64) Vec64 {
	return Vec64{
		a[0] + b[0],
		a[1] + b[1],
		a[2] + b[2],
		a[3] + b[3],
	}
}

// Vec64Sub returns the result of element-wise a - b.
var Vec64Sub func(a, b Vec64) Vec64

// Implemented in vec64.s
func sse2Vec64Sub(a, b Vec64) Vec64

// Implemented in vec64.s
func avxVec64Sub(a, b Vec64) Vec64

func goVec64Sub(a, b Vec64) Vec64 {
	return Vec64{
		a[0] - b[0],
		a[1] - b[1],
		a[2] - b[2],
		a[3] - b[3],
	}
}

// Vec64Mul returns the result of element-wise a * b.
var Vec64Mul func(a, b Vec64) Vec64

// Implemented in vec64.s
func sse2Vec64Mul(a, b Vec64) Vec64

// Implemented in vec64.s
func avxVec64Mul(a, b Vec64) Vec64

func goVec64Mul(a, b Vec64) Vec64 {
	return Vec64{
		a[0] * b[0],
		a[1] * b[1],
		a[2] * b[2],
		a[3] * b[3],
	}
}

// Vec64Div returns the result of element-wise a / b.
var Vec64Div func(a, b Vec64) Vec64

// Implemented in vec64.s
func sse2Vec64Div(a, b Vec64) Vec64

// Implemented in vec64.s
func avxVec64Div(a, b Vec64) Vec64

func goVec64Div(a, b Vec64) Vec64 {
	return Vec64{
		a[0] / b[0],
		a[1] / b[1],
		a[2] / b[2],
		a[3] / b[3],
	}
}

// Vec64Eq returns the result of element-wise a == b.
var Vec64Eq func(a, b Vec64) bool

// Implemented in vec64.s
func sse2Vec64Eq(a, b Vec64) bool

// Implemented in vec64.s
func avxVec64Eq(a, b Vec64) bool

func goVec64Eq(a, b Vec64) bool {
	return a[0] == b[0] && a[1] == b[1] && a[2] == b[2] && a[3] == b[3]
}

// Vec64Min returns the result of element-wise math.Min(a, b).
var Vec64Min func(a, b Vec64) Vec64

// Implemented in vec64.s
func sse2Vec64Min(a, b Vec64) Vec64

// Implemented in vec64.s
func avxVec64Min(a, b Vec64) Vec64

func goVec64Min(a, b Vec64) Vec64 {
	return Vec64{
		math.Min(a[0], b[0]),
		math.Min(a[1], b[1]),
		math.Min(a[2], b[2]),
		math.Min(a[3], b[3]),
	}
}

// Vec64Max returns the result of element-wise math.Max(a, b).
var Vec64Max func(a, b Vec64) Vec64

// Implemented in vec64.s
func sse2Vec64Max(a, b Vec64) Vec64

// Implemented in vec64.s
func avxVec64Max(a, b Vec64) Vec64

func goVec64Max(a, b Vec64) Vec64 {
	return Vec64{
		math.Max(a[0], b[0]),
		math.Max(a[1], b[1]),
		math.Max(a[2], b[2]),
		math.Max(a[3], b[3]),
	}
}
