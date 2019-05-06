package cpu

import "github.com/vfreex/gones/pkg/emulator/memory"

type InterruptVector memory.Ptr

const (
	IV_NMI   InterruptVector = 0xFFFA
	IV_RESET InterruptVector = 0xFFFC
	IV_IRQ   InterruptVector = 0xFFFE
	IV_BRK   InterruptVector = 0xFFFE
)

func (cpu *Cpu) ReadInterruptVector(iv InterruptVector) memory.Ptr {
	addrLow := cpu.Memory.Peek(memory.Ptr(iv))
	addrHigh := cpu.Memory.Peek(memory.Ptr(iv + 1))
	return (memory.Ptr(addrHigh) << 8) | memory.Ptr(addrLow)
}

func (cpu *Cpu) Reset() {
	logger.Debug("Reset CPU")
	cpu.P.Set(PFLAG_B, true)
	cpu.SP -= 3
	cpu.P.Set(PFLAG_I, true)
	cpu.PC = cpu.ReadInterruptVector(IV_RESET)
	logger.Debugf("go to: PC=$%4x", cpu.PC)
	// TODO: APU was silenced ($4015 = 0)
}

func (cpu *Cpu) ExecIRQ() {
	logger.Debug("handling IRQ")
	cpu.P.Set(PFLAG_B, false)
	cpu.PushW(cpu.PC)
	cpu.Push(byte(cpu.P))
	cpu.P.Set(PFLAG_I, true)
	cpu.PC = cpu.ReadInterruptVector(IV_IRQ)
}

func (cpu *Cpu) ExecNMI() {
	logger.Debug("handling NMI")
	cpu.P.Set(PFLAG_B, false)
	cpu.PushW(cpu.PC)
	cpu.Push(byte(cpu.P))
	cpu.P.Set(PFLAG_I, true)
	cpu.PC = cpu.ReadInterruptVector(IV_NMI)
	cpu.NMI = false
}