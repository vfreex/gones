package cpu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
)

var opcodeHandlers = [256]*InstructionHandler{
	0x10: {(*Cpu).ExecBPL, REL},
	0x30: {(*Cpu).ExecBMI, REL},
	0x50: {(*Cpu).ExecBVC, REL},
	0x70: {(*Cpu).ExecBVS, REL},
	0x90: {(*Cpu).ExecBCC, REL},
	0xB0: {(*Cpu).ExecBCS, REL},
	0xD0: {(*Cpu).ExecBNE, REL},
	0xF0: {(*Cpu).ExecBEQ, REL},

	0x24: {(*Cpu).ExecBIT, ZP},
	0x2c: {(*Cpu).ExecBIT, ABS},

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
	0x61: {(*Cpu).ExecADC, IZX},
	0x71: {(*Cpu).ExecADC, IZY},
	0x6d: {(*Cpu).ExecADC, ABS},
	0x7d: {(*Cpu).ExecADC, ABX},
	0x79: {(*Cpu).ExecADC, ABY},

	0xe9: {(*Cpu).ExecSBC, IMM},
	0xe5: {(*Cpu).ExecSBC, ZP},
	0xf5: {(*Cpu).ExecSBC, ZPX},
	0xe1: {(*Cpu).ExecSBC, IZX},
	0xf1: {(*Cpu).ExecSBC, IZY},
	0xed: {(*Cpu).ExecSBC, ABS},
	0xfd: {(*Cpu).ExecSBC, ABX},
	0xf9: {(*Cpu).ExecSBC, ABY},

	0x09: {(*Cpu).ExecORA, IMM},
	0x05: {(*Cpu).ExecORA, ZP},
	0x15: {(*Cpu).ExecORA, ZPX},
	0x01: {(*Cpu).ExecORA, IZX},
	0x11: {(*Cpu).ExecORA, IZY},
	0x0d: {(*Cpu).ExecORA, ABS},
	0x1d: {(*Cpu).ExecORA, ABX},
	0x19: {(*Cpu).ExecORA, ABY},

	0x29: {(*Cpu).ExecAND, IMM},
	0x25: {(*Cpu).ExecAND, ZP},
	0x35: {(*Cpu).ExecAND, ZPX},
	0x21: {(*Cpu).ExecAND, IZX},
	0x31: {(*Cpu).ExecAND, IZY},
	0x2d: {(*Cpu).ExecAND, ABS},
	0x3d: {(*Cpu).ExecAND, ABX},
	0x39: {(*Cpu).ExecAND, ABY},

	0x49: {(*Cpu).ExecEOR, IMM},
	0x45: {(*Cpu).ExecEOR, ZP},
	0x55: {(*Cpu).ExecEOR, ZPX},
	0x41: {(*Cpu).ExecEOR, IZX},
	0x51: {(*Cpu).ExecEOR, IZY},
	0x4d: {(*Cpu).ExecEOR, ABS},
	0x5d: {(*Cpu).ExecEOR, ABX},
	0x59: {(*Cpu).ExecEOR, ABY},

	0xc9: {(*Cpu).ExecCMP, IMM},
	0xc5: {(*Cpu).ExecCMP, ZP},
	0xd5: {(*Cpu).ExecCMP, ZPX},
	0xc1: {(*Cpu).ExecCMP, IZX},
	0xd1: {(*Cpu).ExecCMP, IZY},
	0xcd: {(*Cpu).ExecCMP, ABS},
	0xdd: {(*Cpu).ExecCMP, ABX},
	0xd9: {(*Cpu).ExecCMP, ABY},

	0xe0: {(*Cpu).ExecCPX, IMM},
	0xe4: {(*Cpu).ExecCPX, ZP},
	0xec: {(*Cpu).ExecCPX, ABS},

	0xc0: {(*Cpu).ExecCPY, IMM},
	0xc4: {(*Cpu).ExecCPY, ZP},
	0xcc: {(*Cpu).ExecCPY, ABS},

	0xe8: {(*Cpu).ExecINX, IMP},
	0xc8: {(*Cpu).ExecINY, IMP},

	0xca: {(*Cpu).ExecDEX, IMP},
	0x88: {(*Cpu).ExecDEY, IMP},

	0xe6: {(*Cpu).ExecINC, ZP},
	0xf6: {(*Cpu).ExecINC, ZPX},
	0xee: {(*Cpu).ExecINC, ABS},
	0xfe: {(*Cpu).ExecINC, ABX},

	0xc6: {(*Cpu).ExecDEC, ZP},
	0xd6: {(*Cpu).ExecDEC, ZPX},
	0xce: {(*Cpu).ExecDEC, ABS},
	0xde: {(*Cpu).ExecDEC, ABX},

	0x0a: {(*Cpu).ExecASLA, IMP},
	0x06: {(*Cpu).ExecASL, ZP},
	0x16: {(*Cpu).ExecASL, ZPX},
	0x0e: {(*Cpu).ExecASL, ABS},
	0x1e: {(*Cpu).ExecASL, ABX},

	0x2a: {(*Cpu).ExecROLA, IMP},
	0x26: {(*Cpu).ExecROL, ZP},
	0x36: {(*Cpu).ExecROL, ZPX},
	0x2e: {(*Cpu).ExecROL, ABS},
	0x3e: {(*Cpu).ExecROL, ABX},

	0x4a: {(*Cpu).ExecLSRA, IMP},
	0x46: {(*Cpu).ExecLSR, ZP},
	0x56: {(*Cpu).ExecLSR, ZPX},
	0x4e: {(*Cpu).ExecLSR, ABS},
	0x5e: {(*Cpu).ExecLSR, ABX},

	0x6a: {(*Cpu).ExecRORA, IMP},
	0x66: {(*Cpu).ExecROR, ZP},
	0x76: {(*Cpu).ExecROR, ZPX},
	0x6e: {(*Cpu).ExecROR, ABS},
	0x7e: {(*Cpu).ExecROR, ABX},

	0x00: {(*Cpu).ExecBRK, IMP},
	0x40: {(*Cpu).ExecRTI, IMP},
	0x20: {(*Cpu).ExecJSR, ABS},
	0x60: {(*Cpu).ExecRTS, IMP},
	0x4c: {(*Cpu).ExecJMP, ABS},
	0x6c: {(*Cpu).ExecJMP, IND},


	0x18: {(*Cpu).ExecCLC, IMP},
	0x38: {(*Cpu).ExecSEC, IMP},
	0xd8: {(*Cpu).ExecCLD, IMP},
	0xf8: {(*Cpu).ExecSED, IMP},
	0x58: {(*Cpu).ExecCLI, IMP},
	0x78: {(*Cpu).ExecSEI, IMP},
	0xb8: {(*Cpu).ExecCLV, IMP},
	0xea: {(*Cpu).ExecNOP, IMP},
}

