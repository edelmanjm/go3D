package shaders

import (
	"../objects"
	"image"
	"image/color"
	"math"
)

type ShadeEdges struct {
	Color *color.RGBA
}

func (s ShadeEdges) Shade(object objects.VertexObject, canvas *image.RGBA) {
	for _, face := range object.Faces {
		if face.DoDraw {
			x0, y0 := face.TransformedVertices[1].At(0, 0), face.TransformedVertices[1].At(1, 0)
			x1, y1 := face.TransformedVertices[0].At(0, 0), face.TransformedVertices[0].At(1, 0)
			x2, y2 := face.TransformedVertices[2].At(0, 0), face.TransformedVertices[2].At(1, 0)

			s.drawLine(int(math.Round(x0)), int(math.Round(y0)), int(math.Round(x1)), int(math.Round(y1)), canvas)
			s.drawLine(int(math.Round(x1)), int(math.Round(y1)), int(math.Round(x2)), int(math.Round(y2)), canvas)
			s.drawLine(int(math.Round(x2)), int(math.Round(y2)), int(math.Round(x0)), int(math.Round(y0)), canvas)
		}
	}
}

func (s *ShadeEdges) drawLine(x0, y0, x1, y1 int, img *image.RGBA) {
	dx := x1 - x0
	if dx < 0 {
		dx = -dx
	}
	dy := y1 - y0
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}
	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy

	for {
		img.SetRGBA(x0, y0, *s.Color)
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}
