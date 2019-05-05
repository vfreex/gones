package nes

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
	"github.com/vfreex/gones/pkg/emulator/ppu"
	"image"
	"image/color"
	"math/rand"
	"time"
)

// resolution 256x240

const (
	SCREEN_WIDTH  = 256
	SCREEN_HEIGHT = 240
)

type NesDiplay struct {
	screenPixels    *[SCREEN_HEIGHT][SCREEN_WIDTH]ppu.RBGColor
	app             fyne.App
	mainWindow      fyne.Window
	raster          *canvas.Raster
	canvasObj       fyne.CanvasObject
	NextCh          chan int
	StepInstruction bool
	StepFrame       bool
}

var rnd = rand.New(rand.NewSource(time.Now().Unix()))
var temp = int(0)

func NewDisplay(screenPixels *[SCREEN_HEIGHT][SCREEN_WIDTH]ppu.RBGColor) *NesDiplay {
	app := app.New()
	mainWindow := app.NewWindow("GoNES")
	display := &NesDiplay{
		app:          app,
		mainWindow:   mainWindow,
		screenPixels: screenPixels,
		NextCh:       make(chan int, 1),
		StepInstruction: true,
	}
	gameCanvas := display.render()
	mainWindow.SetContent(widget.NewVBox(gameCanvas,
		widget.NewButton(">", func() {
			display.StepInstruction = true
			display.StepFrame = false
			display.NextCh <- 1
		}),
		widget.NewButton(">>", func() {
			display.StepInstruction = false
			display.StepFrame = true
			display.NextCh <- 1
		}),
		widget.NewButton("||", func() {
			display.StepInstruction = true
			display.StepFrame = false
		}),
		widget.NewButton("->", func() {
			display.StepInstruction = false
			display.StepFrame = false
			display.NextCh <- 1
		}),
	))
	//mainWindow.SetFixedSize(true)
	return display
}

func (p *NesDiplay) render() fyne.CanvasObject {
	//p.update()
	p.raster = canvas.NewRaster(func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				//img.Set(x,y, color.RGBA{byte(rnd.Int()), byte(rnd.Int()), byte(rnd.Int()), 0xff})
				pixel := p.screenPixels[y*SCREEN_HEIGHT/h][x*SCREEN_WIDTH/w]
				img.Set(x, y, color.RGBA{R: byte(pixel >> 16), G: byte(pixel >> 8), B: byte(pixel >> 0), A: 0xff})
			}
		}
		return img
	})
	p.raster.SetMinSize(fyne.NewSize(SCREEN_WIDTH, SCREEN_HEIGHT))
	//p.raster.SetMinSize(fyne.NewSize(400, 300))
	//p.canvasObj = fyne.NewContainer(p.raster)
	p.canvasObj = p.raster
	return p.canvasObj
}

func (p *NesDiplay) Show() {
	p.mainWindow.ShowAndRun()
}

func (p *NesDiplay) Refresh() {
	//temp += 0x100000
	p.mainWindow.Canvas().Refresh(p.canvasObj)
}
