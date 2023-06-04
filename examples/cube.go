package main

import (
	"github.com/micah5/earcut-3d"
)

func main() {
	// Input format is flattened [][]float64 ({{x0, y0, z0, x1, y1, z1...}...}
	// Each inner slice is a coplanar polygon/face
	cube := [][]float64{
		{0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0}, // Front
		{1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1}, // Back
		{0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1}, // Bottom
		{0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0}, // Top
		{0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0}, // Left
		{0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1}, // Right
	}
	triangles := earcut3d.Earcut(cube)

	// Now you can do what you want with the triangles (same format as input- [][]float64)
	// For example, create an obj file to visualize it:
	earcut3d.CreateObjFile("cube.obj", triangles)
}
