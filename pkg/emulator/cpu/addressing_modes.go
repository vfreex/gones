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
