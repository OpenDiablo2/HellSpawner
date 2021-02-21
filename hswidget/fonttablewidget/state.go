package fonttablewidget

import (
	"fmt"

	"github.com/ianling/giu"
)

type fontTableWidgetMode int

const (
	modeViewer fontTableWidgetMode = iota
	modeEditRune
	modeAddItem
)

type widgetState struct {
	mode     fontTableWidgetMode
	editRune editRune
	addItem  addItem
}

// Dispose cleans state
func (s *widgetState) Dispose() {
	s.editRune.Dispose()
	s.addItem.Dispose()
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

type addItem struct {
	newRune editRune
	width,
	height int32
}

func (s *addItem) Dispose() {
	s.newRune.Dispose()
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
		mode: modeViewer,
	})
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
