package cpu

import (
	"fmt"
)

/*
This register points the address from which the next instruction
          byte (opcode or parameter) will be fetched. Unlike other
          registers, this one is 16 bits in length. The low and high 8-bit
          halves of the register are called PCL and PCH, respectively. The
          Program Counter may be read by pushing its value on the stack.
          This can be done either by jumping to a subroutine or by causing
          an interrupt.
*/
type ProgramCounter = uint16

/*
The NMOS 65xx processors have 256 bytes of stack memory, ranging
          from $0100 to $01FF. The S register is a 8-bit offset to the stack
          page. In other words, whenever anything is being pushed on the
          stack, it will be stored to the address $0100+S.

          The Stack pointer can be read and written by transfering its value
          to or from the index register X (see below) with the TSX and TXS
          instructions.
*/
type StackPointer = uint8

/*
This 8-bit register stores the state of the processor. The bits in
          this register are called flags. Most of the flags have something
          to do with arithmetic operations.
*/
type ProcessorStatus uint8

const (
	// Carry flag
	PFLAG_C ProcessorStatus = 1 << iota
	// Zero flag
	PFLAG_Z
	// Interrupt disable flag
	PFLAG_I
	// Decimal mode flag.
	// 2A03 does not support BCD mode so although the flag can be set, its value will be ignored.
	PFLAG_D
	// Break flag
	PFLAG_B
	// unused flag, alway 1
	PFLAG_UNUSED
	// oVerflow flag
	PFLAG_V
	// Negative flag
	PFLAG_N
)

/*
The accumulator is the main register for arithmetic and logic
          operations. Unlike the index registers X and Y, it has a direct
          connection to the Arithmetic and Logic Unit (ALU). This is why
          many operations are only available for the accumulator, not the
          index registers.
*/
type Accumulator = uint8

/*
Register for addressing data with indices
*/
type IndexRegister = uint8

func (p ProcessorStatus) String() string {
	flags := []byte("00000000")
	if p&PFLAG_C != 0 {
		flags[7] = 'C'
	}
	if p&PFLAG_Z != 0 {
		flags[6] = 'Z'
	}
	if p&PFLAG_I != 0 {
		flags[5] = 'I'
	}
	if p&PFLAG_D != 0 {
		flags[4] = 'D'
	}
	if p&PFLAG_B != 0 {
		flags[3] = 'B'
	}
	if p&PFLAG_UNUSED != 0 {
		flags[2] = '1'
	}
	if p&PFLAG_V != 0 {
		flags[1] = 'V'
	}
	if p&PFLAG_N != 0 {
		flags[0] = 'N'
	}
	return fmt.Sprintf("%02x(%s)", byte(p), string(flags))
}

func (p *ProcessorStatus) Set(flag ProcessorStatus, val bool) {
	if val {
		*p |= flag
	} else {
		*p &= ^flag
	}
}