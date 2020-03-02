package shaders

import (
	"../objects"
	"image/color"
	"sort"
)

type ShadeFaces struct {
	Colors []*color.RGBA
}

// FIXME when a triangle is turning over to show its other side, the area gets small and it has issues

func (s ShadeFaces) Shade(object objects.VertexObject, canvas Canvas) {
	for faceI, face := range object.Faces {
		if face.DoDraw {

			z := (face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3

			color := s.Colors[faceI%len(s.Colors)]

			// Make a copy just so we're not messing with the original transformed verticies, not that it should really
			// matter but whatever
			sortable := face.TransformedVertices
			sort.Slice(sortable, func(i, j int) bool {
				return sortable[i].At(1, 0) < sortable[j].At(1, 0)
			})
			// For use in the left/right splitting
			y0 := sortable[0].At(1, 0)
			y1 := sortable[1].At(1, 0)
			y2 := sortable[2].At(1, 0)
			x0 := sortable[0].At(0, 0)
			x1 := sortable[1].At(0, 0)
			x2 := sortable[2].At(0, 0)

			// Main left/right splitting
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

			// If the vertex furthest to the left is vertically between the other two, i.e. if the x value of the middle
			// y value is the smallest (vice versa for other case)
			if x1 <= x0 && x1 <= x2 {
				leftXStep = twoSide
				rightXStep = oneSide
			} else if x1 >= x0 && x1 >= x2 {
				leftXStep = oneSide
				rightXStep = twoSide
			} else {
				// The bottom must be flat
				// FIXME causing the weird line glitchy bits
				leftXStep = oneSide
				rightXStep = twoSide
			}

			leftX := x2
			rightX := x2
			for y := int(y2); float64(y) >= y0; y-- {
				s.drawHorizontalLine(int(leftX), int(rightX), y, z, color, canvas)
				leftStep, leftReset := leftXStep(y)
				rightStep, rightReset := rightXStep(y)
				if leftStep > 25 {
					leftStep = 0
				}
				if rightStep > 25 {
					rightStep = 0
				}
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
	for x := x0; x <= x1; x++ {
		if canvas.SetRGBA(x, y, z, *color) {
			return
		}
	}
}
