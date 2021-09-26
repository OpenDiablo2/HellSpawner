package palettegrideditorwidget

import (
	"fmt"
	"image/color"

	"github.com/AllenDang/giu"
)

type widgetMode int32

const (
	widgetModeGrid widgetMode = iota
	widgetModeEdit
)

// PaletteGridState represents palette grid's state
type widgetState struct {
	Mode widgetMode `json:"mode"`
	editEntryState
}

// Dispose cleans palette grids state
func (ws *widgetState) Dispose() {
	ws.Mode = widgetModeGrid
}

type editEntryState struct {
	Idx  int
	RGBA color.RGBA
}

func (ees *editEntryState) Dispose() {
	ees.Idx = 0
}

func (p *PaletteGridEditorWidget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
}

func (p *PaletteGridEditorWidget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.setState(&widgetState{})
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *PaletteGridEditorWidget) initState() {
	state := &widgetState{
		Mode: widgetModeGrid,
	}

	p.setState(state)
}

func (p *PaletteGridEditorWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
