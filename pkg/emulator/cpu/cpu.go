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
	//cpu.P = 0x34
	//cpu.A, cpu.X, cpu.Y = 0, 0, 0
	//cpu.SP = 0xfd
	cpu.PC = 0x8000
}

func (cpu *Cpu) Test() {
	//log.Printf("%v\n", cpu)
	//cpu.Reset()
	//log.Printf("%v\n", cpu)

	//pc := memory.Ptr(0x8000)
	//for pc < 0x8200 {
	//	op := cpu.Memory.Peek(pc)
	//	opName, length := Decode(op)
	//	nextPC := pc + 1
	//	var operand []byte = nil
	//	if length > 1 {
	//		operand = make([]byte, length-1)
	//		switch length {
	//		case 2:
	//			operand[0] = cpu.Memory.Peek(nextPC)
	//			nextPC++
	//		case 3:
	//			operand[0] = cpu.Memory.Peek(nextPC)
	//			nextPC++
	//			operand[1] = cpu.Memory.Peek(nextPC)
	//			nextPC++
	//		}
	//	}
	//	log.Printf("will exec IP=%x: OP=%x %s %v \n", pc, op, opName, operand)
	//	cpu.ExecInstruction(op, operand)
	//	pc = nextPC
	//}
	for cpu.PC < 0x8200 {
		//op := cpu.Memory.Peek(cpu.PC)
		//opName, am := Decode(op)
		//operands, _ := cpu.ReadOperands(am)
		//log.Printf("will exec opcode=%02x %s (%s) %v \n", op, opName, am, operands)
		cycles, _ := cpu.ExecOneInstruction()
		log.Printf("spent %d cycles", cycles)
		//cpu.PC += bytes
	}

}

func (cpu *Cpu) ExecOneInstruction() (cycles int, bytes uint16) {
	opcode, _, cycles0 := cpu.ReadNextInstruction()
	info := &InstructionInfos[opcode]
	length := info.AddressingMode.GetArgumentCount() + 1

	cpu.PC++

	operandAddr, cycles1 := cpu.AddressOperand(info.AddressingMode)

	handler := opcodeHandlers[opcode]
	if handler == nil {
		log.Fatalf("opcode %02x (%s) is not supported", opcode, info.Nemonics)
	}
	cpu.logRegisters()
	log.Printf("will exec opcode=%02x %s (%s) %x \n", opcode, info.Nemonics, info.AddressingMode, operandAddr)
	cycles2 := handler(cpu, operandAddr)
	cpu.logRegisters()

	return cycles0 + cycles1 + cycles2, length
}

func (cpu *Cpu) logRegisters() {
	log.Printf("PC=%04x, P=%s, SP=%02x, A=%02x, X=%02x, Y=%02x", cpu.PC, cpu.P, cpu.SP, cpu.A, cpu.X, cpu.Y)
}
