package cpu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"github.com/vfreex/gones/pkg/emulator/ram"
	"testing"
)

func setupCPU(program memory.Memory) *Cpu {
	as := &memory.AddressSpaceImpl{}
	// 0x0000 - ox1fff RAM
	as.MapMemory(0, 0x2000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, ram.NewMainRAM(), nil)
	// test RAM
	as.MapMemory(0x2000, 0x6000, memory.MMAP_MODE_READ|memory.MMAP_MODE_WRITE, ram.NewRAM(0x6000), func(addr memory.Ptr) memory.Ptr {
		return 0x2000 + addr
	})
	// test ROM
	as.MapMemory(0x8000, 0x8000, memory.MMAP_MODE_READ, program, func(addr memory.Ptr) memory.Ptr {
		return 0x8000 + addr
	})
	return NewCpu(as)
}

func TestADC_PositiveAddNegative(t *testing.T) {
	// test 7f + ff + 0 == 7e
	program := ram.NewRAM(0x8000)
	// adc #ff
	program.Poke(0x0000, 0x69)
	program.Poke(0x0001, 0xff)
	cpu := setupCPU(program)
	cpu.PC = 0x8000
	cpu.A = 0x7f
	cpu.P.Set(PFLAG_C, false)
	cpu.ExecOneInstruction()
	if cpu.A != 0x7e {
		t.Fatalf("got %x, expected %x", cpu.A, 0x7e)
	}
	if cpu.P&PFLAG_C == 0 {
		t.Fatalf("got CF=0, expected CF=1")
	}
	if cpu.P&PFLAG_V != 0 {
		t.Fatalf("got OF=1, expected OF=0")
	}
	if cpu.P&PFLAG_Z != 0 {
		t.Fatalf("got ZF=1, expected ZF=0")
	}
	if cpu.P&PFLAG_N != 0 {
		t.Fatalf("got NF=1, expected NF=0")
	}
}


func TestADC_PositiveAddPositive(t *testing.T) {
	// test 7f + 7f + 1 == ff
	program := ram.NewRAM(0x8000)
	// adc #ff
	program.Poke(0x0000, 0x69)
	program.Poke(0x0001, 0x7f)
	cpu := setupCPU(program)
	cpu.PC = 0x8000
	cpu.A = 0x7f
	cpu.P.Set(PFLAG_C, true)
	cpu.ExecOneInstruction()
	if cpu.A != 0xff {
		t.Fatalf("got %x, expected %x", cpu.A, 0xff)
	}
	if cpu.P&PFLAG_C != 0 {
		t.Fatalf("got CF=0, expected CF=1")
	}
	if cpu.P&PFLAG_V == 0 {
		t.Fatalf("got OF=1, expected OF=0")
	}
	if cpu.P&PFLAG_Z != 0 {
		t.Fatalf("got ZF=1, expected ZF=0")
	}
	if cpu.P&PFLAG_N == 0 {
		t.Fatalf("got NF=1, expected NF=0")
	}
}


func TestADC_NegativeAddNegative(t *testing.T) {
	// test ff + ff + 0 == fe
	program := ram.NewRAM(0x8000)
	// adc #ff
	program.Poke(0x0000, 0x69)
	program.Poke(0x0001, 0xff)
	cpu := setupCPU(program)
	cpu.PC = 0x8000
	cpu.A = 0xff
	cpu.P.Set(PFLAG_C, false)
	cpu.ExecOneInstruction()
	if cpu.A != 0xfe {
		t.Fatalf("got %x, expected %x", cpu.A, 0xfe)
	}
	if cpu.P&PFLAG_C == 0 {
		t.Fatalf("got CF=0, expected CF=1")
	}
	if cpu.P&PFLAG_V != 0 {
		t.Fatalf("got OF=1, expected OF=0")
	}
	if cpu.P&PFLAG_Z != 0 {
		t.Fatalf("got ZF=1, expected ZF=0")
	}
	if cpu.P&PFLAG_N == 0 {
		t.Fatalf("got NF=0, expected NF=1")
	}
}
