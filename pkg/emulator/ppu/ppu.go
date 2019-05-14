package ppu

import (
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/cpu"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

const (
	ScanlinesPerFrame     = 262
	DotsPerScanline       = 341
	VisualDotsPerScanline = 256
	VisualScanlines       = 240
)

type NewFrameHandler func(frame *[VisualScanlines][VisualDotsPerScanline]RBGColor, frameID int)

type PPUImpl struct {
	cpu                 *cpu.Cpu
	registers           Registers
	sprRam              SprRam
	vram                memory.AddressSpace
	Palette             Palette
	RenderedBuffer      [VisualScanlines][VisualDotsPerScanline]RBGColor
	scanline            int
	dotInScanline       int
	frame               int
	spriteCount         int
	sprites             [8]Sprite
	currentSprites      [8]Sprite
	currentSpritesCount int
	NewFrameHandler     NewFrameHandler
}

var logger = logger2.GetLogger()

func NewPPU(vram memory.AddressSpace, cpu *cpu.Cpu) *PPUImpl {
	ppu := &PPUImpl{
		vram: vram,
		cpu:  cpu,
		//secondaryOAM: ram.NewRAM(32),
	}
	ppu.registers = NewPPURegisters(ppu)
	return ppu
}

func (ppu *PPUImpl) MapToCPUAddressSpace(as memory.AddressSpace) {
	as.AddMapping(0x2000, 0x2000,
		memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, &ppu.registers, func(addr memory.Ptr) memory.Ptr {
			return addr & 0x2007
		})
	as.AddMapping(0x4014, 1,
		memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, &ppu.registers, nil)
}
