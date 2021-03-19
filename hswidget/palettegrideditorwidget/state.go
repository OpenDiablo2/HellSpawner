package palettegrideditorwidget

import (
	"fmt"

	"github.com/ianling/giu"
)

type widgetMode int

const (
	widgetModeGrid widgetMode = iota
	widgetModeEdit
)

// PaletteGridState represents palette grid's state
type widgetState struct {
	mode widgetMode
	editEntryState
}

// Dispose cleans palette grids state
func (ws *widgetState) Dispose() {
	ws.mode = widgetModeGrid
}

type editEntryState struct {
	idx     int
	r, g, b uint8
	hex     string // nolint:structcheck // linter's bug
}

func (ees *editEntryState) Dispose() {
	ees.idx = 0
	ees.r = 0
	ees.g = 0
	ees.b = 0
}

func (p *PaletteGridEditorWidget) getStateID() string {
	return fmt.Sprintf("PaletteGridEditorWidget_%s", p.id)
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
		mode: widgetModeGrid,
	}

	p.setState(state)
}

func (p *PaletteGridEditorWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
