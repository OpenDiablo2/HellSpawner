package hswidget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"
)

const (
// gridMaxWidth    = 160
// gridMaxHeight   = 80
// gridDivisionsXY = 5
// subtileHeight   = gridMaxHeight / gridDivisionsXY
// subtileWidth    = gridMaxWidth / gridDivisionsXY
)

type ds1Controls struct {
	tileX, tileY int32
	object       int32
	// nolint:unused,structcheck // will be used
	subgroup int32
	// nolint:unused,structcheck // will be used
	tile struct {
		floor, wall, shadow, sub int32
	}
}

// DS1ViewerState represents ds1 viewers state
type DS1ViewerState struct {
	*ds1Controls
}

// Dispose clears viewers state
func (is *DS1ViewerState) Dispose() {
	// noop
}

// DS1ViewerWidget represents ds1 viewers widget
type DS1ViewerWidget struct {
	id  string
	ds1 *d2ds1.DS1
}

// DS1Viewer creates a new ds1 viewer
func DS1Viewer(id string, ds1 *d2ds1.DS1) *DS1ViewerWidget {
	result := &DS1ViewerWidget{
		id:  id,
		ds1: ds1,
	}

	return result
}

func (p *DS1ViewerWidget) getStateID() string {
	return fmt.Sprintf("DS1ViewerWidget_%s", p.id)
}

func (p *DS1ViewerWidget) getState() *DS1ViewerState {
	var state *DS1ViewerState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*DS1ViewerState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *DS1ViewerWidget) setState(s *DS1ViewerState) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *DS1ViewerWidget) initState() {
	state := &DS1ViewerState{
		ds1Controls: &ds1Controls{},
	}

	p.setState(state)
}

// Build builds a viewer
func (p *DS1ViewerWidget) Build() {
	state := p.getState()

	tabs := giu.Layout{
		giu.TabItem("Files").Layout(p.makeFilesLayout(state)),
		giu.TabItem("Objects").Layout(p.makeObjectsLayout(state)),
		giu.TabItem("Tiles").Layout(p.makeTilesLayout(state)),
	}

	if len(p.ds1.SubstitutionGroups) > 0 {
		tabs = append(tabs, giu.TabItem("Substitutions").Layout(p.makeSubstitutionsLayout(state)))
	}

	giu.Layout{
		p.makeDataLayout(),
		giu.Separator(),
		giu.TabBar("##TabBar_ds1_" + p.id).Layout(tabs),
	}.Build()
}

func (p *DS1ViewerWidget) makeDataLayout() giu.Layout {
	l := giu.Layout{
		giu.Label(fmt.Sprintf("Version: %d", p.ds1.Version)),
		giu.Label(fmt.Sprintf("Size: %d x %d tiles", p.ds1.Width, p.ds1.Height)),
		giu.Label(fmt.Sprintf("Substitution Type: %d", p.ds1.SubstitutionType)),
		giu.Separator(),
		giu.Label("Number of"),
		giu.Label(fmt.Sprintf("\tWall Layers: %d", p.ds1.NumberOfWalls)),
		giu.Label(fmt.Sprintf("\tFloor Layers: %d", p.ds1.NumberOfFloors)),
		giu.Label(fmt.Sprintf("\tShadow Layers: %d", p.ds1.NumberOfShadowLayers)),
		giu.Label(fmt.Sprintf("\tSubstitution Layers: %d", p.ds1.NumberOfSubstitutionLayers)),
	}

	return l
}

func (p *DS1ViewerWidget) makeFilesLayout(_ *DS1ViewerState) giu.Layout {
	l := giu.Layout{}

	// iterating using the value should not be a big deal as
	// we only expect a handful of strings in this slice.
	for _, str := range p.ds1.Files {
		l = append(l, giu.Label(str))
	}

	return l
}

