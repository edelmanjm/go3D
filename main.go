package main

import (
	"./display"
	"./objects"
	"./shaders"
	"github.com/pkg/profile"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/color"
	"time"
)

//const w = 1024
//const h = 512
const fov = 30

func main() {

	defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

	imageChannel := make(chan *image.RGBA, 60)
	myDisplay := display.CreateDisplay(imageChannel)

	myScene := Scene{
		//projection: mat.NewDense(4, 4, []float64{
		//	1, 0, 0, 0,
		//	0, 1, 0, 0,
		//	0, 0, 1, 0,
		//	0, 0, 0, 1,
		//}),
		projection: mat.NewDense(4, 4, []float64{
			1, 0, 0, 0,
			0, 1, 0, 0,
			0, 0, 1, 0,
			0, 0, -1, 0,
		}),
		scale: mat.NewDense(4, 4, []float64{
			50, 0, 0, 0,
			0, 50, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}),
		camera: Camera{
			mat.NewVecDense(3, []float64{
				0,
				0,
				-1,
			}),
		},
		shaders:           []shaders.Shader{shaders.ShadeEdges{&color.RGBA{0, 0, 0, 0}}},
		nearClippingPlane: -1, farClippingPlane: -10,
	}
	myScene.projectionUpdater = myScene.PerspectiveTransform
	obj := objects.ReadFromObj("/Users/jonathan/Desktop/objs/cube.obj")
	//obj.Transformations = append(obj.Transformations, mat.NewDense(4, 4, []float64{
	//	1, 0, 0, 2,
	//	0, 1, 0, 2,
	//	0, 0, 1, 0,
	//	0, 0, 0, 1,
	//}))
	obj.Transformations = mat.NewDense(4, 4, []float64{
		1, 0, 0, 5,
		0, 1, 0, 5,
		0, 0, 1, 0,
		0, 0, 0, 1,
	})
	myScene.objects = append(myScene.objects, obj)

	go func() {
		var lastTime time.Time
		for {
			lastTime = time.Now()
			x, y := myDisplay.GetDimensions()
			frame := genBlankCanvas(x, y)
			// TODO this probably doesn't actually need to occur every frame
			myScene.projectionUpdater(float64(x), float64(y), fov)
			myScene.drawObjects(frame)
			imageChannel <- frame
			myDisplay.Refresh()
			myDisplay.SetFrametime(float64(time.Since(lastTime).Nanoseconds() / 1e6))
			//fmt.Printf("%d\n", time.Since(lastTime).Nanoseconds() / 1e6)
		}
	}()
	myDisplay.Start()
}

func genBlankCanvas(w, h int) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := range canvas.Pix {
		canvas.Pix[i] = 255
	}
	return canvas
}
