package shaders

import (
	"../objects"
	"image/color"
	"math"
	"sort"
)

type ShadeFacesInt struct {
	Colors []*color.RGBA
}

func (s ShadeFacesInt) Shade(object objects.VertexObject, canvas Canvas) {
	for faceI, face := range object.Faces {
		if face.DoDraw {

			z := (face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3

			randomColor := s.Colors[faceI%len(s.Colors)]

			// Make a copy just so we're not messing with the original transformed verticies, not that it should really
			// matter but whatever
			sortable := face.TransformedVertices
			sort.Slice(sortable, func(i, j int) bool {
				return sortable[i].At(1, 0) < sortable[j].At(1, 0)
			})
			// For use in the left/right splitting
			y0 := int(math.Round(sortable[0].At(1, 0)))
			y1 := int(math.Round(sortable[1].At(1, 0)))
			y2 := int(math.Round(sortable[2].At(1, 0)))
			x0 := int(math.Round(sortable[0].At(0, 0)))
			x1 := int(math.Round(sortable[1].At(0, 0)))
			x2 := int(math.Round(sortable[2].At(0, 0)))

			var cutoffX int
			slope := y2 - y0
			if slope == 0 {
				cutoffX = x2
			} else {
				cutoffX = (y1-y0)*(x2-x0)/(y2-y0) + x0
			}

			s.drawTriangle(x0, y0, x1, cutoffX, y1, z, randomColor, canvas)
			s.drawTriangle(x2, y2, x1, cutoffX, y1, z, randomColor, canvas)
			//s.drawTriangle(x0, y0, x1, x2, y1, z, randomColor, canvas)
			//s.drawTriangle(x2, y2, x0, y0, x1, y1, z, randomColor, canvas)
			//s.drawTriangle(x2, y2, x1, y1, x0, y0, z, randomColor, canvas)
		}
	}
}

func (s *ShadeFacesInt) drawTriangle(x0, y0, x1, x2, yk int, z float64, color *color.RGBA, canvas Canvas) {
	xA := x0
	xB := x0
	yN := y0

	dxA := x1 - xA
	if dxA < 0 {
		dxA = -dxA
	}
	dyA := yk - y0
	if dyA < 0 {
		dyA = -dyA
	}
	var sxA int
	if xA < x1 {
		sxA = 1
	} else {
		sxA = -1
	}
	errA := dxA - dyA

	dxB := x2 - xB
	if dxB < 0 {
		dxB = -dxB
	}
	dyB := yk - y0
	if dyB < 0 {
		dyB = -dyB
	}
	var sxB int
	if xB < x2 {
		sxB = 1
	} else {
		sxB = -1
	}
	errB := dxB - dyB

	var sy int
	if y0 < yk {
		sy = 1
	} else {
		sy = -1
	}

	for {

		for {
			canvas.SetRGBA(xA, yN, z, *color)
			if xA == x1 && yN == yk {
				return
			}
			e2A := 2 * errA
			if e2A > -dyA {
				errA -= dyA
				xA += sxA
			}
			if e2A < dxA {
				errA += dxA
				break
			}
		}

		for {
			canvas.SetRGBA(xB, yN, z, *color)
			if xB == x2 && yN == yk {
				return
			}
			e2B := 2 * errB
			if e2B > -dyB {
				errB -= dyB
				xB += sxB
			}
			if e2B < dxB {
				errB += dxB
				break
			}
		}

		yN += sy
		if xA < xB {
			s.drawHorizontalLine(xA, xB, yN, z, color, canvas)
		} else if xA > xB {
			s.drawHorizontalLine(xB, xA, yN, z, color, canvas)
		}
	}

}

func (s *ShadeFacesInt) drawHorizontalLine(x0, x1, y int, z float64, color *color.RGBA, canvas Canvas) {
	for x := x0; x <= x1; x++ {
		if canvas.SetRGBA(x, y, z, *color) {
			return
		}
	}
}
