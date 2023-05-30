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

func distance(a, b Vector3D) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	dz := a.Z - b.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
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

func normalize(v Vector3D) Vector3D {
	length := math.Sqrt(dotProduct(v, v))
	if length == 0 { // Check for zero length
		return v
	}
	return multiplyByScalar(v, 1.0/length)
}

func isNearlyParallel(v1, v2 Vector3D, tolerance float64) bool {
	// Calculate the dot product of normalized vectors
	v1Normalized := normalize(v1)
	v2Normalized := normalize(v2)
	cosineOfAngle := dotProduct(v1Normalized, v2Normalized)

	// Cosine of angle will be nearly 1 for nearly parallel vectors
	return math.Abs(cosineOfAngle-1) < tolerance || math.Abs(cosineOfAngle+1) < tolerance
}

func findInitialVectors(points3D []Vector3D) []Vector3D {
	centroid := Vector3D{}
	for _, point := range points3D {
		centroid = add(centroid, point)
	}
	centroid = multiplyByScalar(centroid, 1.0/float64(len(points3D)))

	translatedPoints := []Vector3D{}
	for _, point := range points3D {
		translatedPoints = append(translatedPoints, subtract(point, centroid))
	}

	for i := 0; i < len(translatedPoints); i++ {
		for j := i + 1; j < len(translatedPoints); j++ {
			if isNearlyParallel(translatedPoints[i], translatedPoints[j], 0.01) {
				continue
			}
			return []Vector3D{translatedPoints[i], translatedPoints[j]}
		}
	}

	return nil
}

func isZeroVector(v Vector3D) bool {
	const epsilon = 1e-10
	return math.Abs(v.X) < epsilon && math.Abs(v.Y) < epsilon && math.Abs(v.Z) < epsilon
}

func projectPointTo2D(point, refPoint, dir1, dir2 Vector3D) Vector2D {
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

func mergeClosePoints(faces [][]Vector3D, threshold float64) [][]Vector3D {
	type Cluster struct {
		points []Vector3D
	}

	// List of all clusters
	clusters := []*Cluster{}

	for _, face := range faces {
		for _, point := range face {
			found := false
			for _, cluster := range clusters {
				for _, clusterPoint := range cluster.points {
					if distance(point, clusterPoint) < threshold {
						cluster.points = append(cluster.points, point)
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				newCluster := &Cluster{[]Vector3D{point}}
				clusters = append(clusters, newCluster)
			}
		}
	}

	// Build a map from original points to their corresponding centroid
	pointToCentroid := map[Vector3D]Vector3D{}
	for _, cluster := range clusters {
		centroid := calculateCentroid(cluster.points)
		for _, point := range cluster.points {
			pointToCentroid[point] = centroid
		}
	}

	// Build the new set of faces
	newFaces := make([][]Vector3D, len(faces))
	for i, face := range faces {
		newFace := make([]Vector3D, len(face))
		for j, point := range face {
			newFace[j] = pointToCentroid[point]
		}
		newFaces[i] = newFace
	}

	return newFaces
}

func calculateCentroid(points []Vector3D) Vector3D {
	var sumX, sumY, sumZ float64
	for _, point := range points {
		sumX += point.X
		sumY += point.Y
		sumZ += point.Z
	}
	length := float64(len(points))
	return Vector3D{sumX / length, sumY / length, sumZ / length}
}

func triangulate(p []Vector2D, holes ...[]Vector2D) [][]Vector2D {
	if len(p) < 3 {
		return nil
	} else if len(p) == 3 {
		return [][]Vector2D{p}
	}
	var triangles [][]Vector2D

	// Create flat points array
	flatPoints := flatten(p)
	holeIndices := []int{}
	for _, hole := range holes {
		holeIndices = append(holeIndices, len(flatPoints)/2)
		flatPoints = append(flatPoints, flatten(hole)...)
	}

	indices, _ := earcut.Earcut(flatPoints, holeIndices, 2)
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
func transform(points3D []Vector3D, holes3D ...[]Vector3D) [][]Vector3D {
	if len(points3D) < 3 {
		println("Not enough points to transform")
		return nil
	}
	// And two direction vectors dir1 and dir2 in the plane
	vectors := findInitialVectors(points3D)
	orthonormalBasis := gramSchmidt(vectors...)
	if len(orthonormalBasis) < 2 {
		println("Cannot find orthonormal basis")
		return nil
	}
	dir1, dir2 := orthonormalBasis[0], orthonormalBasis[1]

	// Let's project all points to 2D
	points2D := []Vector2D{}
	for _, point3D := range points3D {
		point2D := projectPointTo2D(point3D, points3D[0], dir1, dir2)
		points2D = append(points2D, point2D)
	}

	// Project holes to 2D
	holes2D := [][]Vector2D{}
	for _, hole3D := range holes3D {
		hole2D := []Vector2D{}
		for _, point3D := range hole3D {
			point2D := projectPointTo2D(point3D, points3D[0], dir1, dir2)
			hole2D = append(hole2D, point2D)
		}
		holes2D = append(holes2D, hole2D)
	}

	// Triangulate the 2D points
	triangles2D := triangulate(points2D, holes2D...)

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

func FindBasis(points []Vector3D) []Vector3D {
	vectors := findInitialVectors(points)
	orthonormalBasis := gramSchmidt(vectors...)
	if len(orthonormalBasis) < 2 {
		println("Cannot find orthonormal basis")
		return nil
	}
	return orthonormalBasis
}

func ProjectShapeTo2D(plane []Vector3D, basis []Vector3D) []Vector2D {
	points2D := []Vector2D{}
	for _, point3D := range plane {
		point2D := projectPointTo2D(point3D, plane[0], basis[0], basis[1])
		points2D = append(points2D, point2D)
	}
	return points2D
}

func ProjectShapeTo3D(shape []Vector2D, orthonormalBasis []Vector3D, refPoint Vector3D) []Vector3D {
	// Project 2D points back to 3D
	points3D := []Vector3D{}
	for _, point2D := range shape {
		point3D := projectPointTo3D(point2D, refPoint, orthonormalBasis[0], orthonormalBasis[1])
		points3D = append(points3D, point3D)
	}
	return points3D
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

func Earcut(faces [][]Vector3D, holes ...[][]Vector3D) [][]Vector3D {
	triangles := [][]Vector3D{}
	for i, face := range faces {
		_holes := [][]Vector3D{}
		if len(holes) > i {
			_holes = holes[i]
		}
		triangles = append(triangles, transform(face, _holes...)...)
	}
	triangles = mergeClosePoints(triangles, 0.01) // merge points closer than 0.01
	return triangles
}