type InstructionExecutor func(cpu *Cpu, operandAddr memory.Ptr) (cyclesTook int)

type InstructionHandler struct {
	Executor       InstructionExecutor
	AddressingMode AddressingMode
}

func (cpu *Cpu) ExecLDA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec LDA")
	cpu.A = cpu.Memory.Peek(operandAddr)
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 1
}

func (cpu *Cpu) ExecSTA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec STA")
	cpu.Memory.Poke(operandAddr, cpu.A)
	return 1
}

func (cpu *Cpu) ExecLDX(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec LDX")
	cpu.X = cpu.Memory.Peek(operandAddr)
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}

func (cpu *Cpu) ExecSTX(operandAddr memory.Ptr) int {
	//logger.Debug(";; cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	logger.Debug(";; Exec STX")
	cpu.Memory.Poke(operandAddr, cpu.X)
	//logger.Debug(";; cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	return 1
}

func (cpu *Cpu) ExecLDY(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec LDY")
	cpu.Y = cpu.Memory.Peek(operandAddr)
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}

func (cpu *Cpu) ExecSTY(operandAddr memory.Ptr) int {
	//logger.Debug(";; cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	logger.Debug(";; Exec STY")
	cpu.Memory.Poke(operandAddr, cpu.Y)
	//logger.Debug(";; cpu memory %04x: %02x", operandAddr, cpu.Memory.Peek(operandAddr))
	return 1
}

func (cpu *Cpu) ExecTAX(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec TAX")
	cpu.X = cpu.A
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}
func (cpu *Cpu) ExecTXA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec TXA")
	cpu.A = cpu.X
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 1
}
func (cpu *Cpu) ExecTAY(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec TAY")
	cpu.Y = cpu.A
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}
func (cpu *Cpu) ExecTYA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec TYA")
	cpu.A = cpu.Y
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 1
}
func (cpu *Cpu) ExecTSX(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec TSX")
	cpu.X = cpu.SP
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}
func (cpu *Cpu) ExecTXS(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec TXS")
	cpu.SP = cpu.X
	// don't set P flags
	return 1
}

func (cpu *Cpu) ExecINX(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec INX")
	cpu.X++
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}

