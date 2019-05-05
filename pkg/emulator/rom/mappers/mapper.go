package mappers

import "github.com/vfreex/gones/pkg/emulator/memory"

const (
	PrgBankSize = 16 * 1024 // bytes in a PRG/ROM bank
	ChrBankSize = 8 * 1024  // bytes in a CHR/VROM bank
)

type Mapper interface {
	Map() (prg memory.Memory, chr memory.Memory)
}
