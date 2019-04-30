package memory

import "fmt"

const (
	OAMDMA_ADDR = 0x4014
)

type OamDma struct {
	cpuAs Memory
	oam   Memory
}

func NewOamDma(cpuAs Memory, oam Memory) *OamDma {
	return &OamDma{
		cpuAs: cpuAs,
		oam:   oam,
	}
}

func (p *OamDma) Peek(addr Ptr) byte {
	panic(fmt.Errorf("OAMDAM register is not readable"))
}

func (p *OamDma) Poke(addr Ptr, val byte) {
	if addr != OAMDMA_ADDR {
		panic(fmt.Errorf("OAMDAM address is %04x", OAMDMA_ADDR))
	}
	cpuStartAddr := Ptr(val) << 8
	for oamAddr := Ptr(0); oamAddr < 0x100; oamAddr++ {
		b := p.cpuAs.Peek(cpuStartAddr + oamAddr)
		p.oam.Poke(oamAddr, b)
	}
}
