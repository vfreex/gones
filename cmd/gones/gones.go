package main

import (
	"flag"
	"fmt"
	logger2 "github.com/vfreex/gones/pkg/emulator/common/logger"
	"github.com/vfreex/gones/pkg/emulator/nes"
	"github.com/vfreex/gones/pkg/emulator/rom/ines"
	"os"
)

var logger = logger2.GetLogger()

func main() {
	var fileName string
	flag.Parse()
	if flag.NArg() > 0 {
		fileName = flag.Arg(0)
	}

	if len(fileName) == 0 {
		fmt.Fprintf(os.Stderr, "GoNES v0.3.0-beta\n\nUsage:\n\t<rom-file>\n")
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
