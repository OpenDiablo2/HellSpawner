package dt1widget

const (
	subTileFlagBlockWalk = iota
	subTileFlagBlockLOS
	subTileFlagBlockJump
	subTileFlagBlockPlayerWalk
	subTileFlagUnknown1
	subTileFlagBlockLight
	subTileFlagUnknown2
	subTileFlagUnknown3
)

func subTileString(subtile int32) string {
	lookup := map[byte]string{
		1 << 0: "block walk",
		1 << 1: "block light and line of sight",
		1 << 2: "block jump/teleport",
		1 << 3: "block player walk, allow merc walk",
		1 << 4: "unknown #4",
		1 << 5: "block light only",
		1 << 6: "unknown #6",
		1 << 7: "unknown #7",
	}

	str, found := lookup[byte(1<<subtile)]
	if !found {
		return "unknown"
	}

	return str
}

func getFlagFromPos(x, y int) int {
	var subtileLookup = [5][5]int{
		{20, 21, 22, 23, 24},
		{15, 16, 17, 18, 19},
		{10, 11, 12, 13, 14},
		{5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4},
	}

	return subtileLookup[y][x]
}

func (p *widget) getSubTileFieldToEdit(idx int) *bool {
	state := p.getState()

	tileIdx := state.tileGroup

	switch state.subtileFlag {
	case subTileFlagBlockWalk:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].BlockWalk
	case subTileFlagBlockLOS:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].BlockLOS
	case subTileFlagBlockJump:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].BlockJump
	case subTileFlagBlockPlayerWalk:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].BlockPlayerWalk
	case subTileFlagUnknown1:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].Unknown1
	case subTileFlagBlockLight:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].BlockLight
	case subTileFlagUnknown2:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].Unknown2
	case subTileFlagUnknown3:
		return &p.dt1.Tiles[tileIdx].SubTileFlags[idx].Unknown3
	}

	return nil
}
