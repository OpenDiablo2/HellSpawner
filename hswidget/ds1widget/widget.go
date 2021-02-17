package ds1widget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math/d2vector"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2path"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
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
	imageW, imageH                       = 32, 32
)

const (
	maxByteSize = 255
)

const (
// gridMaxWidth    = 160
// gridMaxHeight   = 80
// gridDivisionsXY = 5
// subtileHeight   = gridMaxHeight / gridDivisionsXY
// subtileWidth    = gridMaxWidth / gridDivisionsXY
)

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

func (p *DS1Widget) getState() *DS1State {
	var state *DS1State

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*DS1State)
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
	state := &DS1State{
		ds1Controls: &ds1Controls{},
	}

	p.setState(state)
}

// Build builds a viewer
func (p *DS1Widget) Build() {
	state := p.getState()

	switch state.mode {
	case ds1EditorModeViewer:
		p.makeViewerLayout().Build()
	case ds1EditorModeAddFile:
		p.makeAddFileLayout().Build()
	case ds1EditorModeAddObject:
		p.makeAddObjectLayout().Build()
	case ds1EditorModeAddPath:
		p.makeAddPathLayout().Build()
	case ds1EditorModeAddFloorShadow:
		p.makeAddFloorShadowLayout(&state.addFloorShadowState).Build()
	case ds1EditorModeAddWall:
		p.makeAddWallLayout().Build()
	case ds1EditorModeConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

func (p *DS1Widget) makeViewerLayout() giu.Layout {
	state := p.getState()

	tabs := giu.Layout{
		giu.TabItem("Files").Layout(p.makeFilesLayout()),
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
				state.confirmDialog = hswidget.NewPopUpConfirmDialog(
					"##"+p.id+"confirmVersionChange",
					"Are you sure, you want to change DS1 Version?",
					"This value is used while decoding and encoding ds1 file\n"+
						"Please see github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1/ds1.go\n"+
						"to get more informations.\n\n"+
						"Continue?",
					func() {
						p.ds1.Version = version
						state.mode = ds1EditorModeViewer
					},
					func() {
						state.mode = ds1EditorModeViewer
					},
				)
				state.mode = ds1EditorModeConfirm
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

func (p *DS1Widget) makeFilesLayout() giu.Layout {
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
			state.mode = ds1EditorModeAddFile
		}),
	}
}

func (p *DS1Widget) makeObjectsLayout(state *DS1State) giu.Layout {
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
			state.mode = ds1EditorModeAddObject
		}),
		giu.Button("Add new path...##"+p.id+"AddPath").Size(actionButtonW, actionButtonH).OnClick(func() {
			state.mode = ds1EditorModeAddPath
		}),
	)

	return l
}

func (p *DS1Widget) makeObjectLayout(state *DS1State) giu.Layout {
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

func (p *DS1Widget) makeTilesLayout(state *DS1State) giu.Layout {
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

// nolint:funlen // cannot reduce
func (p *DS1Widget) makeTileLayout(state *DS1State, t *d2ds1.TileRecord) giu.Layout {
	tabs := giu.Layout{}
	editionButtons := giu.Layout{}

	if len(t.Floors) > 0 {
		tabs = append(
			tabs,
			giu.TabItem("Floors").Layout(giu.Layout{
				p.makeTileFloorsLayout(state, t.Floors),
				giu.Separator(),
				giu.Line(
					giu.Button("Edit floor##"+p.id+"editFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
						p.editFloor()
					}),
					giu.Button("Add floor##"+p.id+"addFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
						p.addFloor()
					}),
					hsutil.MakeImageButton(
						"##"+p.id+"deleteFloor",
						layerDeleteButtonSize, layerDeleteButtonSize,
						p.deleteButtonTexture,
						func() {
							p.deleteFloorRecord()
						},
					),
				),
			}),
		)
	} else {
		editionButtons = append(editionButtons,
			giu.Button("Add floor##"+p.id+"addFloor").Size(actionButtonW, actionButtonH).OnClick(func() {
				p.addFloor()
			}),
		)
	}

	if len(t.Walls) > 0 {
		tabs = append(
			tabs,
			giu.TabItem("Walls").Layout(giu.Layout{
				p.makeTileWallsLayout(state, t.Walls),
				giu.Line(
					giu.Button("Edit wall##"+p.id+"editWall").Size(actionButtonW, actionButtonH).OnClick(func() {
						p.editWall()
					}),
					giu.Button("Add wall##"+p.id+"addWallIn").Size(actionButtonW, actionButtonH).OnClick(func() {
						p.addWall()
					}),
					hsutil.MakeImageButton(
						"##"+p.id+"deleteWall",
						layerDeleteButtonSize, layerDeleteButtonSize,
						p.deleteButtonTexture,
						func() {
							p.deleteWall()
						},
					),
				),
			}),
		)
	} else {
		editionButtons = append(editionButtons, giu.Layout{
			giu.Button("Add wall##"+p.id+"addWallOut").Size(actionButtonW, actionButtonH).OnClick(func() {
				p.addWall()
			}),
		})
	}

	if len(t.Shadows) > 0 {
		tabs = append(
			tabs,
			giu.TabItem("Shadows").Layout(giu.Layout{
				p.makeTileShadowsLayout(state, t.Shadows),
				giu.Line(
					giu.Button("Edit shadow##"+p.id+"editShadow").Size(actionButtonW, actionButtonH).OnClick(func() {
						p.editShadow()
					}),
					hsutil.MakeImageButton(
						"##"+p.id+"deleteFloor",
						layerDeleteButtonSize, layerDeleteButtonSize,
						p.deleteButtonTexture,
						func() {
							p.deleteShadow()
						},
					),
				),
			}),
		)
	} else {
		editionButtons = append(editionButtons,
			giu.Button("Add shadow##"+p.id+"addShadow").Size(actionButtonW, actionButtonH).OnClick(func() {
				p.addShadow()
			}),
		)
	}

	if len(t.Substitutions) > 0 {
		tabs = append(tabs, giu.TabItem("Subs").Layout(p.makeTileSubsLayout(state, t.Substitutions)))
	}

	return giu.Layout{
		giu.TabBar("##TabBar_ds1_tiles" + p.id).Layout(tabs),
		giu.Custom(func() {
			if len(editionButtons) > 0 {
				giu.Layout{
					giu.Separator(),
					giu.Label("Edition tools:"),
					editionButtons,
				}.Build()
			}
		}),
	}
}

