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
	mode          fontTableWidgetMode
	editRuneState editRuneState
	addItemState  addItemState
}

// Dispose cleans state
func (s *widgetState) Dispose() {
	s.editRuneState.Dispose()
	s.addItemState.Dispose()
}

type editRuneState struct {
	editedRune int32
	runeBefore rune
}

// Dispose disposes a rune state
func (e *editRuneState) Dispose() {
	e.editedRune = rune(0)
	e.runeBefore = rune(0)
}

type addItemState struct {
	newRune,
	width,
	height int32
}

func (s *addItemState) Dispose() {
	s.newRune = rune(0)
	s.height = 0
	s.width = 0
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
