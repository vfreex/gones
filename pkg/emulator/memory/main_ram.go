package memory

const (
	MAIN_RAM_SIZE = 1024 * 2 // 2 kiB main RAM
)

type MainRAM struct {
	RAM
}

func NewMainRAM() *MainRAM {
	ram := &MainRAM{}
	ram.data = make([]byte, MAIN_RAM_SIZE)
	return ram
}

func (ram *MainRAM) Peek(addr Ptr) byte {
	return ram.RAM.Peek(addr & 0x7FF)
}

func (ram *MainRAM) Poke(addr Ptr, val byte) {
	ram.RAM.Poke(addr&0x7FF, val)
}
