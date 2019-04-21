package memory

import (
	"fmt"
	"sort"
)

type MMapMode uint

const (
	MMAP_MODE_READ = 1 << iota
	MMAP_MODE_WRITE
	MMAP_MODE_EXEC
)

type MMapEntry struct {
	Offset    Ptr
	Length    PtrDist
	Mode      MMapMode
	Memory    Memory
	Tanslator AddressTranslator
}

type AddressTranslator func(addr Ptr) Ptr

type AddressSpace struct {
	MMapEntries []MMapEntry
}

func (as *AddressSpace) findMMapEntry(addr Ptr) *MMapEntry {
	index := sort.Search(len(as.MMapEntries), func(i int) bool {
		return as.MMapEntries[i].Offset > addr
	}) - 1
	if index < 0 || int(addr)-int(as.MMapEntries[index].Offset) >= int(as.MMapEntries[index].Length) {
		panic(fmt.Errorf("trying to access unmapped address 0x%x", addr))
	}
	return &as.MMapEntries[index]
}

func (as *AddressSpace) Peek(addr Ptr) byte {
	entry := as.findMMapEntry(addr)
	if entry.Mode&MMAP_MODE_READ == 0 {
		panic(fmt.Errorf("permmission denied when trying to read 0x%x", addr))
	}
	mappedAddr := addr
	if entry.Tanslator != nil {
		mappedAddr = entry.Tanslator(addr)
	}
	return entry.Memory.Peek(mappedAddr)
}

func (as *AddressSpace) Poke(addr Ptr, val byte) {
	entry := as.findMMapEntry(addr)
	if entry.Mode&MMAP_MODE_WRITE == 0 {
		panic(fmt.Errorf("permmission denied when trying to write %x", addr))
	}
	mappedAddr := addr
	if entry.Tanslator != nil {
		mappedAddr = entry.Tanslator(addr)
	}
	entry.Memory.Poke(mappedAddr, val)
}
