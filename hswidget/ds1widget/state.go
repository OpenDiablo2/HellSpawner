package ds1widget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

type widgetMode int32

const (
	widgetModeViewer widgetMode = iota
	widgetModeAddFile
	widgetModeAddObject
	widgetModeAddPath
	widgetModeConfirm
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

// widgetState represents ds1 viewers state
type widgetState struct {
	*ds1Controls
	mode           widgetMode
	confirmDialog  *hswidget.PopUpConfirmDialog
	newFilePath    string
	addObjectState ds1AddObjectState
	addPathState   ds1AddPathState
}

// Dispose clears viewers state
func (is *widgetState) Dispose() {
	is.addObjectState.Dispose()
	is.addPathState.Dispose()
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("DS1Widget_%s", p.id)
}

func (p *widget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) initState() {
	state := &widgetState{
		ds1Controls: &ds1Controls{},
	}

	p.setState(state)
}
