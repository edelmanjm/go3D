package main

import (
	"./display"
	"./objects"
	"./objects/transformations"
	"./objects/transformations/animations"
	"./shaders"
	"gonum.org/v1/gonum/mat"
	"image"
	"image/color"
	"time"
)

//const w = 1024
//const h = 512
const fov = 30

func main() {

	//defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()

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
			768, 0, 0, 0,
			0, -768, 0, 0,
			0, 0, 1, 0,
			0, 0, 0, 1,
		}),
		translation: mat.NewDense(4, 4, []float64{
			1, 0, 0, 0,
			0, 1, 0, 0,
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
		shaders: []shaders.Shader{
			shaders.ShadeFaces{Colors: []*color.RGBA{
				//{255, 0, 0, 255},
				//{0, 255, 0, 255},
				//{0, 0, 255, 255},
				{43, 142, 198, 255},
				{41, 35, 114, 255},
				{181, 23, 37, 255},
				{33, 48, 52, 255},
				{163, 183, 190, 255},
			}},
			//shaders.ShadeEdges{&color.RGBA{0, 0, 0, 255}},
		},
		nearClippingPlane: -1, farClippingPlane: -10,
		viewType: PERSPECTIVE,
	}
	myScene.updater = myScene.SceneUpdater

	obj := objects.ReadFromObj("/Users/jonathan/Desktop/objs/bmw2.obj")
	//obj.Transformations = append(obj.Transformations, (&transformations.Translation{0, -3, 0}).GetTransformation)
	obj.Transformations = append(obj.Transformations,
		(&animations.SingleAxisSpin{
			Axis: mat.NewVecDense(3, []float64{
				0, 1, 0,
			}),
			Operator: 10,
			GetClock: func() float64 {
				return float64(time.Now().UnixNano()) / 1e9
			}}).GetTransformation)
	//obj.Transformations = append(obj.Transformations, (transformations.RotationFromEulerAngles(0, math.Pi/7, 0)).GetTransformation)
	obj.Transformations = append(obj.Transformations, (&transformations.Translation{0, -2, -6}).GetTransformation)
	myScene.objects = append(myScene.objects, obj)

	//obj := objects.ReadFromObj("/Users/jonathan/Desktop/objs/fox.obj")
	//obj.Transformations = append(obj.Transformations, (&transformations.Translation{0, -2, -10}).GetTransformation)
	//obj.Transformations = append(obj.Transformations, (transformations.RotationFromAxisAngle(mat.NewVecDense(3, []float64{0, 1, 0}), math.Pi / 2)).GetTransformation)
	//obj.Transformations = append(obj.Transformations, (&transformations.Translation{10, 0, 0}).GetTransformation)
	//myScene.objects = append(myScene.objects, obj)

	//planet1 := objects.ReadFromObj("/Users/jonathan/Desktop/objs/icosphere.obj")
	//planet1.Transformations = append(planet1.Transformations, (&animations.SingleAxisSpin{
	//	mat.NewVecDense(3, []float64{0, 1, 0}),
	//	10,
	//	func() float64 {
	//		return float64(time.Now().UnixNano()) / 1e9
	//	},
	//}).GetTransformation)
	//planet1.Transformations = append(planet1.Transformations, (&transformations.Translation{0, 0, -2}).GetTransformation)
	//planet1.Transformations = append(planet1.Transformations, (&animations.SingleAxisSpin{
	//	mat.NewVecDense(3, []float64{0, 1, 0}),
	//	10,
	//	func() float64 {
	//		return float64(time.Now().UnixNano()) / 1e9
	//	},
	//}).GetTransformation)
	//myScene.objects = append(myScene.objects, planet1)

	go func() {
		var lastTime time.Time
		for {
			lastTime = time.Now()
			w, h := myDisplay.GetDimensions()
			frame := genBlankCanvas(w, h)
			zBuffer := genBlankZBuffer(w, h)
			// TODO this probably doesn't actually need to occur every frame
			myScene.updater(float64(w), float64(h), fov)
			myScene.drawObjects(shaders.Canvas{
				Image:      frame,
				ZBuffer:    zBuffer,
				UseZBuffer: true,
			})
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

func genBlankZBuffer(w, h int) [][]float64 {
	zBuffer := make([][]float64, w)
	for c := 0; c < w; c++ {
		zBuffer[c] = make([]float64, h)
		for r := 0; r < h; r++ {
			zBuffer[c][r] = -1
		}
	}
	return zBuffer
}
