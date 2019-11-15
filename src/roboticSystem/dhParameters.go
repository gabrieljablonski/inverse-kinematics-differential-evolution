package roboticSystem

import (
	"arrays"
	"math"
)

// https://en.wikipedia.org/wiki/Denavit%E2%80%93Hartenberg_parameters
type DHParameters struct {
	Theta float64
	D     float64
	R     float64 // alternate notation for `a`
	Alpha float64
}

func (dh DHParameters) TransformationMatrix() *arrays.Array2D {
	cosTheta := math.Cos(dh.Theta)
	sinTheta := math.Sin(dh.Theta)
	cosAlpha := math.Cos(dh.Alpha)
	sinAlpha := math.Sin(dh.Alpha)
	return &arrays.Array2D{
		{cosTheta, -sinTheta * cosAlpha, sinTheta * sinAlpha, dh.R * cosTheta},
		{sinTheta, cosTheta * cosAlpha, -cosTheta * sinAlpha, dh.R * sinTheta},
		{0, sinAlpha, cosAlpha, dh.D},
		{0, 0, 0, 1},
	}
}

func ParametersToTransformationMatrix(dhParams []DHParameters) *arrays.Array2D {
	baseMatrix := arrays.Identity2D(4)
	for i := 0; i < len(dhParams); i++ {
		baseMatrix = baseMatrix.Multiply(dhParams[i].TransformationMatrix())
	}
	return baseMatrix
}
