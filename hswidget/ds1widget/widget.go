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
	actionButtonW, actionButtonH         = 170, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
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

// widget represents ds1 viewers widget
type widget struct {
	id                  string
	ds1                 *d2ds1.DS1
	deleteButtonTexture *giu.Texture
}

// DS1Viewer creates a new ds1 viewer
func DS1Viewer(id string, ds1 *d2ds1.DS1, dbt *giu.Texture) *widget {
	result := &widget{
		id:                  id,
		ds1:                 ds1,
		deleteButtonTexture: dbt,
	}

	return result
}

func (p *widget) Build() {
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

func (p *widget) makeViewerLayout() giu.Layout {
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

func (p *widget) makeDataLayout() giu.Layout {
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

func (p *widget) makeFilesLayout() giu.Layout {
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

func (p *widget) makeObjectsLayout(state *DS1State) giu.Layout {
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

	l = append(
		l,
		giu.Separator(),
		giu.Line(
			giu.Button("Add new object...##"+p.id+"AddObject").Size(actionButtonW, actionButtonH).OnClick(func() {
				state.mode = ds1EditorModeAddObject
			}),
			giu.Button("Add new path...##"+p.id+"AddPath").Size(actionButtonW, actionButtonH).OnClick(func() {
				state.mode = ds1EditorModeAddPath
			}),
			hsutil.MakeImageButton(
				"##"+p.id+"deleteObject",
				layerDeleteButtonSize, layerDeleteButtonSize,
				p.deleteButtonTexture,
				func() {
					p.deleteObject(state.object)
				},
			),
		),
	)

	return l
}

func (p *widget) makeObjectLayout(state *DS1State) giu.Layout {
	objIdx := int(state.object)

	if objIdx >= len(p.ds1.Objects) {
		state.ds1Controls.object = int32(len(p.ds1.Objects) - 1)
		p.setState(state)
	} else if objIdx < 0 {
		state.ds1Controls.object = 0
		p.setState(state)
	}

	obj := &p.ds1.Objects[int(state.ds1Controls.object)]

	l := giu.Layout{
		giu.Line(
			giu.Label("Type: "),
			makeInputIntFromInt("##"+p.id+"objType", &obj.Type),
		),
		giu.Line(
			giu.Label("ID: "),
			makeInputIntFromInt("##"+p.id+"objID", &obj.ID),
		),
		giu.Label("Position (tiles): "),
		giu.Line(
			giu.Label("\tX: "),
			makeInputIntFromInt("##"+p.id+"objX", &obj.X),
		),
		giu.Line(
			giu.Label("\tY: "),
			makeInputIntFromInt("##"+p.id+"objY", &obj.Y),
		),
		giu.Line(
			giu.Label("Flags: 0x"),
			makeInputIntFromInt("##"+p.id+"objFlags", &obj.Flags),
		),
	}

	if len(obj.Paths) > 0 {
		l = append(
			l,
			giu.Dummy(1, 16),
			p.makePathLayout(obj),
		)
	}

	return l
}

func (p *widget) makePathLayout(obj *d2ds1.Object) giu.Layout {
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
			hsutil.MakeImageButton(
				"##"+p.id+"deletePath"+strconv.Itoa(currentIdx),
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

func (p *widget) makeTilesLayout(state *DS1State) giu.Layout {
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
		l,
		giu.Line(
			giu.SliderInt("Tile X", &state.ds1Controls.tileX, 0, p.ds1.Width-1),
			giu.Button("Add...##"+p.id+"addTileRow"),
			hsutil.MakeImageButton(
				"##"+p.id+"deleteTileRow",
				deleteButtonSize, deleteButtonSize,
				p.deleteButtonTexture,
				func() {},
			),
		),
		giu.Line(
			giu.SliderInt("Tile Y", &state.ds1Controls.tileY, 0, p.ds1.Height-1),
			giu.Button("Add...##"+p.id+"addTileCol"),
			hsutil.MakeImageButton(
				"##"+p.id+"deleteTileCol",
				deleteButtonSize, deleteButtonSize,
				p.deleteButtonTexture,
				func() {},
			),
		),
		p.makeTileLayout(state, &p.ds1.Tiles[ty][tx]),
	)

	return l
}

// nolint:funlen // cannot reduce
func (p *widget) makeTileLayout(state *DS1State, t *d2ds1.TileRecord) giu.Layout {
	tabs := giu.Layout{}
	editionButtons := giu.Layout{}

	if len(t.Floors) > 0 {
		tabs = append(
			tabs,
			giu.TabItem("Floors").Layout(giu.Layout{
				p.makeTileFloorsLayout(state, t.Floors),
				giu.Separator(),
				giu.Line(
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
func (p *widget) makeTileFloorsLayout(state *DS1State, records []d2ds1.FloorShadowRecord) giu.Layout {
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

func (p *widget) makeTileFloorLayout(record *d2ds1.FloorShadowRecord) giu.Layout {
	return giu.Layout{
		giu.Line(
			giu.Label("Prop1: "),
			makeInputIntFromByte("##"+p.id+"floorProp1", &record.Prop1),
		),
		giu.Line(
			giu.Label("Sequence: "),
			makeInputIntFromByte("##"+p.id+"floorSequence", &record.Sequence),
		),
		giu.Line(
			giu.Label("Unknown1: "),
			makeInputIntFromByte("##"+p.id+"floorUnknown1", &record.Unknown1),
		),
		giu.Line(
			giu.Label("Style: "),
			makeInputIntFromByte("##"+p.id+"floorStyle", &record.Style),
		),
		giu.Line(
			giu.Label("Unknown2: "),
			makeInputIntFromByte("##"+p.id+"floorUnknown2", &record.Unknown2),
		),
		giu.Line(
			giu.Label("Hidden: "),
			makeCheckboxFromByte("##"+p.id+"floorHidden", &record.HiddenBytes),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
		),
	}
}

// nolint:dupl // could be changed
func (p *widget) makeTileWallsLayout(state *DS1State, records []d2ds1.WallRecord) giu.Layout {
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

func (p *widget) makeTileWallLayout(record *d2ds1.WallRecord) giu.Layout {
	return giu.Layout{
		giu.Line(
			giu.Label("Prop1: "),
			makeInputIntFromByte("##"+p.id+"wallProp1", &record.Prop1),
		),
		giu.Line(
			giu.Label("Zero: "),
			makeInputIntFromByte("##"+p.id+"wallZero", &record.Zero),
		),
		giu.Line(
			giu.Label("Sequence: "),
			makeInputIntFromByte("##"+p.id+"wallSequence", &record.Sequence),
		),
		giu.Line(
			giu.Label("Unknown1: "),
			makeInputIntFromByte("##"+p.id+"wallUnknown1", &record.Unknown1),
		),
		giu.Line(
			giu.Label("Style: "),
			makeInputIntFromByte("##"+p.id+"wallStyle", &record.Style),
		),
		giu.Line(
			giu.Label("Unknown2: "),
			makeInputIntFromByte("##"+p.id+"wallUnknown2", &record.Unknown2),
		),
		giu.Line(
			giu.Label("Hidden: "),
			makeCheckboxFromByte("##"+p.id+"wallHidden", &record.HiddenBytes),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
		),
	}
}

// nolint:dupl // no need to change
func (p *widget) makeTileShadowsLayout(state *DS1State, records []d2ds1.FloorShadowRecord) giu.Layout {
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

func (p *widget) makeTileShadowLayout(record *d2ds1.FloorShadowRecord) giu.Layout {
	return giu.Layout{
		giu.Line(
			giu.Label("Prop1: "),
			makeInputIntFromByte("##"+p.id+"floorProp1", &record.Prop1),
		),
		giu.Line(
			giu.Label("Sequence: "),
			makeInputIntFromByte("##"+p.id+"floorSequence", &record.Sequence),
		),
		giu.Line(
			giu.Label("Unknown1: "),
			makeInputIntFromByte("##"+p.id+"floorUnknown1", &record.Unknown1),
		),
		giu.Line(
			giu.Label("Style: "),
			makeInputIntFromByte("##"+p.id+"floorStyle", &record.Style),
		),
		giu.Line(
			giu.Label("Unknown2: "),
			makeInputIntFromByte("##"+p.id+"floorUnknown2", &record.Unknown2),
		),
		giu.Line(
			giu.Label("Hidden: "),
			makeCheckboxFromByte("##"+p.id+"floorHidden", &record.HiddenBytes),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
		),
	}
}

// nolint:dupl // it is ok
func (p *widget) makeTileSubsLayout(state *DS1State, records []d2ds1.SubstitutionRecord) giu.Layout {
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

func (p *widget) makeTileSubLayout(record *d2ds1.SubstitutionRecord) giu.Layout {
	unknown32 := int32(record.Unknown)

	return giu.Layout{
		giu.Line(
			giu.Label("Unknown: "),
			giu.InputInt("##"+p.id+"subUnknown", &unknown32).Size(inputIntW).OnChange(func() {
				record.Unknown = uint32(unknown32)
			}),
		),
	}
}

func (p *widget) makeSubstitutionsLayout(state *DS1State) giu.Layout {
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

func (p *widget) makeSubstitutionLayout(group *d2ds1.SubstitutionGroup) giu.Layout {
	l := giu.Layout{
		giu.Label(fmt.Sprintf("TileX: %d", group.TileX)),
		giu.Label(fmt.Sprintf("TileY: %d", group.TileY)),
		giu.Label(fmt.Sprintf("WidthInTiles: %d", group.WidthInTiles)),
		giu.Label(fmt.Sprintf("HeightInTiles: %d", group.HeightInTiles)),
		giu.Label(fmt.Sprintf("Unknown: 0x%x", group.Unknown)),
	}

	return l
}

func (p *widget) makeAddFileLayout() giu.Layout {
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

func (p *widget) makeAddObjectLayout() giu.Layout {
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

func (p *widget) makeAddPathLayout() giu.Layout {
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
func (p *widget) makeAddFloorShadowLayout(output *ds1AddFloorShadowState) giu.Layout {
	state := p.getState()

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
			makeCheckboxFromByte("##"+p.id+"addFloorShadowHidden",
				&output.hidden,
			),
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

func (p *widget) makeAddWallLayout() giu.Layout {
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

func (p *widget) deleteFile(idx int) {
	newFiles := make([]string, 0)

	for n, file := range p.ds1.Files {
		if n != idx {
			newFiles = append(newFiles, file)
		}
	}

	p.ds1.Files = newFiles
}

func (p *widget) addPath() {
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

func (p *widget) deletePath(idx int) {
	state := p.getState()

	newPaths := make([]d2path.Path, 0)

	for n, path := range p.ds1.Objects[state.object].Paths {
		if n != idx {
			newPaths = append(newPaths, path)
		}
	}

	p.ds1.Objects[state.object].Paths = newPaths
}

func (p *widget) deleteObject(idx int32) {
	// first, wee check if index (idx) exist in NpcIndexes
	for n, i := range p.ds1.NpcIndexes {
		if i == int(idx) {
			p.ds1.NpcIndexes = append(p.ds1.NpcIndexes[:n], p.ds1.NpcIndexes[n+1:]...)

			// decrease all indexes in npc list
			for n, i := range p.ds1.NpcIndexes {
				if i > int(idx) {
					p.ds1.NpcIndexes[n]--
				}
			}

			break
		}
	}

	// delete object
	p.ds1.Objects = append(p.ds1.Objects[:idx], p.ds1.Objects[idx+1:]...)
}
