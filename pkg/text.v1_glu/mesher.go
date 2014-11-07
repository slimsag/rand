// Copyright 2014 The Azul3D Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package text

import (
	"log"
	"runtime"
	"sync"

	"azul3d.org/native/tess.v1"
	"azul3d.org/gfx.v1"
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
// Internally the mesher uses a cache to make appends more efficient.
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

	metrics, err := m.f.Measure(i)
	if err != nil {
		return err
	}

	data, err := m.f.Lookup(i)
	glyph := data.(QuadGlyph)

	contours := make([][]float64, glyph.NumContours())
	for c := 0; c < glyph.NumContours(); c++ {
		contour := glyph.Contour(c)
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


/*
			for (i = 0; i < nelems; ++i)
			{
				const int* p = &elems[i*nvp];
				glBegin(GL_TRIANGLE_FAN);
				for (j = 0; j < nvp && p[j] != TESS_UNDEF; ++j)
					glVertex2f(verts[p[j]*2], verts[p[j]*2+1]);
				glEnd();
			}
*/

	var triFan []gfx.Vec3
	var s float32 = 0.015
	addVert := func(x, y float32) {
		x *= s
		y *= s
		triFan = append(triFan, gfx.Vec3{X:x, Z:y})
		if len(triFan) == 3 {
			mesh.Vertices = append(mesh.Vertices, triFan[0])
			mesh.Vertices = append(mesh.Vertices, triFan[1])
			mesh.Vertices = append(mesh.Vertices, triFan[2])
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

/*
	for i := 0; i < t.ElementCount; i++ {
		p := t.Elements[i*input.PolySize:(i*input.PolySize)+input.PolySize]
		start0 := float32(t.Vertices[p[0]*2])
		start1 := float32(t.Vertices[(p[0]*2)+1])
		log.Println(p[1]*2, len(t.Vertices))
		last0 := float32(t.Vertices[p[1]*2])
		last1 := float32(t.Vertices[(p[1]*2)+1])
		for v := 0; v+1 < len(p[2:]); v += 2 {
			v0 := float32(t.Vertices[v*2])
			v1 := float32(t.Vertices[(v*2)+1])
			mesh.Vertices = append(mesh.Vertices, []gfx.Vec3{
				{X: start0, Z: start1},
				{X: last0, Z: last1},
				{X: v0, Z: v1},
			}...)
			last0 = v0
			last1 = v1
		}
	}
*/

	/*
	polySize := 2
	log.Println("elems", t.ElementCount)
	for i := 0; i < t.ElementCount; i++ {
		poly := t.Elements[i * polySize:(i * polySize)+2]
		for j := 0; j < len(poly); j++ {
			if j == -1 { break }
			//log.Println("verts", len(t.Vertices), "mix", j*polySize, "j", j, "polySize", polySize)
			v0 := t.Vertices[j*polySize]
			v1 := t.Vertices[(j*polySize)+1]
			s := 0.004
			v0 *= s
			v1 *= s
			log.Println(v0, v1)
			mesh.Vertices = append(mesh.Vertices, gfx.Vec3{
				float32(v0), 0, float32(v1),
			})
		}
	}
	*/
	for len(mesh.Vertices) % 3 != 0 {
		log.Println("add", len(mesh.Vertices) % 3, "len", len(mesh.Vertices))
		mesh.Vertices = append(mesh.Vertices, gfx.Vec3{
			float32(0), 0, float32(0),
		})
	}
	log.Println("add", len(mesh.Vertices) % 3, "len", len(mesh.Vertices))

/*
  const int nelems = tessGetElementCount(tess);
  const TESSindex* elems = tessGetElements(tess);
  for (int i = 0; i < nelems; i++) {
      const TESSindex* poly = &elems[i * polySize];
      glBegin(GL_POLYGON);
      for (int j = 0; j < polySize; j++) {
          if (poly[j] == TESS_UNDEF) break;
          glVertex2fv(&verts[poly[j]*vertexSize]);
      }
      glEnd();
  }
*/	

	return nil
	// FIXME: remove all below this.
	var (
		lastHits, hits []Point
		hSweep         hSweeper
	)
	for c := 0; c < glyph.NumContours(); c++ {
		contour := glyph.Contour(c)
		// Extent the horizontal sweep point buffer so that it can hold the
		// contour's points.
		if len(contour) > len(hSweep.points) {
			// Copy the first half.
			copy(hSweep.points, contour)

			// Extend by appending the section half.
			hSweep.points = append(hSweep.points, contour[len(hSweep.points):len(contour)]...)
		} else {
			// Clamp the size.
			hSweep.points = hSweep.points[:len(contour)]

			// Copy the contour.
			copy(hSweep.points, contour)
		}

		// Sweep the contour's points horizontally.
		hSweep.sweep()

		// For each swept point from left-to-right, perform intersection.
		for _, pt := range hSweep.points {
			// Perform intersection of the vertical sweep line against all of
			// the glyph's uncurved contour lines.
			sweepLine := lineSegment{
				Start: Point{
					X: pt.X,
					Y: metrics.Bounds.Min.Y,
				},
				End: Point{
					X: pt.X,
					Y: metrics.Bounds.Max.Y,
				},
			}
			hits, mesh.Vertices = m.intersect(hits[:0], sweepLine, glyph, mesh.Vertices)

			// Generate vertices for the intersection points.
			//mesh.Vertices = m.generate(mesh.Vertices, lastHits, hits)

			// Swap the old and new buffers, allowing us to avoid a copy.
			lastHits, hits = hits, lastHits
		}
	}

	//appendSq(mesh, gfx.Vec3{0, 0, 0}, scale)
	return nil
}

func (m *GlyphMesher) generate(vertices []gfx.Vec3, lastHits, hits []Point) []gfx.Vec3 {
	var scale float32 = 0.0005

	var ptSize float32 = 0.02
	for _, hit := range hits {
		p0 := gfx.Vec3{float32(hit.X) * scale, 0, float32(hit.Y) * scale}
		vertices = appendSq(vertices, p0, ptSize)
	}

	/*
		add := func(p Point) {
			vertices = append(vertices, gfx.Vec3{
				X: float32(p.X) * scale,
				Z: float32(p.Y) * scale,
			})
		}

		// FIXME: no append/copy
		points := append(lastHits, hits...)
		sweeper := vSweeper{points}
		sweeper.sweep()
		for i := 0; i < len(points)-2; i+=2 {
			add(points[i])
			add(points[i+1])
			add(points[i+2])
			break
		}
	*/
	return vertices

	/*
		p0 := gfx.Vec3{0, 0, float32(metrics.Bounds.Min.Y) * lnScale}
		p1 := gfx.Vec3{0, 0, float32(metrics.Bounds.Max.Y) * lnScale}
		appendLine(mesh, p0, p1, ptSize)
		log.Println(p0, p1)
	*/
}

var swtch int // FIXME: remove

// intersect performs intersection against every line segments composed of the
// first and third points (i.e. skipping the control points) of every contour
// of the given glyph.
//
// The results are appended to the given point buffer and returned.
func (m *GlyphMesher) intersect(buf []Point, l lineSegment, g QuadGlyph, verts []gfx.Vec3) (hitBuffer []Point, vertsOut []gfx.Vec3) {
	hitBuffer = buf

	var scale float32 = 0.0005
	var ptSize float32 = 0.01

	switch swtch {
	case 0:
		verts = m.debugPoints(verts, g, scale, ptSize*2, true)
	case 1:
		verts = m.debugPoints(verts, g, scale, ptSize*2, true)
		verts = m.debugPoints(verts, g, scale, ptSize, false)
	case 2:
		verts = m.debugPoints(verts, g, scale, ptSize*2, true)
		verts = m.debugPoints(verts, g, scale, ptSize, false)
	}
	return nil, verts

	testLine := func(a, b Point) {
		p0 := gfx.Vec3{float32(a.X) * scale, 0, float32(a.Y) * scale}
		p1 := gfx.Vec3{float32(b.X) * scale, 0, float32(b.Y) * scale}
		verts = appendLine(verts, p0, p1, ptSize)

		/*
			cntLine := lineSegment{
				Start: points[i],
				End: points[i+2],
			}
			p, status := l.intersect(cntLine)
			if status == hit {
				hitBuffer = append(hitBuffer, p)
			}/* else if status == collinear {
				hitBuffer = append(hitBuffer, cntLine.Start)
				hitBuffer = append(hitBuffer, cntLine.End)
			}*/
	}

	// Test each line of each contour, while ignoring control points.
	for c := 0; c < g.NumContours(); c++ {
		points := g.Contour(c)
		log.Println(len(points))
		for i := 3; i < len(points); i += 2 {
			testLine(points[i], points[i-2])
		}

		// FIXME: remove
		switch swtch {
		case 0:
			// Contour closing.
			testLine(points[1], points[len(points)-1])
		}
		break
	}
	vertsOut = verts
	return
}

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
