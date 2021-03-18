package dt1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
)

// we want to render the isometric (floor) and rle (wall) pixel buffers separately
func decodeTileGfxData(blocks []d2dt1.Block, floorPixBuf, wallPixBuf *[]byte, tileYOffset, tileWidth int32) {
	for i := range blocks {
		switch blocks[i].Format() {
		case d2dt1.BlockFormatIsometric:
			d2dt1.DecodeTileGfxData([]d2dt1.Block{blocks[i]}, floorPixBuf, tileYOffset, tileWidth)
		case d2dt1.BlockFormatRLE:
			d2dt1.DecodeTileGfxData([]d2dt1.Block{blocks[i]}, wallPixBuf, tileYOffset, tileWidth)
		}
	}
}
