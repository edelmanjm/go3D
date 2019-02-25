package transformations

import "gonum.org/v1/gonum/mat"

type Translation struct {
	X, Y, Z float64
}

func (translation *Translation) GetTransformation() *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		1, 0, 0, translation.X,
		0, 1, 0, translation.Y,
		0, 0, 1, translation.Z,
		0, 0, 0, 1,
	})
}
