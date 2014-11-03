package simd

import (
	"testing"
)

// Logs generic CPU info, go test with the verbose flag (-v) will show this
// information.
func TestCPUInfo(t *testing.T) {
	t.Log("CPU has SSE2 Support?", CPU.SSE2)
	t.Log("CPU has AVX Support?", CPU.AVX)
}

// Benchmarks the overhead of a single function call. Invokes three functions
// in a chain to avoid the inliner (could also just use the gc flag to disable
// inlining if in doubt).
func emptyc() {}
func emptyb() { emptyc() }
func emptya() { emptyb() }
func Benchmark3EmptyFuncs(bench *testing.B) {
	for n := 0; n < bench.N; n++ {
		emptya()
	}
}
