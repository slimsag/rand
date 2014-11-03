// Package simd implements SIMD accelerated vector math.
//
// This package uses assembly code to implement hardware accelerated vector
// math usable from Go.
//
// As a general rule of thumb:
//  All AMD64 processors support at least SSE2.
//  Intel processors since 2005 support AVX instructions.
//
// Fallbacks are implemented in Go for architectures not supporting such
// extensions, so this code should work on any processor regardless of it's
// features. This package chooses the most (generally) fast method, so it may
// choose Go's generated assembler, SSE2, or AVX depending on feature support.
//
// If desired you can explicitly disable the use of the assembler versions by
// compiling with specific build tags:
//  'nosimd' - Disable all assembler versions (i.e. use only pure-Go versions).
//  'nosse2' - Operate as if SSE2 was not available.
//  'noavx' - Operate as if AVX was not available.
package simd
