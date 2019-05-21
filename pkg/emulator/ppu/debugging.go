package ppu

import (
	"encoding/hex"
	"github.com/vfreex/gones/pkg/emulator/memory"
)

func (ppu *PPUImpl) dumpNametable() {
	var data [0x400]byte
	nt := memory.Ptr(0x2000) + memory.Ptr(ppu.registers.ctrl&PPUCtrl_NameTable)*0x400
	for i := memory.Ptr(0); i < 0x400; i++ {
		data[i] = ppu.vram.Peek(nt + i)
	}
	dump := hex.Dump(data[:])
	logger.Infof("v=%v, x=%v", ppu.registers.v.String(), ppu.registers.x)
	logger.Infof("from nametable: %04x", nt)
	logger.Info(dump)
	logger.Sync()
}
