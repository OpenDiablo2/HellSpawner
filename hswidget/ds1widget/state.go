package ds1widget

import (
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

type ds1EditorMode int

const (
	ds1EditorModeViewer ds1EditorMode = iota
	ds1EditorModeAddFile
	ds1EditorModeAddObject
	ds1EditorModeAddPath
	ds1EditorModeAddFloorShadow
	ds1EditorModeAddWall
	ds1EditorModeConfirm
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

// ds1AddObjectState represents state of new object
type ds1AddObjectState struct {
	objType  int32
	objID    int32
	objX     int32
	objY     int32
	objFlags int32
}

// Dispose clears state
func (t *ds1AddObjectState) Dispose() {
	// noop
}

// ds1AddPathState contains data about new path
type ds1AddPathState struct {
	pathAction int32
	pathX      int32
	pathY      int32
}

// Dispose clears state
func (t *ds1AddPathState) Dispose() {
	// noop
}

// ds1AddFloorShadowState contains data used in
// add floor-shadow record dialog
type ds1AddFloorShadowState struct {
	prop1    int32
	sequence int32
	unknown1 int32
	style    int32
	unknown2 int32
	hidden   int32
	cb       func()
}

// Dispose resets DS1AddFloorShadowState
func (t ds1AddFloorShadowState) Dispose() {
	t.prop1 = 0
	t.sequence = 0
	t.unknown1 = 0
	t.style = 0
	t.unknown2 = 0
	t.hidden = 0
}

// ds1AddWallState contains data used in add wall dialog
type ds1AddWallState struct {
	tileType int32
	zero     int32
	ds1AddFloorShadowState
}

// Dispose cleans DS1AddWallState
func (t *ds1AddWallState) Dispose() {
	t.ds1AddFloorShadowState.Dispose()
}

// DS1State represents ds1 viewers state
type DS1State struct {
	*ds1Controls
	mode                ds1EditorMode
	confirmDialog       *hswidget.PopUpConfirmDialog
	newFilePath         string
	addObjectState      ds1AddObjectState
	addPathState        ds1AddPathState
	addFloorShadowState ds1AddFloorShadowState
	addWallState        ds1AddWallState
}

// Dispose clears viewers state
func (is *DS1State) Dispose() {
	is.addObjectState.Dispose()
	is.addPathState.Dispose()
	is.addFloorShadowState.Dispose()
	is.addWallState.Dispose()
}
