package ppu

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

// http://wiki.nesdev.com/w/index.php/PPU_registers

// PPU control ($2000) register - write only
type PPUCtrl byte

const (
	PPUCtrl_NameTableLow           PPUCtrl = 1 << iota                            // See PPUCtrl_NameTable
	PPUCtrl_NameTableHigh                                                         // See PPUCtrl_NameTable
	PPUCtrl_PPUDataIncrement                                                      // (0=Increment by 1, 1=Increment by 32)
	PPUCtrl_SpritePatternTable                                                    // (0: $0000; 1: $1000; ignored in 8x16 mode)
	PPUCtrl_BackgroundPatternTable                                                // (0: $0000; 1: $1000)
	PPUCtrl_SpriteSize                                                            // (0: 8x8 pixels; 1: 8x16 pixels)
	PPUCtrl_PPUMasterSlave                                                        // (Not used in NES, 0: read backdrop from EXT pins; 1: output color on EXT pins)
	PPUCtrl_NMIOnVBlank                                                           // (0=Disabled, 1=Enabled)
	PPUCtrl_NameTable              = PPUCtrl_NameTableLow | PPUCtrl_NameTableHigh // (0-3=VRAM 2000h,2400h,2800h,2C00h)
)

// PPU mask ($2001) register
type PPUMask byte

const (
	PPUMask_Greyscale            PPUMask = 1 << iota // (0: normal color, 1: produce a greyscale display)
	PPUMask_NoBackgroundClipping                     // 1: Show background in leftmost 8 pixels of screen, 0: Hide
	PPUMask_NoSpriteClipping                         // 1: Show sprites in leftmost 8 pixels of screen, 0: Hide
	PPUMask_BackgroundVisibility                     // 1: Show background
	PPUMask_SpriteVisibility                         // 1: Show sprites
	PPUMask_ColorEmphasizeRed
	PPUMask_ColorEmphasizeGreen
	PPUMask_ColorEmphasizeBlue
)

// PPU status ($2002) register
type PPUStatus byte

const (
	PPUStatus_NotUsed0 PPUStatus = 1 << iota
	PPUStatus_NotUsed1
	PPUStatus_NotUsed2
	PPUStatus_NotUsed3
	PPUStatus_NotUsed4
	PPUStatus_SpriteOverflow
	PPUStatus_Sprite0Hit
	PPUStatus_VBlank
)

type PPURegister = byte
type RegisterACL byte

const (
	PPUCTRL   = 0x2000 // W
	PPUMASK   = 0x2001 // W
	PPUSTATUS = 0x2002 // R
	OAMADDR   = 0x2003 // W
	OAMDATA   = 0x2004 // RW
	PPUSCROLL = 0x2005 // W
	PPUADDR   = 0x2006 // W
	PPUDATA   = 0x2007 // RW
	OAMDMA    = 0x4014 // W
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
	ppu *PPUImpl
	//registers  map[memory.Ptr]*Register
	latch      bool
	latchCache byte
	vramAddr   memory.Ptr
	hscroll    byte
	vscroll    byte
	ctrl       PPUCtrl
	mask       PPUMask
	status     PPUStatus
	oamAddr    byte
}

func NewPPURegisters(ppu *PPUImpl) Registers {
	registers := Registers{
		ppu: ppu,
		//registers: map[memory.Ptr]*Register{
		//	PPUCTRL:   {0, ACL_Write},
		//	PPUMASK:   {0, ACL_Write},
		//	PPUSTATUS: {0, ACL_Read},
		//	OAMADDR:   {0, ACL_Write},
		//	OAMDATA:   {0, ACL_Read | ACL_Write},
		//	PPUSCROLL: {0, ACL_Write},
		//	PPUADDR:   {0, ACL_Write},
		//	PPUDATA:   {0, ACL_Read | ACL_Write},
		//	OAMDMA:    {0, ACL_Write},
		//},
	}
	return registers
}

