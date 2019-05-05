package ppu

import (
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/cpu"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

const (
	SCREEN_WIDTH  = 256
	SCREEN_HEIGHT = 240
)

type PPUImpl struct {
	cpu            *cpu.Cpu
	registers      Registers
	SprRam         SprRam
	vram           memory.AddressSpace
	Palette        Palette
	//cycles         int64
	//Frames         int64
	RenderedBuffer [SCREEN_HEIGHT][SCREEN_WIDTH]RBGColor
	scanline       int
	dotInScanline  int
	frame          int
}

var logger = logger2.GetLogger()

func NewPPU(vram memory.AddressSpace, cpu *cpu.Cpu) *PPUImpl {
	ppu := &PPUImpl{
		vram: vram,
		cpu:  cpu,
	}
	ppu.registers = NewPPURegisters(ppu)
	return ppu
}

func (ppu *PPUImpl) MapToCPUAddressSpace(as memory.AddressSpace) {
	as.AddMapping(0x2000, 0x2000,
		memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, &ppu.registers, nil)
}
