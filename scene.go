package main

import (
	"./objects"
	"./shaders"
	"gonum.org/v1/gonum/mat"
	"image"
	"math"
)

type Scene struct {
	projection, scale                   *mat.Dense
	objects                             []objects.VertexObject
	shaders                             []shaders.Shader
	camera                              Camera
	projectionUpdater                   func(width, height, fov float64)
	nearClippingPlane, farClippingPlane float64
}

func (s Scene) addObject(o objects.VertexObject) {
	s.objects = append(s.objects, o)
}

func (s Scene) drawObjects(canvas *image.RGBA) {
	for _, object := range s.objects {
		// FIXME wait to apply perspective until after culling
		objectMat := mat.NewDense(4, 4, nil)
		fastDenseMatMul4x4_4x4(objectMat, object.Transformations, s.projection)
		for _, face := range object.Faces {
			for i, vertex := range face.RawVertices {
				//face.TransformedVertices[i].Mul(objectMat, vertex)
				fastDenseMatMul4x4_4x1(face.TransformedVertices[i], objectMat, vertex)
				w := face.TransformedVertices[i].At(3, 0)
				face.TransformedVertices[i].Set(0, 0, face.TransformedVertices[i].At(0, 0)/w)
				face.TransformedVertices[i].Set(1, 0, face.TransformedVertices[i].At(1, 0)/w)
				face.TransformedVertices[i].Set(2, 0, face.TransformedVertices[i].At(2, 0)/w)
				// TODO minimize scale computation overhead
				fastDenseMatMul4x4_4x1(face.TransformedVertices[i], s.scale, face.TransformedVertices[i])
			}
			// Check culling; if it should be culled, set the first boi in transformedveriticies to nil
			if face.Visibility != objects.BOTH {
				//face.FaceNormal.MulVec(object.Transformations.Slice(0, 3, 0, 3), face.FaceNormal)

				u := mat.NewDense(4, 1, nil)
				u.Sub(face.TransformedVertices[1], face.TransformedVertices[0])
				v := mat.NewDense(4, 1, nil)
				v.Sub(face.TransformedVertices[2], face.TransformedVertices[1])
				face.FaceNormal = mat.NewVecDense(3, []float64{
					u.At(1, 0)*v.At(2, 0) - u.At(2, 0)*v.At(1, 0),
					u.At(2, 0)*v.At(0, 0) - u.At(0, 0)*v.At(2, 0),
					u.At(0, 0)*v.At(1, 0) - u.At(1, 0)*v.At(0, 0),
				})

				faceVector := mat.NewVecDense(3, []float64{
					(face.TransformedVertices[0].At(0, 0) + face.TransformedVertices[1].At(0, 0) + face.TransformedVertices[2].At(0, 0)) / 3,
					(face.TransformedVertices[0].At(1, 0) + face.TransformedVertices[1].At(1, 0) + face.TransformedVertices[2].At(1, 0)) / 3,
					(face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3,
				})
				if mat.Dot(face.FaceNormal, faceVector) < 0 {
					if face.Visibility == objects.FRONT {
						face.DoDraw = false
					} else {
						face.DoDraw = true
					}
				} else {
					if face.Visibility == objects.BACK {
						face.DoDraw = false
					} else {
						face.DoDraw = true
					}
				}
			} else {
				face.DoDraw = true
			}
		}
		// TODO don't draw if not within the clipping plane bounds
		for _, shader := range s.shaders {
			shader.Shade(object, canvas)
		}
	}
}

func (s Scene) PerspectiveTransform(width, height, fov float64) {
	// TODO figure out where the clipping plane really should be
	// Aspect ratio seems to break things for some reason, but everything works fine without it
	//aspectRatio := width / height
	top := math.Tan((fov/360*2*math.Pi)/2) * s.nearClippingPlane
	bottom := -top
	right := top
	left := -top

	s.projection.Set(0, 0, (2*s.nearClippingPlane)/(right-left))
	s.projection.Set(0, 2, (right+left)/(right-left))
	s.projection.Set(1, 1, (2*s.nearClippingPlane)/(top-bottom))
	s.projection.Set(1, 2, (top+bottom)/(top-bottom))
	s.projection.Set(2, 2, -(s.farClippingPlane+s.nearClippingPlane)/(s.farClippingPlane-s.nearClippingPlane))
	s.projection.Set(2, 3, -(2*s.farClippingPlane*s.nearClippingPlane)/(s.farClippingPlane-s.nearClippingPlane))
}

func fastDenseMatMul4x4_4x4(receiver, a, b *mat.Dense) {
	for r := 0; r <= 3; r++ {
		for c := 0; c <= 3; c++ {
			receiver.Set(r, c,
				a.At(r, 0)*b.At(0, c)+
					a.At(r, 1)*b.At(1, c)+
					a.At(r, 2)*b.At(2, c)+
					a.At(r, 3)*b.At(3, c))
		}
	}
}

func fastDenseMatMul4x4_4x1(receiver, a, b *mat.Dense) {
	for r := 0; r <= 3; r++ {
		receiver.Set(r, 0,
			a.At(r, 0)*b.At(0, 0)+
				a.At(r, 1)*b.At(1, 0)+
				a.At(r, 2)*b.At(2, 0)+
				a.At(r, 3)*b.At(3, 0))
	}
}
