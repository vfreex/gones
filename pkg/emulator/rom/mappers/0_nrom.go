package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
)

type NROMMapper struct {
	mapperBase
	prgBankMapping [2]int
}

func init() {
	MapperConstructors[0] = NewNROMMapper
}

func NewNROMMapper(rom *ines.INesRom) Mapper {
	p := &NROMMapper{}
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
