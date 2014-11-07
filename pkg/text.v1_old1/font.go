// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

// Point represents a single point measured in 1/72th inch units.
type Point struct {
	X, Y int
}

// Pixels converts this point in 1/72th inch units to pixel units, given the
// DPI (dots per inch).
func (p Point) Pixels(dpi int) (x, y int) {
	x = (p.X * dpi) / 72
	y = (p.X * dpi) / 72
	return
}

// Points converts the (x, y) point in pixels to 1/72th inch units, given the
// DPI (dots per inch).
func Points(x, y, dpi int) Point {
	return Point{
		X: (x / dpi) * 72.0,
		Y: (y / dpi) * 72.0,
	}
}

// QuadGlyph represents a glyph composed of 2D quadratic curves and utilizing
// the non-zero winding rule to determine fill.
type QuadGlyph struct {
	// Points is a slice of points representing quadratic bezier curves that
	// fully define the shape of the glyph.
	//
	// The points are in the order of:
	//  1. Start Point
	//  2. Control Point
	//  3. End Point
	Points []Point

	// Indices is a slice of indices into the Points slice such that each index
	// contours[last:index] forms the complete set of start-to-end points that
	// completely make up the contour, where every second point is an off-curve
	// quadratic bezier curve control point. For example:
	//  last := 0
	//  for _, index := range glyph.Indices {
	//      drawCurves(glyph.Points[last:index])
	//      last = index
	//  }
	//
	// Consider using the NumContours() and Contour() methods to give code more
	// clarity.
	Indices []int
}

// NumContours is short-handed for:
//  len(g.Indices)
func (g QuadGlyph) NumContours() int {
	return len(g.Indices)
}

// Contour returns the subslice of g.Points that represents the i'th contour of
// the glyph.
//
// To iterate over all contours one could write:
//  for i := 0; i < g.NumContours(); i++ {
//      points := g.Contour(i)
//  }
func (g QuadGlyph) Contour(i int) []Point {
	var start int
	if i > 0 {
		start = g.Indices[i-1]
	}
	end := g.Indices[i]
	return g.Points[start:end]
}

// GlyphBounds represents the bounding box (composed of a minimum and maximum
// point) of a single glyph.
type GlyphBounds struct {
	Min, Max Point
}

// GlyphMetrics represents the measurements of a single glyph. Many of the
// metrics listed here are described n great detail at:
//  http://www.freetype.org/freetype2/docs/glyphs/glyphs-3.html
type GlyphMetrics struct {
	// Bounding box of the glyph
	Bounds GlyphBounds

	// Horizontal and vertical bearings.
	BearingX, BearingY int

	// Horizontal and vertical advances.
	AdvanceX, AdvanceY int
}

// FontIndex represents the index in a font that a glyph can be located at.
type FontIndex uint

// Font is a generic font source provider.
type Font interface {
	// Index locates the font index for the given rune. If there is no index in
	// the font relating to the given rune, then FontIndex(0) is returned.
	Index(r rune) FontIndex

	// Lookup looks up the glyph data associated with the given index and
	// returns it.
	//
	// If any error occured during lookup, glyphData will be nil and the error
	// will be returned.
	//
	// The returned interface (the glyph data) solely represents the shape of
	// the glyph. If the glyph data is not one of the types listed below then
	// any operation using the data may cause a panic:
	//  QuadGlyph
	Lookup(i FontIndex) (glyphData interface{}, err error)

	// Measure measures the glyph associated with the given index and returns
	// it's measurements.
	//
	// If any error occured during lookup, GlyphSize will be nil and the error
	// will be returned.
	Measure(i FontIndex) (*GlyphMetrics, error)

	// Bounds returns a bounding box that describes the maximum (i.e. union) of
	// all glyphs in the font, such that the returned bounding box is equal
	// to or larger than any arbitrary glyph's size within this font.
	//
	// It is useful for performing quick measurements, for instance (although
	// it will overshoot).
	Bounds() GlyphBounds

	// Kerning returns the amount of horizontal and vertical kerning that is
	// between the given two glyphs associated with the given font indices. If
	// the kerning amount is not known for any axis, -1 is returned for that
	// axis.
	Kerning(a, b FontIndex) (x, y int)
}
