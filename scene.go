package main

import (
	"./objects"
	"./shaders"
	"gonum.org/v1/gonum/mat"
	"image"
	"math"
)

type Scene struct {
	projection, scale, translation      *mat.Dense
	objects                             []objects.VertexObject
	shaders                             []shaders.Shader
	camera                              Camera
	updater                             func(width, height, fov float64)
	nearClippingPlane, farClippingPlane float64
}

func (s Scene) addObject(o objects.VertexObject) {
	s.objects = append(s.objects, o)
}

func (s Scene) drawObjects(canvas *image.RGBA) {
	screenspace := mat.NewDense(4, 4, nil)
	fastDenseMatMul4x4_4x4(screenspace, s.translation, s.scale)
	for _, object := range s.objects {
		// FIXME z-normalization messing with faceVector?
		objectMat := mat.NewDense(4, 4, []float64{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		})
		for _, transformationFunc := range object.Transformations {
			fastDenseMatMul4x4_4x4(objectMat, objectMat, transformationFunc())
		}
		fastDenseMatMul4x4_4x4(objectMat, objectMat, s.projection)
		//objectMat := mat.DenseCopyOf(object.Transformations)
		for _, face := range object.Faces {
			for i, vertex := range face.RawVertices {
				fastDenseMatMul4x4_4x1(face.TransformedVertices[i], objectMat, vertex)
				w := face.TransformedVertices[i].At(3, 0)
				face.TransformedVertices[i].Set(0, 0, face.TransformedVertices[i].At(0, 0)/w)
				face.TransformedVertices[i].Set(1, 0, face.TransformedVertices[i].At(1, 0)/w)
				face.TransformedVertices[i].Set(2, 0, face.TransformedVertices[i].At(2, 0)/w)
			}
			if face.Visibility != objects.BOTH {
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
					// Perspective
					(face.TransformedVertices[0].At(0, 0) + face.TransformedVertices[1].At(0, 0) + face.TransformedVertices[2].At(0, 0)) / 3,
					(face.TransformedVertices[0].At(1, 0) + face.TransformedVertices[1].At(1, 0) + face.TransformedVertices[2].At(1, 0)) / 3,
					(face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3,
					// Orthographic
					//0, 0,
					//-(face.TransformedVertices[0].At(2, 0) + face.TransformedVertices[1].At(2, 0) + face.TransformedVertices[2].At(2, 0)) / 3,
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

			if face.DoDraw {
				for _, vertex := range face.TransformedVertices {
					fastDenseMatMul4x4_4x1(vertex, screenspace, vertex)
					w := vertex.At(3, 0)
					vertex.Set(0, 0, vertex.At(0, 0)/w)
					vertex.Set(1, 0, vertex.At(1, 0)/w)
					vertex.Set(2, 0, vertex.At(2, 0)/w)
				}
			}
		}
		// TODO don't draw if not within the clipping plane bounds
		for _, shader := range s.shaders {
			shader.Shade(object, canvas)
		}
	}
}

func (s Scene) PerspectiveTransformUpdate(width, height, fov float64) {
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

func (s Scene) DisplayUpdate(width, height float64) {
	// TODO set scale
	s.translation.Set(0, 3, width/2)
	s.translation.Set(1, 3, height/2)
}

func (s Scene) SceneUpdater(width, height, fov float64) {
	s.PerspectiveTransformUpdate(width, height, fov)
	s.DisplayUpdate(width, height)
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
