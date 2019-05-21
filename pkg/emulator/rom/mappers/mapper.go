package mappers

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
)

const (
	PrgBankSize = 16 * 1024 // bytes in a PRG/ROM bank
	ChrBankSize = 8 * 1024  // bytes in a CHR/VROM bank
)

//type Mapper interface {
//	Map() (prg memory.Memory, chr memory.Memory)
//}

type Mapper interface {
	PeekPrg(addr memory.Ptr) byte
	PokePrg(addr memory.Ptr, val byte)
	PeekChr(addr memory.Ptr) byte
	PokeChr(addr memory.Ptr, val byte)
	AddNametableMirroringChangeListener(listener NametableMirroringChangeListener)
}

type MapperINesConstructor func(rom *ines.INesRom) Mapper

var MapperConstructors map[int]MapperINesConstructor = make(map[int]MapperINesConstructor)

type NametableMirroringChangeListener func(logical, physical int)

type mapperBase struct {
	prgBin                            []byte
	chrBin                            []byte
	useChrRam                         bool
	prgRam                            [0x3fe0]byte
	nametableMirroringChangeListeners []NametableMirroringChangeListener
}

func (p *mapperBase) AddNametableMirroringChangeListener(listener NametableMirroringChangeListener) {
	p.nametableMirroringChangeListeners = append(p.nametableMirroringChangeListeners, listener)
}

func (p *mapperBase) notifyNametableMirroringChangeListener(logical, physical int) {
	for _, listener := range p.nametableMirroringChangeListeners {
		listener(logical, physical)
	}
}

type MapperPrgMemoryAdapter struct {
	mapper Mapper
}

func (p *MapperPrgMemoryAdapter) Peek(addr memory.Ptr) byte {
	return p.mapper.PeekPrg(addr)
}

func (p *MapperPrgMemoryAdapter) Poke(addr memory.Ptr, val byte) {
	p.mapper.PokePrg(addr, val)
}

type MapperChrMemoryAdapter struct {
	mapper Mapper
}

func (p *MapperChrMemoryAdapter) Peek(addr memory.Ptr) byte {
	return p.mapper.PeekChr(addr)
}

func (p *MapperChrMemoryAdapter) Poke(addr memory.Ptr, val byte) {
	p.mapper.PokeChr(addr, val)
}

func MapAddressSpaces(p Mapper, cpuAS, ppuAS memory.AddressSpace) {
	prg := &MapperPrgMemoryAdapter{p}
	chr := &MapperChrMemoryAdapter{p}
	cpuAS.AddMapping(0x4020, 0xbfe0, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		prg, nil)
	ppuAS.AddMapping(0, 0x2000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE,
		chr, nil)
}
