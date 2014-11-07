package text

import (
	"log"
	"math"

	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
	"code.google.com/p/freetype-go/freetype/truetype"
)

/*
	m.Vertices = append(m.Vertices, []gfx.Vec3{
		// Bottom-left triangle.
		{0, 0, 1},
		{-1, 0, -1},
		{1, 0, -1},
	}...)
	m.Vertices = append(m.Vertices, gfx.Vec3{x, 0, y})
*/
//var s float32 = 0.0005

func scaled(p truetype.Point, s float32) gfx.Vec3 {
	return gfx.Vec3{
		X: float32(p.X) * s,
		Z: float32(p.Y) * s,
	}
}

func appendSq(m *gfx.Mesh, a gfx.Vec3, s float32) {
	v := func(x, y float32) {
		m.Vertices = append(m.Vertices, gfx.Vec3{x, 0, y})
	}
	left := a.X - s
	right := a.X + s
	bottom := a.Z - s
	top := a.Z + s

	v(left, bottom)
	v(left, top)
	v(right, top)

	v(left, bottom)
	v(right, top)
	v(right, bottom)
}

func appendLine(m *gfx.Mesh, a, b gfx.Vec3, w float32) {
	v := func(x, y float32) {
		m.Vertices = append(m.Vertices, gfx.Vec3{x, 0, y})
	}

	v(a.X, a.Z)
	v(b.X, b.Z)
	v(a.X+w, a.Z+w)

	v(b.X, b.Z)
	v(a.X, a.Z)
	v(b.X+w, b.Z+w)

	/*
		left := a.X
		bottom := a.Z
		right := b.X
		top := b.Z

		w := float32(0.007)
		bottom -= w
		top += w
		left -= w
		right += w

		v(left, bottom)
		v(left, top)
		v(right, top)

		v(left, bottom)
		v(right, top)
		v(right, bottom)
	*/
}

func absInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func minDist(points []truetype.Point) float32 {
	var fmin float64 = math.MaxInt32
	for pi, p := range points {
		for ti, t := range points {
			if pi == ti {
				continue
			}
			x := t.X - p.X
			y := t.Y - p.Y
			dist := math.Sqrt(float64(x*x + y*y))
			if dist < fmin {
				fmin = dist
			}
		}
	}
	return float32(fmin)
}

func inside(x, y, xMin, yMin, xMax, yMax float32, points []truetype.Point) int {
	w := xMin - xMax
	h := yMin - yMax
	l := lineSegment{
		Start: lmath.Vec2{float64(x - w), float64(y - h)},
		End:   lmath.Vec2{float64(x + w), float64(y + h)},
	}
	var count int
	for i := 0; i < len(points)-2; i += 2 {
		p0 := points[i]
		p1 := points[i+1]
		k := lineSegment{
			Start: lmath.Vec2{float64(p0.X), float64(p0.Y)},
			End:   lmath.Vec2{float64(p1.X), float64(p1.Y)},
		}
		_, hit := l.intersect(k)
		if hit {
			if l.sideOfPoint(k.Start) >= 0 {
				count++
			} else {
				count--
			}
		}
	}
	//log.Println(count)
	return count
}

var swtch bool

