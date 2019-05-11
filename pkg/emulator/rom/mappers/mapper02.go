package mappers

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
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

type Mapper02 struct {
	bankSelectRegister byte // CNROM only implements the lowest 2 bits
	Prg                Mapper02PrgRom
	Chr                Mapper02ChrRom
}

func NewMapper02(prgBin, chrBin []byte) *Mapper02 {
	mapper := &Mapper02{}
	mapper.Prg.bin = prgBin
	mapper.Prg.mapper = mapper
	if len(chrBin) > 0 {
		mapper.Chr.bin = chrBin
	} else {
		mapper.Chr.bin = make([]byte, ChrBankSize)
		mapper.Chr.isRam = true
	}
	mapper.Chr.mapper = mapper
	return mapper
}

func (p *Mapper02) Map() (prg memory.Memory, chr memory.Memory) {
	return &p.Prg, &p.Chr
}

type Mapper02PrgRom struct {
	bin    []byte
	mapper *Mapper02
}

func (p *Mapper02PrgRom) Peek(addr memory.Ptr) byte {
	if addr < 0x8000 {
		panic(fmt.Errorf("program trying to read from Mapper 02 via invalid ROM address 0x%x", addr))
	}
	var bank int
	if addr >= 0xc000 {
		bank = len(p.bin)/PrgBankSize - 1
	} else {
		bank = int(p.mapper.bankSelectRegister)
	}
	return p.bin[bank*PrgBankSize|int(addr)&0x3fff]
}

func (p *Mapper02PrgRom) Poke(addr memory.Ptr, val byte) {
	if addr < 0x8000 {
		panic(fmt.Errorf("mapper 02 Program ROM address 0x%x is not writable", addr))
	}
	p.mapper.bankSelectRegister = val
}

type Mapper02ChrRom struct {
	bin    []byte
	mapper *Mapper02
	isRam  bool
}

func (p *Mapper02ChrRom) Peek(addr memory.Ptr) byte {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 02 Character ROM address 0x%x is not readable", addr))
	}
	return p.bin[addr]
}

func (p *Mapper02ChrRom) Poke(addr memory.Ptr, val byte) {
	if addr >= 0x2000 {
		panic(fmt.Errorf("mapper 02 Character ROM address 0x%x is not readable", addr))
	}
	if !p.isRam {
		panic(fmt.Errorf("mapper 02 Character ROM address 0x%x is not writable", addr))
	}
	p.bin[addr] = val
}
