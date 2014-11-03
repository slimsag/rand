package simd

import (
	"testing"
)

func TestVec32Add(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Add(a, b)
		want := goVec32Add(a, b)
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
		got := sse2Vec32Add(a, b)
		want := goVec32Add(a, b)
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
	got := Vec32Add(a, b)
	want := goVec32Add(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32AddGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Add(a, b)
	}
}
func BenchmarkVec32AddSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Add(a, b)
	}
}
func BenchmarkVec32AddAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Add(a, b)
	}
}
func BenchmarkVec32AddFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Add(a, b)
	}
}

/*
func TestVec32Sub(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Sub(a, b)
		want := goVec32Sub(a, b)
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
		got := sse2Vec32Sub(a, b)
		want := goVec32Sub(a, b)
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
	got := Vec32Sub(a, b)
	want := goVec32Sub(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32SubGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Sub(a, b)
	}
}
func BenchmarkVec32SubSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Sub(a, b)
	}
}
func BenchmarkVec32SubAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Sub(a, b)
	}
}
func BenchmarkVec32SubFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Sub(a, b)
	}
}

func TestVec32Mul(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Mul(a, b)
		want := goVec32Mul(a, b)
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
		got := sse2Vec32Mul(a, b)
		want := goVec32Mul(a, b)
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
	got := Vec32Mul(a, b)
	want := goVec32Mul(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32MulGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Mul(a, b)
	}
}
func BenchmarkVec32MulSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Mul(a, b)
	}
}
func BenchmarkVec32MulAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Mul(a, b)
	}
}
func BenchmarkVec32MulFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Mul(a, b)
	}
}

func TestVec32Div(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Div(a, b)
		want := goVec32Div(a, b)
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
		got := sse2Vec32Div(a, b)
		want := goVec32Div(a, b)
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
	got := Vec32Div(a, b)
	want := goVec32Div(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32DivGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Div(a, b)
	}
}
func BenchmarkVec32DivSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Div(a, b)
	}
}
func BenchmarkVec32DivAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Div(a, b)
	}
}
func BenchmarkVec32DivFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Div(a, b)
	}
}


func TestVec32Eq(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Eq(a, b)
		want := goVec32Eq(a, b)
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
		got := sse2Vec32Eq(a, b)
		want := goVec32Eq(a, b)
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
	got := Vec32Eq(a, b)
	want := goVec32Eq(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32EqGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Eq(a, b)
	}
}
func BenchmarkVec32EqSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Eq(a, b)
	}
}
func BenchmarkVec32EqAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Eq(a, b)
	}
}
func BenchmarkVec32EqFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Eq(a, b)
	}
}


func TestVec32Min(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Min(a, b)
		want := goVec32Min(a, b)
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
		got := sse2Vec32Min(a, b)
		want := goVec32Min(a, b)
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
	got := Vec32Min(a, b)
	want := goVec32Min(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32MinGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Min(a, b)
	}
}
func BenchmarkVec32MinSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Min(a, b)
	}
}
func BenchmarkVec32MinAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Min(a, b)
	}
}
func BenchmarkVec32MinFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Min(a, b)
	}
}

func TestVec32Max(t *testing.T) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}

	if CPU.AVX {
		// Test AVX-specific version.
		got := avxVec32Max(a, b)
		want := goVec32Max(a, b)
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
		got := sse2Vec32Max(a, b)
		want := goVec32Max(a, b)
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
	got := Vec32Max(a, b)
	want := goVec32Max(a, b)
	if got != want {
		t.Log("Chosen:")
		t.Log("a", a)
		t.Log("b", b)
		t.Log("got", got)
		t.Log("want", want)
		t.Fail()
	}
}
func BenchmarkVec32MaxGo(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		goVec32Max(a, b)
	}
}
func BenchmarkVec32MaxSSE2(bench *testing.B) {
	if !CPU.SSE2 {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		sse2Vec32Max(a, b)
	}
}
func BenchmarkVec32MaxAVX(bench *testing.B) {
	if !CPU.AVX {
		bench.SkipNow()
		return
	}
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		avxVec32Max(a, b)
	}
}
func BenchmarkVec32MaxFast(bench *testing.B) {
	a := Vec32{1, 2, 3, 4}
	b := Vec32{2, 3, 4, 5}
	for n := 0; n < bench.N; n++ {
		Vec32Max(a, b)
	}
}
*/
