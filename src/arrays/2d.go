package arrays

import (
	"fmt"
	"log"
	"strings"
)

type Array2D []*Array1D

func NewArray2D(rows, cols int) *Array2D {
	array := make(Array2D, rows)
	for i := 0; i < rows; i++ {
		line := make(Array1D, cols)
		array[i] = &line
	}
	return &array
}

func (a *Array2D) String() string {
	s := make([]string, a.NRows())
	for i := range a.Items() {
		s[i] = a.GetRow(i).String()
	}
	return fmt.Sprintf("[%s]", strings.Join(s, "\n"))
}

func (a *Array2D) Items() [][]float64 {
	items := make([][]float64, a.NRows())
	for i := range items {
		items[i] = append(items[i], a.GetRow(i).Items()...)
	}
	return items
}

func Identity2D(order int) *Array2D {
	identity := NewArray2D(order, order)
	for i := 0; i < order; i++ {
		identity.SetValue(i, i, 1)
	}
	return identity
}

func (a *Array2D) NRows() int {
	if a == nil {
		return 0
	}
	return len(*a)
}

func (a *Array2D) NColumns() int {
	if a.NRows() == 0 {
		return 0
	}
	return a.GetRow(0).Length()
}

func (a *Array2D) NElements() int {
	return a.NRows() * a.NColumns()
}

func (a *Array2D) Append(v1d *Array1D) {
	*a = append(*a, v1d)
}

func (a *Array2D) Copy() *Array2D {
	array := make(Array2D, a.NRows())
	copy(array, *a)
	return &array
}

func (a *Array2D) GetRow(position int) *Array1D {
	return (*a)[position]
}

func (a *Array2D) GetValue(row, col int) float64 {
	return a.GetRow(row).Get(col)
}

func (a *Array2D) SetRow(position int, value Array1D) {
	if row := a.GetRow(position); row != nil {
		*row = value
	} else {
		(*a)[position] = &value
	}
}

func (a *Array2D) SetValue(rown, coln int, value float64) error {
	if row := a.GetRow(rown); row != nil {
		row.Set(coln, value)
		return nil
	}
	return fmt.Errorf("row %d is nil", rown)
}

func (a *Array2D) Multiply(otherArray *Array2D) *Array2D {
	if a.NColumns() != otherArray.NRows() {
		log.Fatalf("Arrays are not compatible for multiplication:\n%s\n%s", a, otherArray)
	}
	nRows := a.NRows()
	nCols := otherArray.NColumns()
	result := NewArray2D(nRows, nCols)
	for i := 0; i < nRows; i++ {
		for j := 0; j < nCols; j++ {
			for k := 0; k < a.NColumns(); k++ {
				value := result.GetValue(i, j) + a.GetValue(i, k)*otherArray.GetValue(k, j)
				result.SetValue(i, j, value)
			}
		}
	}
	return result
}
