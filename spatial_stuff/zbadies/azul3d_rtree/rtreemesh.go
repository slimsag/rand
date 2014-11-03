package main

import (
	"azul3d.org/v1/gfx"
	"azul3d.org/v1/math"
	"azul3d.org/v1/rtree"
	"fmt"
)

func CubeVerts(start, end math.Vec3, result []gfx.Vec3) []gfx.Vec3 {
	addV := func(x, y, z float32) {
		result = append(result, gfx.Vec3{x, y, z})
	}

	x0 := float32(start.X)
	y0 := float32(start.Y)
	z0 := float32(start.Z)

	x1 := float32(end.X)
	y1 := float32(end.Y)
	z1 := float32(end.Z)

	addV(x0, y0, z0)
	addV(x0, y0, z1)
	addV(x0, y1, z1)

	addV(x1, y1, z0)
	addV(x0, y0, z0)
	addV(x0, y1, z0)

	addV(x1, y0, z1)
	addV(x0, y0, z0)
	addV(x1, y0, z0)

	addV(x1, y1, z0)
	addV(x1, y0, z0)
	addV(x0, y0, z0)

	addV(x0, y0, z0)
	addV(x0, y1, z1)
	addV(x0, y1, z0)

	addV(x1, y0, z1)
	addV(x0, y0, z1)
	addV(x0, y0, z0)

	addV(x0, y1, z1)
	addV(x0, y0, z1)
	addV(x1, y0, z1)

	addV(x1, y1, z1)
	addV(x1, y0, z0)
	addV(x1, y1, z0)

	addV(x1, y0, z0)
	addV(x1, y1, z1)
	addV(x1, y0, z1)

	addV(x1, y1, z1)
	addV(x1, y1, z0)
	addV(x0, y1, z0)

	addV(x1, y1, z1)
	addV(x0, y1, z0)
	addV(x0, y1, z1)

	addV(x1, y1, z1)
	addV(x0, y1, z1)
	addV(x1, y0, z1)
	return result
}

func LineVerts(start, end math.Vec3, width float64, result []gfx.Vec3) []gfx.Vec3 {
	hw := math.Vec3Zero.AddScalar(width).DivScalar(2.0)
	start = start.Sub(hw)
	end = end.Add(hw)
	return CubeVerts(start, end, result)
}

func LinesVerts(points []math.Vec3, connect bool, width float64, result []gfx.Vec3) []gfx.Vec3 {
	if len(points) < 0 {
		panic("LinesVerts(): must provide at least 2 points")
	}
	for i := 0; i < len(points); i += 2 {
		result = LineVerts(points[i], points[i+1], width, result)
	}
	return result
}

func Rect3Verts(r math.Rect3, result []gfx.Vec3) []gfx.Vec3 {
	return CubeVerts(r.Min, r.Max, result)
}

func Rect3LineVerts(r math.Rect3, result []gfx.Vec3) []gfx.Vec3 {
	leftBackBottom := math.Vec3{r.Min.X, r.Min.Y, r.Min.Z}
	leftBackTop := math.Vec3{r.Min.X, r.Min.Y, r.Max.Z}
	leftFrontBottom := math.Vec3{r.Min.X, r.Max.Y, r.Min.Z}
	leftFrontTop := math.Vec3{r.Min.X, r.Max.Y, r.Max.Z}

	rightBackBottom := math.Vec3{r.Max.X, r.Min.Y, r.Min.Z}
	rightBackTop := math.Vec3{r.Max.X, r.Min.Y, r.Max.Z}
	rightFrontBottom := math.Vec3{r.Max.X, r.Max.Y, r.Min.Z}
	rightFrontTop := math.Vec3{r.Max.X, r.Max.Y, r.Max.Z}

	width := 0.005
	return LinesVerts([]math.Vec3{
		// Left-to-right
		leftFrontBottom,
		rightFrontBottom,
		leftBackBottom,
		rightBackBottom,
		leftFrontTop,
		rightFrontTop,
		leftBackTop,
		rightBackTop,

		// Back-to-front
		leftBackBottom,
		leftFrontBottom,
		rightBackBottom,
		rightFrontBottom,
		leftBackTop,
		leftFrontTop,
		rightBackTop,
		rightFrontTop,

		// Bottom-to-top
		leftBackBottom,
		leftBackTop,
		leftFrontBottom,
		leftFrontTop,
		rightBackBottom,
		rightBackTop,
		rightFrontBottom,
		rightFrontTop,
	}, true, width, result)
}

func RTreeMesh(tree *rtree.Tree, level int) *gfx.Mesh {
	m := new(gfx.Mesh)

	bci := -1
	nextBC := func() gfx.Color {
		bci++
		switch bci % 3 {
		case 0:
			return gfx.Color{1, 0, 0, 0}
		case 1:
			return gfx.Color{0, 1, 0, 0}
		case 2:
			return gfx.Color{0, 0, 1, 0}
		}
		panic("never here.")
	}

	nodes := 0
	var add func(n *rtree.Node)
	add = func(n *rtree.Node) {
		//if nodes > 100 {
		//	return
		//}
		nodes++
		if n == nil {
			return
		}
		if level == 0 || nodes == level {
			// Add vertices.
			m.Vertices = Rect3LineVerts(n.Bounds(), m.Vertices)
		}

		/*for _, s := range n.Objects {
			c := s.Bounds().Center()
			sz := 0.005
			sb := math.Rect3{
				Min: math.Vec3Zero.AddScalar(-sz).Add(c),
				Max: math.Vec3Zero.AddScalar(sz).Add(c),
			}
			m.Vertices = Rect3LineVerts(sb, m.Vertices)
		}*/

		for _, c := range n.Children {
			add(c)
		}
	}
	add(tree.Root())

	fmt.Println(tree.Count(), "objects", nodes, "nodes")

	for _ = range m.Vertices {
		// Add barycentric coordinates.
		m.Colors = append(m.Colors, nextBC())
	}
	return m
}
