package dt1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

func getTileTypeImage(t d2enum.TileType) string {
	switch t {
	case d2enum.TileFloor:
		return "floor.png"
	case d2enum.TileLeftWall:
		return "wall_west.png"
	case d2enum.TileRightWall:
		return "wall_north.png"
	case d2enum.TileRightPartOfNorthCornerWall:
		return "corner_upper_north.png"
	case d2enum.TileLeftPartOfNorthCornerWall:
		return "corner_upper_west.png"
	case d2enum.TileLeftEndWall:
		return "corner_upper_east.png"
	case d2enum.TileRightEndWall:
		return "corner_lower_south.png"
	case d2enum.TileSouthCornerWall:
		return "corner_lower_east.png"
	case d2enum.TileLeftWallWithDoor:
		return "door_west.png"
	case d2enum.TileRightWallWithDoor:
		return "door_north.png"
	default:
		return ""
	}
}
