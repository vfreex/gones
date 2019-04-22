package cpu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"log"
)

var opcodeHandlers = [256]*InstructionHandler{
	0x10: {(*Cpu).ExecBPL, REL},

	0x24: {(*Cpu).ExecBIT, ZP},
	0x2c: {(*Cpu).ExecBIT, ABS},

	0x78: {(*Cpu).ExecSEI, IMP},
	0xd8: {(*Cpu).ExecCLD, IMP},

	0xaa: {(*Cpu).ExecTAX, IMP},
	0x8a: {(*Cpu).ExecTXA, IMP},
	0xa8: {(*Cpu).ExecTAY, IMP},
	0x98: {(*Cpu).ExecTYA, IMP},
	0xba: {(*Cpu).ExecTSX, IMP},
	0x9a: {(*Cpu).ExecTXS, IMP},


	0xa2: {(*Cpu).ExecLDX, IMM},
	0xa6: {(*Cpu).ExecLDX, ZP},
	0xb6: {(*Cpu).ExecLDX, ZPY},
	0xae: {(*Cpu).ExecLDX, ABS},
	0xbe: {(*Cpu).ExecLDX, ABY},

	0x86: {(*Cpu).ExecSTX, ZP},
	0x96: {(*Cpu).ExecSTX, ZPY},
	0x8e: {(*Cpu).ExecSTX, ABS},

	0xa0: {(*Cpu).ExecLDY, IMM},
	0xa4: {(*Cpu).ExecLDY, ZP},
	0xb4: {(*Cpu).ExecLDY, ZPX},
	0xac: {(*Cpu).ExecLDY, ABS},
	0xbc: {(*Cpu).ExecLDY, ABX},

	0x84: {(*Cpu).ExecSTY, ZP},
	0x94: {(*Cpu).ExecSTY, ZPX},
	0x8c: {(*Cpu).ExecSTY, ABS},

	0xc8: {(*Cpu).ExecINY, IMP},
	0xe8: {(*Cpu).ExecINX, IMP},

	0xa9: {(*Cpu).ExecLDA, IMM},
	0xa5: {(*Cpu).ExecLDA, ZP},
	0xb5: {(*Cpu).ExecLDA, ZPX},
	0xad: {(*Cpu).ExecLDA, ABS},
	0xbd: {(*Cpu).ExecLDA, ABX},
	0xb9: {(*Cpu).ExecLDA, ABY},
	0xa1: {(*Cpu).ExecLDA, IZX},
	0xb1: {(*Cpu).ExecLDA, IZY},

	0x85: {(*Cpu).ExecSTA, ZP},
	0x95: {(*Cpu).ExecSTA, ZPX},
	0x81: {(*Cpu).ExecSTA, IZX},
	0x91: {(*Cpu).ExecSTA, IZY},
	0x8d: {(*Cpu).ExecSTA, ABS},
	0x9d: {(*Cpu).ExecSTA, ABX},
	0x99: {(*Cpu).ExecSTA, ABY},

	0x68: {(*Cpu).ExecPLA, IMP},
	0x48: {(*Cpu).ExecPHA, IMP},
	0x28: {(*Cpu).ExecPLP, IMP},
	0x08: {(*Cpu).ExecPHP, IMP},

	0x69: {(*Cpu).ExecADC, IMM},
	0x65: {(*Cpu).ExecADC, ZP},
	0x75: {(*Cpu).ExecADC, ZPX},
	0x6d: {(*Cpu).ExecADC, ABS},
	0x7d: {(*Cpu).ExecADC, ABX},
	0x79: {(*Cpu).ExecADC, ABY},
	0x61: {(*Cpu).ExecADC, IZX},
	0x71: {(*Cpu).ExecADC, IZY},
}

type InstructionExecutor func(cpu *Cpu, operandAddr memory.Ptr) (cyclesTook int)

type InstructionHandler struct {
	Executor       InstructionExecutor
	AddressingMode AddressingMode
}

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

func (cpu *Cpu) ExecSTA(operandAddr memory.Ptr) int {
	log.Printf("Exec STA")
	cpu.Memory.Poke(operandAddr, cpu.A)
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

func (cpu *Cpu) ExecLDY(operandAddr memory.Ptr) int {
	log.Printf("Exec LDY")
	cpu.Y = cpu.Memory.Peek(operandAddr)
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}

func (cpu *Cpu) ExecSTY(operandAddr memory.Ptr) int {
	log.Printf("cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	log.Printf("Exec STY")
	cpu.Memory.Poke(operandAddr, cpu.Y)
	log.Printf("cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	return 1
}

func (cpu *Cpu) ExecTAX(operandAddr memory.Ptr) int {
	log.Printf("Exec TAX")
	cpu.X = cpu.A
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}
func (cpu *Cpu) ExecTXA(operandAddr memory.Ptr) int {
	log.Printf("Exec TXA")
	cpu.A = cpu.X
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 1
}
func (cpu *Cpu) ExecTAY(operandAddr memory.Ptr) int {
	log.Printf("Exec TAY")
	cpu.Y = cpu.A
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}
func (cpu *Cpu) ExecTYA(operandAddr memory.Ptr) int {
	log.Printf("Exec TYA")
	cpu.A = cpu.Y
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 1
}
func (cpu *Cpu) ExecTSX(operandAddr memory.Ptr) int {
	log.Printf("Exec TSX")
	cpu.X = cpu.SP
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}
func (cpu *Cpu) ExecTXS(operandAddr memory.Ptr) int {
	log.Printf("Exec TXS")
	cpu.SP = cpu.X
	// don't set P flags
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

func (cpu *Cpu) ExecPLA(operandAddr memory.Ptr) int {
	log.Printf("Exec PLA")
	cpu.SP++
	cpu.A = cpu.Memory.Peek(0x100 | memory.Ptr(cpu.SP))
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 3
}

func (cpu *Cpu) ExecPHA(operandAddr memory.Ptr) int {
	log.Printf("Exec PHA")
	cpu.Memory.Poke(0x100|memory.Ptr(cpu.SP), cpu.A)
	cpu.SP--
	return 2
}

func (cpu *Cpu) ExecPLP(operandAddr memory.Ptr) int {
	log.Printf("Exec PLP")
	cpu.SP++
	cpu.P = ProcessorStatus(cpu.Memory.Peek(0x100 | memory.Ptr(cpu.SP)))
	return 3
}

func (cpu *Cpu) ExecPHP(operandAddr memory.Ptr) int {
	log.Printf("Exec PHP")
	cpu.Memory.Poke(0x100|memory.Ptr(cpu.SP), byte(cpu.P))
	cpu.SP--
	return 2
}

func (cpu *Cpu) ExecADC(operandAddr memory.Ptr) int {
	log.Printf("Exec ADC")
	operand := cpu.Memory.Peek(operandAddr)
	r := uint16(cpu.A) + uint16(operand)
	if cpu.P&PFLAG_C != 0 {
		r++
	}
	r2 := uint8(r)
	cpu.P.Set(PFLAG_C, r > 0xFF)
	// https://en.wikipedia.org/wiki/Overflow_flag
	cpu.P.Set(PFLAG_V, (cpu.A^operand)&0x80 == 0 && (cpu.A^r2)&0x80 != 0)
	cpu.P.Set(PFLAG_Z, r2 == 0)
	cpu.P.Set(PFLAG_N, r2 > 0x7f)
	cpu.A = r2
	return 1
}
