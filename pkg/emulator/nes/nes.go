package nes

import (
	"github.com/vfreex/gones/pkg/emulator/cpu"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/ppu"
	"github.com/vfreex/gones/pkg/emulator/ram"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
	"log"
	"time"
)

type NES interface {
	LoadCartridge(cartridge *ines.INesRom) error
	Start() error
}

type NESImpl struct {
	ticker *time.Ticker
	cpu    *cpu.Cpu
	cpuAS  memory.AddressSpace
	ram    memory.Memory
	ppu    *ppu.PPUImpl
	ppuAS  memory.AddressSpace
	vram   memory.Memory
}

func NewNes() NES {
	nes := &NESImpl{
		cpuAS:  &memory.AddressSpaceImpl{},
		ram:   ram.NewMainRAM(),
		ppu:   ppu.NewPPU(),
		ppuAS: &memory.AddressSpaceImpl{},
		vram:  ram.NewRAM(0x800),
	}
	nes.cpu = cpu.NewCpu(nes.cpuAS)

	// setting up CPU memory map
	// 0x0000 - ox1fff RAM
	nes.cpuAS.AddMapping(0, 0x2000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		nes.ram, nil)
	// fake memory map range
	nes.cpuAS.AddMapping(0x4000, 0x14, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		ram.NewRAM(0x14), func(addr memory.Ptr) memory.Ptr {
			return addr - 0x4000
		})
	nes.ppu.MapToCPUAddressSpace(nes.cpuAS)
	// fake memory map range
	nes.cpuAS.AddMapping(0x4015, 0x3, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		ram.NewRAM(0x03), func(addr memory.Ptr) memory.Ptr {
			return addr - 0x4015
		})

	// setting up PPU memory map
	// https://wiki.nesdev.com/w/index.php/PPU_memory_map
	nes.ppuAS.AddMapping(0x2000, 0x800, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
	nes.vram, func(addr memory.Ptr) memory.Ptr {
			return addr & 0xf7ff
		})

	return nes
}

func (nes *NESImpl) LoadCartridge(cartridge *ines.INesRom) error {
	// load PRG-ROM
	nes.cpuAS.AddMapping(0x8000, 0x8000, memory.MMAP_MODE_READ,
		cartridge.Prg, nil)

	// load CHR-ROM
	nes.ppuAS.AddMapping(0, 0x2000, memory.MMAP_MODE_READ,
		cartridge.Chr, nil)
	return nil
}

func (nes *NESImpl) Start() error {
	nes.cpuAS.Map()
	nes.ppuAS.Map()
	nes.ticker = time.NewTicker(1 * time.Second)
	cpu := nes.cpu
	//go func() {
	for tick := range nes.ticker.C {
		log.Printf("At time %v", tick)
		spentCycles := int64(0)
		for spentCycles < int64(CpuClockRate) {
			spentCycles += int64(cpu.ExecOneInstruction())
		}
		log.Println("realtime CPU clock rate: %v", spentCycles)
	}

	return nil
}
