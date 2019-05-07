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
	return ppu.vram.Peek(baseAddr + 0x3c0 + offset)
}

func (ppu *PPUImpl) ReadPatternTableByte(offset memory.Ptr) byte {
	offset &= 0xfff
	if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
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

func (ppu *PPUImpl) renderSprites() {
	// http://wiki.nesdev.com/w/index.php/PPU_sprite_evaluation
	x := ppu.dotInScanline
	y := ppu.scanline - 21
	if x == 0 {
		// TODO: Cycles 1-64: Secondary OAM (32-byte buffer for current sprites on scanline) is initialized to $FF,
		//  attempting to read $2004 will return $FF
		//for i := 0; i < 32; i++ {
		//	ppu.secondaryOAM.Poke(memory.Ptr(i), 0xff)
		//}
		//} else if x == 1 {
		// Sprite evaluation
		ppu.spriteCount = 0
		ppu.registers.status &= ^PPUStatus_SpriteOverflow
		for spriteIndex := 0; spriteIndex < 64; spriteIndex++ {
			spriteY := int(ppu.SprRam.Peek(memory.Ptr(spriteIndex * 4)))
			deltaY := y - spriteY
			spriteHeight := 8
			if ppu.registers.ctrl&PPUCtrl_SpriteSize != 0 {
				spriteHeight = 16
			}
			if deltaY < 0 || deltaY >= spriteHeight {
				// sprite is not in range
				continue
			}
			// the sprite is in this scanline
			if ppu.spriteCount >= 8 {
				// sprite overflow
				// TODO: implement hardware bug
				ppu.registers.status |= PPUStatus_SpriteOverflow
				break
			}

			var spriteBytes [4]byte
			for i := 0; i < 4; i++ {
				spriteBytes[i] = ppu.SprRam.Peek(memory.Ptr(spriteIndex*4 + i))
				//ppu.secondaryOAM.Poke(memory.Ptr(spriteCount*4+i), toCopy)
			}
			// evaluate sprite
			sprite := &ppu.sprites[ppu.spriteCount]
			sprite.Id = spriteIndex
			sprite.Y = spriteY
			sprite.TileId = int(ppu.SprRam.Peek(memory.Ptr(spriteIndex*4 + 1)))
			sprite.X = int(ppu.SprRam.Peek(memory.Ptr(spriteIndex*4 + 3)))
			sprite.Attr.Unmarshal(ppu.SprRam.Peek(memory.Ptr(spriteIndex*4 + 2)))
			//sprite.Unmarshal(spriteIndex, &spriteBytes)

			ppu.spriteCount++
		}
		if ppu.spriteCount > 0 {
			logger.Debugf("renderSprites: Scanline #%d has %d sprites.", y, ppu.spriteCount)
		}
		// sprite fetches
		for i := 0; i < ppu.spriteCount; i++ {
			// load tile row to sprite
			sprite := &ppu.sprites[i]
			var addr memory.Ptr
			var spriteHeight int
			if ppu.registers.ctrl&PPUCtrl_SpriteSize == 0 {
				// 8*8 sprite
				spriteHeight = 8
				if ppu.registers.ctrl&PPUCtrl_SpritePatternTable != 0 {
					addr = 0x1000
				}
				addr += memory.Ptr(sprite.TileId * 16)
			} else {
				// 8 * 16 sprite
				spriteHeight = 16
				if sprite.TileId&1 != 0 {
					addr = 0x1000
				}
				addr += memory.Ptr(sprite.TileId & ^1 * 16)
			}
			deltaY := y - sprite.Y
			if sprite.Attr.VerticalFlip {
				deltaY ^= spriteHeight - 1 // i.e. deltaY = spriteHeight - 1 - deltaY
			}
			addr += memory.Ptr(deltaY + deltaY&0x8) // if deltaY >=8, select the 3rd 8-byte

			sprite.TileRowLow = ppu.vram.Peek(addr)
			sprite.TileRowHigh = ppu.vram.Peek(addr + 8)
		}
	} else if x == 319 {
		// Sprite fetches (8 sprites total, 8 cycles per sprite)

	}
}
func (ppu *PPUImpl) Render() {
	scanline := ppu.scanline
	x := ppu.dotInScanline
	//y := scanline - 21
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
		ppu.registers.status &= ^PPUStatus_VBlank
	} else if scanline <= 260 {
		// rendering
		ppu.renderSprites()
		if x < 256 {
			// BG Fetch
			ppu.DrawPixel()
		} else if x < 320 {
			//ppu.DrawPixel(patternEntry0, patternEntry1, attrTableEntry)
			// Sprite Fetch
			// Fetches 4x8 bytes; two dummy Name Table entris, and two Pattern Table bytes; for 1st..8th sprite in NEXT scanline (fetches dummy patterns if the scanline contains less than 8 sprites).
			// http://wiki.nesdev.com/w/index.php/PPU_sprite_evaluation
			// First, it clears the list of sprites to draw
			switch x % 8 {
			case 0:
				// Garbage nametable byte
			case 2:
				// Garbage nametable byte
			case 4:
				// Tile bitmap low
			case 6:
				// Tile bitmap high
			}
		} else if x < 336 {
			// BG Fetch
			//ppu.DrawPixel(patternEntry0, patternEntry1, attrTableEntry)
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

func (ppu *PPUImpl) DrawPixel() {
	x := ppu.dotInScanline
	y := ppu.scanline - 21

	if y >= 0 && y < SCREEN_HEIGHT && x >= 0 && x < SCREEN_WIDTH {
		var currentPalette byte
		// Draw background
		if ppu.registers.mask&PPUMask_BackgroundVisibility != 0 {
			tileId := y/8*32 + x/8
			groupId := y/32*8 + x/32
			offsetY := y % 8
			offsetX := x % 8
			fieldY := y % 32 / 16
			fieldX := x % 32 / 16
			field := fieldY*2 + fieldX
			patternId := ppu.ReadNameTableByte(memory.Ptr(tileId))
			low := ppu.ReadPatternTableByte(memory.Ptr(patternId)*16 + memory.Ptr(offsetY))
			high := ppu.ReadPatternTableByte(memory.Ptr(patternId)*16 + 8 + memory.Ptr(offsetY))
			currentPalette = high>>byte(7-offsetX)&1<<1 | low>>byte(7-offsetX)&1
			if currentPalette > 0 {
				// non-global background color
				attr := ppu.ReadAttributeTableByte(memory.Ptr(groupId))
				paletteId := attr >> (byte(field) * 2) & 3
				currentPalette |= paletteId << 2
			}
		}
		// Draw sprites
		if ppu.registers.mask&PPUMask_SpriteVisibility != 0 {
			// Each four bytes in SPR-RAM define attributes for one sprite
			if x == 0 {
				ppu.spriteShown = 0
				ppu.registers.status &= ^PPUStatus_Sprite0Hit
			}
			// TODO: sprite 0 hit flag, clipping, etc
			for spriteIndex := 0; spriteIndex < ppu.spriteCount; spriteIndex++ {
				sprite := &ppu.sprites[spriteIndex]
				deltaX := x - sprite.X
				if deltaX < 0 || deltaX >= 8 {
					// sprite is not in range
					continue
				}
				if sprite.Attr.HorizontalFlip {
					deltaX ^= 7 // i.e. delta = 7 - delta
				}
				colorLow := sprite.TileRowLow >> byte(7-deltaX) & 1
				colorHigh := sprite.TileRowHigh >> byte(7-deltaX) & 1
				spritePalette := byte(colorLow | colorHigh<<1)
				if spritePalette == 0 {
					// transparent pixel
					continue
				}
				if sprite.Id == 0 && currentPalette != 0 {
					// set sprite 0 hit flag
					ppu.registers.status |= PPUStatus_Sprite0Hit
				}
				if currentPalette != 0 && sprite.Attr.BackgroundPriority {
					// background pixel covers this sprite pixel
					continue
				}
				spritePalette |= byte(sprite.Attr.PaletteId << 2)
				currentPalette = spritePalette + 0x10
				if deltaX == 0 {
					ppu.spriteShown++
				}
				break
			}
			if x == SCREEN_WIDTH-1 {
				if ppu.spriteCount > 0 {
					logger.Debugf("Draw: Scanline #%d shows %d/%d sprites.", y, ppu.spriteShown, ppu.spriteCount)
				}
			}
		}
		ppu.RenderedBuffer[y][x] = Color(ppu.Palette.Peek(0x3F00 + memory.Ptr(currentPalette))).ToGRBColor()
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
