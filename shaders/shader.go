package shaders

import (
	"../objects"
	"image"
)

type Shader interface {
	Shade(object objects.VertexObject, canvas *image.RGBA)
}
