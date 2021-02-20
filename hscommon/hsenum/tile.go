package hsenum

// Tile types
const (
	TileFloor = iota
	TileLeftWall
	TileRightWall
	TileRightPartOfNorthCornerWall
	TileLeftPartOfNorthCornerWall
	TileLeftEndWall
	TileRightEndWall
	TileSouthCornerWall
	TileLeftWallWithDoor
	TileRightWallWithDoor
	TileSpecialTile1
	TileSpecialTile2
	TilePillarsColumnsAndStandaloneObjects
	TileShadow
	TileTree
	TileRoof
	TileLowerWallsEquivalentToLeftWall
	TileLowerWallsEquivalentToRightWall
	TileLowerWallsEquivalentToRightLeftNorthCornerWall
	TileLowerWallsEquivalentToSouthCornerwall
)

// GetTileTypeString returns string of tile type
// nolint:gocyclo // can't reduce
func GetTileTypeString(t int32) string {
	switch t {
	case TileFloor:
		return "floor"
	case TileSpecialTile1, TileSpecialTile2:
		return "special"
	case TileShadow:
		return "shadow"
	case TileTree:
		return "wall/object"
	case TileRoof:
		return "roof"
	case TileLeftWall:
		return "Left Wall"
	case TileRightWall:
		return "Upper Wall"
	case TileRightPartOfNorthCornerWall:
		return "Upper part of an Upper-Left corner"
	case TileLeftPartOfNorthCornerWall:
		return "Left part of an Upper-Left corner"
	case TileLeftEndWall:
		return "Upper-Right corner"
	case TileRightEndWall:
		return "Lower-Left corner"
	case TileSouthCornerWall:
		return "Lower-Right corner"
	case TileLeftWallWithDoor:
		return "Left Wall with Door object, but not always"
	case TileRightWallWithDoor:
		return "Upper Wall with Door object, but not always"
	default:
		return "lower wall ?"
	}
}
