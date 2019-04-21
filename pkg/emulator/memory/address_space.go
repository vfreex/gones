package memory

import (
	"fmt"
	"sort"
)

type AddressSpace interface {
	Memory
	MapMemory(offset Ptr, length PtrDist, mode MMapMode, mappedMemory Memory, translator AddressTranslator)
}

type AddressTranslator func(addr Ptr) Ptr

type MMapMode uint

const (
	MMAP_MODE_READ = 1 << iota
	MMAP_MODE_WRITE
	MMAP_MODE_EXEC
)

type MMapEntry struct {
	Offset     Ptr
	Length     PtrDist
	Mode       MMapMode
	Memory     Memory
	Translator AddressTranslator
}

type AddressSpaceImpl struct {
	mMapEntries []MMapEntry
}

func (as *AddressSpaceImpl) MapMemory(offset Ptr, length PtrDist, mode MMapMode, mappedMemory Memory, translator AddressTranslator) {
	as.mMapEntries = append(as.mMapEntries, MMapEntry{
		Offset:     offset,
		Length:     length,
		Mode:       mode,
		Memory:     mappedMemory,
		Translator: translator,
	})
}

func (as *AddressSpaceImpl) lookupMappedMemory(addr Ptr) (*MMapEntry, Ptr) {
	index := sort.Search(len(as.mMapEntries), func(i int) bool {
		return as.mMapEntries[i].Offset > addr
	}) - 1
	if index < 0 || int(addr)-int(as.mMapEntries[index].Offset) >= int(as.mMapEntries[index].Length) {
		panic(fmt.Errorf("trying to access unmapped address 0x%x", addr))
	}
	mappedAddr := addr
	if as.mMapEntries[index].Translator != nil {
		mappedAddr = as.mMapEntries[index].Translator(addr)
	}
	return &as.mMapEntries[index], mappedAddr
}

func (as *AddressSpaceImpl) Peek(addr Ptr) byte {
	entry, mappedAddr := as.lookupMappedMemory(addr)
	if entry.Mode&MMAP_MODE_READ == 0 {
		panic(fmt.Errorf("permmission denied when trying to read 0x%x", addr))
	}
	return entry.Memory.Peek(mappedAddr)
}

func (as *AddressSpaceImpl) Poke(addr Ptr, val byte) {
	entry, mappedAddr := as.lookupMappedMemory(addr)
	if entry.Mode&MMAP_MODE_WRITE == 0 {
		panic(fmt.Errorf("permmission denied when trying to write %x", addr))
	}
	entry.Memory.Poke(mappedAddr, val)
}
