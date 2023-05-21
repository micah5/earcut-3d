package earcut3d

import (
	"fmt"
	"github.com/rclancey/go-earcut"
	"math"
	"os"
)

type Vector3D struct {
	X, Y, Z float64
}

type Vector2D struct {
	X, Y float64
}

func flatten(vectors []Vector2D) []float64 {
	flat := []float64{}
	for _, v := range vectors {
		flat = append(flat, v.X, v.Y)
	}
	return flat
}

func subtract(v1, v2 Vector3D) Vector3D {
	return Vector3D{v1.X - v2.X, v1.Y - v2.Y, v1.Z - v2.Z}
}

func add(v1, v2 Vector3D) Vector3D {
	return Vector3D{v1.X + v2.X, v1.Y + v2.Y, v1.Z + v2.Z}
}

func dotProduct(v1, v2 Vector3D) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func multiplyByScalar(v Vector3D, scalar float64) Vector3D {
	return Vector3D{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func gramSchmidt(vectors ...Vector3D) []Vector3D {
	orthonormalBasis := []Vector3D{}

	for _, vi := range vectors {
		if isZeroVector(vi) { // Check for zero vectors
			continue
		}
		for _, vj := range orthonormalBasis {
			vi = subtract(vi, multiplyByScalar(vj, dotProduct(vi, vj)))
		}

		// Normalize
		length := math.Sqrt(dotProduct(vi, vi))
		if length == 0 { // Check for zero length
			continue
		}
		vi = multiplyByScalar(vi, 1.0/length)

		orthonormalBasis = append(orthonormalBasis, vi)
	}

	return orthonormalBasis
}

func isZeroVector(v Vector3D) bool {
	const epsilon = 1e-10
	return math.Abs(v.X) < epsilon && math.Abs(v.Y) < epsilon && math.Abs(v.Z) < epsilon
}

func projectPointToPlane(point, refPoint, dir1, dir2 Vector3D) Vector2D {
	translatedPoint := subtract(point, refPoint) // Translate to origin
	x := dotProduct(translatedPoint, dir1)
	y := dotProduct(translatedPoint, dir2)
	return Vector2D{x, y}
}

func projectPointTo3D(point2D Vector2D, refPoint, dir1, dir2 Vector3D) Vector3D {
	xComponent := multiplyByScalar(dir1, point2D.X)
	yComponent := multiplyByScalar(dir2, point2D.Y)

	return add(refPoint, add(xComponent, yComponent))
}

func triangulate(p []Vector2D) [][]Vector2D {
	if len(p) < 3 {
		return nil
	} else if len(p) == 3 {
		return [][]Vector2D{p}
	}
	var triangles [][]Vector2D
	flatPoints := flatten(p)
	indices, _ := earcut.Earcut(flatPoints, nil, 2)
	for i := 0; i < len(indices); i += 3 {
		vertices := make([]Vector2D, 3)
		for j := 0; j < 3; j++ {
			index := indices[i+j] * 2
			vertices[j] = Vector2D{flatPoints[index], flatPoints[index+1]}
		}
		triangles = append(triangles, vertices)
	}
	return triangles
}

// Assuming we have 3D points a, b, c, d lying on the same plane
func transform(points3D []Vector3D) [][]Vector3D {
	// And two direction vectors dir1 and dir2 in the plane
	vectors := []Vector3D{subtract(points3D[1], points3D[0]), subtract(points3D[2], points3D[0])}
	orthonormalBasis := gramSchmidt(vectors...)
	if len(orthonormalBasis) < 2 {
		return nil
	}
	dir1, dir2 := orthonormalBasis[0], orthonormalBasis[1]

	// Let's project all points to 2D
	points2D := []Vector2D{}

	for _, point3D := range points3D {
		point2D := projectPointToPlane(point3D, points3D[0], dir1, dir2)
		points2D = append(points2D, point2D)
	}

	// Triangulate the 2D points
	triangles2D := triangulate(points2D)

	// Convert 2D points back to 3D
	triangles3D := [][]Vector3D{}
	for _, triangle := range triangles2D {
		points3DTransformed := []Vector3D{}
		for _, point2D := range triangle {
			point3D := projectPointTo3D(point2D, points3D[0], dir1, dir2)
			points3DTransformed = append(points3DTransformed, point3D)
		}
		triangles3D = append(triangles3D, points3DTransformed)
	}
	return triangles3D
}

func CreateObjFile(name string, triangles [][]Vector3D) {
	// Create file
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a map to store unique vertices and their indices
	vertexIndices := make(map[Vector3D]int)
	currentIndex := 1

	// Write triangles
	for _, triangle := range triangles {
		for _, point := range triangle {
			// If the vertex hasn't been seen before, write it and store its index
			if _, seen := vertexIndices[point]; !seen {
				f.WriteString(fmt.Sprintf("v %f %f %f\n", point.X, point.Y, point.Z))
				vertexIndices[point] = currentIndex
				currentIndex++
			}
		}
	}

	// Write faces
	for _, triangle := range triangles {
		f.WriteString("f")
		for _, point := range triangle {
			f.WriteString(fmt.Sprintf(" %d", vertexIndices[point]))
		}
		f.WriteString("\n")
	}
}

func Earcut(faces [][]Vector3D) [][]Vector3D {
	triangles := [][]Vector3D{}
	for _, face := range faces {
		triangles = append(triangles, transform(face)...)
	}
	return triangles
}
