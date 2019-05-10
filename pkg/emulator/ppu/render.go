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
			spriteY := int(ppu.SprRam.Peek(memory.Ptr(spriteIndex*4)) + 1)
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

func (ppu *PPUImpl) fillShifters() {
	ppu.registers.bgHighShift = ppu.registers.bgHighShift&0xff00 | uint16(ppu.registers.bgHighLatch)
	ppu.registers.bgLowShift = ppu.registers.bgLowShift&0xff00 | uint16(ppu.registers.bgLowLatch)
	ppu.registers.attrHighShift = ppu.registers.attrHighShift&0xff00 | uint16(ppu.registers.attrHighLatch)
	ppu.registers.attrLowShift = ppu.registers.attrLowShift&0xff00 | uint16(ppu.registers.attrLowLatch)
}

func (ppu *PPUImpl) Render() {
	// http://wiki.nesdev.com/w/index.php/File:Ntsc_timing.png
	scanline := ppu.scanline
	x := ppu.dotInScanline
	y := scanline - 21
	lx := x + 16 + int(ppu.registers.hscroll)
	coarseX := lx / 8
	nt := int(ppu.registers.ctrl & PPUCtrl_NameTable)
	if coarseX >= 32 {
		nt ^= 1
		coarseX %= 32
	}
	ly := y + int(ppu.registers.vscroll)
	//if vy < 0 {
	//	vy += 262
	//}
	coarseY := ly / 8
	if coarseY >= 30 {
		nt ^= 2
		coarseY %= 30
	}
	fineY := ly % 8
	ppu.registers.v = uint16(fineY&7<<12 | nt&3<<10 | coarseY&0X1f<<5 | coarseX&0X1f)

	switch {
	case scanline <= 19:
		if scanline == 0 && x == 0 {
			ppu.registers.status |= PPUStatus_VBlank
			if ppu.registers.ctrl&PPUCtrl_NMIOnVBlank != 0 {
				ppu.cpu.NMI = true
			}
		}
	case scanline == 20: // dummy
		if x == 0 {
			ppu.registers.status &= ^(PPUStatus_Sprite0Hit | PPUStatus_SpriteOverflow | PPUStatus_VBlank)
		}
		if x == 339 && ppu.frame&1 != 0 {
			// on every odd frame, the dead cycle at the end is removed
			ppu.dotInScanline = 0
			ppu.scanline++
			break
		}
		//fallthrough
	case scanline <= 260: // visible
		ppu.renderSprites()
		ppu.DrawPixel()
		switch {
		case x == 0: // dot 0: idle
			// ppu.registers.v = ppu.getCurrentNametableAddr()
		case x <= 256: // dot 1-256: fetch background
			switch x % 8 {
			case 1:
				// The shifters are reloaded during ticks 9, 17, 25, ..., 257
				if x >= 9 {
					ppu.fillShifters()
				}
			case 2:
				// fetch nametable
				//nametableEntry := ppu.registers.v&0xfff | 0x2000
				nametableEntry := 0x2000 + nt * 0x400 + coarseY * 32 + coarseX
				ppu.registers.bgNameLatch = ppu.vram.Peek(memory.Ptr(nametableEntry))
			case 4:
				// fetch attrtable
				groupId := coarseY/4*8 + coarseX/4
				attrEntry := ppu.registers.v&0xc00 + 0x23c0 + memory.Ptr(groupId)
				//attrEntry := ntaddr + 0x3c0 + memory.Ptr(groupId)
				attr := ppu.vram.Peek(attrEntry)
				paletteId := attr
				if coarseY%4 >= 2 {
					paletteId >>= 4
				}
				if coarseX%4 >= 2 {
					paletteId >>= 2
				}
				paletteId &= 3
				//if paletteId != paletteId1 {
				//	panic("attr not equal")
				//}
				//if paletteId1&1 != 0 {
				//	ppu.registers.attrLowLatch = 0xff
				//} else {
				//	ppu.registers.attrLowLatch = 0
				//}
				//if paletteId1&2 != 0 {
				//	ppu.registers.attrHighLatch = 0xff
				//} else {
				//	ppu.registers.attrHighLatch = 0
				//}
			case 6:
				// fetch bitmap low from pattern table
				lowAddr := memory.Ptr(ppu.registers.bgNameLatch)*16 + memory.Ptr(fineY) //ppu.registers.v>>12&0x7
				if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
					lowAddr |= 0x1000
				}
				ppu.registers.bgLowLatch = ppu.vram.Peek(lowAddr)
			case 0:
				// fetch bitmap high from pattern table
				highAddr := memory.Ptr(ppu.registers.bgNameLatch)*16 + 8 + memory.Ptr(fineY)// ppu.registers.v>>12&0x7
				if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
					highAddr |= 0x1000
				}
				ppu.registers.bgHighLatch = ppu.vram.Peek(highAddr)
			}
		case x == 257: // 257
			ppu.fillShifters()
			fallthrough
		case x <= 320: // 258-320
		case x <= 336: // 321-336
		default: // 337-340
		}
		//ppu.renderSprites()
		//if scanline > 20 && x < SCREEN_WIDTH {
		//	// BG Fetch
		//	ppu.DrawPixel()
		//}

	case scanline == 261: // post

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
			fineX := ppu.registers.hscroll % 8
			//if fineX < 0 {
			//	panic(fmt.Errorf("fineX = %v", fineX))
			//}
			currentPalette = byte(ppu.registers.bgHighShift>>byte(15-fineX)&1<<1 |
				ppu.registers.bgLowShift>>byte(15-fineX)&1)

			//attr1 := uint8(ppu.registers.attrLowShift>>15 | ppu.registers.attrHighShift>>15<<1)
			//if currentPalette > 0 {
			//	// non-global background color
			//	currentPalette |= attr1
			//}

		}
		// Draw sprites
		if ppu.registers.mask&PPUMask_SpriteVisibility != 0 {
			// Each four bytes in SPR-RAM define attributes for one sprite
			if x == 0 {
				ppu.spriteShown = 0
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

	ppu.registers.bgHighShift <<= 1
	ppu.registers.bgLowShift <<= 1
	//ppu.registers.attrHighShift <<= 1
	//ppu.registers.attrLowShift <<= 1
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
