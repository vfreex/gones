package cpu

import (
	"encoding/hex"
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

const (
	SP_BASE uint16 = 0x100
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

var logger = logger2.GetLogger()

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
}

func (cpu *Cpu) Test() {
	for cpu.PC < 0x810f {
		cycles := cpu.ExecOneInstruction()
		logger.Infof("spent %d cycles", cycles)
	}

}

func (cpu *Cpu) Push(b byte) {
	cpu.Memory.Poke(0x100|memory.Ptr(cpu.SP), b)
	cpu.SP--
}

func (cpu *Cpu) PushW(w uint16) {
	cpu.Push(byte(w >> 8))
	cpu.Push(byte(w & 0xff))
}

func (cpu *Cpu) Pop() byte {
	cpu.SP++
	return cpu.Memory.Peek(0x100 | memory.Ptr(cpu.SP))
}

func (cpu *Cpu) PopW() uint16 {
	low := cpu.Pop()
	high := cpu.Pop()
	return uint16(high)<<8 | uint16(low)
}

func (cpu *Cpu) ExecOneInstruction() (cycles int) {
	opcode := cpu.Memory.Peek(cpu.PC)
	info := &InstructionInfos[opcode]

	arguments := make([]byte, info.AddressingMode.GetArgumentCount())
	switch info.AddressingMode.GetArgumentCount() {
	case 2:
		arguments[1] = cpu.Memory.Peek(cpu.PC + 2)
		fallthrough
	case 1:
		arguments[0] = cpu.Memory.Peek(cpu.PC + 1)
	}
	logger.Debugf("got instruction at %04x: %02x(%s %s) %s",
		cpu.PC, opcode, info.Nemonics, info.AddressingMode, hex.EncodeToString(arguments))
	handler := opcodeHandlers[opcode]
	if handler == nil {
		//logger.Fatalf("opcode %02x (%s) is not supported", opcode, info.Nemonics)
		logger.Fatalf("opcode %02x is not supported", opcode)
	}

	cpu.PC++
	operandAddr, cycles1 := cpu.AddressOperand(handler.AddressingMode)
	cpu.logRegisters()
	logger.Debugf("will exec opcode=%02x %s (%s) %x \n", opcode, info.Nemonics, handler.AddressingMode, operandAddr)
	cycles2 := handler.Executor(cpu, operandAddr)
	cpu.logRegisters()

	return 1 + cycles1 + cycles2
}

func (cpu *Cpu) logRegisters() {
	logger.Debugf("PC=%04x, P=%s, SP=%02x, A=%02x, X=%02x, Y=%02x", cpu.PC, cpu.P, cpu.SP, cpu.A, cpu.X, cpu.Y)
}
