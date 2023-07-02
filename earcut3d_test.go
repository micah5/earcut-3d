package earcut3d

import (
	"math"
	"os"
	"reflect"
	"testing"
)

var cube = [][]float64{
	{0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0},
	{1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1},
	{0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 1},
	{0, 1, 0, 0, 1, 1, 1, 1, 1, 1, 1, 0},
	{0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0},
	{0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 1},
}

func convertToInts(vectors [][]float64) [][]float64 {
	rounded := [][]float64{}
	for _, face := range vectors {
		roundedFace := []float64{}
		for _, vertex := range face {
			roundedFace = append(roundedFace, float64(math.Round(vertex)))
		}
		rounded = append(rounded, roundedFace)
	}
	return rounded
}

func TestEarcut(t *testing.T) {
	triangles := Earcut(cube)
	triangles = convertToInts(triangles)

	// Check if output is correct
	expectedTriangles := [][]float64{
		{0, 1, 1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 1, 1},
		{1, 1, 1, 1, 0, 1, 1, 0, 0},
		{1, 0, 0, 1, 1, 0, 1, 1, 1},
		{1, 0, 1, 0, 0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0, 0, 1, 0, 1},
		{1, 1, 1, 1, 1, 0, 0, 1, 0},
		{0, 1, 0, 0, 1, 1, 1, 1, 1},
		{1, 1, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 1, 1, 0},
		{1, 1, 1, 0, 1, 1, 0, 0, 1},
		{0, 0, 1, 1, 0, 1, 1, 1, 1},
	}
	// ignore order
	for _, triangle := range triangles {
		found := false
		for _, expectedTriangle := range expectedTriangles {
			if reflect.DeepEqual(triangle, expectedTriangle) {
				found = true
			}
		}
		if !found {
			t.Fatalf("Output was incorrect, got: %+v, want: %+v.", triangles, expectedTriangles)
		}
	}
}

func TestEarcutFaces(t *testing.T) {
	faces := EarcutFaces(cube)

	// check that size is correct
	if len(faces) != 6 {
		t.Fatalf("Output was incorrect, got: %+v, want: %+v.", len(faces), 6)
	}

	// check that each face has 2 triangles
	for _, face := range faces {
		if len(face) != 2 {
			t.Fatalf("Output was incorrect, got: %+v, want: %+v.", len(face), 2)
		}
	}
}

func TestProjection(t *testing.T) {
	inputFace := cube[0]
	basis := FindBasis(inputFace)
	points2D := ProjectShapeTo2D(inputFace, basis)
	points3D := ProjectShapeTo3D(points2D, basis, []float64{0, 0, 0})
	points3D = convertToInts([][]float64{points3D})[0]

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
