package fonttablewidget

import (
	"fmt"

	"github.com/ianling/giu"
)

type fontTableWidgetMode int

const (
	// viewer state
	fontTableWidgetViewer fontTableWidgetMode = iota
	// edit rune
	fontTableWidgetEditRune
	// opens a dialog and adds a new item on first free index
	fontTableWidgetAddItem
)

type widgetState struct {
	mode     fontTableWidgetMode
	editRune editRune
}

// Dispose cleans state
func (s *widgetState) Dispose() {
	s.editRune.Dispose()
}

type editRune struct {
	editedRune int32
	startRune  rune
}

// Dispose disposes a rune state
func (e *editRune) Dispose() {
	e.editedRune = rune(0)
	e.startRune = rune(0)
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("COFWidget_%s", p.id)
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

func (p *widget) initState() {
	p.setState(&widgetState{
		mode: fontTableWidgetViewer,
	})
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
