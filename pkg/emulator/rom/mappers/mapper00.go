package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

type Mapper00 struct {
	Prg Mapper00PrgRom
	Chr Mapper00Chr
}

func (p *Mapper00) Map() (prg memory.Memory, chr memory.Memory) {
	return &p.Prg, &p.Chr
}

func NewMapper00(prgBin, chrBin []byte) *Mapper00 {
	mapper := &Mapper00{}
	mapper.Prg.bin = prgBin
	if len(chrBin) > 0 {
		mapper.Chr.bin = chrBin
	} else {
		// cartridge use CHR-RAM rather than CHR-ROM
		mapper.Chr.bin = make([]byte, ChrBankSize)
		mapper.Chr.isRam = true
	}
	return mapper
}

type Mapper00PrgRom struct {
	bin []byte
}

func (p *Mapper00PrgRom) Peek(addr memory.Ptr) byte {
	if addr < 0x8000 {
		panic(fmt.Errorf("program trying to read from Mapper 03 via invalid ROM address 0x%x", addr))
	}
	if len(p.bin) == 2*PrgBankSize {
		return p.bin[addr-0x8000]
	} else {
		return p.bin[(addr-0x8000)&0x3fff]
	}
}

func (p *Mapper00PrgRom) Poke(addr memory.Ptr, val byte) {
	if addr < 0x8000 {
		panic(fmt.Errorf("mapper 00 Program ROM address 0x%x is not writable", addr))
	}
	panic(fmt.Errorf("mapper 00 Program ROM address 0x%x is not writable", addr))
}

type Mapper00Chr struct {
	bin   []byte
	isRam bool
}

func (p *Mapper00Chr) Peek(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 00 Character ROM address 0x%x is not readable", addr))
	}
	return p.bin[addr]
}

func (p *Mapper00Chr) Poke(addr memory.Ptr, val byte) {
	if !p.isRam {
		panic(fmt.Errorf("mapper 00 Character ROM address 0x%x is not writable", addr))
	}
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 00 Character RAM address 0x%x is not writable", addr))
	}
	p.bin[addr] = val
}
