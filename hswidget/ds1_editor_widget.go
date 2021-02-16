package hswidget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math/d2vector"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2path"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	layerDeleteButtonSize                = 24
	inputIntW                            = 40
	filePathW                            = 200
	deleteButtonSize                     = 15
	actionButtonW, actionButtonH         = 200, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	trueFalseListW                       = 60
)

const (
	maxByteSize = 255
)

type ds1EditorState int

const (
	ds1EditorStateViewer ds1EditorState = iota
	ds1EditorStateAddFile
	ds1EditorStateAddObject
	ds1EditorStateAddPath
	ds1EditorStateAddFloorShadow
	ds1EditorStateAddWall
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

// DS1AddObjectState represents state of new object
type DS1AddObjectState struct {
	objType  int32
	objID    int32
	objX     int32
	objY     int32
	objFlags int32
}

// Dispose clears state
func (t *DS1AddObjectState) Dispose() {
	// noop
}

// DS1AddPathState contains data about new path
type DS1AddPathState struct {
	pathAction int32
	pathX      int32
	pathY      int32
}

// Dispose clears state
func (t *DS1AddPathState) Dispose() {
	// noop
}

// DS1AddFloorShadowState contains data used in
// add floor-shadow record dialog
type DS1AddFloorShadowState struct {
	prop1    int32
	sequence int32
	unknown1 int32
	style    int32
	unknown2 int32
	hidden   int32
	cb       func()
}

// Dispose resets DS1AddFloorShadowState
func (t DS1AddFloorShadowState) Dispose() {
	t.prop1 = 0
	t.sequence = 0
	t.unknown1 = 0
	t.style = 0
	t.unknown2 = 0
	t.hidden = 0
}

// DS1AddWallState contains data used in add wall dialog
type DS1AddWallState struct {
	tileType int32
	zero     int32
	DS1AddFloorShadowState
}

// Dispose cleans DS1AddWallState
func (t *DS1AddWallState) Dispose() {
	t.DS1AddFloorShadowState.Dispose()
}

// DS1ViewerState represents ds1 viewers state
type DS1ViewerState struct {
	*ds1Controls
	state               ds1EditorState
	confirmDialog       *PopUpConfirmDialog
	newFilePath         string
	addObjectState      DS1AddObjectState
	addPathState        DS1AddPathState
	addFloorShadowState DS1AddFloorShadowState
	addWallState        DS1AddWallState
}

// Dispose clears viewers state
func (is *DS1ViewerState) Dispose() {
	is.addObjectState.Dispose()
	is.addPathState.Dispose()
	is.addFloorShadowState.Dispose()
	is.addWallState.Dispose()
}

// DS1Widget represents ds1 viewers widget
type DS1Widget struct {
	id                  string
	ds1                 *d2ds1.DS1
	deleteButtonTexture *giu.Texture
}

// DS1Viewer creates a new ds1 viewer
func DS1Viewer(id string, ds1 *d2ds1.DS1, dbt *giu.Texture) *DS1Widget {
	result := &DS1Widget{
		id:                  id,
		ds1:                 ds1,
		deleteButtonTexture: dbt,
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
	case ds1EditorStateAddPath:
		p.makeAddPathLayout().Build()
	case ds1EditorStateAddFloorShadow:
		p.makeAddFloorShadowLayout().Build()
	case ds1EditorStateAddWall:
		p.makeAddWallLayout().Build()
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
	var version int32 = p.ds1.Version

	state := p.getState()

	l := giu.Layout{
		giu.Line(
			giu.Label("Version: "),
			giu.InputInt("##"+p.id+"version", &version).Size(inputIntW).OnChange(func() {
				state.confirmDialog = NewPopUpConfirmDialog(
					"##"+p.id+"confirmVersionChange",
					"Are you sure, you want to change DS1 Version?",
					"This value is used while decoding and encoding ds1 file\n"+
						"Please see github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1/ds1.go\n"+
						"to get more informations.\n\n"+
						"Continue?",
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
		currentIdx := n

		l = append(l, giu.Layout{
			giu.Line(
				hsutil.MakeImageButton(
					"##"+p.id+"DeleteFile"+strconv.Itoa(currentIdx),
					deleteButtonSize, deleteButtonSize,
					p.deleteButtonTexture,
					func() {
						p.deleteFile(currentIdx)
					},
				),
				giu.Label(str),
			),
		})
	}

	return giu.Layout{
		l,
		giu.Separator(),
		giu.Button("Add File##"+p.id+"AddFile").Size(actionButtonW, actionButtonH).OnClick(func() {
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
		giu.Button("Add new object...##"+p.id+"AddObject").Size(actionButtonW, actionButtonH).OnClick(func() {
			state.state = ds1EditorStateAddObject
		}),
		giu.Button("Add new path...##"+p.id+"AddPath").Size(actionButtonW, actionButtonH).OnClick(func() {
			state.state = ds1EditorStateAddPath
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
		l = append(
			l,
			giu.Dummy(1, 16),
			p.makePathLayout(&obj),
		)
	}

	return l
}

func (p *DS1Widget) makePathLayout(obj *d2ds1.Object) giu.Layout {
	rowWidgets := make([]*giu.RowWidget, 0)

	rowWidgets = append(rowWidgets, giu.Row(
		giu.Label("Index"),
		giu.Label("Position"),
		giu.Label("Action"),
		giu.Label(""),
	))

	for idx := range obj.Paths {
		currentIdx := idx
		x, y := obj.Paths[idx].Position.X(), obj.Paths[idx].Position.Y()
		rowWidgets = append(rowWidgets, giu.Row(
			giu.Label(fmt.Sprintf("%d", idx)),
			giu.Label(fmt.Sprintf("(%d, %d)", int(x), int(y))),
			giu.Label(fmt.Sprintf("%d", obj.Paths[idx].Action)),
			hsutil.MakeImageButton("##"+p.id+"deletePath",
				deleteButtonSize, deleteButtonSize,
				p.deleteButtonTexture,
				func() {
					p.deletePath(currentIdx)
				},
			),
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

	numRows, numCols := 0, 0

	numRows = len(p.ds1.Tiles)
	if numRows == 0 {
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
		tabs = append(
			tabs,
			giu.TabItem("Floors").Layout(giu.Layout{
				p.makeTileFloorsLayout(state, t.Floors),
				giu.Separator(),
				giu.Line(
					giu.Button("Edit floor##"+p.id+"editFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
						state.addFloorShadowState.cb = func() {
							newFloor := p.createFloorShadowRecord()

							p.ds1.Tiles[state.tileY][state.tileY].Floors[state.object] = newFloor
						}
						state.state = ds1EditorStateAddFloorShadow
					}),
					giu.Button("Add floor##"+p.id+"addFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
						state.addFloorShadowState.cb = func() {
							newFloor := p.createFloorShadowRecord()

							p.ds1.Tiles[state.tileY][state.tileY].Floors = append(p.ds1.Tiles[state.tileY][state.tileY].Floors, newFloor)

							p.ds1.NumberOfFloors++
						}
						state.state = ds1EditorStateAddFloorShadow
					}),
					hsutil.MakeImageButton(
						"##"+p.id+"deleteFloor",
						layerDeleteButtonSize, layerDeleteButtonSize,
						p.deleteButtonTexture,
						func() {
							p.deleteFloorRecord()
							p.ds1.NumberOfFloors--
							p.recreateLayerStreamTypes()
						},
					),
				),
			}),
		)
	}

	if len(t.Walls) > 0 {
		tabs = append(
			tabs,
			giu.TabItem("Walls").Layout(giu.Layout{
				p.makeTileWallsLayout(state, t.Walls),
				giu.Button("Edit wall##"+p.id+"addFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
					state.addFloorShadowState.cb = func() {
						newWall := p.createWallRecord()

						p.ds1.Tiles[state.tileY][state.tileY].Walls[state.object] = newWall
					}
					state.state = ds1EditorStateAddWall
				}),
			}),
		)
	}

	if len(t.Shadows) > 0 {
		tabs = append(
			tabs,
			giu.TabItem("Shadows").Layout(giu.Layout{
				p.makeTileShadowsLayout(state, t.Shadows),
				giu.Button("Edit shadow##"+p.id+"addFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
					state.addFloorShadowState.cb = func() {
						newShadow := p.createFloorShadowRecord()

						p.ds1.Tiles[state.tileY][state.tileY].Shadows[state.object] = newShadow
					}

					state.state = ds1EditorStateAddFloorShadow
				}),
			}),
		)
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
		giu.InputText("##"+p.id+"newFilePath", &state.newFilePath).Size(filePathW),
		giu.Separator(),
		giu.Line(
			giu.Button("Add##"+p.id+"addFileAdd").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				p.ds1.Files = append(p.ds1.Files, state.newFilePath)
				state.state = ds1EditorStateViewer
			}),
			giu.Button("Cancel##"+p.id+"addFileCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.state = ds1EditorStateViewer
			}),
		),
	}
}

func (p *DS1Widget) makeAddObjectLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Line(
			giu.Label("Type: "),
			giu.InputInt("##"+p.id+"AddObjectType", &state.addObjectState.objType).Size(inputIntW),
		),
		giu.Line(
			giu.Label("ID: "),
			giu.InputInt("##"+p.id+"AddObjectID", &state.addObjectState.objID).Size(inputIntW),
		),
		giu.Line(
			giu.Label("X: "),
			giu.InputInt("##"+p.id+"AddObjectX", &state.addObjectState.objX).Size(inputIntW),
		),
		giu.Line(
			giu.Label("Y: "),
			giu.InputInt("##"+p.id+"AddObjectY", &state.addObjectState.objY).Size(inputIntW),
		),
		giu.Line(
			giu.Label("Flags: "),
			giu.InputInt("##"+p.id+"AddObjectFlags", &state.addObjectState.objFlags).Size(inputIntW),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##"+p.id+"AddObjectSave").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
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
			giu.Button("Cancel##"+p.id+"AddObjectCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.state = ds1EditorStateViewer
			}),
		),
	}
}

func (p *DS1Widget) makeAddPathLayout() giu.Layout {
	state := p.getState()

	// https://github.com/OpenDiablo2/OpenDiablo2/issues/811
	// this list should be created like in COFWidget.makeAddLayerLayout
	actionsList := []string{"1", "2", "3"}

	return giu.Layout{
		giu.Line(
			giu.Label("Action: "),
			giu.Combo("##"+p.id+"newPathAction",
				actionsList[state.addPathState.pathAction],
				actionsList, &state.addPathState.pathAction,
			).Size(bigListW),
		),
		giu.Label("Vector:"),
		giu.Line(
			giu.Label("\tX: "),
			giu.InputInt("##"+p.id+"newPathX", &state.addPathState.pathX).Size(inputIntW),
		),
		giu.Line(
			giu.Label("\tY: "),
			giu.InputInt("##"+p.id+"newPathY", &state.addPathState.pathY).Size(inputIntW),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##"+p.id+"AddPathSave").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				newPath := d2path.Path{
					// nolint:gomnd // npc actions starts from 1
					Action: int(state.addPathState.pathAction) + 1,
					Position: d2vector.NewPosition(
						float64(state.addPathState.pathX),
						float64(state.addPathState.pathY),
					),
				}

				p.ds1.Objects[state.object].Paths = append(p.ds1.Objects[state.object].Paths, newPath)

				state.state = ds1EditorStateViewer
			}),
			giu.Button("Cancel##"+p.id+"AddPathCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.state = ds1EditorStateViewer
			}),
		),
	}
}

