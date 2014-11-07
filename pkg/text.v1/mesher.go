// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"log"
	"runtime"
	"sync"

	"azul3d.org/gfx.v1"
	"azul3d.org/native/tess.v1"
)

// cacheEntry represents a single cache entry.
type cacheEntry struct {
	// Number of uses (i.e. importance) of this entry.
	use int
}

// GlyphMesher is capable of generating meshes from glyph data.
type GlyphMesher struct {
	f     Font
	cache map[FontIndex]*cacheEntry
	byUse []*cacheEntry
}

// Append appends an appropriate mesh for the given font index to the given
// mesh object.
//
// There are cases (spaces, etc) where no vertices will be appended to the
// mesh, thus making it invalid. You should explicitly check for this validity.
//
// Internally the mesher makes use of a cache to make appends more efficient.
func (m *GlyphMesher) Append(mesh *gfx.Mesh, i FontIndex) error {
	// FIXME: sync
	// FIXME: use m.cache
	// FIXME: egg m.byUse
	// FIXME: determine scale

	// FIXME: remove
	swtch++
	if swtch > 1 {
		swtch = 0
	}
	log.Println(swtch)

	// Lookup the glyph data for the font index.
	data, err := m.f.Lookup(i)
	if err != nil {
		return err
	}
	glyph := data.(QuadGlyph)

	mesh.Vertices = m.appendTris(glyph, mesh.Vertices)

	// FIXME: uncomment
	/*
	// Generate curve data for the glyph.
	var interps []gfx.Vec3
	mesh.Vertices, interps = m.appendCurves(glyph, mesh.Vertices, interps)

	// Add the Interp vertex attribute to the mesh.
	if mesh.Attribs == nil {
		mesh.Attribs = make(map[string]gfx.VertexAttrib)
	}
	if len(interps) > 0 {
		mesh.Attribs["Interp"] = gfx.VertexAttrib{
			Data: interps,
		}
	}
	*/
	return nil
}

const ptScale = 10

// appendCurves appends all of the concave and convex quadratic bezier curves
// found in the glyph to the given verts and interps slices.
func (m *GlyphMesher) appendCurves(g QuadGlyph, verts, interps []gfx.Vec3) (vertsOut, interpsOut []gfx.Vec3) {
	// FIXME: put under a DEBUG constant.
	switch swtch {
	//case 0:
	//	verts = m.debugPoints(verts, g, 1, ptScale*2, true)
	case 1:
		verts = m.debugPoints(verts, g, 1, ptScale*2, true)
		verts = m.debugPoints(verts, g, 1, ptScale, false)
		return verts, interps
	}

	addCurve := func(a, b, c Point) {
		// If the control point, b, is in-between the line segment a-c, then
		// this is not a curve but simply a linear line segment -- thus we can
		// omit generation of it.
		ac := lineSegment{
			Start: a,
			End: c,
		}
		if ac.isBetween(b) {
			// It's a linear line segment.
			return
		}

		// Determine if the curve is concave (otherwise it's convex) by simply
		// examaning which side the point is on.
		var concave float32 = 0
		l := lineSegment{Start: a, End: c}
		if l.sideOfPoint(b) > 0 {
			concave = 1
		}

		// Append the actual triangle.
		verts = append(verts, gfx.Vec3{float32(a.X), 0, float32(a.Y)})
		verts = append(verts, gfx.Vec3{float32(b.X), 0, float32(b.Y)})
		verts = append(verts, gfx.Vec3{float32(c.X), 0, float32(c.Y)})

		// Append the interpretation values for the bezier curve shader. The Y
		// value serves just to identify concave/vs/convex curves.
		interps = append(interps, gfx.Vec3{0, concave, 0})
		interps = append(interps, gfx.Vec3{.5, concave, 0})
		interps = append(interps, gfx.Vec3{1, concave, 1})
	}

	for c := 0; c < g.NumContours(); c++ {
		points := g.Contour(c)
		for i := 2; i < len(points); i += 2 {
			addCurve(points[i], points[i-1], points[i-2])
		}

		// Contour closing.
		addCurve(points[0], points[len(points)-1], points[len(points)-2])
	}

	return verts, interps
}

