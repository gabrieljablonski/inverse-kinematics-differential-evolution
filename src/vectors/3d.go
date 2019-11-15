package vectors

import (
	"arrays"
	"fmt"
	"math"
)

const (
	DefaultW = 1
)

type Vector3D struct {
	X, Y, Z float64
	W       float64
}

func (v Vector3D) String() string {
	return fmt.Sprintf("%.5f,%.5f,%.5f", v.X, v.Y, v.Z)
}

func NewVector3D(x, y, z float64) Vector3D {
	return Vector3D{
		X: x,
		Y: y,
		Z: z,
		W: DefaultW,
	}
}

func (v Vector3D) Distance(otherVector Vector3D) float64 {
	dX := v.X - otherVector.X
	dY := v.Y - otherVector.Y
	dZ := v.Z - otherVector.Z
	return math.Sqrt(dX*dX + dY*dY + dZ*dZ)
}

func (v Vector3D) AsRow() *arrays.Array2D {
	return &arrays.Array2D{{v.X, v.Y, v.Z, v.W}}
}

func (v Vector3D) AsColumn() *arrays.Array2D {
	return &arrays.Array2D{{v.X}, {v.Y}, {v.Z}, {v.W}}
}

func (v Vector3D) Transform(matrix *arrays.Array2D) Vector3D {
	transformed := matrix.Multiply(v.AsColumn())
	x := transformed.GetValue(0, 0)
	y := transformed.GetValue(1, 0)
	z := transformed.GetValue(2, 0)
	w := transformed.GetValue(3, 0)
	return Vector3D{x, y, z, w}
}
