package utils

import (
	"arrays"
	"math"
	"math/rand"
)

type Range1D struct {
	LowerBound float64
	UpperBound float64
}

func RandomInRange(r Range1D) float64 {
	return r.LowerBound + (r.UpperBound-r.LowerBound)*rand.Float64()
}

func PickRandom(setSize, howMany, exclude int) *[]int {
	alreadyPicked := make([]bool, setSize)
	picked := make([]int, howMany)
	for i := range picked {
		for {
			agentNumber := rand.Intn(setSize)
			if agentNumber != exclude && !alreadyPicked[agentNumber] {
				picked[i] = agentNumber
				alreadyPicked[agentNumber] = true
				break
			}
		}
	}
	return &picked
}

func ConstrainValue(value float64, constraint Range1D) float64 {
	return math.Max(constraint.LowerBound,
		math.Min(value, constraint.UpperBound))
}

func ConstrainArray(array arrays.Array1D, constraint Range1D) *arrays.Array1D {
	for i, value := range array {
		array.Set(i, ConstrainValue(value, constraint))
	}
	return &array
}
