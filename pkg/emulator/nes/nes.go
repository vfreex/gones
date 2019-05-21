package nes

import (
	"fmt"
	pkgLogger "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/cpu"
	"github.com/vfreex/gones/pkg/emulator/joypad"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/ppu"
	"github.com/vfreex/gones/pkg/emulator/ram"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
	"github.com/vfreex/gones/pkg/emulator/rom/mappers"
	"time"
)

const (
	FPS = 60
)

var logger = pkgLogger.GetLogger()

type NES interface {
	LoadCartridge(cartridge *ines.INesRom) error
	Start() error
}

type NESImpl struct {
	ticker  *time.Ticker
	cpu     *cpu.Cpu
	cpuAS   memory.AddressSpace
	ram     memory.Memory
	ppu     *ppu.PPUImpl
	ppuAS   memory.AddressSpace
	vram    *ram.CIRam
	display *NesDiplay
	joypads *joypad.Joypads
}

func NewNes() NES {
	nes := &NESImpl{
		cpuAS:   &memory.AddressSpaceImpl{},
		ram:     ram.NewMainRAM(),
		ppuAS:   &memory.AddressSpaceImpl{},
		vram:    ram.NewCIRam(),
		joypads: joypad.NewJoypads(),
	}
	nes.cpu = cpu.NewCpu(nes.cpuAS)
	nes.ppu = ppu.NewPPU(nes.ppuAS, nes.cpu)
	nes.display = NewDisplay(&nes.ppu.RenderedBuffer)

	// setting up CPU memory map
	// 0x0000 - ox1fff RAM
	nes.cpuAS.AddMapping(0, 0x2000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		nes.ram, nil)
	// fake memory map range
	nes.cpuAS.AddMapping(0x4000, 0x14, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		ram.NewRAM(0x14), func(addr memory.Ptr) memory.Ptr {
			return addr - 0x4000
		})
	//nes.cpuAS.AddMapping(0x4014, 1, memory.MMAP_MODE_WRITE,
	nes.ppu.MapToCPUAddressSpace(nes.cpuAS)
	// fake memory map range
	nes.cpuAS.AddMapping(0x4015, 1, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		ram.NewRAM(0x01), func(addr memory.Ptr) memory.Ptr {
			return addr - 0x4015
		})
	nes.cpuAS.AddMapping(0x4016, 2, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		nes.joypads, nil)

	// setting up PPU memory map
	// https://wiki.nesdev.com/w/index.php/PPU_memory_map
	nes.ppuAS.AddMapping(0x3F00, 0x100,
		memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, &nes.ppu.Palette, nil)

	return nes
}

func (nes *NESImpl) LoadCartridge(cartridge *ines.INesRom) error {
	if cartridge.Header.Flags6&ines.FLAGS6_FOUR_SCREEN_VRAM_ON != 0 {
		nes.vram.SetNametableMirroring(0,0)
		nes.vram.SetNametableMirroring(1,1)
		nes.vram.SetNametableMirroring(2,2)
		nes.vram.SetNametableMirroring(3,3)
	} else if cartridge.Header.Flags6&ines.FLAGS6_VERTICAL_MIRRORING != 0 {
		nes.vram.SetNametableMirroring(0,0)
		nes.vram.SetNametableMirroring(1,1)
		nes.vram.SetNametableMirroring(2,0)
		nes.vram.SetNametableMirroring(3,1)
	} else {
		nes.vram.SetNametableMirroring(0,0)
		nes.vram.SetNametableMirroring(1,0)
		nes.vram.SetNametableMirroring(2,1)
		nes.vram.SetNametableMirroring(3,1)
	}
	var mapper mappers.Mapper
	mapperConstructor := mappers.MapperConstructors[cartridge.Header.GetMapperType()]
	if mapperConstructor != nil {
		mapper = mapperConstructor(cartridge)
	} else {
		panic(fmt.Errorf("cartridge uses unsupported mapper %v", cartridge.Header.GetMapperType()))
	}
	mappers.MapAddressSpaces(mapper, nes.cpuAS, nes.ppuAS)

	nes.ppuAS.AddMapping(0x2000, 0x1f00, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		nes.vram, nil)

	// mapper may change nametable mirroring at runtime
	mapper.AddNametableMirroringChangeListener(func(logical, physical int) {
		nes.vram.SetNametableMirroring(logical, physical)
	})

	return nil
}

func (nes *NESImpl) Start() error {
	nes.cpuAS.Map()
	nes.ppuAS.Map()

	const fps = 60
	interval := 1 * time.Second / fps
	cpuCyclesPerFrame := 29780
	nes.ticker = time.NewTicker(interval)
	nes.ppu.NewFrameHandler = func(frame *[240][256]ppu.RBGColor, frameID int) {
		nes.display.Refresh()
	}
	cpu := nes.cpu
	cpu.PowerUp()

	frames := 0
	go func() {
		for tick := range nes.ticker.C {
			//tick:=time.Now()
			logger.Infof("At time %v", tick)

			spentCycles := int64(0)
			loop := 0
			for spentCycles < int64(cpuCyclesPerFrame) {
				if nes.display.RequestReset {
					nes.cpu.Reset()
					nes.display.RequestReset = false
				}
				if nes.display.StepInstruction {
					<-nes.display.NextCh
				}
				cycles := int64(cpu.ExecOneInstruction())
				//cycles := int64(1)
				if cycles <= 0 {
					panic("invalid cycle")
				}
				for pp := int64(0); pp < cycles*3; pp++ {
					nes.ppu.Step()
				}
				spentCycles += cycles
				loop++
				//logger.Debug("")
				//logger.Infof("spent %d/%d CPU cycles", spentCycles, cpuCyclesPerFrame)
			}
			//nes.display.Refresh()
			// update joypad
			nes.joypads.Joypads[0].Buttons = nes.display.Keys
			//logger.SetOutput(os.Stderr)
			logger.Info("----------------------------------------------------------")
			now := time.Now()
			actualTime := now.Sub(tick)
			logger.Infof("spent %v/%v to render frame #%d after running %v loops / %v cycles",
				actualTime, interval, frames, loop, spentCycles)
			frames++
			//nes.ticker.Stop()
			//close(stopCh)
			if nes.display.StepFrame {
				ch := <-nes.display.NextCh
				if ch == 0xff {
					nes.cpu.Reset()
					// TODO: also reset PPU?
				}
			}
		}
	}()
	nes.display.Show()
	//<-stopCh
	nes.ticker.Stop()
	return nil
}
