package animations

import (
	".."
	"gonum.org/v1/gonum/mat"
	"math"
)

type SingleAxisSpin struct {
	Axis     *mat.VecDense
	Operator float64
	GetClock func() float64
}

func (a *SingleAxisSpin) GetTransformation() *mat.Dense {
	return transformations.RotationFromAxisAngle(a.Axis, math.Mod(a.GetClock(), a.Operator)/a.Operator*2*math.Pi).GetTransformation()
}
