package ram

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

// NES only has 2 kiB on-chip VRAM, but certain mappers can provide additional VRAM up to 4kiB.
// We implement CIRAM as 4 kiB for simplicity.
// http://wiki.nesdev.com/w/index.php/Mirroring#Nametable_Mirroring

type CIRam struct {
	ram          [0x1000]byte
	mirroringMap [4]int
}

func NewCIRam() *CIRam {
	ram := &CIRam{}
	return ram
}

func (p *CIRam) SetNametableMirroring(logical, physical int) {
	if logical < 0 || logical >= len(p.mirroringMap) {
		panic(fmt.Errorf("logical nametable ID %v is out of bound [%v, %v)", logical, 0, len(p.mirroringMap)))
	}
	if physical < 0 || physical >= len(p.mirroringMap) {
		panic(fmt.Errorf("logical nametable ID %v is out of bound [%v, %v)", physical, 0, len(p.mirroringMap)))
	}
	p.mirroringMap[logical] = physical
}

func (p *CIRam) mapAddr(addr memory.Ptr) memory.Ptr {
	logical := (addr & 0xfff) / 0x400
	physical := p.mirroringMap[logical]
	return memory.Ptr(physical*0x400) | addr&0x3ff
}

func (p *CIRam) Peek(addr memory.Ptr) byte {
	if addr < 0x2000 {
		panic(fmt.Errorf("error reading CIRAM via invalid nametable address %04x", addr))
	}
	return p.ram[p.mapAddr(addr)]
}

func (p *CIRam) Poke(addr memory.Ptr, val byte) {
	if addr < 0x2000 {
		panic(fmt.Errorf("error writing CIRAM via invalid nametable address %04x", addr))
	}
	p.ram[p.mapAddr(addr)] = val
}
