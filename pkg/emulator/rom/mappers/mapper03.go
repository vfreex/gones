package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
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

type Mapper03 struct {
	bankSelectRegister byte // CNROM only implements the lowest 2 bits
	Prg                Mapper03PrgRom
	Chr                Mapper03ChrRom
}

func NewMapper03(prgBin, chrBin []byte) *Mapper03 {
	mapper := &Mapper03{}
	mapper.Prg.bin = prgBin
	mapper.Prg.mapper = mapper
	mapper.Chr.bin = chrBin
	mapper.Chr.mapper = mapper
	return mapper
}

func (p *Mapper03) Map() (prg memory.Memory, chr memory.Memory) {
	return &p.Prg, &p.Chr
}

type Mapper03PrgRom struct {
	bin []byte
	mapper *Mapper03
}

func (p *Mapper03PrgRom) Peek(addr memory.Ptr) byte {
	if addr < 0x8000 {
		panic(fmt.Errorf("program trying to read from Mapper 03 via invalid ROM address 0x%x", addr))
	}
	if len(p.bin) == 2*PrgBankSize {
		return p.bin[addr-0x8000]
	} else {
		return p.bin[(addr-0x8000)&0x3fff]
	}
}

func (p *Mapper03PrgRom) Poke(addr memory.Ptr, val byte) {
	if addr < 0x8000 {
		panic(fmt.Errorf("mapper 03 Program ROM address 0x%x is not writable", addr))
	}
	p.mapper.bankSelectRegister = val & 0x3 // CNROM only implements the lowest 2 bits, capping it at 32 KiB
}

type Mapper03ChrRom struct {
	bin []byte
	mapper *Mapper03
}

func (p *Mapper03ChrRom) Peek(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 03 Character ROM address 0x%x is not readable", addr))
	}
	newBank := int(p.mapper.bankSelectRegister)
	return p.bin[newBank * ChrBankSize | int(addr)]
}

func (*Mapper03ChrRom) Poke(addr memory.Ptr, val byte) {
	panic(fmt.Errorf("mapper 03 Character ROM address 0x%x is not writable", addr))
}
