package display

import (
	"github.com/fyne-io/fyne"
	"github.com/fyne-io/fyne/canvas"
	"github.com/fyne-io/fyne/desktop"
	"image"
	"image/color"
	"strconv"
)

type DisplayWindow struct {
	window           fyne.Window
	canvas           fyne.CanvasObject
	latestImage      *image.RGBA
	frametimeDisplay *canvas.Text
}

func (d *DisplayWindow) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	d.canvas.Resize(size)
}

func (d *DisplayWindow) MinSize(objects []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(768, 768)
}

func (d *DisplayWindow) Refresh() {
	d.window.Canvas().Refresh(d.canvas)
}

func (d *DisplayWindow) RenderImage(px, py, w, h int) color.Color {
	return d.latestImage.At(px, py)
}

func (d *DisplayWindow) Start() {
	d.window.ShowAndRun()
}

func (d *DisplayWindow) GetDimensions() (x, y int) {
	return d.canvas.CurrentSize().Width, d.canvas.CurrentSize().Height
}

func (d *DisplayWindow) SetFrametime(frametime float64) {
	d.frametimeDisplay.Text = strconv.FormatFloat(frametime, 'f', 0, 64) + "ms"
}

// TODO any way to specify this in the interface?
func CreateDisplay(imageChannel chan *image.RGBA) *DisplayWindow {
	window := desktop.NewApp().NewWindow("Display")
	display := &DisplayWindow{window: window}
	display.canvas = canvas.NewRaster(display.RenderImage)
	go func() {
		for {
			display.latestImage = <-imageChannel
		}
	}()

	i := image.NewRGBA(image.Rect(0, 0, 500, 500))
	for x := 0; x < i.Bounds().Dx(); x++ {
		for y := 0; y < i.Bounds().Dy(); y++ {
			i.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
		}
	}
	display.latestImage = i

	display.frametimeDisplay = canvas.NewText("asdf", color.RGBA{255, 0, 0, 255})
	display.frametimeDisplay.Move(fyne.NewPos(24, 24))
	display.window.SetContent(display.frametimeDisplay)

	window.SetContent(fyne.NewContainerWithLayout(display, display.canvas, display.frametimeDisplay))

	return display
}
