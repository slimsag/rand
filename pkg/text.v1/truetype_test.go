// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import "testing"

type expTestCase struct {
	r rune
	contours [][]Point
}

var expTests = []expTestCase{
	{
		r: 'i',
		contours: [][]Point{
			{{377, 1120}, {377, 560}, {377, 0}, {285, 0}, {193, 0}, {193, 560}, {193, 1120}, {285, 1120}},
			{{377, 1556}, {377, 1439}, {377, 1323}, {285, 1323}, {193, 1323}, {193, 1439}, {193, 1556}, {285, 1556}},
		},
	},{
		r: 'j',
		contours: [][]Point{
			{{377, 1120}, {377, 550}, {377, -20}, {377, -234}, {295, -330}, {214, -426}, {33, -426}, {-2, -426}, {-37, -426}, {-37, -348}, {-37, -270}, {-12, -270}, {12, -270}, {117, -270}, {155, -221}, {193, -173}, {193, -20}, {193, 550}, {193, 1120}, {285, 1120}},
			{{377, 1556}, {377, 1439}, {377, 1323}, {285, 1323}, {193, 1323}, {193, 1439}, {193, 1556}, {285, 1556}},
		},
	},
}

func ptsEqual(a, b []Point) bool {
	if len(a) != len(b) {
		return false
	}
	for i, p := range a {
		if b[i] != p {
			return false
		}
	}
	return true
}

// Test for proper expansion of truetype font points.
func TestTruetypeExp(t *testing.T) {
	// Open font file.
	f, err := LoadFontFile("testdata/Vera.ttf")
	if err != nil {
		t.Fatal(err)
	}

	for _, cs := range expTests {
		// Find the font index for the rune.
		idx, ok := f.Index(cs.r)
		if !ok {
			t.Fatalf("%q: can't find in font\n", cs.r)
		}

		// Lookup the glyph data.
		data, err := f.Lookup(idx)
		if err != nil {
			t.Fatal(err)
		}
		glyph := data.(QuadGlyph)
		if len(cs.contours) != glyph.NumContours() {
			t.Logf("%q want %d contours, got %d contours\n", cs.r, len(cs.contours), glyph.NumContours())
		}
		for i, expect := range cs.contours {
			contour := glyph.Contour(i)
			if !ptsEqual(contour, expect) {
				t.Logf("%q want len=%d, got len=%d\n", cs.r, len(contour), len(expect))
				t.Log("got", contour)
				t.Log("expect", expect)
			}
		}
	}
}
