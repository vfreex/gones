package ppu

import (
	"bytes"
	"fmt"
	"github.com/vfreex/gones/pkg/emulator/memory"
	"math/rand"
	"time"
)

// The logical screen resolution processed by the PPU is 256x240 pixels
// The PPU renders 262 scanlines per frame.
// Each scanline lasts for 341 PPU clock cycles (113.667 CPU clock cycles; 1 CPU cycle = 3 PPU cycles),
// with each clock cycle producing one pixel
//func (ppu *PPUImpl) RenderFrame() {
//	for scanlineId := 0; scanlineId < SCREEN_HEIGHT; scanlineId++ {
//		for x := 0; x < SCREEN_WIDTH; x++ {
//			//ppu.RenderedBuffer[scanlineId][x] = RBGColor(rnd.Int())
//			//ppu.RenderPixel(scanlineId, x)
//		}
//	}
//	ppu.registers.status |= PPUStatus_VBlank
//}
func (ppu *PPUImpl) getCurrentNametableAddr() memory.Ptr {
	// The tiles are fetched from Pattern Table 0 or 1 (depending on Bit 4 in PPU Control Register 1)
	var baseAddr memory.Ptr
	switch ppu.registers.ctrl & PPUCtrl_NameTable {
	case 0:
		baseAddr = 0x2000
	case 1:
		baseAddr = 0x2400
	case 2:
		baseAddr = 0x2800
	case 3:
		baseAddr = 0x2c00
	}
	return baseAddr
}
func (ppu *PPUImpl) ReadNameTableByte(offset memory.PtrDist) byte {
	baseAddr := ppu.getCurrentNametableAddr()
	return ppu.vram.Peek(baseAddr + offset)
}

func (ppu *PPUImpl) ReadAttributeTableByte(offset memory.PtrDist) byte {
	baseAddr := ppu.getCurrentNametableAddr()
	return ppu.vram.Peek(baseAddr + 0x3c0 + memory.Ptr(offset))
}

func (ppu *PPUImpl) ReadPatternTableByte(offset memory.Ptr) byte {
	if ppu.registers.ctrl & PPUCtrl_BackgroundPatternTable != 0 {
		offset |= 0x1000
	}
	return ppu.vram.Peek(offset)
}

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

func (ppu *PPUImpl) TestDisplay(scanlineId int, x int, colorId byte) {
	color := ppu.Palette.Peek(memory.Ptr(colorId) + 0x3F00)
	//log.Infof("(%d, %d)=%x,%4x", scanlineId, x, colorId, color)
	if scanlineId < SCREEN_HEIGHT && x < SCREEN_WIDTH {
		//ppu.RenderedBuffer[x][scanlineId] = RGBMap[rnd.Int63()%64]
		ppu.RenderedBuffer[scanlineId][x] = RGBMap[color]
	}
}
func (ppu *PPUImpl) DrawPixel(patternEntry0, patternEntry1, attrTableEntry byte) {
	x := ppu.dotInScanline
	y := ppu.scanline - 21
	tileX := x / 8
	tileY := y / 8

	if ppu.scanline >= 0 && ppu.scanline < SCREEN_HEIGHT && x >= 0 && x<SCREEN_WIDTH {
		// Draw Background
		if ppu.registers.mask & PPUMask_BackgroundVisibility != 0 {
			field := tileX%4/2 + tileY%4/2*2
			paletteIndex := ((attrTableEntry >> uint(field*2) & 0x3) << 2) |
				(patternEntry0 >> uint(x%8) & 0x1) |
				((patternEntry1 >> uint(x%8) & 0x1) << 1)
			ppu.RenderedBuffer[y][x] = Color(ppu.Palette.Peek(0x3F00 + memory.Ptr(paletteIndex))).ToGRBColor()
		}
	}

}
func (ppu *PPUImpl) Render() {
	scanline := ppu.scanline
	x := ppu.dotInScanline
	y := scanline - 21
	if scanline == 0 {
		if x == 0 {
			ppu.registers.status |= PPUStatus_VBlank
			if ppu.registers.ctrl&PPUCtrl_NMIOnVBlank != 0 {
				ppu.cpu.NMI = true
			}
		}
		// VINT pulled down, nops
	} else if scanline <= 20 {
		// dummy scanline
	} else if scanline <= 260 {
		// rendering
		var nametableEntry byte
		var attrTableEntry byte
		var patternEntry0 byte
		var patternEntry1 byte
		//var paletteIndex byte

		var nametableEntryF byte
		var attrTableEntryF byte
		var patternEntry0F byte
		var patternEntry1F byte
		if x < 256 {
			//tileX := x / 8
			//tileY := y / 8
			//field := tileX%4/2 + tileY%4/2*2
			//paletteIndex = ((attrTableEntry >> uint(field*2) & 0x3) << 2) |
			//	(patternEntry0 >> uint(x%8) & 0x1) |
			//	((patternEntry1 >> uint(x%8) & 0x1) << 1)

			//ppu.RenderedBuffer[y][x] = Color(ppu.Palette.Peek(0x3F00 + memory.Ptr(paletteIndex))).ToGRBColor()
			ppu.DrawPixel(patternEntry0, patternEntry1, attrTableEntry)

			tileXF := x/8 + 2
			tileYF := y / 8
			tileIndex := tileYF*32 + tileXF
			areaX := tileXF / 4
			areaY := tileYF / 4
			areaIndex := areaY*8 + areaX
			switch x % 8 {
			case 0:
				nametableEntryF = ppu.ReadNameTableByte(memory.PtrDist(tileIndex))
			case 1:
				attrTableEntryF = ppu.ReadAttributeTableByte(memory.PtrDist(areaIndex))
			case 2:
				patternTableOffset := memory.Ptr(nametableEntry)*16 + memory.Ptr(y)%8
				patternEntry0F = ppu.ReadPatternTableByte(patternTableOffset)
			case 3:
				patternTableOffset := memory.Ptr(nametableEntry)*16 + 8 + memory.Ptr(y)%8
				patternEntry1F = ppu.ReadPatternTableByte(patternTableOffset)

				nametableEntry = nametableEntryF
				attrTableEntry = attrTableEntryF
				patternEntry0 = patternEntry0F
				patternEntry1 = patternEntry1F
			}
		}
	} else {
		// when this scanline finishes, the VINT flag is set
		//ppu.dumpVRAM()
	}
	ppu.dotInScanline++
	if ppu.dotInScanline >= 341 {
		ppu.dotInScanline %= 341
		ppu.scanline++
		if ppu.scanline >= 262 {
			ppu.scanline %= 262
			ppu.frame++
		}
	}
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
