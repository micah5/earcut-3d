package main

import (
	"github.com/micah5/earcut-3d"
	"math"
)

func main() {
	radius := 1.0
	divisions := 9

	vertices := make([]earcut3d.Vector3D, (divisions+1)*(divisions+1))
	faces := make([][]int, divisions*divisions)

	index := 0
	for i := 0; i <= divisions; i++ {
		phi := math.Pi * float64(i) / float64(divisions)
		for j := 0; j <= divisions; j++ {
			theta := 2.0 * math.Pi * float64(j) / float64(divisions)

			x := radius * math.Sin(phi) * math.Cos(theta)
			y := radius * math.Sin(phi) * math.Sin(theta)
			z := radius * math.Cos(phi)

			vertices[index] = earcut3d.Vector3D{x, y, z}
			index++
		}
	}

	// Create faces
	for i := 0; i < divisions; i++ {
		for j := 0; j < divisions; j++ {
			// check if inner part of sphere
			// excluding the poles because they're triangles and too lazy to deal with them separately
			if i != 0 && i != divisions-1 {
				lowerLeft := i*(divisions+1) + j + 1
				lowerRight := lowerLeft + 1
				upperRight := lowerRight + divisions + 1
				upperLeft := lowerLeft + divisions + 1

				face := []int{lowerLeft, upperLeft, upperRight, lowerRight}
				faces[i*divisions+j] = face
			}
		}
	}

	// Prepare for writing
	sphere := make([][]earcut3d.Vector3D, 0)
	for _, face := range faces {
		if len(face) == 4 {
			sphere = append(sphere, []earcut3d.Vector3D{vertices[face[0]-1], vertices[face[1]-1], vertices[face[2]-1], vertices[face[3]-1]})
		}
	}

	triangles := earcut3d.Earcut(sphere)
	earcut3d.CreateObjFile("sphere.obj", triangles)
}
