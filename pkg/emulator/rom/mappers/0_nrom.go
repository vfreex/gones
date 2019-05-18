package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/ram"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
)

type NROMMapper struct {
	mapperBase
	prgBankMapping [2]int
	prgRam         [0x3fe0]byte
	useChrRam      bool
}

func init() {
	MapperConstructors[0] = NewNROMMapper
}

func NewNROMMapper(rom *ines.INesRom, mirroringController ram.NametableMirrorController) Mapper {
	p := &NROMMapper{}
	p.MirroringController = mirroringController
	p.prgBin = rom.PrgBin
	if len(p.prgBin) > PrgBankSize {
		p.prgBankMapping[1] = 1
	}
	if len(rom.ChrBin) > 0 {
		p.chrBin = rom.ChrBin
	} else {
		// cartridge use CHR-RAM rather than CHR-ROM
		p.chrBin = make([]byte, ChrBankSize)
		p.useChrRam = true
	}
	if rom.Header.Flags6&ines.FLAGS6_FOUR_SCREEN_VRAM_ON != 0 {
		mirroringController.SetNametableMirroring(0, 0)
		mirroringController.SetNametableMirroring(1, 1)
		mirroringController.SetNametableMirroring(2, 2)
		mirroringController.SetNametableMirroring(3, 3)
	} else if rom.Header.Flags6&ines.FLAGS6_VERTICAL_MIRRORING != 0 {
		mirroringController.SetVerticalMirroring()
	} else {
		mirroringController.SetHorizontalMirroring()
	}
	return p
}

func (p *NROMMapper) PeekPrg(addr memory.Ptr) byte {
	if addr < 0x4020 {
		panic(fmt.Errorf("mapper 0 PRG-ROM address %04x is not configured", addr))
	}
	if addr < 0x8000 {
		return p.prgRam[addr-0x4020]
	}
	bank := p.prgBankMapping[int(addr-0x8000)/PrgBankSize]
	return p.prgBin[bank*PrgBankSize|int(addr)&0x3fff]
}

func (p *NROMMapper) PokePrg(addr memory.Ptr, val byte) {
	if addr < 0x4020 {
		panic(fmt.Errorf("mapper 0 PRG-ROM address %04x is not configured", addr))
	}
	if addr < 0x8000 {
		p.prgRam[addr-0x4020] = val
		return
	}
	panic(fmt.Errorf("mapper 0 PRG-ROM address %04x is not writable", addr))
}

func (p *NROMMapper) PeekChr(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 0 CHR-ROM/CHR-RAM address %04x is not configured", addr))
	}
	return p.chrBin[addr]
}

func (p *NROMMapper) PokeChr(addr memory.Ptr, val byte) {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 0 CHR-ROM/CHR-RAM %04x is not configured", addr))
	}
	if p.useChrRam {
		panic(fmt.Errorf("this mapper 0 cartridge uses CHR-ROM, writing address %04x is not possible", addr))
	}
	p.chrBin[addr] = val
}
