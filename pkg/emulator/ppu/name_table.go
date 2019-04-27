package ppu

// PPU 2000h-23FFh - Name Table 0 and Attribute Table 0 (1K)
// PPU 2400h-27FFh - Name Table 1 and Attribute Table 1 (1K)
// PPU 2800h-2BFFh - Name Table 2 and Attribute Table 2 (1K)
// PPU 2C00h-2FFFh - Name Table 3 and Attribute Table 3 (1K)
// PPU 3000h-3EFFh - Mirror of 2000h-2EFFh

// A name tables is essentially a matrix of tile numbers, pointing to the tiles stored in the pattern
// tables.
// The tiles are fetched from Pattern Table 0 or 1 (depending on Bit 4 in PPU Control Register 1)..
// The name tables are 32x30 (960/3c0h) tiles and since each tile is 8x8 pixels, the entire name
// table is 256x240 pixels.

// Each Name Table is directly followed by an Attribute Table of 40h bytes,
//  containing 2bit background palette numbers for each 16x16 pixel field (2*2 group of tiles).
//  Each byte in the Attribute table defines palette numbers for a 32x32 pixel area (4*4 group of tiles):
//  Bit0-1  Palette Number for upperleft 16x16 pixels of the 32x32 area
//  Bit2-3  Palette Number for upperright 16x16 pixels of the 32x32 area
//  Bit4-5  Palette Number for lowerleft 16x16 pixels of the 32x32 area
//  Bit6-7  Palette Number for lowerright 16x16 pixels of the 32x32 area

type NameTableAndAttributeTable struct {
	nameTable      [0x3c0]byte
	attributeTable [0x40]byte
}
