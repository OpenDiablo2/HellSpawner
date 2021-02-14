package hswidget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"
)

type ds1EditorState int

const (
	ds1EditorStateViewer ds1EditorState = iota
	ds1EditorStateAddFile
	ds1EditorStateAddObject
	ds1EditorStateConfirm
)

const (
// gridMaxWidth    = 160
// gridMaxHeight   = 80
// gridDivisionsXY = 5
// subtileHeight   = gridMaxHeight / gridDivisionsXY
// subtileWidth    = gridMaxWidth / gridDivisionsXY
)

const (
	imageW, imageH = 32, 32
)

type ds1Controls struct {
	tileX, tileY int32
	object       int32
	// nolint:structcheck // will be used
	subgroup int32
	// nolint:structcheck // will be used
	tile struct {
		floor, wall, shadow, sub int32
	}
}

type DS1AddObjectState struct {
	objType  int32
	objID    int32
	objX     int32
	objY     int32
	objFlags int32
}

// DS1ViewerState represents ds1 viewers state
type DS1ViewerState struct {
	*ds1Controls
	state          ds1EditorState
	confirmDialog  *PopUpConfirmDialog
	newFilePath    string
	addObjectState DS1AddObjectState
}

// Dispose clears viewers state
func (is *DS1ViewerState) Dispose() {
	// noop
}

// DS1Widget represents ds1 viewers widget
type DS1Widget struct {
	id  string
	ds1 *d2ds1.DS1
}

// DS1Viewer creates a new ds1 viewer
func DS1Viewer(id string, ds1 *d2ds1.DS1) *DS1Widget {
	result := &DS1Widget{
		id:  id,
		ds1: ds1,
	}

	return result
}

func (p *DS1Widget) getStateID() string {
	return fmt.Sprintf("DS1Widget_%s", p.id)
}

