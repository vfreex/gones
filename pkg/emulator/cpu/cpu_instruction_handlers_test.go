package cpu

import (
	"github.com/vfreex/gones/pkg/emulator/ram"
	"log"
	"testing"
)

func TestOpcodeHandlers(t *testing.T) {
	for opcode, info := range InstructionInfos {
		handler := opcodeHandlers[opcode]
		if handler == nil {
			log.Printf("Opcode %02x (%s %d) is not implmeneted yet.", opcode, info.Nemonics, info.AddressingMode)
			continue
		}
		if handler.AddressingMode != info.AddressingMode {
			t.Fatalf("BUG: Addressing mode for %02x opcodeHandlers doesn't match the info in InstructionInfos: got %d, expected: %d",
				opcode, handler.AddressingMode, info.AddressingMode)
		}
	}
}

func execAdc(a byte, x byte, c bool) (ans byte, status ProcessorStatus) {
	ram := ram.NewRAM(2)
	cpu := NewCpu(ram)
	ram.Poke(0, 0x69)
	ram.Poke(1, x)
	cpu.PC = 0
	cpu.A = a
	cpu.P.Set(PFLAG_C, c)
	cpu.ExecOneInstruction()
	return cpu.A, cpu.P
}

func TestADC_Signed(t *testing.T) {
	for a := -128; a < 128; a++ {
		for x := -128; x < 128; x++ {
			for c := 0; c < 2; c++ {
				r, p := execAdc(byte(a), byte(x), c != 0)
				actual := int8(r)
				actualOF := p&PFLAG_V != 0
				expected := a + x + c
				expectedOF := expected >= 128 || expected < -128
				expected2 := int8(expected)
				if expectedOF == actualOF && expected2 == actual {
					continue
				}
				t.Fatalf("error computing %02x+%02x+%x, got %02x %v, excpeted %02x %v", a, x, c,
					actual, actualOF, expected2, expectedOF)
			}
		}
	}
}

func TestADC_Unsigned(t *testing.T) {
	for a := 0; a < 256; a++ {
		for x := 0; x < 256; x++ {
			for c := 0; c < 2; c++ {
				r, p := execAdc(byte(a), byte(x), c != 0)
				actual := uint8(r)
				actualCF := p&PFLAG_C != 0
				expected := a + x + c
				expectedCF := expected >= 256
				expected2 := uint8(expected)
				if expectedCF == actualCF && expected2 == actual {
					continue
				}
				t.Fatalf("error computing %02x+%02x+%x, got %02x %v, excpeted %02x %v", a, x, c,
					actual, actualCF, expected2, expectedCF)
			}
		}
	}
}

func execSbc(a byte, x byte, c bool) (ans byte, status ProcessorStatus) {
	ram := ram.NewRAM(2)
	cpu := NewCpu(ram)
	ram.Poke(0, 0xe9)
	ram.Poke(1, x)
	cpu.PC = 0
	cpu.A = a
	cpu.P.Set(PFLAG_C, c)
	cpu.ExecOneInstruction()
	return cpu.A, cpu.P
}

func TestSBC_Signed(t *testing.T) {
	for a := -128; a < 128; a++ {
		for x := -128; x < 128; x++ {
			for c := 0; c < 2; c++ {
				r, p := execSbc(byte(a), byte(x), c != 0)
				actual := int8(r)
				actualOF := p&PFLAG_V != 0
				expected := a - x - (1 - c)
				expected2 := int8(expected)
				expectedOF := expected >= 128 || expected < -128
				if expectedOF == actualOF && expected2 == actual {
					continue
				}
				t.Fatalf("error computing %02x+%02x+%x, got %02x %v, excpeted %02x %v", a, x, c,
					actual, actualOF, expected2, expectedOF)
			}
		}
	}
}

func TestSBC_Unsigned(t *testing.T) {
	for a := 0; a < 256; a++ {
		for x := 0; x < 256; x++ {
			for c := 0; c < 2; c++ {
				r, p := execSbc(byte(a), byte(x), c != 0)
				actual := r
				actualCF := p&PFLAG_C != 0
				expected := uint8(a - x - (1 - c))
				expectedCF := a-x-(1-c) >= 0
				if expectedCF == actualCF && expected == actual {
					continue
				}
				t.Fatalf("error computing %02x-%02x-1+%x, got %02x %v, excpeted %02x %v", a, x, c,
					actual, actualCF, expected, expectedCF)
			}
		}
	}
}
