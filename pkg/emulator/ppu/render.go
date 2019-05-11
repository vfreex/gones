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
	dot := ppu.dotInScanline
	y := ppu.scanline //- 21
	switch {
	case dot == 0:
		ppu.currentSpritesCount = ppu.spriteCount
		ppu.currentSprites = ppu.sprites
	case dot == 64:
		ppu.spriteCount = 0
		ppu.registers.status &= ^PPUStatus_SpriteOverflow
	case dot == 256:
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
	case dot >= 257 && dot <= 320:
		// sprite fetches
		i := (dot - 257) / 8
		if i >= ppu.spriteCount {
			break
		}
		sprite := &ppu.sprites[i]
		switch (dot - 257) % 8 {
		case 5:
			addr := ppu.spriteTileAddr(sprite)
			sprite.TileRowLow = ppu.vram.Peek(addr)
		case 7:
			addr := ppu.spriteTileAddr(sprite)
			sprite.TileRowHigh = ppu.vram.Peek(addr + 8)
		}
	}
}

func (ppu *PPUImpl) spriteTileAddr(sprite *Sprite) memory.Ptr {
	y := ppu.scanline
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
	addr += memory.Ptr(deltaY + deltaY&0x8)
	return addr
}

func (ppu *PPUImpl) fillShifters() {
	ppu.registers.bgHighShift = ppu.registers.bgHighShift&0xff00 | uint16(ppu.registers.bgHighLatch)
	ppu.registers.bgLowShift = ppu.registers.bgLowShift&0xff00 | uint16(ppu.registers.bgLowLatch)
	ppu.registers.attrHighShift = ppu.registers.attrHighShift&0xff00 | uint16(ppu.registers.attrHighLatch)
	ppu.registers.attrLowShift = ppu.registers.attrLowShift&0xff00 | uint16(ppu.registers.attrLowLatch)
}

func (ppu *PPUImpl) DrawPixel() {
	x := ppu.dotInScanline - 2
	y := ppu.scanline
	//ppu.registers.bgHighShift <<= 1
	//ppu.registers.bgLowShift <<= 1
	//ppu.registers.attrHighShift <<= 1
	//ppu.registers.attrLowShift <<= 1
	if y >= 0 && y < SCREEN_HEIGHT && x >= 0 && x < SCREEN_WIDTH {
		var currentPalette byte
		// Draw background
		if ppu.registers.mask&PPUMask_BackgroundVisibility != 0 &&
			(ppu.registers.mask&PPUMask_NoBackgroundClipping != 0 || x >= 8) {
			fineX := ppu.registers.hscroll % 8
			currentPalette = byte(ppu.registers.bgHighShift>>byte(15-fineX)&1<<1 |
				ppu.registers.bgLowShift>>byte(15-fineX)&1)
			if currentPalette > 0 {
				attr := byte(ppu.registers.attrHighShift>>byte(15-fineX)&1<<1 |
					ppu.registers.attrLowShift>>byte(15-fineX)&1)
				currentPalette |= attr & 3 << 2
			}
		}
		// Draw sprites
		if ppu.registers.mask&PPUMask_SpriteVisibility != 0 &&
			(ppu.registers.mask&PPUMask_NoSpriteClipping != 0 || x >= 8) {
			// Each four bytes in SPR-RAM define attributes for one sprite
			if x == 0 {
				ppu.spriteShown = 0
			}
			for spriteIndex := 0; spriteIndex < ppu.currentSpritesCount; spriteIndex++ {
				sprite := &ppu.currentSprites[spriteIndex]
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
			//if x == SCREEN_WIDTH-1 {
			//	if ppu.currentSpritesCount > 0 {
			//		logger.Warnf("Draw: Scanline #%d shows %d/%d sprites.", y, ppu.spriteShown, ppu.currentSpritesCount)
			//	}
			//}
		}
		ppu.RenderedBuffer[y][x] = Color(ppu.Palette.Peek(0x3F00 + memory.Ptr(currentPalette))).ToGRBColor()
	}
	ppu.registers.bgHighShift <<= 1
	ppu.registers.bgLowShift <<= 1
	ppu.registers.attrHighShift <<= 1
	ppu.registers.attrLowShift <<= 1
}

const (
	ScanlinesPerFrame     = 262
	DotsPerScanline       = 341
	DotsPerTileRow        = 8
	VisualDotsPerScanline = 256
	VisualScanlines       = 240
)

func (ppu *PPUImpl) fetchBgTileRow(nt int, coarseX, coarseY, fineY int, cycle int) {
	switch cycle {
	case 1:
		// fetch nametable
		ntAddr := 0x2000 | nt*0x400 | coarseY*32 | coarseX
		ppu.registers.bgNameLatch = ppu.vram.Peek(memory.Ptr(ntAddr))
	case 3:
		// fetch attrtable
		attrAddr := 0x23C0 | nt*0x400 | coarseY/4<<3 | coarseX/4
		attr := ppu.vram.Peek(memory.Ptr(attrAddr))
		if coarseY%4 >= 2 {
			attr >>= 4
		}
		if coarseX%4 >= 2 {
			attr >>= 2
		}
		paletteId := attr & 3
		if paletteId&1 != 0 {
			ppu.registers.attrLowLatch = 0xff
		} else {
			ppu.registers.attrLowLatch = 0
		}
		if paletteId&2 != 0 {
			ppu.registers.attrHighLatch = 0xff
		} else {
			ppu.registers.attrHighLatch = 0
		}
	case 5:
		// fetch bitmap low from pattern table
		lowAddr := memory.Ptr(ppu.registers.bgNameLatch)*16 + memory.Ptr(fineY) //ppu.registers.v>>12&0x7
		if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
			lowAddr |= 0x1000
		}
		ppu.registers.bgLowLatch = ppu.vram.Peek(lowAddr)
	case 7:
		// fetch bitmap high from pattern table
		highAddr := memory.Ptr(ppu.registers.bgNameLatch)*16 + 8 + memory.Ptr(fineY) // ppu.registers.v>>12&0x7
		if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
			highAddr |= 0x1000
		}
		ppu.registers.bgHighLatch = ppu.vram.Peek(highAddr)
	}
}

