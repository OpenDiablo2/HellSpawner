package palettegridwidget

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
	// nolint:unused,structcheck // will be used
	loading bool
	// nolint:unused,structcheck // will be used
	failure bool
	texture [256]*giu.Texture
	editEntryState
}

// Dispose cleans palette grids state
func (ws *widgetState) Dispose() {
	ws.mode = widgetModeGrid
}

type editEntryState struct {
	idx     int
	r, g, b uint8
}

func (ees *editEntryState) Dispose() {
	ees.idx = 0
	ees.r = 0
	ees.g = 0
	ees.b = 0
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("PaletteGridWidget_%s", p.id)
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
	state := &widgetState{
		mode: widgetModeGrid,
	}

	p.reloadTextures()

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) reloadTextures() {
	// Prevent multiple invocation to LoadImage.
	p.setState(&widgetState{})

	for x := 0; x < 256; x++ {
		p.loadTexture(x)
	}
}
