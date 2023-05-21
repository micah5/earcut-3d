package earcut3d

import (
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

func TestEarcut(t *testing.T) {
	triangles := Earcut(cube)

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
