package main

import (
	"github.com/micah5/earcut-3d"
)

func main() {
	plane := [][]float64{
		{0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0},
	}
	holes := [][]float64{
		{0, 0.25, 0.25, 0, 0.25, 0.75, 0, 0.45, 0.75, 0, 0.45, 0.25},
		{0, 0.55, 0.25, 0, 0.55, 0.75, 0, 0.75, 0.75, 0, 0.75, 0.25},
	}
	triangles := earcut3d.Earcut(plane, holes)

	earcut3d.CreateObjFile("hole.obj", triangles)
}
