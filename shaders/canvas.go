package shaders

import (
	"image"
	"image/color"
)

type Canvas struct {
	Image      *image.RGBA
	ZBuffer    [][]float64
	UseZBuffer bool
}

func (c *Canvas) SetRGBA(x int, y int, z float64, color color.RGBA) bool {
	if 0 < x && x < len(c.ZBuffer) && 0 < y && y < len(c.ZBuffer[0]) {
		// FIXME because line drawing issues
		if y > len(c.ZBuffer[0]) {
			return true
		}
		if !c.UseZBuffer || z > c.ZBuffer[x][y] {
			c.ZBuffer[x][y] = z
			c.Image.SetRGBA(x, y, color)
		}
	}
	return false
}
