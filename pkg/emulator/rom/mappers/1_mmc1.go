package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
)

type MMC1Mapper struct {
	mapperBase
	shiftRegister byte
	writeCounter  int
	registers     [4]byte
	useChrRam     bool
}

func init() {
	MapperConstructors[1] = NewMMC1Mapper
}

func NewMMC1Mapper(rom *ines.INesRom) Mapper {
	p := &MMC1Mapper{}
	p.prgBin = rom.PrgBin
	if len(rom.ChrBin) > 0 {
		p.chrBin = rom.ChrBin
	} else {
		// cartridge use CHR-RAM rather than CHR-ROM
		p.chrBin = make([]byte, ChrBankSize)
		p.useChrRam = true
	}
	p.registers[0] = 0x0c
	return p
}

func (p *MMC1Mapper) mapPrgAddr(addr memory.Ptr) int {
	offset := int(addr) & 0x3fff
	bank := int(p.registers[3] & 0x0f)
	switch p.registers[0] >> 2 & 0x3 {
	case 0: // Switchable 32K Area at 8000h-FFFFh
		fallthrough
	case 1: // Switchable 32K Area at 8000h-FFFFh
		bank &= 0x0e
		if addr >= 0xc000 {
			bank++
		}
	case 2:
		// Switchable 16K Area at C000h-FFFFh (via Register 3)
		// And Fixed  16K Area at 8000h-BFFFh (always 1st 16K)
		if addr < 0xc000 {
			bank = 0
		}
	default:
		// Switchable 16K Area at 8000h-BFFFh (via Register 3)
		// And Fixed  16K Area at C000h-FFFFh (always last 16K)
		if addr >= 0xc000 {
			bank = len(p.prgBin)/PrgBankSize - 1
		}
	}
	physicalAddr := bank*PrgBankSize | offset
	if physicalAddr >= len(p.prgBin) {
		panic(fmt.Sprintf("error accessing Mapper 1 PRG-ROM with address %04x (%04x/%04x)",
			addr, physicalAddr, len(p.prgBin)))
	}
	return physicalAddr
}

func (p *MMC1Mapper) PeekPrg(addr memory.Ptr) byte {
	if addr < 0x4020 {
		panic(fmt.Errorf("program trying to read PRG-ROM from Mapper 1 via invalid ROM address 0x%x", addr))
	}
	if addr < 0x8000 {
		return p.prgRam[addr-0x4020]
	}
	return p.prgBin[p.mapPrgAddr(addr)]
}

func (p *MMC1Mapper) PokePrg(addr memory.Ptr, val byte) {
	if addr < 0x4020 {
		panic(fmt.Errorf("mapper 1 PRG-ROM address 0x%x is not configured", addr))
	}
	if addr < 0x8000 {
		// write to PRG-RAM
		p.prgRam[addr-0x4020] = val
		return
	}
	// write to mapper register
	if val&0x80 != 0 {
		// reset SR
		p.shiftRegister = 0
		p.writeCounter = 0
	} else {
		p.writeCounter++
		p.shiftRegister >>= 1
		p.shiftRegister |= val & 1 << 4
		if p.writeCounter == 5 {
			p.registers[addr>>13&3] = p.shiftRegister
			p.shiftRegister = 0
			p.writeCounter = 0
			// TODO: cartridge set nametable mirroring
			switch p.registers[0] & 0x3 {
			case 0: // one-screen, lower bank;
				p.notifyNametableMirroringChangeListener(0, 0)
				p.notifyNametableMirroringChangeListener(1, 0)
				p.notifyNametableMirroringChangeListener(2, 0)
				p.notifyNametableMirroringChangeListener(3, 0)
			case 1: // one-screen, upper bank;
				p.notifyNametableMirroringChangeListener(0, 1)
				p.notifyNametableMirroringChangeListener(1, 1)
				p.notifyNametableMirroringChangeListener(2, 1)
				p.notifyNametableMirroringChangeListener(3, 1)
			case 2: // Two-Screen Vertical Mirroring
				p.notifyNametableMirroringChangeListener(0, 0)
				p.notifyNametableMirroringChangeListener(1, 1)
				p.notifyNametableMirroringChangeListener(2, 0)
				p.notifyNametableMirroringChangeListener(3, 1)
			case 3: // Two-Screen Horizontal Mirroring
				p.notifyNametableMirroringChangeListener(0, 0)
				p.notifyNametableMirroringChangeListener(1, 0)
				p.notifyNametableMirroringChangeListener(2, 1)
				p.notifyNametableMirroringChangeListener(3, 1)
			}
		}
	}
}
func (p *MMC1Mapper) mapChrAddr(addr memory.Ptr) int {
	bank := 0
	if p.registers[0]&0x10 != 0 {
		// Swap 4K of VROM at PPU 0000h and 1000h
		if addr >= 0x1000 {
			bank = int(p.registers[2] & 0x1f)
		} else {
			bank = int(p.registers[1] & 0x1f)
		}
		//return bank*ChrBankSize/2 | int(addr&0x0fff)
	} else {
		// Swap 8K of VROM at PPU 0000h
		bank = int(p.registers[1] & 0x1e)
		if addr >= 0x1000 {
			bank++
		}
	}
	return ChrBankSize/2*bank | int(addr&0x0fff)
}
func (p *MMC1Mapper) PeekChr(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 1 CHR-ROM/CHR-RAM address 0x%x is not configured", addr))
	}
	return p.chrBin[p.mapChrAddr(addr)]
}

func (p *MMC1Mapper) PokeChr(addr memory.Ptr, val byte) {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 1 CHR-ROM/CHR-RAM address 0x%x is not configured", addr))
	}
	if !p.useChrRam {
		panic(fmt.Errorf("this mapper 1 cartridge uses CHR-ROM, writing address %04x is not possible", addr))
	}
	p.chrBin[p.mapChrAddr(addr)] = val
}
