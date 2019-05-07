package main

import (
	"flag"
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/nes"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
	"os"
	"fmt"
)

var logger = logger2.GetLogger()

func main() {
	var fileName string

	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/cpu_reset/ram_after_reset.nes" // passed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/cpu_reset/registers.nes" // passed
	//
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/instr_misc/instr_misc.nes"
	//
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_nes_cpu_test5/official.nes" // passed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_nes_cpu_test5/cpu.nes" // illegal opcode
	//
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_ppu_tests_2005.09.15b/palette_ram.nes" // passed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_ppu_tests_2005.09.15b/power_up_palette.nes" //failed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_ppu_tests_2005.09.15b/sprite_ram.nes" // passed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_ppu_tests_2005.09.15b/vbl_clear_time.nes" // failed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/blargg_ppu_tests_2005.09.15b/vram_access.nes" // passed

	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/cpu_exec_space/test_cpu_exec_space_ppuio.nes" // failed
	//fileName = "/Users/vfreex/Documents/hack/NES/tests/nes-test-roms/cpu_interrupts_v2/cpu_interrupts.nes" // failed

	//fileName = "/Users/vfreex/Documents/hack/NES/color_test_src/color_test.nes" // passed



	flag.Parse()
	if flag.NArg() > 0 {
		fileName = flag.Arg(0)
	}

	if len(fileName) == 0 {
		fmt.Fprintf(os.Stderr, "GoNES v0.1-alpha\n\nUsage:\n\t<rom-file>\n")
		os.Exit(1)
		return
	}

	romFile, err := os.Open(fileName)
	if err != nil {
		panic(fmt.Errorf("error opening ROM file: %v - %v", fileName, err))
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
