package cpu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"log"
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

func (cpu *Cpu) ExecDEX(operandAddr memory.Ptr) int {
	log.Printf("Exec DEX")
	cpu.X--
	cpu.P.Set(PFLAG_Z, cpu.X == 0)
	cpu.P.Set(PFLAG_N, cpu.X >= 128)
	return 1
}

func (cpu *Cpu) ExecDEY(operandAddr memory.Ptr) int {
	log.Printf("Exec DEY")
	cpu.Y--
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
	cycles := 1
	if cpu.P&PFLAG_N == 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBMI(operandAddr memory.Ptr) int {
	log.Printf("Exec BMI")
	cycles := 1
	if cpu.P&PFLAG_N != 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBVC(operandAddr memory.Ptr) int {
	log.Printf("Exec BVC")
	cycles := 1
	if cpu.P&PFLAG_V == 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBVS(operandAddr memory.Ptr) int {
	log.Printf("Exec BVS")
	cycles := 1
	if cpu.P&PFLAG_V != 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBCC(operandAddr memory.Ptr) int {
	log.Printf("Exec BCC")
	cycles := 1
	if cpu.P&PFLAG_C == 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBCS(operandAddr memory.Ptr) int {
	log.Printf("Exec BCS")
	cycles := 1
	if cpu.P&PFLAG_C != 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBNE(operandAddr memory.Ptr) int {
	log.Printf("Exec BNE")
	cycles := 1
	if cpu.P&PFLAG_Z == 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
}

func (cpu *Cpu) ExecBEQ(operandAddr memory.Ptr) int {
	log.Printf("Exec BEQ")
	cycles := 1
	if cpu.P&PFLAG_Z != 0 {
		cycles++
		log.Printf("before jump: PC=%2x", cpu.PC)
		cpu.PC = operandAddr
		log.Printf("jump to PC=%2x", operandAddr)
	}
	return cycles
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
	// http://www.6502.org/tutorials/vflag.html#2.4
	cpu.P.Set(PFLAG_V, (cpu.A^operand)&0x80 == 0 && (cpu.A^r2)&0x80 != 0)
	cpu.P.Set(PFLAG_Z, r2 == 0)
	cpu.P.Set(PFLAG_N, r2 > 0x7f)
	cpu.A = r2
	return 1
}

func (cpu *Cpu) ExecSBC(operandAddr memory.Ptr) int {
	log.Printf("Exec SBC")
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
	log.Printf("Exec ORA")
	operand := cpu.Memory.Peek(operandAddr)
	cpu.A |= operand
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecAND(operandAddr memory.Ptr) int {
	log.Printf("Exec AND")
	operand := cpu.Memory.Peek(operandAddr)
	cpu.A &= operand
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecEOR(operandAddr memory.Ptr) int {
	log.Printf("Exec EOR")
	operand := cpu.Memory.Peek(operandAddr)
	cpu.A ^= operand
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecCMP(operandAddr memory.Ptr) int {
	log.Printf("Exec CMP")
	operand := cpu.Memory.Peek(operandAddr)
	r := cpu.A - operand
	cpu.P.Set(PFLAG_C, cpu.A >= operand)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 1
}

func (cpu *Cpu) ExecCPX(operandAddr memory.Ptr) int {
	log.Printf("Exec CPX")
	operand := cpu.Memory.Peek(operandAddr)
	r := cpu.X - operand
	cpu.P.Set(PFLAG_C, cpu.X >= operand)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 1
}

func (cpu *Cpu) ExecCPY(operandAddr memory.Ptr) int {
	log.Printf("Exec CPY")
	operand := cpu.Memory.Peek(operandAddr)
	r := cpu.Y - operand
	cpu.P.Set(PFLAG_C, cpu.Y >= operand)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 1
}

func (cpu *Cpu) ExecINC(operandAddr memory.Ptr) int {
	log.Printf("Exec INC")
	r := cpu.Memory.Peek(operandAddr) + 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecDEC(operandAddr memory.Ptr) int {
	log.Printf("Exec DEC")
	r := cpu.Memory.Peek(operandAddr) - 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecASLA(operandAddr memory.Ptr) int {
	log.Printf("Exec ASLA")
	cpu.P.Set(PFLAG_C, cpu.A > 0x7f)
	cpu.A <<= 1
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, cpu.A > 0x7f)
	return 1
}

func (cpu *Cpu) ExecASL(operandAddr memory.Ptr) int {
	log.Printf("Exec ASL")
	operand := cpu.Memory.Peek(operandAddr)
	r := operand << 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_C, operand > 0x7f)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, r > 0x7f)
	return 3
}

func (cpu *Cpu) ExecROLA(operandAddr memory.Ptr) int {
	log.Printf("Exec ROLA")
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
	log.Printf("Exec ROL")
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
	log.Printf("Exec LSRA")
	cpu.P.Set(PFLAG_C, cpu.A&1 != 0)
	cpu.A >>= 1
	cpu.P.Set(PFLAG_Z, cpu.A == 0)
	cpu.P.Set(PFLAG_N, false)
	return 1
}

func (cpu *Cpu) ExecLSR(operandAddr memory.Ptr) int {
	log.Printf("Exec LSR")
	operand := cpu.Memory.Peek(operandAddr)
	r := operand >> 1
	cpu.Memory.Poke(operandAddr, r)
	cpu.P.Set(PFLAG_C, operand&1 != 0)
	cpu.P.Set(PFLAG_Z, r == 0)
	cpu.P.Set(PFLAG_N, false)
	return 3
}

func (cpu *Cpu) ExecRORA(operandAddr memory.Ptr) int {
	log.Printf("Exec RORA")
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
	log.Printf("Exec ROR")
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
