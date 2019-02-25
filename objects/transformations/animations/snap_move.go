package animations

import (
	"gonum.org/v1/gonum/mat"
	"math"
)
import ".."

type SnapMove struct {
	X0, Y0, Z0,
	X1, Y1, Z1,
	Operator float64
	GetClock func() float64
}

func (a *SnapMove) GetTransformation() *mat.Dense {
	return (&transformations.Translation{
		a.X0 + (a.X1-a.X0)*(math.Mod(a.GetClock(), a.Operator))/a.Operator,
		a.Y0 + (a.Y1-a.Y0)*(math.Mod(a.GetClock(), a.Operator))/a.Operator,
		a.Z0 + (a.Z1-a.Z0)*(math.Mod(a.GetClock(), a.Operator))/a.Operator,
	}).GetTransformation()
}
