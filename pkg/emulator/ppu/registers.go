package ppu

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

type PPURegister = byte
type RegisterACL byte

const (
	PPUCTRL   = 0x2000
	PPUMASK   = 0x2001
	PPUSTATUS = 0x2002
	OAMADDR   = 0x2003
	OAMDATA   = 0x2004
	PPUSCROLL = 0x2005
	PPUADDR   = 0x2006
	PPUDATA   = 0x2007
	OAMDMA    = 0x4014
)

const (
	ACL_Read RegisterACL = 1 << iota
	ACL_Write
)

type Register struct {
	value byte
	acl   RegisterACL
}

type Registers struct {
	ppu       *PPUImpl
	registers map[memory.Ptr]*Register
}

func NewPPURegisters(ppu *PPUImpl) Registers {
	registers := Registers{
		ppu: ppu,
		registers: map[memory.Ptr]*Register{
			PPUCTRL:   {0, ACL_Write},
			PPUMASK:   {0, ACL_Write},
			PPUSTATUS: {0, ACL_Read},
			OAMADDR:   {0, ACL_Write},
			OAMDATA:   {0, ACL_Read | ACL_Write},
			PPUSCROLL: {0, ACL_Write},
			PPUADDR:   {0, ACL_Write},
			PPUDATA:   {0, ACL_Read | ACL_Write},
			OAMDMA:    {0, ACL_Write},
		},
	}
	return registers
}

func (p Registers) Peek(addr memory.Ptr) byte {
	r, ok := p.registers[addr]
	if !ok {
		panic(fmt.Errorf("no PPU register at address %04x", addr))
	}
	if r.acl&ACL_Read == 0 {
		panic(fmt.Errorf("no read access to PPU register %04x", addr))
	}
	return r.value
}

func (p Registers) Poke(addr memory.Ptr, val byte) {
	r, ok := p.registers[addr]
	if !ok {
		panic(fmt.Errorf("no PPU register at address %04x", addr))
	}
	if r.acl&ACL_Write == 0 {
		panic(fmt.Errorf("no write access to PPU register %04x", addr))
	}
	r.value = val
	switch addr {
	case OAMDMA:
		p.onOAMDMAWrite()
	}
}

func (p Registers) onOAMDMAWrite() {
	// Transfers 256 bytes from CPU Memory area into SPR-RAM. The transfer takes 512 CPU clock cycles, two cycles per byte, the transfer starts about immediately after writing to 4014h: The CPU either fetches the first byte of the next instruction, and then begins DMA, or fetches and executes the next instruction, and then begins DMA. The CPU is halted during transfer.
	// Bit7-0  Upper 8bit of source address (Source=N*100h) (Lower bits are zero)
	// Data is written to Port 2004h. The destination address in SPR-RAM is thus [2003h], which should be normally initialized to zero - unless one wants to "rotate" the target area, which may be useful when implementing more than eight (flickering) sprites per scanline.
	//srcAddr := memory.Ptr(p.registers[OAMDMA].value) << 8
	//destAddr := memory.Ptr(p.registers[OAMADDR].value)
	//p.ppu.sprRam.data[destAddr] = ?
	panic("not implemented")
}
