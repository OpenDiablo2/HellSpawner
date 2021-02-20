package dt1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
)

func decodeWallBlock(block *d2dt1.Block, wallPixBuf *[]byte, tileYOffset, tileWidth int32) {
	// RLE Encoding
	blockX := int32(block.X)
	blockY := int32(block.Y)
	x := int32(0)
	y := int32(0)
	idx := 0
	length := block.Length

	for length > 0 {
		b1 := block.EncodedData[idx]
		b2 := block.EncodedData[idx+1]
		idx += 2
		length -= 2

		if (b1 | b2) == 0 {
			x = 0
			y++

			continue
		}

		x += int32(b1)
		length -= int32(b2)

		for b2 > 0 {
			offset := ((blockY + y + tileYOffset) * tileWidth) + (blockX + x)
			(*wallPixBuf)[offset] = block.EncodedData[idx]
			idx++
			x++
			b2--
		}
	}
}

// nolint:gomnd // 3D isometric decoding
func decodeFloorBlock(block *d2dt1.Block, floorPixBuf *[]byte, tileYOffset, tileWidth int32) {
	xjump := []int32{14, 12, 10, 8, 6, 4, 2, 0, 2, 4, 6, 8, 10, 12, 14}
	nbpix := []int32{4, 8, 12, 16, 20, 24, 28, 32, 28, 24, 20, 16, 12, 8, 4}
	blockX := int32(block.X)
	blockY := int32(block.Y)
	length := int32(256)
	x := int32(0)
	y := int32(0)
	idx := 0

	for length > 0 {
		x = xjump[y]
		n := nbpix[y]
		length -= n

		for n > 0 {
			offset := ((blockY + y + tileYOffset) * tileWidth) + (blockX + x)
			(*floorPixBuf)[offset] = block.EncodedData[idx]
			x++
			n--
			idx++
		}
		y++
	}
}

// this is copied from `OpenDiablo2/d2common/d2fileformats/d2dt1`,
// we want to render the isometric (floor) and rle (wall) pixel buffers separately
func decodeTileGfxData(blocks []d2dt1.Block, floorPixBuf, wallPixBuf *[]byte, tileYOffset, tileWidth int32) {
	for i := range blocks {
		switch blocks[i].Format() {
		case d2dt1.BlockFormatIsometric:
			decodeFloorBlock(&blocks[i], floorPixBuf, tileYOffset, tileWidth)
		case d2dt1.BlockFormatRLE:
			decodeWallBlock(&blocks[i], wallPixBuf, tileYOffset, tileWidth)
		}
	}
}
