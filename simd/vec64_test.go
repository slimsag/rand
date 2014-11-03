package simd

import (
	"testing"
)

func TestVec64Add(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Add(a, b)
		want := goVec64Add(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Add(a, b)
		want := goVec64Add(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Add(a, b)
	want := goVec64Add(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64AddGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Add(a, b)
	}
}
func BenchmarkVec64AddSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Add(a, b)
	}
}
func BenchmarkVec64AddAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Add(a, b)
	}
}
func BenchmarkVec64AddFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Add(a, b)
	}
}

func TestVec64Sub(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Sub(a, b)
		want := goVec64Sub(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Sub(a, b)
		want := goVec64Sub(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Sub(a, b)
	want := goVec64Sub(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64SubGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Sub(a, b)
	}
}
func BenchmarkVec64SubSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Sub(a, b)
	}
}
func BenchmarkVec64SubAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Sub(a, b)
	}
}
func BenchmarkVec64SubFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Sub(a, b)
	}
}

func TestVec64Mul(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Mul(a, b)
		want := goVec64Mul(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Mul(a, b)
		want := goVec64Mul(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Mul(a, b)
	want := goVec64Mul(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64MulGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Mul(a, b)
	}
}
func BenchmarkVec64MulSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Mul(a, b)
	}
}
func BenchmarkVec64MulAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Mul(a, b)
	}
}
func BenchmarkVec64MulFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Mul(a, b)
	}
}

func TestVec64Div(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Div(a, b)
		want := goVec64Div(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Div(a, b)
		want := goVec64Div(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Div(a, b)
	want := goVec64Div(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64DivGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Div(a, b)
	}
}
func BenchmarkVec64DivSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Div(a, b)
	}
}
func BenchmarkVec64DivAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Div(a, b)
	}
}
func BenchmarkVec64DivFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Div(a, b)
	}
}

func TestVec64Eq(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Eq(a, b)
		want := goVec64Eq(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Eq(a, b)
		want := goVec64Eq(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Eq(a, b)
	want := goVec64Eq(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64EqGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Eq(a, b)
	}
}
func BenchmarkVec64EqSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Eq(a, b)
	}
}
func BenchmarkVec64EqAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Eq(a, b)
	}
}
func BenchmarkVec64EqFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Eq(a, b)
	}
}

func TestVec64Min(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Min(a, b)
		want := goVec64Min(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Min(a, b)
		want := goVec64Min(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Min(a, b)
	want := goVec64Min(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64MinGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Min(a, b)
	}
}
func BenchmarkVec64MinSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Min(a, b)
	}
}
func BenchmarkVec64MinAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Min(a, b)
	}
}
func BenchmarkVec64MinFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Min(a, b)
	}
}

func TestVec64Max(t *testing.T) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec64Max(a, b)
		want := goVec64Max(a, b)
		if got != want {
			t.Log("AVX:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	if CPU.SSE2 {
		// Test SSE2-specific version.
		got := sse2Vec64Max(a, b)
		want := goVec64Max(a, b)
		if got != want {
			t.Log("SSE2:")
			t.Log("a", a)
			t.Log("b", b)
			t.Log("got", got)
			t.Log("want", want)
			t.Fail()
		}
	}
	// Test *chosen* version.
	got := Vec64Max(a, b)
	want := goVec64Max(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec64MaxGo(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec64Max(a, b)
	}
}
func BenchmarkVec64MaxSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec64Max(a, b)
	}
}
func BenchmarkVec64MaxAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec64Max(a, b)
	}
}
func BenchmarkVec64MaxFast(bench *testing.B) {
	a := Vec64{1, 2, 3, 4}
	b := Vec64{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec64Max(a, b)
	}
}
