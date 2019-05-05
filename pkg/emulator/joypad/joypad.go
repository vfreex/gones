package joypad

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

const (
	Button_A byte = 1 << iota
	Button_B
	Button_Select
	Button_Start
	Button_Up
	Button_Down
	Button_Left
	Button_Right
)

const (
	Joypad_1 = 0x4016
	Joypad_2 = 0x4017
)

type Joypad struct {
	Buttons byte
	Shift   byte
}

//func NewJoypad(port memory.Ptr) *Joypad {
//	p := &Joypad{Port: port}
//	return p
//}

type Joypads struct {
	Joypads [2]Joypad
	Reset   bool
}

func NewJoypads() *Joypads {
	return &Joypads{}
}

func (p *Joypads) getJoypad(addr memory.Ptr) *Joypad {
	var joypad *Joypad
	switch addr {
	case Joypad_1:
		joypad = &p.Joypads[0]
	case Joypad_2:
		joypad = &p.Joypads[1]
	default:
		panic(fmt.Errorf("invalid Joypad port address: %02x", addr))
	}
	return joypad
}
func (p *Joypads) Peek(addr memory.Ptr) byte {
	var r byte
	joypad := p.getJoypad(addr)
	if joypad.Shift < 8 {
		r = (joypad.Buttons >> joypad.Shift) & 1
		if !p.Reset {
			joypad.Shift++
		}
	} else {
		r = 1
	}
	return r
}

func (p *Joypads) Poke(addr memory.Ptr, val byte) {
	if val&1 == 1 {
		p.Reset = true
		p.Joypads[0].Shift = 0
		p.Joypads[1].Shift = 0
	} else {
		p.Reset = false
	}
}
