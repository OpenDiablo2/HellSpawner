package ds1widget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math/d2vector"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2path"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
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
	textureLoader       hscommon.TextureLoader
}

// Create creates a new ds1 viewer
func Create(textureLoader hscommon.TextureLoader, id string, ds1 *d2ds1.DS1, dbt *giu.Texture) giu.Widget {
	result := &widget{
		id:                  id,
		ds1:                 ds1,
		deleteButtonTexture: dbt,
		textureLoader:       textureLoader,
	}

	return result
}

func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case widgetModeViewer:
		p.makeViewerLayout().Build()
	case widgetModeAddFile:
		p.makeAddFileLayout().Build()
	case widgetModeAddObject:
		p.makeAddObjectLayout().Build()
	case widgetModeAddPath:
		p.makeAddPathLayout().Build()
	case widgetModeConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

// creates standard viewer/editor layout
func (p *widget) makeViewerLayout() giu.Layout {
	state := p.getState()

	tabs := giu.Layout{
		giu.TabItem("Files").Layout(p.makeFilesLayout()),
		giu.TabItem("Objects").Layout(p.makeObjectsLayout(state)),
		giu.TabItem("Tiles").Layout(p.makeTilesTabLayout(state)),
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

// makeDataLayout creates basic data layout
// used in p.makeViewerLayout
func (p *widget) makeDataLayout() giu.Layout {
	version := int32(p.ds1.Version())

	state := p.getState()

	w, h := int32(p.ds1.Width()), int32(p.ds1.Height())
	l := giu.Layout{
		giu.Line(
			giu.Label("Version: "),
			giu.InputInt("##"+p.id+"version", &version).Size(inputIntW).OnChange(func() {
				state.confirmDialog = hswidget.NewPopUpConfirmDialog(
					"##"+p.id+"confirmVersionChange",
					"Are you sure, you want to change DS1 Version?",
					"This value is used while decoding and encoding ds1 file\n"+
						"Please check github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1/ds1_version.go\n"+
						"to get more informations what does version determinates.\n\n"+
						"Continue?",
					func() {
						p.ds1.SetVersion(int(version))
						state.mode = widgetModeViewer
					},
					func() {
						state.mode = widgetModeViewer
					},
				)
				state.mode = widgetModeConfirm
			}),
		),
		// giu.Label(fmt.Sprintf("Size: %d x %d tiles", p.ds1.Width, p.ds1.Height)),
		giu.Label("Size:"),
		giu.Line(
			giu.Label("\tWidth: "),
			giu.InputInt("##"+p.id+"width", &w).Size(inputIntW).OnChange(func() {
				state.confirmDialog = hswidget.NewPopUpConfirmDialog(
					"##"+p.id+"confirmWidthChange",
					"Are you really sure, you want to change size of DS1 tiles?",
					"This will affect all your tiles in Tile tab.\n"+
						"Continue?",
					func() {
						p.ds1.SetWidth(int(w))
						state.mode = widgetModeViewer
					},
					func() {
						state.mode = widgetModeViewer
					},
				)
				state.mode = widgetModeConfirm
			}),
		),
		giu.Line(
			giu.Label("\tHeight: "),
			giu.InputInt("##"+p.id+"height", &h).Size(inputIntW).OnChange(func() {
				state.confirmDialog = hswidget.NewPopUpConfirmDialog(
					"##"+p.id+"confirmWidthChange",
					"Are you really sure, you want to change size of DS1 tiles?",
					"This will affect all your tiles in Tile tab.\n"+
						"Continue?",
					func() {
						p.ds1.SetHeight(int(h))
						state.mode = widgetModeViewer
					},
					func() {
						state.mode = widgetModeViewer
					},
				)
				state.mode = widgetModeConfirm
			}),
		),
		giu.Label(fmt.Sprintf("Substitution Type: %d", p.ds1.SubstitutionType)),
		giu.Separator(),
		giu.Label("Number of"),
		giu.Label(fmt.Sprintf("\tWall Layers: %d", len(p.ds1.Walls))),
		giu.Label(fmt.Sprintf("\tFloor Layers: %d", len(p.ds1.Floors))),
		giu.Label(fmt.Sprintf("\tShadow Layers: %d", len(p.ds1.Shadows))),
		giu.Label(fmt.Sprintf("\tSubstitution Layers: %d", len(p.ds1.Substitutions))),
	}

	return l
}

// makeFilesLayout creates files list
// used in p.makeViewerLayout (files tab)
func (p *widget) makeFilesLayout() giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	// iterating using the value should not be a big deal as
	// we only expect a handful of strings in this slice.
	for n, str := range p.ds1.Files {
		currentIdx := n

		l = append(l, giu.Layout{
			giu.Line(
				hswidget.MakeImageButton(
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
			state.mode = widgetModeAddFile
		}),
	}
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
				state.mode = widgetModeViewer
			}),
			giu.Button("Cancel##"+p.id+"addFileCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = widgetModeViewer
			}),
		),
	}
}

