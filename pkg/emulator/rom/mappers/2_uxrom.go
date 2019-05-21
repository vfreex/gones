package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
)

/*
http://wiki.nesdev.com/w/index.php/UxROM
PRG ROM capacity	256K/4096K
PRG ROM window	16K + 16K fixed
PRG RAM capacity	None
CHR capacity	8K
CHR window	n/a

CPU $8000-$BFFF: 16 KB switchable PRG ROM bank
CPU $C000-$FFFF: 16 KB PRG ROM bank, fixed to the last bank
*/

type UxRomMapper struct {
	mapperBase
	bankSelect byte
}

func init() {
	MapperConstructors[2] = NewUxRomMapper
}

func NewUxRomMapper(rom *ines.INesRom) Mapper {
	p := &UxRomMapper{}
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

func (p *UxRomMapper) PeekPrg(addr memory.Ptr) byte {
	if addr < 0x4020 {
		panic(fmt.Errorf("program trying to read from Mapper 2 via invalid ROM address %04x", addr))
	}
	if addr < 0x8000 {
		return p.prgRam[addr-0x4020]
	}
	var bank int
	if addr >= 0xc000 {
		bank = len(p.prgBin)/PrgBankSize - 1
	} else {
		bank = int(p.bankSelect)
	}
	return p.prgBin[bank*PrgBankSize|int(addr)&0x3fff]
}

func (p *UxRomMapper) PokePrg(addr memory.Ptr, val byte) {
	if addr < 0x4020 {
		panic(fmt.Errorf("mapper 2 PRG-ROM address 0x%x is not configured", addr))
	}
	if addr < 0x8000 {
		// write to PRG-RAM
		p.prgRam[addr-0x4020] = val
		return
	}
	p.bankSelect = val
}

func (p *UxRomMapper) PeekChr(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 2 CHR-ROM/CHR-RAM address %04x is not configured", addr))
	}
	return p.chrBin[addr]
}

func (p *UxRomMapper) PokeChr(addr memory.Ptr, val byte) {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 2 CHR-ROM/CHR-RAM address %04x is not configured", addr))
	}
	if !p.useChrRam {
		panic(fmt.Errorf("this mapper 2 cartridge uses CHR-ROM, writing address %04x is not possible", addr))
	}
	p.chrBin[addr] = val
}
