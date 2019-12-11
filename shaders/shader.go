package shaders

import (
	"../objects"
)

type Shader interface {
	Shade(object objects.VertexObject, canvas Canvas)
}
