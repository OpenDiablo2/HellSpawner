package dt1widget

import (
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
)

func getTileTypeImage(t int32) string {
	switch t {
	case hsenum.TileFloor:
		return "floor.png"
	case hsenum.TileLeftWall:
		return "wall_west.png"
	case hsenum.TileRightWall:
		return "wall_north.png"
	case hsenum.TileRightPartOfNorthCornerWall:
		return "corner_upper_north.png"
	case hsenum.TileLeftPartOfNorthCornerWall:
		return "corner_upper_west.png"
	case hsenum.TileLeftEndWall:
		return "corner_upper_east.png"
	case hsenum.TileRightEndWall:
		return "corner_lower_south.png"
	case hsenum.TileSouthCornerWall:
		return "corner_lower_east.png"
	case hsenum.TileLeftWallWithDoor:
		return "door_west.png"
	case hsenum.TileRightWallWithDoor:
		return "door_north.png"
	default:
		return ""
	}
}
