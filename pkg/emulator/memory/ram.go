package memory

type RAM struct {
	data []byte
}

func NewRAM(size int) *RAM {
	return &RAM{data: make([]byte, size)}
}

func (r *RAM) Peek(addr Ptr) byte {
	return r.data[addr]
}

func (r *RAM) Poke(addr Ptr, val byte) {
	r.data[addr] = val
}
