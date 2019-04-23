package ppu

import "github.com/vfreex/gones/pkg/emulator/memory"

// SPR-RAM Memory Map (8bit buswidth, 0-FFh)
//  00-FF         Sprite Attributes (256 bytes, for 64 sprites / 4 bytes each)
// Sprite RAM is directly built-in in the PPU chip. SPR-RAM is not connected to CPU or PPU bus, and can be accessed via I/O Ports only.
type SprRam struct {
	data [0x100]byte
}

func (p *SprRam) Peek(addr memory.Ptr) byte {
	return p.data[addr]
}

func (p *SprRam) Poke(addr memory.Ptr, val byte) {
	p.data[addr] = val
}
