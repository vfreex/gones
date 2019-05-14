package ppu

import (
	"github.com/vfreex/gones/pkg/emulator/memory"
)

// The logical screen resolution processed by the PPU is 256x240 pixels
// The PPU renders 262 scanlines per frame.
// Each scanline lasts for 341 PPU clock cycles (113.667 CPU clock cycles; 1 CPU cycle = 3 PPU cycles),
// with each clock cycle producing one pixel

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
			spriteY := int(ppu.sprRam.Peek(memory.Ptr(spriteIndex * 4)))
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
				spriteBytes[i] = ppu.sprRam.Peek(memory.Ptr(spriteIndex*4 + i))
				//ppu.secondaryOAM.Poke(memory.Ptr(spriteCount*4+i), toCopy)
			}
			// evaluate sprite
			sprite := &ppu.sprites[ppu.spriteCount]
			sprite.Id = spriteIndex
			sprite.Y = spriteY
			sprite.TileId = int(ppu.sprRam.Peek(memory.Ptr(spriteIndex*4 + 1)))
			sprite.X = int(ppu.sprRam.Peek(memory.Ptr(spriteIndex*4 + 3)))
			sprite.Attr.Unmarshal(ppu.sprRam.Peek(memory.Ptr(spriteIndex*4 + 2)))
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

func (ppu *PPUImpl) drawPixel() {
	x := ppu.dotInScanline - 2
	y := ppu.scanline
	if y >= 0 && y < VisualScanlines && x >= 0 && x < VisualDotsPerScanline {
		var currentPalette byte
		// Draw background
		if ppu.registers.mask&PPUMask_BackgroundVisibility != 0 &&
			(ppu.registers.mask&PPUMask_NoBackgroundClipping != 0 || x >= 8) {
			fineX := ppu.registers.x
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
				break
			}
		}
		color := ppu.Palette.Peek(0x3F00 + memory.Ptr(currentPalette))
		if ppu.registers.mask & PPUMask_Greyscale != 0 {
			color &= 0x30
		}
		ppu.RenderedBuffer[y][x] = Color(color).ToGRBColor()
	}
	ppu.registers.bgHighShift <<= 1
	ppu.registers.bgLowShift <<= 1
	ppu.registers.attrHighShift <<= 1
	ppu.registers.attrLowShift <<= 1
}

func (ppu *PPUImpl) fetchBgTileRow(step int) {
	switch step {
	case 1:
		// fetch nametable
		ntAddr := 0x2000 | ppu.registers.v.Address()&0xfff
		ppu.registers.bgNameLatch = ppu.vram.Peek(ntAddr)
	case 3:
		// fetch attrtable
		v := ppu.registers.v.Address()
		attrAddr := 0x23C0 | v&0x0C00 | v>>4&0x38 | v>>2&0x07
		attr := ppu.vram.Peek(attrAddr)
		if ppu.registers.v.CoarseY()%4 >= 2 {
			attr >>= 4
		}
		if ppu.registers.v.CoarseX()%4 >= 2 {
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
		// fetch bitmap low from pattern table\
		lowAddr := memory.Ptr(ppu.registers.bgNameLatch)*16 + memory.Ptr(ppu.registers.v.FineY())
		if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
			lowAddr |= 0x1000
		}
		ppu.registers.bgLowLatch = ppu.vram.Peek(lowAddr)
	case 7:
		// fetch bitmap high from pattern table
		highAddr := memory.Ptr(ppu.registers.bgNameLatch)*16 + 8 + memory.Ptr(ppu.registers.v.FineY())
		if ppu.registers.ctrl&PPUCtrl_BackgroundPatternTable != 0 {
			highAddr |= 0x1000
		}
		ppu.registers.bgHighLatch = ppu.vram.Peek(highAddr)
		logger.Debugf("at (%v, %v): v=%v", ppu.scanline, ppu.dotInScanline, ppu.registers.v.String())
	}
}

func (ppu *PPUImpl) Step() {
	// http://wiki.nesdev.com/w/index.php/PPU_rendering
	// http://wiki.nesdev.com/w/index.php/File:Ntsc_timing.png
	scanline := ppu.scanline
	dot := ppu.dotInScanline
	switch {
	case scanline == 261: // pre
		switch {
		case dot == 1:
			ppu.registers.status &= ^(PPUStatus_Sprite0Hit | PPUStatus_SpriteOverflow | PPUStatus_VBlank)
		case dot >= 280 && dot <= 304:
			if ppu.registers.mask&(PPUMask_BackgroundVisibility|PPUMask_SpriteVisibility) != 0 {
				ppu.registers.v.SetCoarseY(ppu.registers.t.CoarseY())
				ppu.registers.v.SetFineY(ppu.registers.t.FineY())
				ppu.registers.v.SetNametable(ppu.registers.v.Nametable()&1 | ppu.registers.t.Nametable()&2)
			}
		case dot == 340 && ppu.frame&1 != 0 &&
			ppu.registers.mask&(PPUMask_BackgroundVisibility|PPUMask_SpriteVisibility) != 0:
			//on every odd frame, scanline 0, dot 0 is skipped
			ppu.dotInScanline = 0
			ppu.scanline = 0
			goto end
		}
		fallthrough
	case scanline >= 0 && scanline <= 239: // Visible scanlines
		ppu.renderSprites()
		if dot >= 2 && dot <= 257 || dot >= 322 && dot <= 337 {
			ppu.drawPixel()
			if dot%8 == 1 {
				ppu.fillShifters()
			}
		}
		switch {
		case dot == 0: // idle
		case dot >= 1 && dot <= 256: // fetches 3rd..34th tile in scanline
			if scanline == 261 {
				break
			}
			ppu.fetchBgTileRow((dot - 1) % 8)
			if ppu.registers.mask&(PPUMask_BackgroundVisibility|PPUMask_SpriteVisibility) != 0 {
				if dot%8 == 0 { // dots 8, 16, 24, ..., 256: increase coarseX
					ppu.registers.v.IncreaseCoarseX()
				}
				if dot == 256 { // increase fineY
					ppu.registers.v.IncreaseFineY()
				}
			}
		case dot == 257:
			if ppu.registers.mask&(PPUMask_BackgroundVisibility|PPUMask_SpriteVisibility) != 0 {
				ppu.registers.v.SetCoarseX(ppu.registers.t.CoarseX())
				ppu.registers.v.SetNametable(ppu.registers.v.Nametable()&2 | ppu.registers.t.Nametable()&1)
			}
		case dot >= 257 && dot <= 320: // fetching the sprites on the next scanline
		case dot >= 321 && dot <= 336: // fetching the first two tiles for the next scanline
			ppu.fetchBgTileRow((dot - 321) % 8)
			if ppu.registers.mask&(PPUMask_BackgroundVisibility|PPUMask_SpriteVisibility) != 0 {
				if dot%8 == 0 { // dots 328, 336: increase coarseX
					ppu.registers.v.IncreaseCoarseX()
				}
			}
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
	}
end:
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