// makeObjectsLayout creates objects info tab
// used in p.makeViewerLayout (in objects tab)
func (p *widget) makeObjectsLayout(state *widgetState) giu.Layout {
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
				state.mode = widgetModeAddObject
			}),
			giu.Button("Add path to this object...##"+p.id+"AddPath").Size(actionButtonW, actionButtonH).OnClick(func() {
				state.mode = widgetModeAddPath
			}),
			hswidget.MakeImageButton(
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

// makeObjectLayout creates informations about single object
// used in p.makeObjectsLayout
func (p *widget) makeObjectLayout(state *widgetState) giu.Layout {
	if objIdx := int(state.object); objIdx >= len(p.ds1.Objects) {
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
			hswidget.MakeInputInt(
				"##"+p.id+"objType",
				inputIntW,
				&obj.Type,
				nil,
			),
		),
		giu.Line(
			giu.Label("ID: "),
			hswidget.MakeInputInt(
				"##"+p.id+"objID",
				inputIntW,
				&obj.ID,
				nil,
			),
		),
		giu.Label("Position (tiles): "),
		giu.Line(
			giu.Label("\tX: "),
			hswidget.MakeInputInt(
				"##"+p.id+"objX",
				inputIntW,
				&obj.X,
				nil,
			),
		),
		giu.Line(
			giu.Label("\tY: "),
			hswidget.MakeInputInt(
				"##"+p.id+"objY",
				inputIntW,
				&obj.Y,
				nil,
			),
		),
		giu.Line(
			giu.Label("Flags: 0x"),
			hswidget.MakeInputInt(
				"##"+p.id+"objFlags",
				inputIntW,
				&obj.Flags,
				nil,
			),
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

// makePathLayout creates paths table
// used in p.makeObjectLayout
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
			hswidget.MakeImageButton(
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

// makeTilesTabLayout creates tiles layout (tile x, y)
func (p *widget) makeTilesTabLayout(state *widgetState) giu.Layout {
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

	numRows := p.ds1.Height()
	if numRows == 0 {
		return l
	}

	if ty >= numRows {
		state.ds1Controls.tileY = int32(numRows - 1)
		p.setState(state)
	}

	if numCols := p.ds1.Width(); tx >= numCols {
		state.ds1Controls.tileX = int32(numCols - 1)
		p.setState(state)
	}

	tx, ty = int(state.tileX), int(state.tileY)

	l = append(
		l,
		giu.SliderInt("Tile X", &state.ds1Controls.tileX, 0, int32(p.ds1.Width()-1)),
		giu.SliderInt("Tile Y", &state.ds1Controls.tileY, 0, int32(p.ds1.Height()-1)),
		p.makeTileTabLayout(state, tx, ty),
	)

	return l
}

// makeTileTabLayout creates tabs for tile types
// used in p.makeTilesLayout
func (p *widget) makeTileTabLayout(state *widgetState, x, y int) giu.Layout {
	tabs := giu.Layout{}
	editionButtons := giu.Layout{}

	tabs = append(
		tabs,
		p.makeTileLayout(state, x, y, d2ds1.FloorLayerGroup),
		p.makeTileLayout(state, x, y, d2ds1.WallLayerGroup),
		p.makeTileLayout(state, x, y, d2ds1.ShadowLayerGroup),
		p.makeTileLayout(state, x, y, d2ds1.SubstitutionLayerGroup),
	)

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

func (p *widget) makeTileLayout(state *widgetState, x, y int, t d2ds1.LayerGroupType) giu.Layout {
	l := giu.Layout{}
	group := p.ds1.GetLayersGroup(t)
	numRecords := len(*group)

	// this is a pointer to appropriate record index
	var recordIdx *int32
	// addCb is a callback for layer-add button
	var addCb func(int32)
	// delCb is a callback for layer-delete button
	var deleteCb func(int32)

	// sets "everything" ;-)
	switch t {
	case d2ds1.FloorLayerGroup:
		recordIdx = &state.tile.floor
		addCb = p.addFloor
		deleteCb = p.deleteFloor
	case d2ds1.WallLayerGroup:
		recordIdx = &state.tile.wall
		addCb = p.addWall
		deleteCb = p.deleteWall
	case d2ds1.ShadowLayerGroup:
		recordIdx = &state.tile.shadow
	case d2ds1.SubstitutionLayerGroup:
		recordIdx = &state.tile.sub
	}

	var addBtn *giu.ButtonWidget
	if addCb != nil {
		addBtn = giu.Button("Add "+t.String()+" ##"+p.id+"addButton").
			Size(actionButtonW, actionButtonH).
			OnClick(func() { addCb(*recordIdx) })
	}

	var deleteBtn giu.Layout
	if deleteCb != nil {
		deleteBtn = hswidget.MakeImageButton(
			"##"+p.id+"delete"+t.String(),
			layerDeleteButtonSize, layerDeleteButtonSize,
			p.deleteButtonTexture,
			func() {
				deleteCb(*recordIdx)
			},
		)
	}

	if numRecords > 0 {
		// checks, if record index is correct
		if int(*recordIdx) >= numRecords {
			*recordIdx = int32(numRecords - 1)

			p.setState(state)
		} else if *recordIdx < 0 {
			*recordIdx = 0

			p.setState(state)
		}

		if numRecords > 1 {
			l = append(l, giu.SliderInt(t.String(), recordIdx, 0, int32(numRecords-1)))
		}

		l = append(l, p.makeTabTileLayout((*group)[*recordIdx].Tile(x, y), t))
	}

	return giu.Layout{giu.TabItem(t.String()).Layout(giu.Layout{
		l,
		giu.Separator(),
		giu.Custom(func() {
			var l giu.Layout
			if btn := addBtn; btn != nil {
				l = append(l, btn)
			}
			if btn := deleteBtn; btn != nil && numRecords > 0 {
				l = append(l, btn)
			}
			giu.Line(l...).Build()
		}),
	})}
}

func (p *widget) makeTabTileLayout(record *d2ds1.Tile, t d2ds1.LayerGroupType) giu.Layout {
	// for substitutions, only unknown bytes should be displayed
	if t == d2ds1.SubstitutionLayerGroup {
		unknown32 := int32(record.Substitution)

		return giu.Layout{
			giu.Line(
				giu.Label("Unknown: "),
				giu.InputInt("##"+p.id+"subUnknown", &unknown32).Size(inputIntW).OnChange(func() {
					record.Substitution = uint32(unknown32)
				}),
			),
		}
	}

	// common for shadows/walls/floors (like d2ds1.tileCommonFields)
	l := giu.Layout{
		giu.Line(
			giu.Label("Prop1: "),
			hswidget.MakeInputInt(
				"##"+p.id+"floorProp1",
				inputIntW,
				&record.Prop1,
				nil,
			),
		),
		giu.Line(
			giu.Label("Sequence: "),
			hswidget.MakeInputInt(
				"##"+p.id+"floorSequence",
				inputIntW,
				&record.Sequence,
				nil,
			),
		),
		giu.Line(
			giu.Label("Unknown1: "),
			hswidget.MakeInputInt(
				"##"+p.id+"floorUnknown1",
				inputIntW,
				&record.Unknown1,
				nil,
			),
		),
		giu.Line(
			giu.Label("Style: "),
			hswidget.MakeInputInt(
				"##"+p.id+"floorStyle",
				inputIntW,
				&record.Style,
				nil,
			),
		),
		giu.Line(
			giu.Label("Unknown2: "),
			hswidget.MakeInputInt(
				"##"+p.id+"floorUnknown2",
				inputIntW,
				&record.Unknown2,
				nil,
			),
		),
		giu.Line(
			giu.Label("Hidden: "),
			hswidget.MakeCheckboxFromByte(
				"##"+p.id+"floorHidden",
				&record.HiddenBytes,
			),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("RandomIndex: %v", record.RandomIndex)),
		),
		giu.Line(
			giu.Label(fmt.Sprintf("YAdjust: %v", record.YAdjust)),
		),
	}

	if t == d2ds1.WallLayerGroup {
		l = append(l,
			giu.Line(
				giu.Label("Zero: "),
				hswidget.MakeInputInt(
					"##"+p.id+"wallZero",
					inputIntW,
					&record.Zero,
					nil,
				),
			),
		)
	} else if t == d2ds1.FloorLayerGroup || t == d2ds1.ShadowLayerGroup {
		l = append(l,
			giu.Line(
				giu.Label(fmt.Sprintf("Animated: %v", record.Animated)),
			),
		)
	}

	return l
}

func (p *widget) makeSubstitutionsLayout(state *widgetState) giu.Layout {
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

				state.mode = widgetModeViewer
			}),
			giu.Button("Cancel##"+p.id+"AddObjectCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = widgetModeViewer
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
				state.mode = widgetModeViewer
			}),
			giu.Button("Cancel##"+p.id+"AddPathCancel").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.mode = widgetModeViewer
			}),
		),
	}
}

func (p *widget) deleteFile(idx int) {
	p.ds1.Files = append(p.ds1.Files[:idx], p.ds1.Files[idx+1:]...)
}

func (p *widget) addPath() {
	state := p.getState()

	newPath := d2path.Path{
		// npc actions starts from 1
		Action: int(state.addPathState.pathAction) + 1,
		Position: d2vector.NewPosition(
			float64(state.addPathState.pathX),
			float64(state.addPathState.pathY),
		),
	}

	p.ds1.Objects[state.object].Paths = append(p.ds1.Objects[state.object].Paths, newPath)
}

func (p *widget) deletePath(idx int) {
	state := p.getState()
	p.ds1.Objects[state.object].Paths = append(p.ds1.Objects[state.object].Paths[:idx], p.ds1.Objects[state.object].Paths[idx+1:]...)
}

func (p *widget) deleteObject(idx int32) {
	p.ds1.Objects = append(p.ds1.Objects[:idx], p.ds1.Objects[idx+1:]...)
}
