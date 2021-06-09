package dt1widget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/HellSpawner/hswidget/dt1widget/tiletypeimage"
)

func drawTileTypeImage(t d2enum.TileType) giu.Widget {
	tileImageFile := getTileTypeImage(t)

	switch t {
	case d2enum.TileFloor:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor()
		})
	case d2enum.TileLeftWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().WestWall(true)
		})
	case d2enum.TileRightWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().NorthWall(true)
		})
	case d2enum.TileRightPartOfNorthCornerWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().WestWall(false).NorthWall(true)
		})
	case d2enum.TileLeftPartOfNorthCornerWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().WestWall(true).NorthWall(false)
		})
	case d2enum.TileLeftEndWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().EastWall()
		})
	case d2enum.TileRightEndWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().SoathWall()
		})
	case d2enum.TileSouthCornerWall:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().Corner()
		})
	case d2enum.TileLeftWallWithDoor:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().WestDoor()
		})
	case d2enum.TileRightWallWithDoor:
		return giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()
			b := tiletypeimage.TileTypeImage(canvas, pos)
			b.Floor().NorthDoor()
		})
	default:
		return giu.ImageWithFile("./hsassets/images/" + tileImageFile)
	}
}

func getTileTypeImage(t d2enum.TileType) string {
	switch t {
	case d2enum.TileSouthCornerWall:
		return "corner_lower_east.png"
	default:
		return ""
	}
}
