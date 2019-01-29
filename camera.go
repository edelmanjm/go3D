package main

import "gonum.org/v1/gonum/mat"

type Camera struct {
	// Camera is always at the origin
	Normal *mat.VecDense
}