// nolint:dupl // yah, thats duplication of makeTileWallLayout but it isn't complete and can be changed
func (p *DS1Widget) makeTileFloorsLayout(state *DS1State, records []d2ds1.FloorShadowRecord) giu.Layout {
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
func (p *DS1Widget) makeTileWallsLayout(state *DS1State, records []d2ds1.WallRecord) giu.Layout {
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
func (p *DS1Widget) makeTileShadowsLayout(state *DS1State, records []d2ds1.FloorShadowRecord) giu.Layout {
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
func (p *DS1Widget) makeTileSubsLayout(state *DS1State, records []d2ds1.SubstitutionRecord) giu.Layout {
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

func (p *DS1Widget) makeSubstitutionsLayout(state *DS1State) giu.Layout {
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
				state.mode = ds1EditorModeViewer
			}),
			giu.Button("Cancel##"+p.id+"addFileCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = ds1EditorModeViewer
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

				state.mode = ds1EditorModeViewer
			}),
			giu.Button("Cancel##"+p.id+"AddObjectCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = ds1EditorModeViewer
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
				p.addPath()
			}),
			giu.Button("Cancel##"+p.id+"AddPathCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = ds1EditorModeViewer
			}),
		),
	}
}

// output in argument, because we're using this method in a two cases:
// first for adding floor an shadow
// second to in p.makeAddWallLayout, so we should specify, wher does
// we want to save results (in DS1State.addFloorShadowState,
// or in DS1State.addWallState)
func (p *DS1Widget) makeAddFloorShadowLayout(output *ds1AddFloorShadowState) giu.Layout {
	state := p.getState()

	trueFalseList := []string{"false", "true"}

	return giu.Layout{
		giu.Line(
			giu.Label("Prop 1: "),
			giu.InputInt("##"+p.id+"addFloorShadowProp1", &output.prop1).Size(inputIntW).OnChange(func() {
				if output.prop1 > maxByteSize {
					output.prop1 = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Sequence: "),
			giu.InputInt("##"+p.id+"addFloorShadowSequence", &output.sequence).Size(inputIntW).OnChange(func() {
				if output.sequence > maxByteSize {
					output.sequence = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Unknown 1: "),
			giu.InputInt("##"+p.id+"addFloorShadowUnknown1", &output.unknown1).Size(inputIntW).OnChange(func() {
				if output.unknown1 > maxByteSize {
					output.unknown1 = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Style: "),
			giu.InputInt("##"+p.id+"addFloorShadowStyle", &output.style).Size(inputIntW).OnChange(func() {
				if output.style > maxByteSize {
					output.style = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Unknown 2: "),
			giu.InputInt("##"+p.id+"addFloorShadowUnknown2", &output.unknown2).Size(inputIntW).OnChange(func() {
				if output.unknown2 > maxByteSize {
					output.unknown2 = maxByteSize
				}
			}),
		),
		giu.Line(
			giu.Label("Hidden: "),
			giu.Combo(
				"##"+p.id+"addFloorShadowHidden",
				trueFalseList[output.hidden],
				trueFalseList, &output.hidden,
			).Size(trueFalseListW),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##"+p.id+"AddFloorShadowSave").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				output.cb()
				state.mode = ds1EditorModeViewer
			}),
			giu.Button("Cancel##"+p.id+"AddFloorShadowCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = ds1EditorModeViewer
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
		p.makeAddFloorShadowLayout(&state.addWallState.ds1AddFloorShadowState),
	}
}

func (p *DS1Widget) deleteFile(idx int) {
	newFiles := make([]string, 0)

	for n, file := range p.ds1.Files {
		if n != idx {
			newFiles = append(newFiles, file)
		}
	}

	p.ds1.Files = newFiles
}

func (p *DS1Widget) addPath() {
	state := p.getState()

	newPath := d2path.Path{
		// nolint:gomnd // npc actions starts from 1
		Action: int(state.addPathState.pathAction) + 1,
		Position: d2vector.NewPosition(
			float64(state.addPathState.pathX),
			float64(state.addPathState.pathY),
		),
	}

	p.ds1.Objects[state.object].Paths = append(p.ds1.Objects[state.object].Paths, newPath)

	state.mode = ds1EditorModeViewer
}

func (p *DS1Widget) deletePath(idx int) {
	state := p.getState()

	newPaths := make([]d2path.Path, 0)

	for n, path := range p.ds1.Objects[state.object].Paths {
		if n != idx {
			newPaths = append(newPaths, path)
		}
	}

	p.ds1.Objects[state.object].Paths = newPaths
}
