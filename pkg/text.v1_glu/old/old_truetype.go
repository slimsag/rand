// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"io/ioutil"
	"log"

	"code.google.com/p/freetype-go/freetype/raster"
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

	// Here we express implicit truetype points into a complete set of
	// quadratic bezier curves.
	//
	// Truetype fonts have implicit points in them. See figure 3 here:
	//  https://developer.apple.com/fonts/TTRefMan/RM01/Chap1.html#direction
	// or the discussion here:
	//  http://stackoverflow.com/questions/20733790/truetype-fonts-glyph-are-made-of-quadratic-bezier-why-do-more-than-one-consecu
	var (
		glyph  QuadGlyph
		lastPt = f.buf.Point[len(f.buf.Point)-1]
	)
	convertPoint := func(p truetype.Point) Point {
		return Point{
			X: int(p.X),
			Y: int(p.Y),
		}
	}
	e0 := 0
	for _, e := range f.buf.End {
		startIndex := len(glyph.Points)
		points := f.buf.Point[e0:e]
		for pi, p := range points {
			onCurve := (p.Flags & 0x1) != 0
			lastPtOnCurve := (lastPt.Flags & 0x1) != 0

			if pi == 0 && !onCurve {
				panic("first point of contour is not on-curve")
			}
			if pi == len(points)-1 && !onCurve {
				log.Println("last point of contour is not on-curve")
				//panic("last point of contour is not on-curve")
			}

			if !onCurve && !lastPtOnCurve {
				// This is an implied on-curve point.
				implied := Point{
					X: int(p.X+lastPt.X) / 2,
					Y: int(p.Y+lastPt.Y) / 2,
				}
				glyph.Points = append(glyph.Points, implied)
			} else if onCurve && lastPtOnCurve {
				// Two on-curve points -- insert midpoint.
				midpoint := Point{
					X: int(p.X+lastPt.X) / 2,
					Y: int(p.Y+lastPt.Y) / 2,
				}
				glyph.Points = append(glyph.Points, midpoint)
			}
			glyph.Points = append(glyph.Points, convertPoint(p))
			lastPt = p
		}
		endIndex := len(glyph.Points)
		newPoints := glyph.Points[startIndex:endIndex]
		if len(newPoints)%2 != 0 {
			log.Println("Failed to expand point set.")
			//panic("Failed to expand point set.")
		}
		glyph.Indices = append(glyph.Indices, endIndex)
		e0 = e
	}
	return glyph, nil
}

func (f *TruetypeFont) expandContour(points []truetype.Point)

// drawContour draws the given closed contour with the given offset.
func (c *Context) drawContour(ps []truetype.Point, dx, dy raster.Fix32) {
	if len(ps) == 0 {
		return
	}
	// ps[0] is a truetype.Point measured in FUnits and positive Y going upwards.
	// start is the same thing measured in fixed point units and positive Y
	// going downwards, and offset by (dx, dy)
	start := raster.Point{
		X: dx + raster.Fix32(ps[0].X<<2),
		Y: dy - raster.Fix32(ps[0].Y<<2),
	}
	c.r.Start(start)
	q0, on0 := start, true
	for _, p := range ps[1:] {
		q := raster.Point{
			X: dx + raster.Fix32(p.X<<2),
			Y: dy - raster.Fix32(p.Y<<2),
		}
		on := p.Flags&0x01 != 0
		if on {
			if on0 {
				c.r.Add1(q)
			} else {
				c.r.Add2(q0, q)
			}
		} else {
			if on0 {
				// No-op.
			} else {
				mid := raster.Point{
					X: (q0.X + q.X) / 2,
					Y: (q0.Y + q.Y) / 2,
				}
				c.r.Add2(q0, mid)
			}
		}
		q0, on0 = q, on
	}
	// Close the curve.
	if on0 {
		c.r.Add1(start)
	} else {
		c.r.Add2(q0, start)
	}
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
