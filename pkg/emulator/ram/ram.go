package ram

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
)

type RAM struct {
	data []byte
}

func NewRAM(size int) *RAM {
	return &RAM{data: make([]byte, size)}
}

func (r *RAM) Peek(addr memory.Ptr) byte {
	return r.data[addr]
}

func (r *RAM) Poke(addr memory.Ptr, val byte) {
	r.data[addr] = val
}
