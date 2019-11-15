package arrays

import (
	"fmt"
	"strings"
)

type Array1D []float64

func (a *Array1D) String() string {
	s := make([]string, a.Length())
	for i, value := range a.Items() {
		s[i] = fmt.Sprintf("%.5f", value)
	}
	return fmt.Sprintf("[%s]", strings.Join(s, ", "))
}

func (a *Array1D) Length() int {
	if a == nil {
		return 0
	}
	return len(*a)
}

func (a *Array1D) Items() []float64 {
	return *a
}

func (a *Array1D) Copy() *Array1D {
	array := make(Array1D, a.Length())
	copy(array, *a)
	return &array
}

func (a *Array1D) Append(value float64) {
	*a = append(*a, value)
}

func (a *Array1D) Get(position int) float64 {
	return (*a)[position]
}

func (a *Array1D) Set(position int, value float64) {
	(*a)[position] = value
}

func (a *Array1D) Sum() float64 {
	sum := 0.0
	for _, value := range a.Items() {
		sum += value
	}
	return sum
}

func (a *Array1D) Add(otherArray *Array1D) *Array1D {
	result := make(Array1D, a.Length())
	copy(result, *a)
	for i := range a.Items() {
		result[i] += otherArray.Get(i)
	}
	return &result
}

func (a *Array1D) Subtract(otherArray *Array1D) *Array1D {
	result := make(Array1D, a.Length())
	copy(result, *a)
	for i := range a.Items() {
		result[i] -= otherArray.Get(i)
	}
	return &result
}

func (a *Array1D) MultiplyByConstant(c float64) *Array1D {
	result := make(Array1D, a.Length())
	copy(result, *a)
	for i := range a.Items() {
		result[i] *= c
	}
	return &result
}
