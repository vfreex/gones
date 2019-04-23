package nes

import (
	"github.com/vfreex/gones/pkg/emulator/cpu"
	"github.com/vfreex/gones/pkg/emulator/memory"
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
}

func NewNes() NES {
	cpuAS := &memory.AddressSpaceImpl{}
	mainRam := ram.NewMainRAM()
	nes := &NESImpl{
		cpu:   cpu.NewCpu(cpuAS),
		cpuAS: cpuAS,
		ram:   mainRam,
	}

	// setting up memory map
	// 0x0000 - ox1fff RAM
	cpuAS.AddMapping(0, 0x2000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		mainRam, nil)
	cpuAS.AddMapping(0x2000, 0x6000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		ram.NewRAM(0x6000), func(addr memory.Ptr) memory.Ptr {
			return addr - 0x2000
		})
	return nes
}

func (nes *NESImpl) LoadCartridge(cartridge *ines.INesRom) error {
	nes.cpuAS.AddMapping(0x8000, 0x8000,
		memory.MMAP_MODE_READ, cartridge, nil)
	return nil
}

func (nes *NESImpl) Start() error {
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
	//}()
	return nil
}