func (p *Registers) Peek(addr memory.Ptr) byte {
	//r, ok := p.registers[addr]
	//if !ok {
	//	panic(fmt.Errorf("no PPU register at address %04x", addr))
	//}
	//if r.acl&ACL_Read == 0 {
	//	panic(fmt.Errorf("no read access to PPU register %04x", addr))
	//}
	//return r.value
	var r byte
	switch addr {
	case PPUSTATUS:
		r = byte(p.status)
		p.status &= ^PPUStatus_VBlank
		p.latch = false
	case OAMDATA:
		// The address is NOT auto-incremented after <reading> from 2004h.
		return p.ppu.SprRam.data[p.oamAddr]
	case PPUDATA:
		// Reading from VRAM 0000h-3EFFh loads the desired value into a latch,
		// and returns the OLD content of the latch to the CPU
		if p.vramAddr < 0x3f00 {
			r = p.latchCache
			p.latchCache = p.ppu.vram.Peek(p.vramAddr)
		} else {
			// reading from Palette memory VRAM 3F00h-3FFFh does directly access the desired address.
			r = p.ppu.vram.Peek(p.vramAddr)
		}
		// The PPU will auto-increment the VRAM address (selected via Port 2006h)
		// after each read/write from/to Port 2007h by 1 or 32 (depending on Bit2 of $2000).
		if p.ctrl&PPUCtrl_PPUDataIncrement != 0 {
			p.vramAddr += 32
		} else {
			p.vramAddr++
		}
	default:
		panic(fmt.Errorf("PPU register %04x is not readable", addr))
	}
	return r
}

func (p *Registers) Poke(addr memory.Ptr, val byte) {
	switch addr {
	case PPUCTRL:
		p.ctrl = PPUCtrl(val)
	case PPUMASK:
		p.mask = PPUMask(val)
	case OAMADDR:
		p.oamAddr = val
	case OAMDATA:
		// The Port 2003h address is auto-incremented by 1 after each <write> to 2004h.
		p.ppu.SprRam.data[p.oamAddr] = val
		p.oamAddr++
	case PPUSCROLL:
		if !p.latch {
			p.hscroll = val
		} else {
			p.vscroll = val
		}
		p.latch = !p.latch
	case PPUADDR:
		// After reading PPUSTATUS to reset the address latch,
		// write the 16-bit address of VRAM you want to access here
		if !p.latch { // first write, upper byte (6bit)
			p.vramAddr = (memory.Ptr(val) & 0x3F) << 8
		} else { // second write, lower byte (8bit)
			p.vramAddr |= memory.Ptr(val)
		}
		p.latch = !p.latch
	case PPUDATA:
		p.ppu.vram.Poke(p.vramAddr, val)
		// The PPU will auto-increment the VRAM address (selected via Port 2006h)
		// after each read/write from/to Port 2007h by 1 or 32 (depending on Bit2 of $2000).
		if p.ctrl&PPUCtrl_PPUDataIncrement != 0 {
			p.vramAddr += 32
		} else {
			p.vramAddr++
		}
	case OAMDMA:
		p.onOAMDMAWrite(val)
	default:
		panic(fmt.Errorf("PPU register %04x is not writable", addr))
	}
}

func (p *Registers) onOAMDMAWrite(val byte) {
	// Transfers 256 bytes from CPU Memory area into SPR-RAM. The transfer takes 512 CPU clock cycles, two cycles per byte, the transfer starts about immediately after writing to 4014h: The CPU either fetches the first byte of the next instruction, and then begins DMA, or fetches and executes the next instruction, and then begins DMA. The CPU is halted during transfer.
	// Bit7-0  Upper 8bit of source address (Source=N*100h) (Lower bits are zero)
	// Data is written to Port 2004h. The destination address in SPR-RAM is thus [2003h], which should be normally initialized to zero - unless one wants to "rotate" the target area, which may be useful when implementing more than eight (flickering) sprites per scanline.
	//srcAddr := memory.Ptr(p.registers[OAMDMA].value) << 8
	//destAddr := memory.Ptr(p.registers[OAMADDR].value)
	//p.ppu.sprRam.data[destAddr] = ?
	src := memory.Ptr(val) << 8
	for i := memory.PtrDist(0); i < 256; i++ {
		b := p.ppu.cpu.Memory.Peek(src + i)
		p.Poke(OAMDATA, b)
	}
}
