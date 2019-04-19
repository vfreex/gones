/*
Byte     Contents
---------------------------------------------------------------------------
0-3      String "NES^Z" used to recognize .NES files.
4        Number of 16kB ROM banks.
5        Number of 8kB VROM banks.
6        bit 0     1 for vertical mirroring, 0 for horizontal mirroring.
bit 1     1 for battery-backed RAM at $6000-$7FFF.
bit 2     1 for a 512-byte trainer at $7000-$71FF.
bit 3     1 for a four-screen VRAM layout.
bit 4-7   Four lower bits of ROM Mapper Type.
7        bit 0     1 for VS-System cartridges.
bit 1-3   Reserved, must be zeroes!
bit 4-7   Four higher bits of ROM Mapper Type.
8        Number of 8kB RAM banks. For compatibility with the previous
versions of the .NES format, assume 1x8kB RAM page when this
byte is zero.
9        bit 0     1 for PAL cartridges, otherwise assume NTSC.
bit 1-7   Reserved, must be zeroes!
10-15    Reserved, must be zeroes!
16-...   ROM banks, in ascending order. If a trainer is present, its
512 bytes precede the ROM bank contents.
...-EOF  VROM banks, in ascending order.
---------------------------------------------------------------------------
*/

package ines

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	INES_FILE_MAGIC   = "NES\x1a"
	PRG_BANK_SIZE     = 16 * 1024 // bytes in a PRG/ROM bank
	CHR_BANK_SIZE     = 8 * 1024  // bytes in a CHR/VROM bank
	PRG_RAM_BANK_SIZE = 8 * 1024  // bytes in a RPG RAM bank
)

const (
	FLAGS6_VERTICAL_MIRRORING  = 1
	FLAGS6_BATTERY_RAM_ON      = 1 << 1
	FLAGS6_TRAINER_ON          = 1 << 2
	FLAGS6_FOUR_SCREEN_VRAM_ON = 1 << 3

	FLAGS7_VS_SYSTEM_ON = 1

	FLAGS9_PAL_ON = 1
)

const (
	MAPPER_NORM = 0,
)

type INesHeader struct {
	Magic      [4]byte
	PrgSize    byte // PRG/ROM in 16 KB units
	ChrSize    byte // CHR/VROM in 8 KB units, value 0 means the board uses CHR RAM
	Flags6     byte
	Flags7     byte
	PrgRamSize byte // in 8kB units, value 0 means 1x8kB for compatibility, see http://wiki.nesdev.com/w/index.php/PRG_RAM_circuit
	Flags9     byte
	_          [6]byte
}
type INesRom struct {
	Header  INesHeader
	Trainer *bytes.Buffer
	Prg     *bytes.Buffer
	Chr     *bytes.Buffer
	Extra   *bytes.Buffer
}

func NewINesRom(reader io.Reader) (*INesRom, error) {
	rom := &INesRom{}
	header := &rom.Header
	if err := binary.Read(reader, binary.LittleEndian, header); err != nil {
		return rom, err
	}
	if string(rom.Header.Magic[:]) != INES_FILE_MAGIC {
		return rom, fmt.Errorf("no valid header is found")
	}

	rom.Trainer = &bytes.Buffer{}

	prgBuf := make([]byte, 0, PRG_BANK_SIZE*int(header.PrgSize))
	rom.Prg = bytes.NewBuffer(prgBuf)
	if _, err := io.CopyN(rom.Prg, reader, int64(rom.Prg.Cap())); err != nil {
		return rom, err
	}

	chrBuf := make([]byte, 0, CHR_BANK_SIZE*int(header.ChrSize))
	rom.Chr = bytes.NewBuffer(chrBuf)
	if _, err := io.CopyN(rom.Chr, reader, int64(rom.Chr.Cap())); err != nil {
		return rom, err
	}
	rom.Extra = &bytes.Buffer{}
	if _, err := io.Copy(rom.Extra, reader); err != nil {
		return rom, err
	}
	return rom, nil
}
func (p *INesRom) String() string {
	return fmt.Sprintf("iNESRom{header: %v, trainer: %d, PRG: %d, CHR: %d, EXTRA: %d}", &p.Header, p.Trainer.Len(), p.Prg.Len(), p.Chr.Len(), p.Extra.Len())
}
func (p *INesRom) MatchesFileMagic(reader io.Reader) (bool, error) {
	magic := make([]byte, 4)
	if n, err := reader.Read(magic); n != 4 || string(magic) != INES_FILE_MAGIC {
		return false, err
	}

	return true, nil
}
func (h *INesHeader) String() string {
	m := map[string]interface{}{
		"type":          "iNES",
		"mapper_type":   h.GetMapperType(),
		"prg_bytes":     int(h.PrgSize) * PRG_BANK_SIZE,
		"chr_bytes":     int(h.ChrSize) * CHR_BANK_SIZE,
		"prg_ram_bytes": PRG_RAM_BANK_SIZE,
	}
	if h.PrgRamSize > 0 {
		m["prg_ram_bytes"] = PRG_RAM_BANK_SIZE * int(h.PrgRamSize)
	}
	r, _ := json.Marshal(m)
	return string(r)
}
func (h *INesHeader) GetMapperType() int {
	return int((h.Flags7 & 0xF0) | (h.Flags6 >> 4))
}
