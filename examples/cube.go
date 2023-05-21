package main

import (
	"github.com/micah5/earcut-3d"
)

func main() {
	// Input format is [][]Vector3D, where each inner slice is a coplanar polygon/face
	cube := [][]earcut3d.Vector3D{
		{{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}}, // Front
		{{1, 0, 0}, {1, 1, 0}, {1, 1, 1}, {1, 0, 1}}, // Back
		{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}}, // Bottom
		{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}}, // Top
		{{0, 0, 0}, {0, 1, 0}, {1, 1, 0}, {1, 0, 0}}, // Left
		{{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}}, // Right
	}
	triangles := earcut3d.Earcut(cube)

	// Now you can do what you want with the triangles (same format as input- [][]Vector3D)
	// For example, create an obj file to visualize it:
	earcut3d.CreateObjFile("cube.obj", triangles)
}
