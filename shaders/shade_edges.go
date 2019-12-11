package shaders

import (
	"../objects"
	"image/color"
)

type ShadeEdges struct {
	Color *color.RGBA
}

func (s ShadeEdges) Shade(object objects.VertexObject, canvas Canvas) {
	for _, face := range object.Faces {
		if face.DoDraw {
			x0, y0 := face.TransformedVertices[1].At(0, 0), face.TransformedVertices[1].At(1, 0)
			x1, y1 := face.TransformedVertices[0].At(0, 0), face.TransformedVertices[0].At(1, 0)
			x2, y2 := face.TransformedVertices[2].At(0, 0), face.TransformedVertices[2].At(1, 0)

			//s.drawLine(int(math.Round(x0)), int(math.Round(y0)), int(math.Round(x1)), int(math.Round(y1)), canvas)
			//s.drawLine(int(math.Round(x1)), int(math.Round(y1)), int(math.Round(x2)), int(math.Round(y2)), canvas)
			//s.drawLine(int(math.Round(x2)), int(math.Round(y2)), int(math.Round(x0)), int(math.Round(y0)), canvas)

			// Use average z-depth
			z := (face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3
			s.drawLine(int(x0), int(y0), int(x1), int(y1), z, canvas)
			s.drawLine(int(x1), int(y1), int(x2), int(y2), z, canvas)
			s.drawLine(int(x2), int(y2), int(x0), int(y0), z, canvas)
		}
	}
}

func (s *ShadeEdges) drawLine(x0, y0, x1, y1 int, z float64, canvas Canvas) {
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
		canvas.SetRGBA(x0, y0, z, *s.Color)
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
