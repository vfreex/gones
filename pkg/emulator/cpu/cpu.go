package cpu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"log"
)

const (
	SP_BASE uint16 = 0x100
)

const (
	INT_VEC_NMI     uint16 = 0xFFFA // 7 cycles
	INT_VEC_RESET   uint16 = 0xFFFC //
	INT_VEC_IRQ_BRK uint16 = 0xFFFE // 7 cycles
)

type Cpu struct {
	// registers
	PC   ProgramCounter
	P    ProcessorStatus
	SP   StackPointer
	A    Accumulator
	X, Y IndexRegister
	// memory
	Memory memory.Memory
}

func NewCpu(memory memory.Memory) *Cpu {
	cpu := &Cpu{Memory: memory}
	cpu.Init()
	return cpu
}

func (cpu *Cpu) Init() {
	/*
		P = $34[1] (IRQ disabled)[2]
		A, X, Y = 0
		S = $FD
		$4017 = $00 (frame irq enabled)
		$4015 = $00 (all channels disabled)
		$4000-$400F = $00 (not sure about $4010-$4013)
		All 15 bits of noise channel LFSR = $0000[3]. The first time the LFSR is clocked from the all-0s state, it will shift in a 1.
		Internal memory ($0000-$07FF) has unreliable startup state. Some machines may have consistent RAM contents at power-on, but others do not.
		Emulators often implement a consistent RAM startup state (e.g. all $00 or $FF, or a particular pattern), and flash carts like the PowerPak may partially or fully initialize RAM before starting a program, so an NES programmer must be careful not to rely on the startup contents of RAM.
	*/
	cpu.P = 0x34
	cpu.A, cpu.X, cpu.Y = 0, 0, 0
	cpu.SP = 0xfd
	cpu.PC = 0x8000
}

func (cpu *Cpu) Test() {
	//log.Printf("%v\n", cpu)
	//cpu.Reset()
	//log.Printf("%v\n", cpu)

	pc := memory.Ptr(0x8000)
	for pc < 0x8200 {
		op := cpu.Memory.Peek(pc)
		opName, length := Decode(op)
		nextPC := pc + 1
		var operand []byte = nil
		if length > 1 {
			operand = make([]byte, length-1)
			switch length {
			case 2:
				operand[0] = cpu.Memory.Peek(nextPC)
				nextPC++
			case 3:
				operand[0] = cpu.Memory.Peek(nextPC)
				nextPC++
				operand[1] = cpu.Memory.Peek(nextPC)
				nextPC++
			}
		}
		log.Printf("IP=%x: OP=%x %s %v \n", pc, op, opName, operand)
		cpu.ExecInstruction(op, operand)
		pc = nextPC
	}

}

func (cpu *Cpu) ExecInstruction(opcode byte, operands []byte) {
	info := &InstructionInfos[opcode]
	length := uint16(0)
	switch info.AddressingMode {
	case IMP:
		length = 1
	case IND, ABS, ABX, ABY:
		length = 3
	default:
		length = 2
	}
	log.Printf("EXEC: PC=%04x, %02x %v - %s %v", cpu.PC, opcode, operands, info.Nemonics, operands)
	cpu.PC += length
}

//func (cpu *Cpu) Push(b byte) {
//	cpu.Memory[SP_BASE|uint16(cpu.SP)] = b
//	cpu.SP--
//}
//
//func (cpu *Cpu) Pop() byte {
//	cpu.SP++
//	return cpu.Memory[SP_BASE|uint16(cpu.SP)]
//}

//func (cpu *Cpu) Dump() {
//	p := 0x8000
//	log.Printf("OP: %x \n", cpu.Memory[p])
//}
