package objects

import "gonum.org/v1/gonum/mat"

type Triangle struct {
	RawVertices         [3]*mat.Dense
	VertNormals         [3]*mat.VecDense
	FaceNormal          *mat.VecDense
	TransformedVertices [3]*mat.Dense
	// TODO more dynamic system
	Visibility Sideness
	DoDraw     bool
}

type VertexObject struct {
	Faces           []*Triangle
	Transformations *mat.Dense
}

type Sideness uint8

const (
	FRONT Sideness = 0
	BACK  Sideness = 1
	BOTH  Sideness = 2
)
