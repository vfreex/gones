package cpu

type AddressingMode int

//go:generate stringer -type=AddressingMode -output addressing_modes.gen.go
const (
	IMP AddressingMode = iota
	IMM
	ZP
	ZPX
	ZPY
	IZX
	IZY
	ABS
	ABX
	ABY
	IND
	REL
)

func (i AddressingMode) GetArgumentCount() uint16 {
	length := uint16(1)
	switch i {
	case IMP:
		length = 0
	case IND, ABS, ABX, ABY:
		length = 2
	}
	return length
}
