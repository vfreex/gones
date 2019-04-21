package memory

type Ptr = uint16
type PtrDiff = int16
type PtrDist = uint16

type Memory interface {
	Peek(addr Ptr) byte
	Poke(addr Ptr, val byte)
}