func (p *DS1ViewerWidget) makeObjectsLayout(state *DS1ViewerState) giu.Layout {
	numObjects := int32(len(p.ds1.Objects))

	l := giu.Layout{}

	if numObjects > 1 {
		l = append(l, giu.SliderInt("Object Index", &state.object, 0, numObjects-1))
	}

	if numObjects > 0 {
		l = append(l, p.makeObjectLayout(state))
	} else {
		line := giu.Line(
			giu.Label("No objects."),
			giu.ImageWithFile("hsassets/images/shrug.png").Size(32, 32),
		)
		l = append(l, line)
	}

	return l
}

func (p *DS1ViewerWidget) makeObjectLayout(state *DS1ViewerState) giu.Layout {
	objIdx := int(state.object)

	if objIdx >= len(p.ds1.Objects) {
		state.ds1Controls.object = int32(len(p.ds1.Objects) - 1)
		p.setState(state)
	} else if objIdx < 0 {
		state.ds1Controls.object = 0
		p.setState(state)
	}

	obj := p.ds1.Objects[int(state.ds1Controls.object)]

	l := giu.Layout{
		giu.Label(fmt.Sprintf("Type: %d", obj.Type)),
		giu.Label(fmt.Sprintf("ID: %d", obj.ID)),
		giu.Label(fmt.Sprintf("Position: (%d, %d) tiles", obj.X, obj.Y)),
		giu.Label(fmt.Sprintf("Flags: 0x%X", obj.Flags)),
	}

	if len(obj.Paths) > 0 {
		l = append(l, giu.Dummy(1, 16), p.makePathLayout(&obj))
	}

	return l
}

func (p *DS1ViewerWidget) makePathLayout(obj *d2ds1.Object) giu.Layout {
	rowWidgets := make([]*giu.RowWidget, 0)

	rowWidgets = append(rowWidgets, giu.Row(
		giu.Label("Index"),
		giu.Label("Position"),
		giu.Label("Action"),
	))

	for idx := range obj.Paths {
		x, y := obj.Paths[idx].Position.X(), obj.Paths[idx].Position.Y()
		rowWidgets = append(rowWidgets, giu.Row(
			giu.Label(fmt.Sprintf("%d", idx)),
			giu.Label(fmt.Sprintf("(%d, %d)", int(x), int(y))),
			giu.Label(fmt.Sprintf("%d", obj.Paths[idx].Action)),
		))
	}

	return giu.Layout{
		giu.Label("Path Points:"),
		giu.FastTable("").Border(true).Rows(rowWidgets),
	}
}

func (p *DS1ViewerWidget) makeTilesLayout(state *DS1ViewerState) giu.Layout {
	l := giu.Layout{}

	tx, ty := int(state.tileX), int(state.tileY)

	if ty < 0 {
		state.ds1Controls.tileY = 0
		p.setState(state)
	}

	if tx < 0 {
		state.ds1Controls.tileX = 0
		p.setState(state)
	}

	numRows, numCols := 0, 0

	numRows = len(p.ds1.Tiles)
	if numRows < 1 {
		return l
	}

	if ty >= numRows {
		state.ds1Controls.tileY = int32(numRows - 1)
		p.setState(state)
	}

	numCols = len(p.ds1.Tiles[0])
	if tx >= numCols {
		state.ds1Controls.tileX = int32(numCols - 1)
		p.setState(state)
	}

	tx, ty = int(state.tileX), int(state.tileY)

	l = append(
		l, giu.SliderInt("Tile X", &state.ds1Controls.tileX, 0, p.ds1.Width-1),
		giu.SliderInt("Tile Y", &state.ds1Controls.tileY, 0, p.ds1.Height-1),
		p.makeTileLayout(state, &p.ds1.Tiles[ty][tx]),
	)

	return l
}

