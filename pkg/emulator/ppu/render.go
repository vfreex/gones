package ppu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
	"math/rand"
	"time"
)

// The logical screen resolution processed by the PPU is 256x240 pixels
// The PPU renders 262 scanlines per frame.
// Each scanline lasts for 341 PPU clock cycles (113.667 CPU clock cycles; 1 CPU cycle = 3 PPU cycles),
// with each clock cycle producing one pixel
func (ppu *PPUImpl) RenderFrame() {
	for scanlineId := 0; scanlineId < SCREEN_HEIGHT; scanlineId++ {
		for x := 0; x < SCREEN_WIDTH; x++ {
			//ppu.RenderedBuffer[scanlineId][x] = RBGColor(rnd.Int())
			ppu.RenderPixel(scanlineId, x)
		}
	}
	ppu.registers.status |= PPUStatus_VBlank
	ppu.cycles++
}

func (ppu *PPUImpl) ReadNameTableByte(tableId, index byte) byte {
	// The tiles are fetched from Pattern Table 0 or 1 (depending on Bit 4 in PPU Control Register 1)
	// assuming table #0
	//logger.Debugf("accessing VRAM %04x", 0x2000+memory.Ptr(index))
	return ppu.vram.Peek(0x2000 + memory.Ptr(index))
}

func (ppu *PPUImpl) ReadAttributeTableByte(tableId, index byte) byte {
	return ppu.vram.Peek(0x23c0 + memory.Ptr(index))
}

func (ppu *PPUImpl) ReadPatternTableByte(tableId, tileId, entryId byte) byte {
	return ppu.vram.Peek(memory.Ptr(entryId))
}

func (ppu *PPUImpl) RenderPixel(scanlineId int, x int) {
	if scanlineId < 20 {
		// VINT unset
		return
	}
	if scanlineId == 20 {
		// dummy scanline
		// For odd frames, the cycle at the end of the scanline is skipped
		if x == 339 && ppu.Frames&1 != 0 {
			ppu.cycles++
			ppu.RenderPixel(0, 0)
			return
		}
	}
	if scanlineId < 261 {
		//At the beginning of each scanline,
		// the data for the first two tiles is already loaded into the shift registers (and ready to be rendered),'
		// so the first tile that gets fetched is Tile 3.
		nameTableEntry := byte(0)
		attrTableEntry := byte(0)
		patternTableEntry0, patternTableEntry1 := byte(0), byte(0)
		if x < 256 {
			pixelX := 2*8 + x
			pixelY := scanlineId
			tileId := byte(pixelY/8*32 + pixelX/8)
			tileX := byte(pixelX / 8)
			tileY := byte(pixelY / 8)
			switch x / 8 {
			case 0:
				nameTableEntry = ppu.ReadNameTableByte(0, tileId)
			case 2:
				attrTableEntry = ppu.ReadAttributeTableByte(0, tileId/4)
			case 4:
				patternTableEntry0 = ppu.ReadPatternTableByte(0, nameTableEntry, (tileY*8+tileX)*2/8)
			case 6:
				patternTableEntry1 = ppu.ReadPatternTableByte(0, nameTableEntry, (tileY*8+tileX)*2/8+1)
			}
			colorID := ((patternTableEntry0 >> byte(pixelX/4)) & 0x3) | (((attrTableEntry >> byte(pixelX/4)) & 0x3) << 2)
			colorID1 := ((patternTableEntry1 >> byte(pixelX/4)) & 0x3) | (((attrTableEntry >> byte(pixelX/4)) & 0x3) << 2)

			ppu.TestDisplay(pixelY, pixelX, colorID)
			ppu.TestDisplay(pixelY, pixelX, colorID1)

		} else if x < 320 {

		} else if x < 336 {
			pixelX := x - 320
			pixelY := scanlineId
			tileId := byte(pixelY/8*32 + pixelX/8)
			tileX := byte(pixelX / 8)
			tileY := byte(pixelY / 8)
			switch x / 8 {
			case 0:
				nameTableEntry = ppu.ReadNameTableByte(0, tileId)
			case 2:
				attrTableEntry = ppu.ReadAttributeTableByte(0, tileId/4)
			case 4:
				patternTableEntry0 = ppu.ReadPatternTableByte(0, tileId, (tileY*8+tileX)*2/8)
			case 6:
				patternTableEntry1 = ppu.ReadPatternTableByte(0, tileId, (tileY*8+tileX)*2/8+1)
			}
			colorID := ((patternTableEntry0 >> byte(pixelX/4)) & 0x3) | (((attrTableEntry >> byte(pixelX/4)) & 0x3) << 2)
			colorID1 := ((patternTableEntry1 >> byte(pixelX/4)) & 0x3) | (((attrTableEntry >> byte(pixelX/4)) & 0x3) << 2)

			ppu.TestDisplay(pixelY, pixelX, colorID)
			ppu.TestDisplay(pixelY, pixelX, colorID1)
		} else if x < 340 {
			ppu.registers.status |= PPUStatus_VBlank
		}
	}

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
