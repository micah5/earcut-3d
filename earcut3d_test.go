package earcut3d

import (
	"math"
	"os"
	"reflect"
	"testing"
)

var cube = [][]Vector3D{
	{{0, 0, 0}, {0, 0, 1}, {0, 1, 1}, {0, 1, 0}},
	{{1, 0, 0}, {1, 1, 0}, {1, 1, 1}, {1, 0, 1}},
	{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}, {0, 0, 1}},
	{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}, {1, 1, 0}},
	{{0, 0, 0}, {0, 1, 0}, {1, 1, 0}, {1, 0, 0}},
	{{0, 0, 1}, {1, 0, 1}, {1, 1, 1}, {0, 1, 1}},
}

func convertToInts(vectors [][]Vector3D) [][]Vector3D {
	rounded := [][]Vector3D{}
	for _, face := range vectors {
		roundedFace := []Vector3D{}
		for _, vertex := range face {
			roundedFace = append(roundedFace, Vector3D{
				float64(math.Round(vertex.X)),
				float64(math.Round(vertex.Y)),
				float64(math.Round(vertex.Z)),
			})
		}
		rounded = append(rounded, roundedFace)
	}
	return rounded
}

func TestEarcut(t *testing.T) {
	triangles := Earcut(cube)
	triangles = convertToInts(triangles)

	// Check if output is correct
	expectedTriangles := [][]Vector3D{
		{{0, 1, 1}, {0, 1, 0}, {0, 0, 0}},
		{{0, 0, 0}, {0, 0, 1}, {0, 1, 1}},
		{{1, 1, 1}, {1, 0, 1}, {1, 0, 0}},
		{{1, 0, 0}, {1, 1, 0}, {1, 1, 1}},
		{{1, 0, 1}, {0, 0, 1}, {0, 0, 0}},
		{{0, 0, 0}, {1, 0, 0}, {1, 0, 1}},
		{{1, 1, 1}, {1, 1, 0}, {0, 1, 0}},
		{{0, 1, 0}, {0, 1, 1}, {1, 1, 1}},
		{{1, 1, 0}, {1, 0, 0}, {0, 0, 0}},
		{{0, 0, 0}, {0, 1, 0}, {1, 1, 0}},
		{{1, 1, 1}, {0, 1, 1}, {0, 0, 1}},
		{{0, 0, 1}, {1, 0, 1}, {1, 1, 1}},
	}
	if !reflect.DeepEqual(triangles, expectedTriangles) {
		t.Fatalf("Output was incorrect, got: %+v, want: %+v.", triangles, expectedTriangles)
	}
}

func TestProjection(t *testing.T) {
	inputFace := cube[0]
	basis := FindBasis(inputFace)
	points2D := ProjectShapeTo2D(inputFace, basis)
	points3D := ProjectShapeTo3D(points2D, basis, inputFace[0])
	points3D = convertToInts([][]Vector3D{points3D})[0]

	// Check that points3D are the same as the original points
	if !reflect.DeepEqual(points3D, inputFace) {
		t.Fatalf("Output was incorrect, got: %+v, want: %+v.", points3D, cube[0])
	}
}

func TestCreateObjFile(t *testing.T) {
	triangles := Earcut(cube)

	// Create file
	filename := "/tmp/cube.obj"
	CreateObjFile(filename, triangles)

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("File %s was not created", filename)
	}
	defer os.Remove(filename)
}
