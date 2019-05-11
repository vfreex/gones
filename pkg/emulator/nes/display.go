package nes

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/driver/desktop"
	"fyne.io/fyne/widget"
	"github.com/vfreex/gones/pkg/emulator/joypad"
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
	PressedKeys     byte
	ReleasedKeys    byte
	Keys            byte
	img             *image.RGBA
}

var rnd = rand.New(rand.NewSource(time.Now().Unix()))
var temp = int(0)

func NewDisplay(screenPixels *[SCREEN_HEIGHT][SCREEN_WIDTH]ppu.RBGColor) *NesDiplay {
	app := app.New()
	mainWindow := app.NewWindow("GoNES")
	display := &NesDiplay{
		app:             app,
		mainWindow:      mainWindow,
		screenPixels:    screenPixels,
		NextCh:          make(chan int, 1),
		StepInstruction: false,
	}
	gameCanvas := display.render()
	mainWindow.SetContent(
		widget.NewVBox(gameCanvas,
			widget.NewHBox(
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
				widget.NewButton("RESET", func() {
					display.StepInstruction = true
					display.StepFrame = false
					display.NextCh <- 0xff
				}),
			),
		))
	mainWindow.Canvas().(desktop.Canvas).SetOnKeyDown(func(event *fyne.KeyEvent) {
		switch event.Name {
		case fyne.KeyReturn:
			display.Keys |= joypad.Button_Start
		case fyne.KeyA:
			fallthrough
		case fyne.KeyLeft:
			display.Keys |= joypad.Button_Left
		case fyne.KeyW:
			fallthrough
		case fyne.KeyUp:
			display.Keys |= joypad.Button_Up
		case fyne.KeyD:
			fallthrough
		case fyne.KeyRight:
			display.Keys |= joypad.Button_Right
		case fyne.KeyS:
			fallthrough
		case fyne.KeyDown:
			display.Keys |= joypad.Button_Down
		case fyne.KeyZ:
			display.Keys |= joypad.Button_B
		case fyne.KeyX:
			display.Keys |= joypad.Button_A
		case "LeftControl":
			display.Keys |= joypad.Button_Select
		}
	})
	mainWindow.Canvas().(desktop.Canvas).SetOnKeyUp(func(event *fyne.KeyEvent) {
		switch event.Name {
		case fyne.KeyReturn:
			display.Keys &= ^joypad.Button_Start
		case fyne.KeyA:
			fallthrough
		case fyne.KeyLeft:
			display.Keys &= ^ joypad.Button_Left
		case fyne.KeyW:
			fallthrough
		case fyne.KeyUp:
			display.Keys &= ^joypad.Button_Up
		case fyne.KeyD:
			fallthrough
		case fyne.KeyRight:
			display.Keys &= ^ joypad.Button_Right
		case fyne.KeyS:
			fallthrough
		case fyne.KeyDown:
			display.Keys &= ^joypad.Button_Down
		case fyne.KeyZ:
			display.Keys &= ^ joypad.Button_B
		case fyne.KeyX:
			display.Keys &= ^ joypad.Button_A
		case "LeftControl":
			display.Keys &= ^ joypad.Button_Select
		}
	})
	//mainWindow.SetFixedSize(true)
	return display
}

func (p *NesDiplay) render() fyne.CanvasObject {
	//p.update()
	lastW, lastH := 0, 0
	p.raster = canvas.NewRaster(func(w, h int) image.Image {
		if p.img == nil || w != lastW || h != lastH {
			p.img = image.NewRGBA(image.Rect(0, 0, w, h))
			lastW = w
			lastH = h
		}
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				pixel := p.screenPixels[y*SCREEN_HEIGHT/h][x*SCREEN_WIDTH/w]
				p.img.SetRGBA(x, y, color.RGBA{R: byte(pixel >> 16), G: byte(pixel >> 8), B: byte(pixel >> 0), A: 0xff})
			}
		}
		return p.img
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
