package cpu

import (
	"encoding/hex"
	"fmt"
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
	// interrupts
	NMI bool
	IRQ bool
	// waitCycles
	Wait int
}

var logger = logger2.GetLogger()

func NewCpu(memory memory.Memory) *Cpu {
	cpu := &Cpu{Memory: memory}
	return cpu
}

func (cpu *Cpu) PowerUp() {
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
	cpu.PC = cpu.ReadInterruptVector(IV_RESET)
	logger.Debugf("entrypoint: PC=$%4x", cpu.PC)
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
	if cpu.NMI {
		cpu.ExecNMI()
	} else if cpu.IRQ && cpu.P&PFLAG_I == 0 {
		cpu.ExecIRQ()
	}
	cpu.logInstruction()
	opcode := cpu.Memory.Peek(cpu.PC)
	handler := opcodeHandlers[opcode]
	if handler == nil {
		logger.Fatalf("opcode %02x is not supported", opcode)
	}

	cpu.PC++
	operandAddr, cycles1 := cpu.AddressOperand(handler.AddressingMode)
	cpu.logRegisters()
	cycles2 := handler.Executor(cpu, operandAddr)
	cpu.logRegisters()

	wait := cpu.Wait
	cpu.Wait = 0
	return 1 + cycles1 + cycles2 + wait
}

func (cpu *Cpu) logInstruction() {
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
	logger.Debugf("L%04x: %s %s ; %02x (%s-%d) %s",
		cpu.PC, info.Nemonics, formatInstructionArgument(info.AddressingMode, arguments),
		opcode, info.Nemonics, info.AddressingMode, hex.EncodeToString(arguments))
}

func formatInstructionArgument(am AddressingMode, args []byte) string {
	r := ""
	switch am {
	case IMM:
		r = fmt.Sprintf("#$%x", args[0])
	case ZP:
		r = fmt.Sprintf("$%02x", args[0])
	case ZPX:
		r = fmt.Sprintf("$%02x,X", args[0])
	case ZPY:
		r = fmt.Sprintf("$%02x,Y", args[0])
	case REL:
		r = fmt.Sprintf("*$%+x", int8(args[0]))
	case ABS:
		r = fmt.Sprintf("$%02x%02x", args[1], args[0])
	case ABX:
		r = fmt.Sprintf("$%02x%02x,X", args[1], args[0])
	case ABY:
		r = fmt.Sprintf("$%02x%02x,Y", args[1], args[0])
	case IND:
		r = fmt.Sprintf("($%02x%02x)", args[1], args[0])
	case IZX:
		r = fmt.Sprintf("($%02x,X)", args[0])
	case IZY:
		r = fmt.Sprintf("($%02x),Y", args[0])
	}
	return r
}

func (cpu *Cpu) logRegisters() {
	logger.Debugf(";; PC=%04x, P=%s, SP=%02x, A=%02x, X=%02x, Y=%02x", cpu.PC, cpu.P, cpu.SP, cpu.A, cpu.X, cpu.Y)
}
