package simd

// CPUInfo represents generic information about the CPU.
type CPUInfo struct {
	// Whether or not the CPU supports AVX.
	AVX bool

	// Whether or not the CPU supports SSE2.
	SSE2 bool
}

// CPU is filled at initialization time with information about the CPU.
var CPU CPUInfo

func init() {
	// Fill CPU with generic information.
	CPU = CPUInfo{
		AVX:  haveAVX,  // see avx.go, noavx.go
		SSE2: haveSSE2, // See sse2.go, nosse2.go
	}

	// Depending on CPU features, install the fastest functions. We don't do
	// this using build tags because some processors support, e.g. SSE2 but not
	// AVX (only way to tell is at runtime).
	if CPU.AVX {
		// Vec64
		Vec64Add = avxVec64Add
		Vec64Sub = avxVec64Sub
		Vec64Mul = avxVec64Mul
		Vec64Div = goVec64Div // Go division is faster than SSE2/AVX ?
		Vec64Min = avxVec64Min
		Vec64Max = avxVec64Max
		Vec64Eq = goVec64Eq // Go equality is faster than SSE2/AVX ?

		// Vec32
		Vec32Add = avxVec32Add
		Vec32Sub = avxVec32Sub
		Vec32Mul = avxVec32Mul
		Vec32Div = goVec32Div // Go division is faster than SSE2/AVX ?
		Vec32Min = avxVec32Min
		Vec32Max = avxVec32Max
		Vec32Eq = goVec32Eq // Go equality is faster than SSE2/AVX ?
	} else if CPU.SSE2 {
		// Vec64
		Vec64Add = sse2Vec64Add
		Vec64Sub = sse2Vec64Sub
		Vec64Mul = sse2Vec64Mul
		Vec64Div = goVec64Div // Go division is faster than SSE2/AVX ?
		Vec64Min = sse2Vec64Min
		Vec64Max = sse2Vec64Max
		Vec64Eq = goVec64Eq // Go equality is faster than SSE2/AVX ?

		// Vec32
		Vec32Add = sse2Vec32Add
		Vec32Sub = sse2Vec32Sub
		Vec32Mul = sse2Vec32Mul
		Vec32Div = goVec32Div // Go division is faster than SSE2/AVX ?
		Vec32Min = sse2Vec32Min
		Vec32Max = sse2Vec32Max
		Vec32Eq = goVec32Eq // Go equality is faster than SSE2/AVX ?
	} else {
		// Vec64
		Vec64Add = goVec64Add
		Vec64Sub = goVec64Sub
		Vec64Mul = goVec64Mul
		Vec64Div = goVec64Div
		Vec64Min = goVec64Min
		Vec64Max = goVec64Max
		Vec64Eq = goVec64Eq

		// Vec32
		Vec32Add = goVec32Add
		Vec32Sub = goVec32Sub
		Vec32Mul = goVec32Mul
		Vec32Div = goVec32Div
		Vec32Min = goVec32Min
		Vec32Max = goVec32Max
		Vec32Eq = goVec32Eq
	}
}
