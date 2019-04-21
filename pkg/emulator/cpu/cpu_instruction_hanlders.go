package cpu

import (
	"encoding/hex"
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"log"
)

var opcodeHandlers = [256]InstructionHandler{
	0x10: (*Cpu).ExecBPL,
	0x24: (*Cpu).ExecBIT,
	0x2c: (*Cpu).ExecBIT,
	0x78: (*Cpu).ExecSEI,
	0x86: (*Cpu).ExecSTX,
	0x8e: (*Cpu).ExecSTX,
	0x96: (*Cpu).ExecSTX,
	0x9a: (*Cpu).ExecTXS,
	0xd8: (*Cpu).ExecCLD,
	0xa2: (*Cpu).ExecLDX,
	0xc8: (*Cpu).ExecINY,
	0xe8: (*Cpu).ExecINX,
	0xa9: (*Cpu).ExecLDA,
	0xa5: (*Cpu).ExecLDA,
	0xb5: (*Cpu).ExecLDA,
	0xad: (*Cpu).ExecLDA,
	0xbd: (*Cpu).ExecLDA,
	0xb9: (*Cpu).ExecLDA,
	0xa1: (*Cpu).ExecLDA,
	0xb1: (*Cpu).ExecLDA,
}

func (cpu *Cpu) AddressOperand(am AddressingMode) (memory.Ptr, int) {
	switch am {
	case IMP:
		return 0, 0
	case IMM:
		return cpu.AddressImm()
	case ZP:
		return cpu.AddressZP()
	case ZPX:
		return cpu.AddressZP()
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
	return addr&0xff00 != (addr+memory.Ptr(offset))&0xff00
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
	high := cpu.Memory.Peek((addr + 1) & 0xff)
	addr2 := memory.Ptr(high)<<8 | memory.Ptr(low)
	return addr2, 4
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

func (cpu *Cpu) ReadNextInstruction() (byte, []byte, int) {
	opcode := cpu.Memory.Peek(cpu.PC)
	info := &InstructionInfos[opcode]
	arguments := make([]byte, info.AddressingMode.GetArgumentCount())
	cycles := 1
	switch info.AddressingMode.GetArgumentCount() {
	case 2:
		arguments[1] = cpu.Memory.Peek(cpu.PC + 2)
		//cycles++
		fallthrough
	case 1:
		arguments[0] = cpu.Memory.Peek(cpu.PC + 1)
		//cycles++
	}
	log.Printf("got instruction at %04x: %02x(%s %s) %s",
		cpu.PC, opcode, info.Nemonics, info.AddressingMode, hex.EncodeToString(arguments))
	return opcode, arguments, cycles
}

type InstructionHandler func(cpu *Cpu, operandAddr memory.Ptr) (cyclesTook int)

func (cpu *Cpu) ExecSEI(operandAddr memory.Ptr) int {
	log.Printf("Exec SEI")
	cpu.P.Set(PFLAG_I, true)
	return 1
}

func (cpu *Cpu) ExecCLD(operandAddr memory.Ptr) int {
	log.Printf("Exec CLD")
	cpu.P.Set(PFLAG_D, false)
	return 1
}

func (cpu *Cpu) ExecLDA(operandAddr memory.Ptr) int {
	log.Printf("Exec LDA")
	cpu.A = cpu.Memory.Peek(operandAddr)
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 1
}

func (cpu *Cpu) ExecLDX(operandAddr memory.Ptr) int {
	log.Printf("Exec LDX")
	cpu.X = cpu.Memory.Peek(operandAddr)
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}

func (cpu *Cpu) ExecSTX(operandAddr memory.Ptr) int {
	log.Printf("cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	log.Printf("Exec STX")
	cpu.Memory.Poke(operandAddr, cpu.X)
	log.Printf("cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	return 1
}

func (cpu *Cpu) ExecTXS(operandAddr memory.Ptr) int {
	log.Printf("Exec TXS")
	cpu.SP = cpu.X
	return 1
}

func (cpu *Cpu) ExecINX(operandAddr memory.Ptr) int {
	log.Printf("Exec INX")
	cpu.X++
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}

func (cpu *Cpu) ExecINY(operandAddr memory.Ptr) int {
	log.Printf("Exec INY")
	cpu.Y++
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}

func (cpu *Cpu) ExecBIT(operandAddr memory.Ptr) int {
	log.Printf("Exec BIT")

	operand := cpu.Memory.Peek(operandAddr)
	result := cpu.A & operand

	cpu.P.Set(PFLAG_Z, result == 0)
	cpu.P.Set(PFLAG_V, operand&0x40 != 0)
	cpu.P.Set(PFLAG_N, operand >= 128)
	return 1
}

func (cpu *Cpu) ExecBPL(operandAddr memory.Ptr) int {
	log.Printf("Exec BPL")
	if cpu.P&PFLAG_N == 0 {
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return 1
}
