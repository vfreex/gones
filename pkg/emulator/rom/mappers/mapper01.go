package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/ram"
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

type Mapper01 struct {
	bankSelectRegister byte // CNROM only implements the lowest 2 bits
	Prg                Mapper01PrgRom
	Chr                Mapper01ChrRom
	PrgRam             *ram.RAM
	ShiftRegister      byte
	WriteCounter       int
	Registers          [4]byte
}

func NewMapper01(prgBin, chrBin []byte) *Mapper01 {
	mapper := &Mapper01{}
	mapper.Prg.bin = prgBin
	mapper.Prg.mapper = mapper
	if len(chrBin) > 0 {
		mapper.Chr.bin = chrBin
	} else {
		// cartridge use CHR-RAM rather than CHR-ROM
		mapper.Chr.bin = make([]byte, ChrBankSize)
		mapper.Chr.isRam = true
	}
	mapper.Chr.mapper = mapper
	mapper.PrgRam = ram.NewRAM(0x4000)
	mapper.Registers[0] = 0x0c
	return mapper
}

func (p *Mapper01) Map() (prg memory.Memory, chr memory.Memory) {
	return &p.Prg, &p.Chr
}

type Mapper01PrgRom struct {
	bin    []byte
	mapper *Mapper01
}

func (p *Mapper01PrgRom) Peek(addr memory.Ptr) byte {
	if addr < 0x4020 {
		panic(fmt.Errorf("program trying to read from Mapper 01 via invalid ROM address 0x%x", addr))
	}
	if addr < 0x8000 {
		return p.mapper.PrgRam.Peek(addr - 0x4000)
	}
	offset := int(addr) & 0x3fff
	bank := int(p.mapper.Registers[3] & 0x0f)
	switch p.mapper.Registers[0] >> 2 & 0x3 {
	case 0: // Switchable 32K Area at 8000h-FFFFh
		fallthrough
	case 1: // Switchable 32K Area at 8000h-FFFFh
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
			bank = len(p.bin)/PrgBankSize - 1
		}
	}
	physicalAddr := bank*PrgBankSize | offset
	return p.bin[physicalAddr]
}

func (p *Mapper01PrgRom) Poke(addr memory.Ptr, val byte) {
	if addr < 0x4020 {
		panic(fmt.Errorf("mapper 01 Program ROM address 0x%x is not writable", addr))
	}
	if addr < 0x8000 {
		// write to PRG-RAM
		p.mapper.PrgRam.Poke(addr-0x4000, val)
		return
	}
	// write to mapper register
	if val&0x80 != 0 {
		// reset SR
		p.mapper.ShiftRegister = 0
		p.mapper.WriteCounter = 0
	} else {
		p.mapper.WriteCounter++
		p.mapper.ShiftRegister >>= 1
		p.mapper.ShiftRegister |= val & 1 << 4
		if p.mapper.WriteCounter == 5 {
			p.mapper.Registers[addr>>13&3] = p.mapper.ShiftRegister
			p.mapper.ShiftRegister = 0
			p.mapper.WriteCounter = 0
			// TODO: cartridge set nametable mirroring
			switch p.mapper.Registers[0] & 0x3 {
			case 2: // Two-Screen Vertical Mirroring
			case 3: // Two-Screen Horizontal Mirroring
			default:
			}
		}
	}

}

type Mapper01ChrRom struct {
	bin    []byte
	mapper *Mapper01
	isRam  bool
}

func (p *Mapper01ChrRom) Peek(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 01 Character ROM address 0x%x is not readable", addr))
	}
	bank := 0
	if p.mapper.Registers[0]&0x10 != 0 {
		// Swap 4K of VROM at PPU 0000h and 1000h
		if addr >= 0x1000 {
			bank = int(p.mapper.Registers[2] & 0x1f)
		} else {
			bank = int(p.mapper.Registers[1] & 0x1f)
		}
		return p.bin[bank*ChrBankSize/2|int(addr&0x0fff)]
	} else {
		// Swap 8K of VROM at PPU 0000h
		bank = int(p.mapper.Registers[1] & 0x1f)
		return p.bin[bank*ChrBankSize/2|int(addr)]
	}
}

func (p *Mapper01ChrRom) Poke(addr memory.Ptr, val byte) {
	if !p.isRam {
		panic(fmt.Errorf("mapper 01 Character ROM address 0x%x is not writable", addr))
	}
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 01 Character RAM address 0x%x is not writable", addr))
	}
	p.bin[addr] = val
}