func (p *DS1ViewerWidget) makeTileLayout(state *DS1ViewerState, t *d2ds1.TileRecord) giu.Layout {
	tabs := giu.Layout{}

	if len(t.Floors) > 0 {
		tabs = append(tabs, giu.TabItem("Floors").Layout(p.makeTileFloorsLayout(state, t.Floors)))
	}

	if len(t.Walls) > 0 {
		tabs = append(tabs, giu.TabItem("Walls").Layout(p.makeTileWallsLayout(state, t.Walls)))
	}

	if len(t.Shadows) > 0 {
		tabs = append(tabs, giu.TabItem("Shadows").Layout(p.makeTileShadowsLayout(state, t.Shadows)))
	}

	if len(t.Substitutions) > 0 {
		tabs = append(tabs, giu.TabItem("Subs").Layout(p.makeTileSubsLayout(state, t.Substitutions)))
	}

	return giu.Layout{
		giu.TabBar("##TabBar_ds1_tiles" + p.id).Layout(tabs),
	}
}

func (p *DS1ViewerWidget) makeTileFloorsLayout(state *DS1ViewerState, records []d2ds1.FloorShadowRecord) giu.Layout {
	l := giu.Layout{}

	if len(records) == 0 {
		return l
	}

	recordIdx := int(state.tile.floor)
	numRecords := len(records)

	if recordIdx >= numRecords {
		recordIdx = numRecords - 1
		state.tile.floor = int32(recordIdx)
		p.setState(state)
	} else if recordIdx < 0 {
		recordIdx = 0
		state.tile.floor = int32(recordIdx)
		p.setState(state)
	}

	if numRecords > 1 {
		l = append(l, giu.SliderInt("Floor", &state.tile.floor, 0, int32(numRecords-1)))
	}

	l = append(l, p.makeTileFloorLayout(&records[recordIdx]))

	return l
}

func (p *DS1ViewerWidget) makeTileFloorLayout(record *d2ds1.FloorShadowRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Prop1: %v", record.Prop1)),
		giu.Label(fmt.Sprintf("Sequence: %v", record.Sequence)),
		giu.Label(fmt.Sprintf("Unknown1: %v", record.Unknown1)),
		giu.Label(fmt.Sprintf("Style: %v", record.Style)),
		giu.Label(fmt.Sprintf("Unknown2: %v", record.Unknown2)),
		giu.Label(fmt.Sprintf("Hidden: %v", record.Hidden)),
		giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
		giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
	}
}

func (p *DS1ViewerWidget) makeTileWallsLayout(state *DS1ViewerState, records []d2ds1.WallRecord) giu.Layout {
	l := giu.Layout{}

	if len(records) == 0 {
		return l
	}

	recordIdx := int(state.tile.wall)
	numRecords := len(records)

	if recordIdx >= numRecords {
		recordIdx = numRecords - 1
		state.tile.wall = int32(recordIdx)
		p.setState(state)
	} else if recordIdx < 0 {
		recordIdx = 0
		state.tile.wall = int32(recordIdx)
		p.setState(state)
	}

	if numRecords > 1 {
		l = append(l, giu.SliderInt("Wall", &state.tile.wall, 0, int32(numRecords-1)))
	}

	l = append(l, p.makeTileWallLayout(&records[recordIdx]))

	return l
}

func (p *DS1ViewerWidget) makeTileWallLayout(record *d2ds1.WallRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Prop1: %v", record.Prop1)),
		giu.Label(fmt.Sprintf("Zero: %v", record.Zero)),
		giu.Label(fmt.Sprintf("Sequence: %v", record.Sequence)),
		giu.Label(fmt.Sprintf("Unknown1: %v", record.Unknown1)),
		giu.Label(fmt.Sprintf("Style: %v", record.Style)),
		giu.Label(fmt.Sprintf("Unknown2: %v", record.Unknown2)),
		giu.Label(fmt.Sprintf("Hidden: %v", record.Hidden)),
		giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
	}
}

