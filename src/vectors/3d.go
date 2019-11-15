package vectors

import (
	"arrays"
	"fmt"
	"math"
	"utils"
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

func RandomVector3D(constraint utils.Range1D) Vector3D {
	x := utils.RandomInRange(constraint)
	y := utils.RandomInRange(constraint)
	z := utils.RandomInRange(constraint)
	return NewVector3D(x, y, z)
}

func Vector3DFromMatrix(array *arrays.Array2D) (Vector3D, error) {
	if array.NRows() == 1 && array.NColumns() == 4 {
		return vector3DFromRow(array), nil
	} else if array.NColumns() == 1 && array.NRows() == 4 {
		return vector3DFromColumn(array), nil
	}
	return Vector3D{}, fmt.Errorf("invalid matrix")
}

func vector3DFromRow(array *arrays.Array2D) Vector3D {
	return Vector3D{
		array.GetValue( 0, 0),
		array.GetValue( 0, 1),
		array.GetValue( 0, 2),
		array.GetValue( 0, 3),
	}
}

func vector3DFromColumn(array *arrays.Array2D) Vector3D {
	return Vector3D{
		array.GetValue( 0, 0),
		array.GetValue( 1, 0),
		array.GetValue( 2, 0),
		array.GetValue( 3, 0),
	}
}

func (v Vector3D) AsRow() *arrays.Array2D {
	return &arrays.Array2D{{v.X, v.Y, v.Z, v.W}}
}

func (v Vector3D) AsColumn() *arrays.Array2D {
	return &arrays.Array2D{{v.X}, {v.Y}, {v.Z}, {v.W}}
}

func (v Vector3D) Transform(matrix *arrays.Array2D) Vector3D {
	transformed := matrix.Multiply(v.AsColumn())
	vector, _ := Vector3DFromMatrix(transformed)
	return vector
}

func (v Vector3D) Distance(otherVector Vector3D) float64 {
	dX := v.X - otherVector.X
	dY := v.Y - otherVector.Y
	dZ := v.Z - otherVector.Z
	return math.Sqrt(dX*dX + dY*dY + dZ*dZ)
}
