package transformations

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

type Rotation struct {
	I, J, K, W float64
}

func RotationFromAxisAngle(axis *mat.VecDense, theta float64) *Rotation {
	return &Rotation{
		axis.At(0, 0) * math.Sin(theta/2),
		axis.At(1, 0) * math.Sin(theta/2),
		axis.At(2, 0) * math.Sin(theta/2),
		math.Cos(theta / 2),
	}
}

func RotationFromEulerAngles(phi, theta, psi float64) *Rotation {
	return &Rotation{
		math.Sin(phi/2)*math.Cos(theta/2)*math.Cos(psi/2) - math.Cos(phi/2)*math.Sin(theta/2)*math.Sin(psi/2),
		math.Cos(phi/2)*math.Sin(theta/2)*math.Cos(psi/2) + math.Sin(phi/2)*math.Cos(theta/2)*math.Sin(psi/2),
		math.Cos(phi/2)*math.Cos(theta/2)*math.Sin(psi/2) - math.Sin(phi/2)*math.Sin(theta/2)*math.Cos(psi/2),
		math.Cos(phi/2)*math.Cos(theta/2)*math.Cos(psi/2) + math.Sin(phi/2)*math.Sin(theta/2)*math.Sin(psi/2),
	}
}

func (r *Rotation) GetTransformation() *mat.Dense {
	return mat.NewDense(4, 4, []float64{
		1 - 2*(r.J*r.J+r.K*r.K), 0 + 2*(r.I*r.J-r.K*r.W), 0 + 2*(r.I*r.K+r.J*r.W), 0,
		0 + 2*(r.I*r.J+r.K*r.W), 1 - 2*(r.I*r.I+r.K*r.K), 0 + 2*(r.J*r.K-r.I*r.W), 0,
		0 + 2*(r.I*r.K-r.J*r.W), 0 + 2*(r.J*r.K+r.I*r.W), 1 - 2*(r.I*r.I+r.J*r.J), 0,
		0, 0, 0, 1,
	})
}
