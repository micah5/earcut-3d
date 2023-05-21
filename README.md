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
```
Make sure every face is coplanar. If not, the library will print a warning.
