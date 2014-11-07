
// isFilled tests whether the given point is within a filled region of the
// glyph, by utilizing the non-zero winding rule.
//
// The test is extremely simple making no consideration of the bezier curve
// itself (it only considers linear segments).
func (g QuadGlyph) isFilled(p Point) bool {
	log.Println("BEGIN")
	count := 0
	far := lineSegment{
		Start: p,
		End: Point{
			X: p.X+10000,
			Y: p.Y+10000,
		},
	}

	for c := 0; c < g.NumContours(); c++ {
		contour := g.Contour(c)
		last := contour[0]
		for _, pt := range contour[1:] {
			line := lineSegment{Start:last, End:pt}
			_, status := line.intersect(far)
			if status == hit {
				if far.sideOfPoint(last) < 0 {
					count++
				} else {
					count--
				}
				log.Println(count, far.sideOfPoint(last))
			}
			last = pt
		}
	}
	log.Println(count)
	return count != 0
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















// intersect performs intersection against every line segments composed of the
// first and third points (i.e. skipping the control points) of every contour
// of the given glyph.
//
// The results are appended to the given point buffer and returned.
func (m *GlyphMesher) intersect(buf []Point, l lineSegment, g QuadGlyph, verts []gfx.Vec3) (hitBuffer []Point, vertsOut []gfx.Vec3) {
	return
	hitBuffer = buf

	var scale float32 = 0.0005
	var ptSize float32 = 0.01

	// FIXME: put under a DEBUG constant.
	switch swtch {
	//case 0:
	//	verts = m.debugPoints(verts, g, scale, ptSize*2, true)
	case 1:
		verts = m.debugPoints(verts, g, scale, ptSize*2, true)
		verts = m.debugPoints(verts, g, scale, ptSize, false)
		return nil, verts
	}

	// FIXME: garbage
	/*
	testLine := func(a, b Point) {
		p0 := gfx.Vec3{float32(a.X) * scale, 0, float32(a.Y) * scale}
		p1 := gfx.Vec3{float32(b.X) * scale, 0, float32(b.Y) * scale}
		verts = appendLine(verts, p0, p1, ptSize)
	}

	// Test each line of each contour, while ignoring control points.
	for c := 0; c < g.NumContours(); c++ {
		points := g.Contour(c)
		log.Println(len(points))
		for i := 2; i < len(points); i += 2 {
			testLine(points[i], points[i-2])
		}

		// Contour closing.
		testLine(points[0], points[len(points)-2])
	}
	*/

	return nil, verts
}




















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











	metrics, err := m.f.Measure(i)
	if err != nil {
		return err
	}


