package dt1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
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

// nolint:gocyclo // can't reduce
// https://github.com/OpenDiablo2/OpenDiablo2/pull/1089
func getTileTypeString(t d2enum.TileType) string {
	switch t {
	case d2enum.TileFloor:
		return "floor"
	case d2enum.TileSpecialTile1, d2enum.TileSpecialTile2:
		return "special"
	case d2enum.TileShadow:
		return "shadow"
	case d2enum.TileTree:
		return "wall/object"
	case d2enum.TileRoof:
		return "roof"
	case d2enum.TileLeftWall:
		return "Left Wall"
	case d2enum.TileRightWall:
		return "Upper Wall"
	case d2enum.TileRightPartOfNorthCornerWall:
		return "Upper part of an Upper-Left corner"
	case d2enum.TileLeftPartOfNorthCornerWall:
		return "Left part of an Upper-Left corner"
	case d2enum.TileLeftEndWall:
		return "Upper-Right corner"
	case d2enum.TileRightEndWall:
		return "Lower-Left corner"
	case d2enum.TileSouthCornerWall:
		return "Lower-Right corner"
	case d2enum.TileLeftWallWithDoor:
		return "Left Wall with Door object, but not always"
	case d2enum.TileRightWallWithDoor:
		return "Upper Wall with Door object, but not always"
	default:
		return "lower wall ?"
	}
}
