package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
)

/*
http://wiki.nesdev.com/w/index.php/INES_Mapper_003

iNES Mapper 003 is used to designate the CNROM board,
generalized to support up to 256 banks (2048 KiB) of CHR ROM.

PRG ROM size: 16 KiB or 32 KiB
PRG ROM bank size: Not bankswitched
PRG RAM: None
CHR capacity: Up to 2048 KiB ROM
CHR bank size: 8 KiB
Nametable mirroring: Fixed vertical or horizontal mirroring
Subject to bus conflicts: Yes (CNROM), but not all compatible boards have bus conflicts.
*/

type CNROMMapper struct {
	mapperBase
	bankSelect byte
}

func init() {
	MapperConstructors[3] = NewCNROMMapper
}

func NewCNROMMapper(rom *ines.INesRom) Mapper {
	p := &CNROMMapper{}
	p.prgBin = rom.PrgBin
	if len(rom.ChrBin) > 0 {
		p.chrBin = rom.ChrBin
	} else {
		// cartridge use CHR-RAM rather than CHR-ROM
		p.chrBin = make([]byte, ChrBankSize)
		p.useChrRam = true
	}
	return p
}

func (p *CNROMMapper) PeekPrg(addr memory.Ptr) byte {
	if addr < 0x4020 {
		panic(fmt.Errorf("program trying to read from Mapper 3 via invalid ROM address %04x", addr))
	}
	if addr < 0x8000 {
		return p.prgRam[addr-0x4020]
	}
	if len(p.prgBin) == 2*PrgBankSize {
		return p.prgBin[addr-0x8000]
	} else {
		return p.prgBin[(addr-0x8000)&0x3fff]
	}
}

func (p *CNROMMapper) PokePrg(addr memory.Ptr, val byte) {
	if addr < 0x4020 {
		panic(fmt.Errorf("mapper 3 PRG-ROM address 0x%x is not configured", addr))
	}
	if addr < 0x8000 {
		// write to PRG-RAM
		p.prgRam[addr-0x4020] = val
		return
	}
	p.bankSelect = val
}

func (p *CNROMMapper) PeekChr(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 3 CHR-ROM/CHR-RAM address %04x is not configured", addr))
	}
	newBank := int(p.bankSelect)
	return p.chrBin[newBank*ChrBankSize|int(addr)]
}

func (p *CNROMMapper) PokeChr(addr memory.Ptr, val byte) {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 3 CHR-ROM/CHR-RAM address %04x is not configured", addr))
	}
	if !p.useChrRam {
		panic(fmt.Errorf("this mapper 3 cartridge uses CHR-ROM, writing address %04x is not possible", addr))
	}
	newBank := int(p.bankSelect)
	p.chrBin[newBank*ChrBankSize|int(addr)] = val
}
