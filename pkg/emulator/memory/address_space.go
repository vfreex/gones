package memory

import (
	"fmt"
	"sort"
)

type AddressSpace interface {
	Memory
	Map()
	AddMapping(offset Ptr, length PtrDist, mode MMapMode, mappedMemory Memory, translator AddressTranslator)
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

type MMapEntries []MMapEntry

func (p MMapEntries) Len() int {
	return len([]MMapEntry(p))
}

func (p MMapEntries) Less(i, j int) bool {
	return []MMapEntry(p)[i].Offset < []MMapEntry(p)[j].Offset
}

func (p MMapEntries) Swap(i, j int) {
	t := []MMapEntry(p)[i]
	[]MMapEntry(p)[i] = []MMapEntry(p)[j]
	[]MMapEntry(p)[j] = t
}

type AddressSpaceImpl struct {
	mMapEntries MMapEntries
}

func (as *AddressSpaceImpl) AddMapping(offset Ptr, length PtrDist, mode MMapMode, mappedMemory Memory, translator AddressTranslator) {
	as.mMapEntries = append(as.mMapEntries, MMapEntry{
		Offset:     offset,
		Length:     length,
		Mode:       mode,
		Memory:     mappedMemory,
		Translator: translator,
	})
}

func (as *AddressSpaceImpl) Map() {
	sort.Sort(as.mMapEntries)
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
		panic(fmt.Errorf("permission denied when trying to read 0x%x", addr))
	}
	return entry.Memory.Peek(mappedAddr)
}

func (as *AddressSpaceImpl) Poke(addr Ptr, val byte) {
	entry, mappedAddr := as.lookupMappedMemory(addr)
	if entry.Mode&MMAP_MODE_WRITE == 0 {
		panic(fmt.Errorf("permission denied when trying to write %x", addr))
	}
	entry.Memory.Poke(mappedAddr, val)
}