func (p *DS1Widget) makeAddFloorShadowLayout() giu.Layout {
	state := p.getState()

	trueFalseList := []string{"false", "true"}

	return giu.Layout{
		giu.Line(
			giu.Label("Prop 1: "),
			giu.InputInt("##"+p.id+"addFloorShadowProp1", &state.addFloorShadowState.prop1).Size(inputIntW).OnChange(func() {
				if state.addFloorShadowState.prop1 > maxByteSize {
					state.addFloorShadowState.prop1 = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Sequence: "),
			giu.InputInt("##"+p.id+"addFloorShadowSequence", &state.addFloorShadowState.sequence).Size(inputIntW).OnChange(func() {
				if state.addFloorShadowState.sequence > maxByteSize {
					state.addFloorShadowState.sequence = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Unknown 1: "),
			giu.InputInt("##"+p.id+"addFloorShadowUnknown1", &state.addFloorShadowState.unknown1).Size(inputIntW).OnChange(func() {
				if state.addFloorShadowState.unknown1 > maxByteSize {
					state.addFloorShadowState.unknown1 = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Style: "),
			giu.InputInt("##"+p.id+"addFloorShadowStyle", &state.addFloorShadowState.style).Size(inputIntW).OnChange(func() {
				if state.addFloorShadowState.style > maxByteSize {
					state.addFloorShadowState.style = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Unknown 2: "),
			giu.InputInt("##"+p.id+"addFloorShadowUnknown2", &state.addFloorShadowState.unknown2).Size(inputIntW).OnChange(func() {
				if state.addFloorShadowState.unknown2 > maxByteSize {
					state.addFloorShadowState.unknown2 = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Hidden: "),
			giu.Combo(
				"##"+p.id+"addFloorShadowHidden",
				trueFalseList[state.addFloorShadowState.hidden],
				trueFalseList, &state.addFloorShadowState.hidden,
			).Size(trueFalseListW),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##"+p.id+"AddFloorShadowSave").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.addFloorShadowState.cb()
				state.state = ds1EditorStateViewer
			}),
			giu.Button("Cancel##"+p.id+"AddFloorShadowCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.state = ds1EditorStateViewer
			}),
		),
	}
}

func (p *DS1Widget) makeAddWallLayout() giu.Layout {
	state := p.getState()

	// nolint:gomnd // enumeration of tile types starts from 0, but we must give length (starts from 1) in argument
	tileTypeList := make([]string, d2enum.TileLowerWallsEquivalentToSouthCornerwall+1)
	for i := d2enum.TileFloor; i <= d2enum.TileLowerWallsEquivalentToSouthCornerwall; i++ {
		// this list should be a group of strings, which describes d2enum.TileType
		tileTypeList[int(i)] = strconv.Itoa(int(i))
	}

	return giu.Layout{
		giu.Line(
			giu.Label("Type: "),
			giu.Combo("##"+p.id+"AddWallType", tileTypeList[int(state.addWallState.tileType)], tileTypeList,
				&state.addWallState.tileType).Size(bigListW),
		),
		giu.Line(
			giu.Label("Zero: "),
			giu.InputInt("##"+p.id+"AddWallZero", &state.addWallState.zero).Size(inputIntW).OnChange(func() {
				if state.addWallState.zero > maxByteSize {
					state.addWallState.zero = maxByteSize
				}
			}),
		),
		// this fields are constant for flor, shadow and wall
		p.makeAddFloorShadowLayout(),
	}
}

func (p *DS1Widget) createFloorShadowRecord() d2ds1.FloorShadowRecord {
	state := p.getState()

	newFloorShadowRecord := d2ds1.FloorShadowRecord{
		Prop1:       byte(state.addFloorShadowState.prop1),
		Sequence:    byte(state.addFloorShadowState.sequence),
		Unknown1:    byte(state.addFloorShadowState.unknown1),
		Style:       byte(state.addFloorShadowState.style),
		Unknown2:    byte(state.addFloorShadowState.unknown2),
		HiddenBytes: byte(state.addFloorShadowState.hidden),
	}

	return newFloorShadowRecord
}

func (p *DS1Widget) createWallRecord() d2ds1.WallRecord {
	state := p.getState()

	newWall := d2ds1.WallRecord{
		Type:        d2enum.TileType(state.addWallState.tileType),
		Zero:        byte(state.addWallState.zero),
		Prop1:       byte(state.addWallState.prop1),
		Sequence:    byte(state.addWallState.sequence),
		Unknown1:    byte(state.addWallState.unknown1),
		Style:       byte(state.addWallState.style),
		Unknown2:    byte(state.addWallState.unknown2),
		HiddenBytes: byte(state.addWallState.hidden),
	}

	return newWall
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

func (p *DS1Widget) deletePath(idx int) {
	state := p.getState()

	newPaths := make([]d2path.Path, 0)

	for n, i := range p.ds1.Objects[state.object].Paths {
		if n != idx {
			newPaths = append(newPaths, i)
		}
	}

	p.ds1.Objects[state.object].Paths = newPaths
}

func (p *DS1Widget) deleteFloorRecord() {
	state := p.getState()

	newFloors := make([]d2ds1.FloorShadowRecord, 0)

	for n, i := range p.ds1.Tiles[state.tileY][state.tileX].Floors {
		if n != int(state.object) {
			newFloors = append(newFloors, i)
		}
	}

	p.ds1.Tiles[state.tileY][state.tileX].Floors = newFloors

}

// Warning: this is 1:1 copy from
// github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1.(*DS1).setupLayerStreamType()
// but this method is unexported for now, so...
func (p *DS1Widget) recreateLayerStreamTypes() {
	var layerStream []d2enum.LayerStreamType

	// nolint:gomnd // this is constant version
	// see in OpenDiablo2
	if p.ds1.Version < 4 {
		layerStream = []d2enum.LayerStreamType{
			d2enum.LayerStreamWall1,
			d2enum.LayerStreamFloor1,
			d2enum.LayerStreamOrientation1,
			d2enum.LayerStreamSubstitute,
			d2enum.LayerStreamShadow,
		}
	} else {
		// nolint:gomnd // constant (each wall layer has d2enum.LayerStreamWall and Orientation)
		layerStream = make([]d2enum.LayerStreamType,
			(p.ds1.NumberOfWalls*2)+p.ds1.NumberOfFloors+p.ds1.NumberOfShadowLayers+p.ds1.NumberOfSubstitutionLayers)

		layerIdx := 0
		for i := 0; i < int(p.ds1.NumberOfWalls); i++ {
			layerStream[layerIdx] = d2enum.LayerStreamType(int(d2enum.LayerStreamWall1) + i)
			layerStream[layerIdx+1] = d2enum.LayerStreamType(int(d2enum.LayerStreamOrientation1) + i)
			layerIdx += 2
		}
		for i := 0; i < int(p.ds1.NumberOfFloors); i++ {
			layerStream[layerIdx] = d2enum.LayerStreamType(int(d2enum.LayerStreamFloor1) + i)
			layerIdx++
		}
		if p.ds1.NumberOfShadowLayers > 0 {
			layerStream[layerIdx] = d2enum.LayerStreamShadow
			layerIdx++
		}
		if p.ds1.NumberOfSubstitutionLayers > 0 {
			layerStream[layerIdx] = d2enum.LayerStreamSubstitute
		}
	}

	p.ds1.LayerStreamTypes = layerStream
}
