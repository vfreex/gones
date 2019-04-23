package ppu

import "github.com/vfreex/gones/pkg/emulator/memory"

// http://wiki.nesdev.com/w/index.php/PPU_palettes

// The palette for the background runs from VRAM $3F00 to $3F0F;
// the palette for the sprites runs from $3F10 to $3F1F.
// Each color takes up one byte.

type Color byte

const (
	ChrominanceMask Color = 0x0f
	LuminanceMask   Color = 0x30
)

type Palette struct {
	entries [0x20]Color
}

func getEntryIndex(addr memory.Ptr) memory.Ptr {
	index := addr & 0x001f
	switch index {
	case 0x10, 0x14, 0x18, 0x1c:
		index -= 0x10
	}
	return index
}

func (p *Palette) GetColor(addr memory.Ptr) Color {
	return p.entries[getEntryIndex(addr)]
}

func (p *Palette) SetColor(addr memory.Ptr, val Color) {
	p.entries[getEntryIndex(addr)] = val
}

func (p *Palette) Peek(addr memory.Ptr) byte {
	return byte(p.GetColor(addr))
}

func (p *Palette) Poke(addr memory.Ptr, val byte) {
	p.SetColor(addr, Color(val))
}

type RBGColor int32

var RGBMap = [0x40]RBGColor{
	0x00: 0x757575, 0x271b8f, 0x0000ab, 0x47009f, 0x8f0077, 0xab0013, 0xa70000, 0x7f0b00,
	0x08: 0x432f00, 0x004700, 0x005100, 0x003f17, 0x1b3f5f, 0x000000, 0x000000, 0x000000,
	0x10: 0xbcbcbc, 0x0073ef, 0x2b3bef, 0x8300f3, 0xbf00bf, 0xe7005b, 0xdb2b00, 0xcb4f0f,
	0x18: 0x8b7300, 0x009700, 0x00ab00, 0x00933b, 0x00838b, 0x000000, 0x000000, 0x000000,

	0x20: 0xffffff, 0x3fbfff, 0x5f97ff, 0xa78bfd, 0xf77bff, 0xff77b7, 0xff7763, 0xff9b3b,
	0x28: 0xf3bf3f, 0x83d313, 0x4fdf4b, 0x58f898, 0x00ebdb, 0x000000, 0x000000, 0x000000,
	0x30: 0xffffff, 0xabe7ff, 0xc7d7ff, 0xd7cbff, 0xffc7ff, 0xffc7db, 0xffbfb3, 0xffdbab,
	0x38: 0xffe7a3, 0xe3ffa3, 0xabf3bf, 0xb3ffcf, 0x9ffff3, 0x000000, 0x000000, 0x000000,
}

func (color Color) ToGRBColor() RBGColor {
	return RGBMap[color&0x3f]
}