func (ppu *PPUImpl) Render() {
	// http://wiki.nesdev.com/w/index.php/PPU_rendering
	scanline := ppu.scanline
	dot := ppu.dotInScanline
	switch {
	case scanline >= 0 && scanline <= 239: // Visible scanlines
		ppu.renderSprites()
		if dot >= 2 && dot <= 257 || dot >= 322 && dot <= 337 {
			ppu.DrawPixel()
			if dot%8 == 1 {
				ppu.fillShifters()
			}
		}
		switch {
		case dot == 0: // idle
		case dot >= 1 && dot <= 256: // fetches 3rd..34th tile in scanline
			// (33th tile may be parts visible if BG scrolled, 34th is never visible)
			// dot 1-8: tile 2, dot 9-18: tile3
			nt := int(ppu.registers.ctrl & PPUCtrl_NameTable)
			fy := scanline + int(ppu.registers.vscroll)
			coarseY := fy / 8
			if coarseY >= 30 {
				nt ^= 2
				coarseY %= 30
			}
			fineY := fy % 8
			coarseX := (dot-1)/8 + int(ppu.registers.hscroll)/8 + 2
			if coarseX >= 32 {
				nt ^= 1
				coarseX %= 32
			}
			// http://wiki.nesdev.com/w/index.php/File:Ntsc_timing.png
			ppu.fetchBgTileRow(nt, coarseX, coarseY, fineY, (dot-1)%8)
		case dot >= 257 && dot <= 320: // fetching the sprites on the next scanline
		case dot >= 321 && dot <= 336: // fetching the first two tiles for the next scanline
			nt := int(ppu.registers.ctrl & PPUCtrl_NameTable)
			fy := scanline + int(ppu.registers.vscroll) + 1
			coarseY := fy / 8
			if coarseY >= 30 {
				nt ^= 2
				coarseY %= 30
			}
			fineY := fy % 8
			coarseX := (dot-321)/8 + int(ppu.registers.hscroll)/8
			if coarseX >= 32 {
				nt ^= 1
				coarseX %= 32
			}
			ppu.fetchBgTileRow(nt, coarseX, coarseY, fineY, (dot-321)%8)
		case dot >= 337 && dot <= 340:
		}
	case scanline == 240: // post scanline
	case scanline == 241: // VINT
		if dot == 1 {
			ppu.registers.status |= PPUStatus_VBlank
			if ppu.registers.ctrl&PPUCtrl_NMIOnVBlank != 0 {
				ppu.cpu.NMI = true
			}
		}
	case scanline == 261: // pre
		if dot == 1 {
			ppu.registers.status &= ^(PPUStatus_Sprite0Hit | PPUStatus_SpriteOverflow | PPUStatus_VBlank)
		}
		if dot >= 2 && dot <= 257 || dot >= 322 && dot <= 337 {
			ppu.DrawPixel()
			if dot%8 == 1 {
				ppu.fillShifters()
			}
		}
		switch {
		case dot >= 321 && dot <= 336: // fetching the first two tiles for the next scanline
			nt := int(ppu.registers.ctrl & PPUCtrl_NameTable)
			fy := int(ppu.registers.vscroll)
			coarseY := fy / 8
			if coarseY >= 30 {
				nt ^= 2
				coarseY %= 30
			}
			fineY := fy % 8
			coarseX := (dot-321)/8 + int(ppu.registers.hscroll)/8
			if coarseX >= 32 {
				nt ^= 1
				coarseX %= 32
			}
			ppu.fetchBgTileRow(nt, coarseX, coarseY, fineY, (dot-321)%8)
		case dot == 340 && ppu.frame&1 != 0:
			// on every odd frame, the dead cycle at the end is removed
			ppu.dotInScanline = 0
			ppu.scanline = 0
		}
	}
	ppu.dotInScanline++
	if ppu.dotInScanline >= DotsPerScanline {
		ppu.dotInScanline %= DotsPerScanline
		ppu.scanline++
		if ppu.scanline >= ScanlinesPerFrame {
			ppu.scanline %= ScanlinesPerFrame
			ppu.frame++
		}
	}
}