func (p *DS1ViewerWidget) makeTileShadowsLayout(state *DS1ViewerState, records []d2ds1.FloorShadowRecord) giu.Layout {
	l := giu.Layout{}

	if len(records) == 0 {
		return l
	}

	recordIdx := int(state.tile.shadow)
	numRecords := len(records)

	if recordIdx >= numRecords {
		recordIdx = numRecords - 1
		state.tile.shadow = int32(recordIdx)
		p.setState(state)
	} else if recordIdx < 0 {
		recordIdx = 0
		state.tile.shadow = int32(recordIdx)
		p.setState(state)
	}

	if numRecords > 1 {
		l = append(l, giu.SliderInt("Shadow", &state.tile.shadow, 0, int32(numRecords-1)))
	}

	l = append(l, p.makeTileShadowLayout(&records[recordIdx]))

	return l
}

func (p *DS1ViewerWidget) makeTileShadowLayout(record *d2ds1.FloorShadowRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Prop1: %v", record.Prop1)),
		giu.Label(fmt.Sprintf("Sequence: %v", record.Sequence)),
		giu.Label(fmt.Sprintf("Unknown1: %v", record.Unknown1)),
		giu.Label(fmt.Sprintf("Style: %v", record.Style)),
		giu.Label(fmt.Sprintf("Unknown2: %v", record.Unknown2)),
		giu.Label(fmt.Sprintf("Hidden: %v", record.Hidden)),
		giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
		giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
	}
}

func (p *DS1ViewerWidget) makeTileSubsLayout(state *DS1ViewerState, records []d2ds1.SubstitutionRecord) giu.Layout {
	l := giu.Layout{}

	if len(records) == 0 {
		return l
	}

	recordIdx := int(state.tile.sub)
	numRecords := len(records)

	if recordIdx >= numRecords {
		recordIdx = numRecords - 1
		state.tile.sub = int32(recordIdx)
		p.setState(state)
	} else if recordIdx < 0 {
		recordIdx = 0
		state.tile.sub = int32(recordIdx)
		p.setState(state)
	}

	if numRecords > 1 {
		l = append(l, giu.SliderInt("Substitution", &state.tile.sub, 0, int32(numRecords-1)))
	}

	l = append(l, p.makeTileSubLayout(&records[recordIdx]))

	return l
}

func (p *DS1ViewerWidget) makeTileSubLayout(record *d2ds1.SubstitutionRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Unknown: %v", record.Unknown)),
	}
}

func (p *DS1ViewerWidget) makeSubstitutionsLayout(state *DS1ViewerState) giu.Layout {
	l := giu.Layout{}

	recordIdx := int(state.subgroup)
	numRecords := len(p.ds1.SubstitutionGroups)

	if p.ds1.SubstitutionGroups == nil || numRecords == 0 {
		return l
	}

	if recordIdx >= numRecords {
		recordIdx = numRecords - 1
		state.subgroup = int32(recordIdx)
		p.setState(state)
	} else if recordIdx < 0 {
		recordIdx = 0
		state.subgroup = int32(recordIdx)
		p.setState(state)
	}

	if numRecords > 1 {
		l = append(l, giu.SliderInt("Substitution", &state.subgroup, 0, int32(numRecords-1)))
	}

	l = append(l, p.makeSubstitutionLayout(&p.ds1.SubstitutionGroups[recordIdx]))

	return l
}

func (p *DS1ViewerWidget) makeSubstitutionLayout(group *d2ds1.SubstitutionGroup) giu.Layout {
	l := giu.Layout{
		giu.Label(fmt.Sprintf("TileX: %d", group.TileX)),
		giu.Label(fmt.Sprintf("TileY: %d", group.TileY)),
		giu.Label(fmt.Sprintf("WidthInTiles: %d", group.WidthInTiles)),
		giu.Label(fmt.Sprintf("HeightInTiles: %d", group.HeightInTiles)),
		giu.Label(fmt.Sprintf("Unknown: 0x%x", group.Unknown)),
	}

	return l
}
