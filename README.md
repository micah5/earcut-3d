# earcut-3d

Go library for 3D earcut triangulation.

It works by projecting every plane to 2D, using [go-earcut](https://github.com/rclancey/go-earcut) to create the triangles and finally transforming it back to 3D.
I made it because of a shortcoming in [Mapbox's earcut library](https://github.com/mapbox/earcut) where it ignores the Z component (i.e. 2D only).

![My project2](https://github.com/micah5/earcut-3d/assets/40206415/2285024c-1bc6-48cc-9e69-d684f9e4f19a)

## Installation

```
go get github.com/micah5/earcut-3d
```

## Usage

```go
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
```

Make sure every face is coplanar. If not, the library will print a warning.
