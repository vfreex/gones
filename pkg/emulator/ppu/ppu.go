package ppu

import "github.com/vfreex/gones/pkg/emulator/memory"

type PPUImpl struct {
	registers Registers
	sprRam SprRam
}

func NewPPU() *PPUImpl {
	ppu := &PPUImpl{
	}
	return ppu
}


func (ppu *PPUImpl) MapToCPUAddressSpace(as memory.AddressSpace) {
	as.AddMapping(0x2000, 0x2000,
		memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, ppu.registers, nil)
	as.AddMapping(0x4014, 1,
		memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, ppu.registers, nil)
}
