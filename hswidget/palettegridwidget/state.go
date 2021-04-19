package palettegridwidget

import (
	"fmt"

	"github.com/AllenDang/giu"
)

// PaletteGridState represents palette grid's state
type widgetState struct {
	// nolint:unused,structcheck // will be used
	loading bool
	// nolint:unused,structcheck // will be used
	failure bool
	texture [256]*giu.Texture
}

// Dispose cleans palette grids state
func (ws *widgetState) Dispose() {
	// noop
}

func (p *PaletteGridWidget) getStateID() string {
	return fmt.Sprintf("PaletteGridWidget_%s", p.id)
}

func (p *PaletteGridWidget) getState() *widgetState {
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

func (p *PaletteGridWidget) initState() {
	state := &widgetState{}

	p.reloadTextures()

	p.setState(state)
}

func (p *PaletteGridWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