func (cpu *Cpu) ExecINY(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec INY")
	cpu.Y++
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}

func (cpu *Cpu) ExecDEX(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec DEX")
	cpu.X--
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}

func (cpu *Cpu) ExecDEY(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec DEY")
	cpu.Y--
	cpu.P.Set(PFLAG_Z, cpu.Y == 0)
	cpu.P.Set(PFLAG_N, cpu.Y >= 128)
	return 1
}

func (cpu *Cpu) ExecBIT(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BIT")

	operand := cpu.Memory.Peek(operandAddr)
	result := cpu.A & operand

	cpu.P.Set(PFLAG_Z, result == 0)
	cpu.P.Set(PFLAG_V, operand&0x40 != 0)
	cpu.P.Set(PFLAG_N, operand >= 128)
	return 1
}

func (cpu *Cpu) ExecBPL(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BPL")
	cycles := 1
	if cpu.P&PFLAG_N == 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBMI(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BMI")
	cycles := 1
	if cpu.P&PFLAG_N != 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBVC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BVC")
	cycles := 1
	if cpu.P&PFLAG_V == 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBVS(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BVS")
	cycles := 1
	if cpu.P&PFLAG_V != 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBCC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BCC")
	cycles := 1
	if cpu.P&PFLAG_C == 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBCS(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BCS")
	cycles := 1
	if cpu.P&PFLAG_C != 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBNE(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BNE")
	cycles := 1
	if cpu.P&PFLAG_Z == 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBEQ(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BEQ")
	cycles := 1
	if cpu.P&PFLAG_Z != 0 {
		cycles++
		logger.Debugf(";; before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		logger.Debugf(";; jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecPLA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec PLA")
	cpu.A = cpu.Pop()
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A >= 128)
	return 3
}

func (cpu *Cpu) ExecPHA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec PHA")
	cpu.Push(cpu.A)
	return 2
}

func (cpu *Cpu) ExecPLP(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec PLP")
	cpu.P = ProcessorStatus(cpu.Pop())
	return 3
}

func (cpu *Cpu) ExecPHP(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec PHP")
	cpu.Push(byte(cpu.P | PFLAG_B))
	return 2
}

func (cpu *Cpu) ExecADC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ADC")
	operand := cpu.Memory.Peek(operandAddr)
	r := uint16(cpu.A) + uint16(operand)
	if cpu.P&PFLAG_C != 0 {
		r++
	}
	r2 := uint8(r)
	cpu.P.Set(PFLAG_C, r > 0xFF)
	// http://www.6502.org/tutorials/vflag.html#2.4
	cpu.P.Set(PFLAG_V, (cpu.A^operand)&0x80 == 0 && (cpu.A^r2)&0x80 != 0)
	cpu.P.Set(PFLAG_Z, r2 == 0)
	cpu.P.Set(PFLAG_N, r2 > 0x7f)
	cpu.A = r2
	return 1
}

func (cpu *Cpu) ExecSBC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec SBC")
	operand := cpu.Memory.Peek(operandAddr)
	operand2 := ^operand
	r := uint16(cpu.A) + uint16(operand2)
	if cpu.P&PFLAG_C != 0 {
		r++
	}
	r2 := uint8(r)
	cpu.P.Set(PFLAG_C, r > 0xFF)
	// http://www.6502.org/tutorials/vflag.html#2.4
	cpu.P.Set(PFLAG_V, (cpu.A^operand2)&0x80 == 0 && (cpu.A^r2)&0x80 != 0)
	cpu.P.Set(PFLAG_Z, r2 == 0)
	cpu.P.Set(PFLAG_N, r2 > 0x7f)
	cpu.A = r2
	return 1
}

func (cpu *Cpu) ExecORA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ORA")
	operand := cpu.Memory.Peek(operandAddr)
	cpu.A |= operand
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecAND(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec AND")
	operand := cpu.Memory.Peek(operandAddr)
	cpu.A &= operand
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecEOR(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec EOR")
	operand := cpu.Memory.Peek(operandAddr)
	cpu.A ^= operand
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecCMP(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CMP")
	operand := cpu.Memory.Peek(operandAddr)
	r := cpu.A - operand
	cpu.P.Set(PFLAG_C, cpu.A >= operand)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 1
}

func (cpu *Cpu) ExecCPX(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CPX")
	operand := cpu.Memory.Peek(operandAddr)
	r := cpu.X - operand
	cpu.P.Set(PFLAG_C, cpu.X >= operand)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 1
}

func (cpu *Cpu) ExecCPY(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CPY")
	operand := cpu.Memory.Peek(operandAddr)
	r := cpu.Y - operand
	cpu.P.Set(PFLAG_C, cpu.Y >= operand)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 1
}

func (cpu *Cpu) ExecINC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec INC")
	r := cpu.Memory.Peek(operandAddr) + 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecDEC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec DEC")
	r := cpu.Memory.Peek(operandAddr) - 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecASLA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ASLA")
	cpu.P.Set(PFLAG_C, cpu.A > 0x7f)
	cpu.A <<= 1
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecASL(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ASL")
	operand := cpu.Memory.Peek(operandAddr)
	r := operand << 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_C, operand > 0x7f)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecROLA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ROLA")
	cf := cpu.P&PFLAG_C != 0
	cpu.P.Set(PFLAG_C, cpu.A > 0x7f)
	cpu.A <<= 1
	if cf {
		cpu.A |= 1
	}
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecROL(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ROL")
	operand := cpu.Memory.Peek(operandAddr)
	r := operand << 1
	if cpu.P&PFLAG_C != 0 {
		r |= 1
	}
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_C, operand > 0x7f)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecLSRA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec LSRA")
	cpu.P.Set(PFLAG_C, cpu.A&1 != 0)
	cpu.A >>= 1
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, false)
	return 1
}

func (cpu *Cpu) ExecLSR(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec LSR")
	operand := cpu.Memory.Peek(operandAddr)
	r := operand >> 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_C, operand&1 != 0)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, false)
	return 3
}

func (cpu *Cpu) ExecRORA(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec RORA")
	cf := cpu.P&PFLAG_C != 0
	cpu.P.Set(PFLAG_C, cpu.A&1 != 0)
	cpu.A >>= 1
	if cf {
		cpu.A |= 0x80
	}
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cf)
	return 1
}

