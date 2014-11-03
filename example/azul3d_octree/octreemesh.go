package main

import (
	"fmt"

	"azul3d.org/gfx.v1"
	"azul3d.org/lmath.v1"
	"azul3d.org/octree.v1"
)

func CubeVerts(start, end lmath.Vec3, result []gfx.Vec3) []gfx.Vec3 {
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

func LineVerts(start, end lmath.Vec3, width float64, result []gfx.Vec3) []gfx.Vec3 {
	hw := lmath.Vec3Zero.AddScalar(width).DivScalar(2.0)
	start = start.Sub(hw)
	end = end.Add(hw)
	return CubeVerts(start, end, result)
}

func LinesVerts(points []lmath.Vec3, connect bool, width float64, result []gfx.Vec3) []gfx.Vec3 {
	if len(points) < 0 {
		panic("LinesVerts(): must provide at least 2 points")
	}
	for i := 0; i < len(points); i += 2 {
		result = LineVerts(points[i], points[i+1], width, result)
	}
	return result
}

func Rect3Verts(r lmath.Rect3, result []gfx.Vec3) []gfx.Vec3 {
	return CubeVerts(r.Min, r.Max, result)
}

func Rect3LineVerts(r lmath.Rect3, result []gfx.Vec3) []gfx.Vec3 {
	leftBackBottom := lmath.Vec3{r.Min.X, r.Min.Y, r.Min.Z}
	leftBackTop := lmath.Vec3{r.Min.X, r.Min.Y, r.Max.Z}
	leftFrontBottom := lmath.Vec3{r.Min.X, r.Max.Y, r.Min.Z}
	leftFrontTop := lmath.Vec3{r.Min.X, r.Max.Y, r.Max.Z}

	rightBackBottom := lmath.Vec3{r.Max.X, r.Min.Y, r.Min.Z}
	rightBackTop := lmath.Vec3{r.Max.X, r.Min.Y, r.Max.Z}
	rightFrontBottom := lmath.Vec3{r.Max.X, r.Max.Y, r.Min.Z}
	rightFrontTop := lmath.Vec3{r.Max.X, r.Max.Y, r.Max.Z}

	width := 0.005
	return LinesVerts([]lmath.Vec3{
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

func OctreeMesh(tree *octree.Tree, level int) *gfx.Mesh {
	m := gfx.NewMesh()

	nc := 0
	var add func(n *octree.Node)
	add = func(n *octree.Node) {
		if n == nil {
			return
		}
		nc++
		if nc > 6000 {
			return
		}
		if level == 0 || n.Level() == level {
			// Add vertices.
			m.Vertices = Rect3LineVerts(n.Bounds(), m.Vertices)
		}

		//fmt.Println("node", nodes, len(n.Objects), "objects")
		/*for oct := 0; oct < 9; oct++ {
			for i := 0; i < n.NumObjects(oct); i++ {
				s := n.Object(oct, i)
				c := s.Bounds().Center()
				sz := 0.005
				sb := math.Rect3{
					Min: math.Vec3Zero.AddScalar(-sz).Add(c),
					Max: math.Vec3Zero.AddScalar(sz).Add(c),
				}
				sb = s.Bounds()
				m.Vertices = Rect3LineVerts(sb, m.Vertices)
			}
		}*/

		for c := 0; c < 8; c++ {
			add(n.Child(octree.ChildIndex(c)))
		}
	}
	add(tree.Root())

	fmt.Printf("%d objects - %d nodes\n", tree.NumObjects(), tree.NumNodes())

	// Generate the barycentric coordinates.
	//m.GenerateBary()
	return m
}
