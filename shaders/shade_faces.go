package shaders

import (
	"../objects"
	"image/color"
	"sort"
)

type ShadeFaces struct {
	Colors []*color.RGBA
}

func (s ShadeFaces) Shade(object objects.VertexObject, canvas Canvas) {
	for faceI, face := range object.Faces {
		if face.DoDraw {

			z := (face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3

			color := s.Colors[faceI%len(s.Colors)]
			sort.Slice(face.TransformedVertices, func(i, j int) bool {
				return face.TransformedVertices[i].At(1, 0) < face.TransformedVertices[j].At(1, 0)
			})
			y0 := face.TransformedVertices[0].At(1, 0)
			y1 := face.TransformedVertices[1].At(1, 0)
			y2 := face.TransformedVertices[2].At(1, 0)
			x0 := face.TransformedVertices[0].At(0, 0)
			x1 := face.TransformedVertices[1].At(0, 0)
			x2 := face.TransformedVertices[2].At(0, 0)

			var leftXStep, rightXStep func(y int) (float64, bool)
			twoSide := func(y int) (float64, bool) {
				if y == int(y1) {
					return 0, true
				} else if y > int(y1) {
					return (x2 - x1) / (y2 - y1), false
				} else {
					return (x1 - x0) / (y1 - y0), false
				}
			}
			oneSide := func(y int) (float64, bool) {
				return (x2 - x0) / (y2 - y0), false
			}
			if x1 < x0 {
				leftXStep = twoSide
				rightXStep = oneSide
			} else {
				leftXStep = oneSide
				rightXStep = twoSide
			}

			leftX := x2
			rightX := x2
			for y := int(y2); float64(y) >= y0; y-- {
				s.drawHorizontalLine(int(leftX), int(rightX), y, z, color, canvas)
				leftStep, leftReset := leftXStep(y)
				rightStep, rightReset := rightXStep(y)
				if leftReset {
					leftX = x1
				} else {
					leftX -= leftStep
				}
				if rightReset {
					rightX = x1
				} else {
					rightX -= rightStep
				}
			}

		}
	}
}

func (s *ShadeFaces) drawHorizontalLine(x0, x1, y int, z float64, color *color.RGBA, canvas Canvas) {
	// FIXME TEMPORARY HARD CODING
	if x0 < 0 {
		x0 = 0
	}
	if x1 > 2880 {
		x1 = 2880
	}
	for x := x0; x <= x1; x++ {
		if canvas.SetRGBA(x, y, z, *color) {
			return
		}
	}
}
