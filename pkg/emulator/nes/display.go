package nes

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/vfreex/gones/pkg/emulator/joypad"
	"github.com/vfreex/gones/pkg/emulator/ppu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
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
	canvasObj       *canvas.Image
	NextCh          chan int
	StepInstruction bool
	StepFrame       bool
	RequestReset    bool
	PressedKeys     byte
	ReleasedKeys    byte
	Keys            byte
	img             *image.NRGBA
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
		container.NewVBox(gameCanvas,
			container.NewHBox(
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
					display.RequestReset = true
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
		case desktop.KeyControlLeft:
			fallthrough
		case desktop.KeyControlRight:
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
		case desktop.KeyControlLeft:
			fallthrough
		case desktop.KeyControlRight:
			display.Keys &= ^ joypad.Button_Select
		}
	})
	mainWindow.SetFixedSize(true)
	return display
}

func (p *NesDiplay) render() fyne.CanvasObject {
	p.img = image.NewNRGBA(image.Rect(0, 0, SCREEN_WIDTH, SCREEN_HEIGHT))
	p.canvasObj = canvas.NewImageFromImage(p.img)
	p.canvasObj.ScaleMode = canvas.ImageScalePixels
	p.canvasObj.SetMinSize(fyne.NewSize(SCREEN_WIDTH*2, SCREEN_HEIGHT*2))
	return p.canvasObj
}

func (p *NesDiplay) Show() {
	p.mainWindow.ShowAndRun()
}

func (p *NesDiplay) Refresh() {
	//temp += 0x100000
	for y := 0; y < SCREEN_HEIGHT; y++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			pixel := p.screenPixels[y][x]
			p.img.SetNRGBA(x, y, color.NRGBA{R: byte(pixel >> 16), G: byte(pixel >> 8), B: byte(pixel >> 0), A: 0xff})
		}
	}

	p.canvasObj.Refresh()
}