func (cpu *Cpu) ExecROR(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec ROR")
	operand := cpu.Memory.Peek(operandAddr)
	r := operand >> 1
	if cpu.P&PFLAG_C != 0 {
		r |= 0x80
	}
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_C, operand&1 != 0)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecBRK(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec BRK")
	cpu.P.Set(PFLAG_B, true)
	cpu.PushW(cpu.PC + 1)
	cpu.Push(byte(cpu.P))
	cpu.P.Set(PFLAG_I, true)
	cpu.PC = cpu.ReadInterruptVector(IV_BRK)
	return 6
}

func (cpu *Cpu) ExecRTI(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec RTI")
	cpu.P = ProcessorStatus(cpu.Pop())
	cpu.PC = cpu.PopW()
	return 5
}

func (cpu *Cpu) ExecJSR(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec JSR")
	cpu.PushW(cpu.PC - 1)
	logger.Debugf(";; before jump: PC=%2x", cpu.PC)
	cpu.PC = operandAddr
	logger.Debugf(";; jump to PC=%2x", cpu.PC)
	return 5
}

func (cpu *Cpu) ExecRTS(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec RTS")
	logger.Debugf(";; before jump: PC=%2x", cpu.PC)
	cpu.PC = cpu.PopW() + 1
	logger.Debugf(";; jump to PC=%2x", cpu.PC)
	return 5
}

func (cpu *Cpu) ExecJMP(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec JMP")
	logger.Debugf(";; before jump: PC=%2x", cpu.PC)
	cpu.PC = operandAddr
	logger.Debugf(";; jump to PC=%2x", cpu.PC)
	return 0
}

func (cpu *Cpu) ExecCLC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CLC")
	cpu.P.Set(PFLAG_C, false)
	return 1
}

func (cpu *Cpu) ExecSEC(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec SEC")
	cpu.P.Set(PFLAG_C, true)
	return 1
}

func (cpu *Cpu) ExecCLD(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CLD")
	cpu.P.Set(PFLAG_D, false)
	return 1
}

func (cpu *Cpu) ExecSED(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec SED")
	cpu.P.Set(PFLAG_D, true)
	return 1
}

func (cpu *Cpu) ExecCLI(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CLI")
	cpu.P.Set(PFLAG_I, false)
	return 1
}

func (cpu *Cpu) ExecSEI(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec SEI")
	cpu.P.Set(PFLAG_I, true)
	return 1
}

func (cpu *Cpu) ExecCLV(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec CLV")
	cpu.P.Set(PFLAG_V, false)
	return 1
}
func (cpu *Cpu) ExecNOP(operandAddr memory.Ptr) int {
	logger.Debug(";; Exec NOP")
	return 1
}