func (p *DS1Widget) getState() *DS1ViewerState {
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

func (p *DS1Widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *DS1Widget) initState() {
	state := &DS1ViewerState{
		ds1Controls:    &ds1Controls{},
		addObjectState: DS1AddObjectState{},
	}

	p.setState(state)
}

// Build builds a viewer
func (p *DS1Widget) Build() {
	state := p.getState()

	switch state.state {
	case ds1EditorStateViewer:
		p.makeViewerLayout().Build()
	case ds1EditorStateAddFile:
		p.makeAddFileLayout().Build()
	case ds1EditorStateAddObject:
		p.makeAddObjectLayout().Build()
	case ds1EditorStateConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

func (p *DS1Widget) makeViewerLayout() giu.Layout {
	state := p.getState()

	tabs := giu.Layout{
		giu.TabItem("Files").Layout(p.makeFilesLayout(state)),
		giu.TabItem("Objects").Layout(p.makeObjectsLayout(state)),
		giu.TabItem("Tiles").Layout(p.makeTilesLayout(state)),
	}

	if len(p.ds1.SubstitutionGroups) > 0 {
		tabs = append(tabs, giu.TabItem("Substitutions").Layout(p.makeSubstitutionsLayout(state)))
	}

	return giu.Layout{
		p.makeDataLayout(),
		giu.Separator(),
		giu.TabBar("##TabBar_ds1_" + p.id).Layout(tabs),
	}
}

func (p *DS1Widget) makeDataLayout() giu.Layout {
	state := p.getState()
	var version int32 = int32(p.ds1.Version)
	l := giu.Layout{
		giu.Line(
			giu.Label("Version: "),
			giu.InputInt("##"+p.id+"version", &version).Size(30).OnChange(func() {
				state.confirmDialog = NewPopUpConfirmDialog(
					"##"+p.id+"confirmVersionChange",
					"Are you sure, you want to change DS1 Version?",
					"This value is used while decoding and encoding ds1 file\nPlease see github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1/ds1.go\nto get more informations.\ncontinue",
					func() {
						p.ds1.Version = version
						state.state = ds1EditorStateViewer
					},
					func() {
						state.state = ds1EditorStateViewer
					},
				)
				state.state = ds1EditorStateConfirm
			}),
		),
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

func (p *DS1Widget) makeFilesLayout(_ *DS1ViewerState) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	// iterating using the value should not be a big deal as
	// we only expect a handful of strings in this slice.
	for n, str := range p.ds1.Files {
		l = append(l, giu.Layout{
			giu.Line(
				giu.Button("Delete##"+p.id+"DeleteFile"+strconv.Itoa(n)).Size(45, 30).OnClick(func() { p.deleteFile(n) }),
				giu.Label(str),
			),
		})
	}

	return giu.Layout{
		l,
		giu.Button("Add File##"+p.id+"AddFile").Size(100, 30).OnClick(func() {
			state.state = ds1EditorStateAddFile
		}),
	}
}

func (p *DS1Widget) makeObjectsLayout(state *DS1ViewerState) giu.Layout {
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
			giu.ImageWithFile("hsassets/images/shrug.png").Size(imageW, imageH),
		)
		l = append(l, line)
	}

	l = append(l, giu.Separator(),
		giu.Button("Add new object...##"+p.id+"AddObject").Size(200, 30).OnClick(func() {
			state.state = ds1EditorStateAddObject
		}),
	)

	return l
}

func (p *DS1Widget) makeObjectLayout(state *DS1ViewerState) giu.Layout {
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

func (p *DS1Widget) makePathLayout(obj *d2ds1.Object) giu.Layout {
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

func (p *DS1Widget) makeTilesLayout(state *DS1ViewerState) giu.Layout {
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

	numRows := len(p.ds1.Tiles)
	if numRows < 1 {
		return l
	}

	if ty >= numRows {
		state.ds1Controls.tileY = int32(numRows - 1)
		p.setState(state)
	}

	if numCols := len(p.ds1.Tiles[0]); tx >= numCols {
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

func (p *DS1Widget) makeTileLayout(state *DS1ViewerState, t *d2ds1.TileRecord) giu.Layout {
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

// nolint:dupl // yah, thats duplication of makeTileWallLayout but it isn't complete and can be changed
func (p *DS1Widget) makeTileFloorsLayout(state *DS1ViewerState, records []d2ds1.FloorShadowRecord) giu.Layout {
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

func (p *DS1Widget) makeTileFloorLayout(record *d2ds1.FloorShadowRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Prop1: %v", record.Prop1)),
		giu.Label(fmt.Sprintf("Sequence: %v", record.Sequence)),
		giu.Label(fmt.Sprintf("Unknown1: %v", record.Unknown1)),
		giu.Label(fmt.Sprintf("Style: %v", record.Style)),
		giu.Label(fmt.Sprintf("Unknown2: %v", record.Unknown2)),
		giu.Label(fmt.Sprintf("Hidden: %v", record.Hidden())),
		giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
		giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
	}
}

// nolint:dupl // could be changed
func (p *DS1Widget) makeTileWallsLayout(state *DS1ViewerState, records []d2ds1.WallRecord) giu.Layout {
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

func (p *DS1Widget) makeTileWallLayout(record *d2ds1.WallRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Prop1: %v", record.Prop1)),
		giu.Label(fmt.Sprintf("Zero: %v", record.Zero)),
		giu.Label(fmt.Sprintf("Sequence: %v", record.Sequence)),
		giu.Label(fmt.Sprintf("Unknown1: %v", record.Unknown1)),
		giu.Label(fmt.Sprintf("Style: %v", record.Style)),
		giu.Label(fmt.Sprintf("Unknown2: %v", record.Unknown2)),
		giu.Label(fmt.Sprintf("Hidden: %v", record.Hidden())),
		giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
	}
}

// nolint:dupl // no need to change
func (p *DS1Widget) makeTileShadowsLayout(state *DS1ViewerState, records []d2ds1.FloorShadowRecord) giu.Layout {
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

func (p *DS1Widget) makeTileShadowLayout(record *d2ds1.FloorShadowRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Prop1: %v", record.Prop1)),
		giu.Label(fmt.Sprintf("Sequence: %v", record.Sequence)),
		giu.Label(fmt.Sprintf("Unknown1: %v", record.Unknown1)),
		giu.Label(fmt.Sprintf("Style: %v", record.Style)),
		giu.Label(fmt.Sprintf("Unknown2: %v", record.Unknown2)),
		giu.Label(fmt.Sprintf("Hidden: %v", record.Hidden())),
		giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
		giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
	}
}

// nolint:dupl // it is ok
func (p *DS1Widget) makeTileSubsLayout(state *DS1ViewerState, records []d2ds1.SubstitutionRecord) giu.Layout {
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

func (p *DS1Widget) makeTileSubLayout(record *d2ds1.SubstitutionRecord) giu.Layout {
	return giu.Layout{
		giu.Label(fmt.Sprintf("Unknown: %v", record.Unknown)),
	}
}

func (p *DS1Widget) makeSubstitutionsLayout(state *DS1ViewerState) giu.Layout {
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

func (p *DS1Widget) makeSubstitutionLayout(group *d2ds1.SubstitutionGroup) giu.Layout {
	l := giu.Layout{
		giu.Label(fmt.Sprintf("TileX: %d", group.TileX)),
		giu.Label(fmt.Sprintf("TileY: %d", group.TileY)),
		giu.Label(fmt.Sprintf("WidthInTiles: %d", group.WidthInTiles)),
		giu.Label(fmt.Sprintf("HeightInTiles: %d", group.HeightInTiles)),
		giu.Label(fmt.Sprintf("Unknown: 0x%x", group.Unknown)),
	}

	return l
}

func (p *DS1Widget) makeAddFileLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Label("File path:"),
		giu.InputText("##"+p.id+"newFilePath", &state.newFilePath).Size(200),
		giu.Separator(),
		giu.Line(
			giu.Button("Add##"+p.id+"addFileAdd").Size(50, 30).OnClick(func() {
				p.ds1.Files = append(p.ds1.Files, state.newFilePath)
				state.state = ds1EditorStateViewer
			}),
			giu.Button("Cancel##"+p.id+"addFileCancel").Size(50, 30).OnClick(func() {
				state.state = ds1EditorStateViewer
			}),
		),
	}
}

func (p *DS1Widget) makeAddObjectLayout() giu.Layout {
	state := p.getState()
	_ = state

	return giu.Layout{
		giu.Line(
			giu.Label("Type: "),
			giu.InputInt("##"+p.id+"AddObjectType", &state.addObjectState.objType).Size(40),
		),
		giu.Line(
			giu.Label("ID: "),
			giu.InputInt("##"+p.id+"AddObjectID", &state.addObjectState.objID).Size(40),
		),
		giu.Line(
			giu.Label("X: "),
			giu.InputInt("##"+p.id+"AddObjectX", &state.addObjectState.objX).Size(40),
		),
		giu.Line(
			giu.Label("Y: "),
			giu.InputInt("##"+p.id+"AddObjectY", &state.addObjectState.objY).Size(40),
		),
		giu.Line(
			giu.Label("Flags: "),
			giu.InputInt("##"+p.id+"AddObjectFlags", &state.addObjectState.objFlags).Size(40),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##"+p.id+"AddObjectSave").Size(50, 30).OnClick(func() {
				newObject := d2ds1.Object{
					Type:  int(state.addObjectState.objType),
					ID:    int(state.addObjectState.objID),
					X:     int(state.addObjectState.objX),
					Y:     int(state.addObjectState.objY),
					Flags: int(state.addObjectState.objFlags),
				}

				p.ds1.Objects = append(p.ds1.Objects, newObject)

				state.state = ds1EditorStateViewer
			}),
			giu.Button("Cancel##"+p.id+"AddObjectCancel").Size(50, 30).OnClick(func() {
				state.state = ds1EditorStateViewer
			}),
		),
	}
}

func (p *DS1Widget) deleteFile(idx int) {
	newFiles := make([]string, 0)

	for n, i := range p.ds1.Files {
		if n != idx {
			newFiles = append(newFiles, i)
		}
	}

	p.ds1.Files = newFiles
}
