package cpu

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
type ProcessorStatus = uint8

const (
	// Carry flag
	C_FLAG = 1
	// Zero flag
	Z_FLAG = 1 << 1
	// Interrupt disable flag
	I_FLAG = 1 << 2
	// Decimal mode flag.
	// 2A03 does not support BCD mode so although the flag can be set, its value will be ignored.
	D_FLAG = 1 << 3
	// Break flag
	B_FLAG = 1 << 4
	// unused flag, alway 1
	UNUSED_FLAG = 1 << 5
	// oVerflow flag
	V_FLAG = 1 << 6
	// Negative flag
	N_FLAG = 1 << 7
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
