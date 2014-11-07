// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"sync"
)

// FontCache caches access to a font source, improving repetitive access to it
// significantly.
type FontCache struct {
	access      sync.RWMutex
	source      FontSource
	indexLookup map[rune]FontIndex
}

// Implements the FontSource interface.
func (f *FontCache) Index(r rune) FontIndex {
	// If we have the index cached already, then simply return it/
	f.access.RLock()
	index, ok := f.indexLookup[r]
	f.access.RUnlock()
	if ok {
		return index
	}

	// Perform a lookup and cache the result for later.
	index = f.source.Index(r)
	f.access.Lock()
	f.indexLookup[r] = index
	f.access.Unlock()
	return index
}

// FontSource is a generic font source provider.
type FontSource interface {
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
	//  []QuadCurve
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
