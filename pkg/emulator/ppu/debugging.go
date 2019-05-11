package ppu

import (
	"bytes"
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

func dumpMemory(mem memory.Memory, start memory.Ptr, length memory.PtrDist) {
	loops := length / 0x10
	for i := memory.PtrDist(0); i < loops; i += 0x10 {
		s := dumpRow(mem, start+i, 0x10)
		logger.Debug(s)
	}
}

func dumpRow(mem memory.Memory, start memory.Ptr, length memory.PtrDist) string {
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("%04x", start))
	for i := memory.PtrDist(0); i < length; i++ {
		buf.WriteByte(' ')
		buf.WriteString(fmt.Sprintf("%02x", mem.Peek(start+i)))
	}
	return buf.String()
}

func (ppu *PPUImpl) dumpVRAM() {
	ntaddr := ppu.getCurrentNametableAddr()
	logger.Debugf("current nametable addr: %04x", ntaddr)
	logger.Debug("current nametable content:")
	dumpMemory(ppu.vram, ntaddr, 0x3c0)
	logger.Debug("current pattern table content:")
	dumpMemory(ppu.vram, 0, 0x1000)
	//logger.Sync()
}