// appendTris appends all of the fill triangles for the glyph to the given
// verts slice.
func (m *GlyphMesher) appendTris(g QuadGlyph, verts []gfx.Vec3) (vertsOut []gfx.Vec3) {
	contours := make([][]float64, g.NumContours())
	for c := 0; c < g.NumContours(); c++ {
		contour := g.Contour(c)
		contours[c] = make([]float64, len(contour)*2)
		for i, pt := range contour {
			contours[c][i] = float64(pt.X)
			contours[c][i+1] = float64(pt.Y)
		}
	}
	log.Println(len(contours), "contours")

	input := tess.Input{
		Contours: contours,
		WindingRule: tess.WindingNonZero,
		PolySize: 6,
	}
	t := tess.New()
	t.Tesselate(input)

	var triFan []gfx.Vec3
	addVert := func(x, y float32) {
		//log.Println(x, y)
		//verts = appendSq(verts, gfx.Vec3{X:x, Z:y}, ptScale*2)
		triFan = append(triFan, gfx.Vec3{X:x, Z:y})
		if len(triFan) == 3 {
			verts = append(verts, triFan[0])
			verts = append(verts, triFan[1])
			verts = append(verts, triFan[2])
			triFan = triFan[:1]
		}
	}

	for i := 0; i < t.ElementCount; i++ {
		p := t.Elements[i*input.PolySize: (i*input.PolySize) + input.PolySize]
		for _, pj := range p {
			if pj == -1 || (pj*2)+1 >= len(t.Vertices) { break }
			log.Println(pj*2, len(t.Vertices))
			v0 := t.Vertices[pj*2]
			v1 := t.Vertices[(pj*2)+1]
			addVert(float32(v0), float32(v1))
		}
	}

	for len(verts) % 3 != 0 {
		log.Println("add", len(verts) % 3, "len", len(verts))
		verts = append(verts, gfx.Vec3{float32(0), 0, float32(0)})
	}
	log.Println("add", len(verts) % 3, "len", len(verts))
	return verts
}

var swtch int // FIXME: remove

var (
	// A limit on the maximum number of glyph meshers that can exist per-font.
	//
	// -1 means that a value of GOMAXPROCS should be used, which effectively
	// allows each OS thread to have it's own mesher.
	MaxMeshers = -1

	// The number of glyphs that can be stored in a glyph mesher's individual
	// cache.
	GlyphCacheSize = 300

	// Explicitly not public because it's very implementation dependent.
	cacheEntryUseLimit = 50

	// Map of meshers per font object.
	fontMeshersAccess sync.RWMutex
	fontMeshers       = make(map[Font][]*GlyphMesher, 32)
)

// FindGlyphMesher finds a glyph mesher for the given font. If a mesher for the
// given font object does not exist, one is created.
//
// The MaxMeshers variable allows for finer control over the maximum number of
// meshers that can exist per font object, although the default value is often
// the best.
//
// In order to maximize distribution of meshers across goroutines, it makes the
// most sense to call FindGlyphMesher somewhat often rather than storing the
// returned one and using it later. The lookup is implemented using a map.
func FindGlyphMesher(f Font) *GlyphMesher {
	// Determine the literal MaxMeshers value.
	maxMeshers := MaxMeshers
	if maxMeshers == -1 {
		maxMeshers = runtime.GOMAXPROCS(-1)
	}
	if maxMeshers <= 0 {
		panic("text: MaxMeshers has invalid value")
	}

	// Sync access to the map.
	fontMeshersAccess.Lock()
	defer fontMeshersAccess.Unlock()

	// Check if we have an existing mesher object and are at our limit of
	// meshers for that font object.
	meshers, ok := fontMeshers[f]
	if ok && len(meshers) < maxMeshers {
		chosen := meshers[0]

		// Rotate meshers so that subsequent calls will return the chosen one
		// last.
		last := meshers[0]
		for i := 1; i < len(meshers); i++ {
			m := meshers[i]
			meshers[i] = last
			last = m
		}
		return chosen
	}

	// Create a new mesher.
	m := &GlyphMesher{
		f: f,
	}
	m.cache = make(map[FontIndex]*cacheEntry, GlyphCacheSize)
	m.byUse = make([]*cacheEntry, 0, GlyphCacheSize)

	// Store the mesher for later use.
	fontMeshers[f] = []*GlyphMesher{m}
	return m
}
