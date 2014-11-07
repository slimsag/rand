// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"io/ioutil"

	"code.google.com/p/freetype-go/freetype/truetype"
)

// TruetypeFont is an truetype.Font that implements the Font interface.
type TruetypeFont struct {
	// The literal truetype font source.
	*truetype.Font

	// The hinting policy to use for the font, NoHinting by default.
	truetype.Hinting

	buf *truetype.GlyphBuf
}

// Implements the Font interface.
func (f *TruetypeFont) Index(r rune) FontIndex {
	return FontIndex(f.Font.Index(r))
}

// Implements the Font interface.
func (f *TruetypeFont) Bounds() GlyphBounds {
	b := f.Font.Bounds(f.Font.FUnitsPerEm())
	return GlyphBounds{
		Min: Point{
			X: int(b.XMin),
			Y: int(b.YMin),
		},
		Max: Point{
			X: int(b.XMax),
			Y: int(b.YMax),
		},
	}
}

// Implements the Font interface.
func (f *TruetypeFont) Kerning(a, b FontIndex) (x, y int) {
	x = int(f.Font.Kerning(f.Font.FUnitsPerEm(), truetype.Index(a), truetype.Index(b)))
	y = -1
	return
}

// Implements the Font interface.
func (f *TruetypeFont) Measure(i FontIndex) (*GlyphMetrics, error) {
	// Load into glyph buffer.
	if f.buf == nil {
		f.buf = truetype.NewGlyphBuf()
	}
	index := truetype.Index(i)
	scale := f.Font.FUnitsPerEm()
	err := f.buf.Load(f.Font, scale, index, f.Hinting)
	if err != nil {
		return nil, err
	}

	hMetric := f.Font.HMetric(scale, index)
	vMetric := f.Font.VMetric(scale, index)
	b := f.buf.B
	return &GlyphMetrics{
		Bounds: GlyphBounds{
			Min: Point{X: int(b.XMin), Y: int(b.YMin)},
			Max: Point{X: int(b.XMax), Y: int(b.YMax)},
		},
		BearingX: int(hMetric.LeftSideBearing),
		AdvanceX: int(hMetric.AdvanceWidth),
		BearingY: int(vMetric.TopSideBearing),
		AdvanceY: int(vMetric.AdvanceHeight),
	}, nil
}

// Implements the Font interface.
func (f *TruetypeFont) Lookup(i FontIndex) (glyphData interface{}, err error) {
	// Load into glyph buffer.
	if f.buf == nil {
		f.buf = truetype.NewGlyphBuf()
	}
	index := truetype.Index(i)
	scale := f.Font.FUnitsPerEm()
	err = f.buf.Load(f.Font, scale, index, f.Hinting)
	if err != nil {
		return nil, err
	}

	// Expand each contour.
	var (
		glyph QuadGlyph
		e0    int
	)
	for _, e1 := range f.buf.End {
		f.expandContour(&glyph, f.buf.Point[e0:e1])
		e0 = e1
	}
	return glyph, nil
}

// Expands the contour defined by points, appending the results to glyph.Points
// and glyph.Indices.
func (f *TruetypeFont) expandContour(glyph *QuadGlyph, points []truetype.Point) {
	if len(points) == 0 {
		return
	}

	// Converts a truetype.Point directly to a Point.
	convertPoint := func(p truetype.Point) Point {
		return Point{
			X: int(p.X),
			Y: int(p.Y),
		}
	}

	//startIndex := len(glyph.Points)
	start := convertPoint(points[0])
	glyph.Points = append(glyph.Points, start)
	q0, on0 := start, true

	// Finds the midpoint of a and b.
	mid := func(a, b Point) Point {
		return Point{
			X: (a.X + b.X) / 2,
			Y: (a.Y + b.Y) / 2,
		}
	}

	// Express a linear curve (a line) using q0 and q1 points.
	expLinear := func(q1 Point) {
		glyph.Points = append(glyph.Points, mid(q0, q1))
		glyph.Points = append(glyph.Points, q1)
	}

	// Express a quadratic curve using start (q0), end(q2), and control (q1)
	// points.
	expQuad := func(q1, q2 Point) {
		glyph.Points = append(glyph.Points, q1)
		glyph.Points = append(glyph.Points, q2)
	}

	for _, p := range points[1:] {
		q := convertPoint(p)
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				expLinear(q)
			} else {
				expQuad(q0, q)
			}
		} else {
			if !on0 {
				expQuad(q0, mid(q0, q))
			}
		}
		q0, on0 = q, on
	}

	// Close the curve.
	if on0 {
		expLinear(start)
	} else {
		expQuad(q0, start)
	}
	//log.Println("faux?", len(glyph.Points[startIndex:]) % 2 != 0)
	//if len(glyph.Points[startIndex:]) % 2 != 0 {
	// Insert faux point to produce a set of N % 2 == 0 points.
	//	glyph.Points = append(glyph.Points, start)
	//panic("expansion of contour failed to produce n % 2 == 0 points.")
	//}

	//if len(glyph.Points[startIndex:]) % 2 != 0 {
	//	panic("expansion of contour failed to produce n % 2 == 0 points.")
	//}
	glyph.Indices = append(glyph.Indices, len(glyph.Points))
}

// LoadFont loads the given truetype font file data and returns a new
// *TruetypeFont. For easier loading, consider using LoadTTFFile instead. This
// method is simply short-hand for:
//  tt, err := truetype.Parse(data)
//  if err != nil {
//      return nil, err
//  }
//  return &TruetypeFont{
//      Font: tt,
//  }
func LoadFont(data []byte) (*TruetypeFont, error) {
	f, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	return &TruetypeFont{
		Font: f,
	}, nil
}

// LoadFontFile loads a truetype font file and returns it or any errors
// encountered during loading. It is short-hand for:
//  data, err := ioutil.ReadFile(path)
//  if err != nil {
//      return nil, err
//  }
//  return LoadFont(data)
func LoadFontFile(path string) (*TruetypeFont, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return LoadFont(data)
}
