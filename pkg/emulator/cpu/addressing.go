package cpu

import (
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

func (cpu *Cpu) AddressOperand(am AddressingMode) (memory.Ptr, int) {
	switch am {
	case IMP:
		return 0, 0
	case IMM:
		return cpu.AddressImm()
	case ZP:
		return cpu.AddressZP()
	case ZPX:
		return cpu.AddressZPX()
	case ZPY:
		return cpu.AddressZPY()
	case ABS:
		return cpu.AddressAbs()
	case ABX:
		return cpu.AddressAbX()
	case ABY:
		return cpu.AddressAbY()
	case REL:
		return cpu.AddressRel()
	case IND:
		return cpu.AddressInd()
	case IZX:
		return cpu.AddressIzx()
	case IZY:
		return cpu.AddressIzy()
	default:
		panic(fmt.Errorf("unsupported addressing mode: %s", am))
	}
}

func isCrossPage(addr memory.Ptr, offset uint8) bool {
	return addr&0xff00 != (addr-memory.Ptr(offset))&0xff00
}

func (cpu *Cpu) AddressImm() (memory.Ptr, int) {
	addr := cpu.PC
	cpu.PC++
	return addr, 0
}

func (cpu *Cpu) AddressZP() (memory.Ptr, int) {
	addr, _ := cpu.AddressImm()
	addr = memory.Ptr(cpu.Memory.Peek(addr))
	return addr, 1
}

func (cpu *Cpu) AddressZPX() (memory.Ptr, int) {
	addr, _ := cpu.AddressImm()
	addr = memory.Ptr(cpu.Memory.Peek(addr)+cpu.X) & 0xff
	return addr, 2
}

func (cpu *Cpu) AddressZPY() (memory.Ptr, int) {
	addr, _ := cpu.AddressImm()
	addr = memory.Ptr(cpu.Memory.Peek(addr)+cpu.Y) & 0xff
	return addr, 2
}

func (cpu *Cpu) AddressAbs() (memory.Ptr, int) {
	low, _ := cpu.AddressImm()
	high, _ := cpu.AddressImm()
	addr := (memory.Ptr(cpu.Memory.Peek(high)) << 8) | memory.Ptr(cpu.Memory.Peek(low))
	return addr, 2
}
func (cpu *Cpu) AddressAbX() (memory.Ptr, int) {
	addr, _ := cpu.AddressAbs()
	addr += memory.Ptr(cpu.X)
	if isCrossPage(addr, cpu.X) {
		return addr, 3
	}
	return addr, 2
}
func (cpu *Cpu) AddressAbY() (memory.Ptr, int) {
	addr, _ := cpu.AddressAbs()
	addr += memory.Ptr(cpu.Y)
	if isCrossPage(addr, cpu.Y) {
		return addr, 3
	}
	return addr, 2
}

func (cpu *Cpu) AddressRel() (memory.Ptr, int) {
	addr, _ := cpu.AddressImm()
	return cpu.PC + memory.PtrDist(int8(cpu.Memory.Peek(addr))), 1
}

func (cpu *Cpu) AddressInd() (memory.Ptr, int) {
	addr, _ := cpu.AddressAbs()
	low := cpu.Memory.Peek(addr)
	// 6502 CPU bug
	addr2 := addr&0xff00 | (addr+1)&0x00ff
	high := cpu.Memory.Peek(addr2)
	addr3 := memory.Ptr(high)<<8 | memory.Ptr(low)
	return addr3, 4
}

func (cpu *Cpu) AddressIzx() (memory.Ptr, int) {
	addr, _ := cpu.AddressZPX()
	low := memory.Ptr(cpu.Memory.Peek(addr))
	high := memory.Ptr(cpu.Memory.Peek((addr + 1) & 0xff))
	return high<<8 | low, 4
}

func (cpu *Cpu) AddressIzy() (memory.Ptr, int) {
	addr, _ := cpu.AddressZP()
	low := memory.Ptr(cpu.Memory.Peek(addr))
	high := memory.Ptr(cpu.Memory.Peek((addr + 1) & 0xff))
	addr2 := high<<8 | low
	if isCrossPage(addr2, cpu.Y) {
		return addr2 + memory.Ptr(cpu.Y), 4
	}
	return addr2 + memory.Ptr(cpu.Y), 3
}
