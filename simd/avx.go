// +build amd64,!noavx,!nosimd

package simd

// Implemented in simd.s
func avxSupport() bool

// See also: noavx.go
var haveAVX = avxSupport()
