package cpu

import "github.com/vfreex/gones/pkg/emulator/memory"

type InterruptVector memory.Ptr
const (
	IV_NMI InterruptVector = 0xFFFA
	IV_RESET InterruptVector = 0xFFFC
	IV_IRQ InterruptVector = 0xFFFE
	IV_BRK InterruptVector = 0xFFFE
)

func (cpu *Cpu) ReadInterruptVector(iv InterruptVector) memory.Ptr {
	addrLow := cpu.Memory.Peek(memory.Ptr(iv))
	addrHigh := cpu.Memory.Peek(memory.Ptr(iv + 1))
	return (memory.Ptr(addrHigh) << 8) | memory.Ptr(addrLow)
}