func appendContour(buf *truetype.GlyphBuf, m *gfx.Mesh, points []truetype.Point) {
	scale := float32(0.0005)
	points = expressPoints(points)

	xMin := float32(buf.B.XMin)
	xMax := float32(buf.B.XMax)
	yMin := float32(buf.B.YMin)
	yMax := float32(buf.B.YMax)
	incr := minDist(points)

	prevX := xMin
	for x := xMin + incr; x < xMax; x += incr {
		sweepLine := lineSegment{
			Start: lmath.Vec2{float64(x), float64(yMin)},
			End:   lmath.Vec2{float64(x), float64(yMax)},
		}
		for i := 0; i < len(points)-2; i += 2 {
			p0 := points[i]
			p1 := points[i+1]
			k := lineSegment{
				Start: lmath.Vec2{float64(p0.X), float64(p0.Y)},
				End:   lmath.Vec2{float64(p1.X), float64(p1.Y)},
			}

			// Intersect the line segment K with the sweep line.
			p, hit := sweepLine.intersect(k)
			if hit {
				appendSq(m, gfx.Vec3{float32(p.X) * scale, 0, float32(p.Y) * scale}, 0.015)
			}
			/*
				appendLine(
					m,
					gfx.Vec3{float32(p0.X)*scale, 0, float32(p0.Y)*scale},
					gfx.Vec3{float32(p1.X)*scale, 0, float32(p1.Y)*scale},
					0.01,
				)
			*/
		}
		prevX = x
	}
	_ = prevX

	// Draw lines.
	for i := 0; i < len(points)-2; i += 2 {
		p0 := points[i]
		p1 := points[i+1]
		appendLine(
			m,
			gfx.Vec3{float32(p0.X) * scale, 0, float32(p0.Y) * scale},
			gfx.Vec3{float32(p1.X) * scale, 0, float32(p1.Y) * scale},
			0.01,
		)
	}
	/*
		incr = 75.0
		pointSize := incr * 0.0001
		for x := xMin; x < xMax; x += incr {
			for y := yMin; y < yMax; y += incr {
				if swtch {
					appendSq(m, gfx.Vec3{x*scale, 0, y*scale}, pointSize)
					continue
				}
				//count := inside(x, y, xMin, yMin, xMax, yMax, points)
				//appendSq(m, gfx.Vec3{x*scale, 0, y*scale}, pointSize*float32(count))
			}
		}
		for i := 0; i < len(points)-2; i+=2 {
			p0 := points[i]
			p1 := points[i+1]
			appendLine(
				m,
				gfx.Vec3{float32(p0.X)*scale, 0, float32(p0.Y)*scale},
				gfx.Vec3{float32(p1.X)*scale, 0, float32(p1.Y)*scale},
				0.01,
			)
		}
	*/
}

type sortedPoints []truetype.Point

func (p sortedPoints) Len() int      { return len(p) }
func (p sortedPoints) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p sortedPoints) Less(ii, jj int) bool {
	i := p[ii].X + p[ii].Y
	j := p[jj].X + p[jj].Y
	return i < j
}

// expressPoints returns a slice with the implicit bezier curve points added.
//
// truetype fonts have implicit points in them. See figure 3 here:
//  https://developer.apple.com/fonts/TTRefMan/RM01/Chap1.html#direction
// or the discussion here:
//  http://stackoverflow.com/questions/20733790/truetype-fonts-glyph-are-made-of-quadratic-bezier-why-do-more-than-one-consecu
func expressPoints(points []truetype.Point) []truetype.Point {
	var (
		exp  []truetype.Point
		last truetype.Point
	)
	for _, p := range points {
		onCurve := (p.Flags & 0x1) > 0
		lastOnCurve := (last.Flags & 0x1) > 0
		if !onCurve && !lastOnCurve {
			// This is an implied on-curve point.
			implied := truetype.Point{
				Flags: 0x1, // on-curve
				X:     (p.X + last.X) / 2,
				Y:     (p.Y + last.Y) / 2,
			}
			exp = append(exp, implied)
		}
		exp = append(exp, p)
		last = p
	}
	return exp
}

func AppendMesh(buf *truetype.GlyphBuf, m *gfx.Mesh) {
	swtch = !swtch
	before := len(m.Vertices)
	if before > 0 {
		m.VerticesChanged = true
	}

	e0 := 0
	for _, e1 := range buf.End {
		appendContour(buf, m, buf.Point[e0:e1])
		e0 = e1
	}

	// Extend interps by the number of added vertices.
	var interp []gfx.Vec3
	attr, ok := m.Attribs["Interp"]
	if ok {
		interp = attr.Data.([]gfx.Vec3)
		attr.Changed = true
	}

	added := len(m.Vertices) - before
	log.Println(added, "vertices")
	for i := 0; i < added; i++ {
		var ip gfx.Vec3
		switch i % 3 {
		case 0:
			ip = gfx.Vec3{.5, 0, 0}
		case 1:
			ip = gfx.Vec3{0, 0, 0}
		case 2:
			ip = gfx.Vec3{1, 0, 1}
		}
		interp = append(interp, ip)
	}

	attr.Data = interp
	m.Attribs["Interp"] = attr
}
