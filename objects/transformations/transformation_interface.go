package transformations

import "gonum.org/v1/gonum/mat"

type Transformation interface {
	GetTransformation() *mat.Dense
}
