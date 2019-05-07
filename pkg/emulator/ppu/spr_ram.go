package ppu

import "github.com/vfreex/gones/pkg/emulator/memory"

// SPR-RAM Memory Map (8bit buswidth, 0-FFh)
//  00-FF         Sprite Attributes (256 bytes, for 64 sprites / 4 bytes each)
// Sprite RAM is directly built-in in the PPU chip. SPR-RAM is not connected to CPU or PPU bus, and can be accessed via I/O Ports only.

type SpriteAttrMask byte

const (
	SpriteAttr_PaletteLow SpriteAttrMask = 1 << iota
	SpriteAttr_PaletteHigh
	SpriteAttr_Unused2
	SpriteAttr_Unused3
	SpriteAttr_Unused4
	SpriteAttr_BackgroundPriority
	SpriteAttr_HorizontalFlip
	SpriteAttr_VerticalFlip
)

type SprRam struct {
	data [0x100]byte
}

func (p *SprRam) Peek(addr memory.Ptr) byte {
	return p.data[addr]
}

func (p *SprRam) Poke(addr memory.Ptr, val byte) {
	p.data[addr] = val
}

type SpriteAttr struct {
	PaletteId          int
	BackgroundPriority bool
	HorizontalFlip     bool
	VerticalFlip       bool
}

type Sprite struct {
	Id     int
	X, Y   int
	TileId int
	Attr   SpriteAttr
	TileRowLow, TileRowHigh byte
}

func (p *Sprite) Unmarshal(id int, spriteBytes *[4]byte) {
	//if len(spriteBytes)<4 {
	//	panic("4 bytes are required to unmarshal a Sprite")
	//}
	p.Id = id
	p.X = int(spriteBytes[3])
	p.Y = int(spriteBytes[0])
	p.TileId = int(spriteBytes[1])
	p.Attr.Unmarshal(spriteBytes[2])
}

func (p *SpriteAttr) Unmarshal(attr byte) {
	*p = SpriteAttr{
		PaletteId:          int(attr & byte(SpriteAttr_PaletteLow|SpriteAttr_PaletteHigh)),
		BackgroundPriority: attr&byte(SpriteAttr_BackgroundPriority) != 0,
		HorizontalFlip:     attr&byte(SpriteAttr_HorizontalFlip) != 0,
		VerticalFlip:       attr&byte(SpriteAttr_VerticalFlip) != 0,
	}
}
