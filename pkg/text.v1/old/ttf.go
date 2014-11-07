package text

import (
	"log"
	"math"
	"sort"

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

func absInt32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func minXDist(points []truetype.Point) int32 {
	var min int32 = math.MaxInt32
	for pi, p := range points {
		for ti, t := range points {
			if pi == ti {
				continue
			}
			dist := absInt32(t.X - p.X)
			if dist < min {
				min = dist
			}
		}
	}
	if min == 0 {
		min = 1
	}
	return min
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
	//points = expressPoints(points)
	scale := float32(0.001)

	lrsweep := &lRSweeper{
		points: points,
	}
	_ = lrsweep

	ltr := make([]truetype.Point, len(points))
	copy(ltr, points)
	sort.Sort(leftToRight(ltr))

	yMin := float32(buf.B.YMin)
	yMax := float32(buf.B.YMax)
	for kstart, start := range ltr {
		if kstart > 10 && swtch {
			break
		}
		sweepLine := lineSegment{
			Start: lmath.Vec2{float64(start.X), float64(yMin)},
			End:   lmath.Vec2{float64(start.X), float64(yMax)},
		}
		// Draw sweep line.
		appendLine(
			m,
			gfx.Vec3{(float32(start.X) * scale) - 0.75, 0, (float32(yMin) * scale) - 0.75},
			gfx.Vec3{(float32(start.X) * scale) - 0.75, 0, (float32(yMax) * scale) - 0.75},
			0.005,
		)

		anyHit := false
		for i := 0; i < len(points); i += 1 {
			p0 := points[i]
			i2 := i + 1
			if i2 == len(points) {
				i2 = 0
			}
			p1 := points[i2]
			k := lineSegment{
				Start: lmath.Vec2{float64(p0.X), float64(p0.Y)},
				End:   lmath.Vec2{float64(p1.X), float64(p1.Y)},
			}

			// Intersect the line segment K with the sweep line.
			p, hit := k.intersect(sweepLine)
			if hit {
				anyHit = true
				appendSq(m, gfx.Vec3{(float32(p.X) * scale) - 0.75, 0, (float32(p.Y) * scale) - 0.75}, 0.015)
			}

			appendLine(
				m,
				gfx.Vec3{(float32(p0.X) * scale) - 0.75, 0, (float32(p0.Y) * scale) - 0.75},
				gfx.Vec3{(float32(p1.X) * scale) - 0.75, 0, (float32(p1.Y) * scale) - 0.75},
				0.005,
			)
		}

		if !anyHit {
			appendSq(m, gfx.Vec3{(float32(start.X) * scale) - 0.75, 0, (float32(yMin) * scale) - 0.8}, 0.015)
		}

		//_ = sweepLine
		//appendSq(m, gfx.Vec3{(float32(p.X)*scale)-0.75, 0, (float32(p.Y)*scale)-0.75}, 0.015)
	}

	/*
		xMin := float32(buf.B.XMin)
		xMax := float32(buf.B.XMax)
		yMin := float32(buf.B.YMin)
		yMax := float32(buf.B.YMax)
		incr := float32(minXDist(points))
		incr = 32

		prevX := xMin
		for x := xMin+incr; x < xMax; x += incr {
			sweepLine := lineSegment{
				Start: lmath.Vec2{float64(x), float64(yMin-100)},
				End: lmath.Vec2{float64(x), float64(yMax+100)},
			}
			appendLine(
				m,
				gfx.Vec3{(float32(x)*scale)-0.75, 0, (float32(yMin)*scale)-0.75},
				gfx.Vec3{(float32(x)*scale)-0.75, 0, (float32(yMax)*scale)-0.75},
				0.005,
			)
			for i := 0; i < len(points)-2; i+=2 {
				p0 := points[i]
				p1 := points[i+1]
				k := lineSegment{
					Start: lmath.Vec2{float64(p0.X), float64(p0.Y)},
					End: lmath.Vec2{float64(p1.X), float64(p1.Y)},
				}

				// Intersect the line segment K with the sweep line.
				p, hit := k.intersect(sweepLine)
				if hit {
					appendSq(m, gfx.Vec3{(float32(p.X)*scale)-0.75, 0, (float32(p.Y)*scale)-0.75}, 0.015)
				} else {
					appendSq(m, gfx.Vec3{(float32(x)*scale)-0.75, 0, (float32(yMin)*scale)-0.8}, 0.015)
				}


				appendLine(
					m,
					gfx.Vec3{float32(p0.X)*scale, 0, float32(p0.Y)*scale},
					gfx.Vec3{float32(p1.X)*scale, 0, float32(p1.Y)*scale},
					0.01,
				)
			}
			prevX = x
		}
		_ = prevX
	*/

	/*
		// Draw lines.
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

type leftToRight []truetype.Point

func (p leftToRight) Len() int      { return len(p) }
func (p leftToRight) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p leftToRight) Less(ii, jj int) bool {
	i := p[ii].X
	j := p[jj].X
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

	if len(buf.End) > 1 {
		return
	}
	e0 := 0
	for _, e1 := range buf.End {
		appendContour(buf, m, buf.Point[e0:e1])
		e0 = e1
		break
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
