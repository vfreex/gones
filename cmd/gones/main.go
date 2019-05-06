package main

import (
	"flag"
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/nes"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
	"os"
)

var logger = logger2.GetLogger()

func main() {
	var fileName string

	fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/cpu_reset/ram_after_reset.nes" // passed
	fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/cpu_reset/registers.nes" // passed

	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/instr_misc/instr_misc.nes"

	fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_nes_cpu_test5/official.nes"

	//fileName = "/Users/vfreex/Downloads/nestests/blargg_ppu_tests_2005.09.15b/sprite_ram.nes" // passed
	//fileName = "/Users/vfreex/Downloads/nestests/blargg_ppu_tests_2005.09.15b/palette_ram.nes" // passed
	//fileName = "/Users/vfreex/Downloads/nestests/blargg_ppu_tests_2005.09.15b/vram_access.nes" // failed

	flag.Parse()
	if flag.NArg() > 1 {
		fileName = flag.Arg(0)
	}

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
