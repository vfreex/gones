package main

import (
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/nes"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
	"os"
)

var logger = logger2.GetLogger()

func main() {
	//fileName := "/Users/vfreex/Documents/hack/NES/NES_Dev_01/ctnes.nes"
	//fileName := "/Users/vfreex/Documents/hack/NES/roms/Balloon Fight (U) .nes"
	fileName := "/Users/vfreex/Documents/hack/NES/tests/branch_timing_tests/1.Branch_Basics.nes"
	romFile, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer romFile.Close()
	var rom *ines.INesRom
	if rom, err = ines.NewINesRom(romFile); err != nil {
		panic(err)
	}
	logger.Warnf("iNES ROM file loaded: %v\n", rom)

	nes := nes.NewNes()
	nes.LoadCartridge(rom)
	nes.Start()
}
