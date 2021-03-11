package palettemapwidget

import (
	"fmt"

	"github.com/ianling/giu"
)

type widgetState struct {
	selection int32
	slider1   int32
	slider2   int32
	textures  map[string]*giu.Texture
}

// Dispose cleans viewer's state
func (p *widgetState) Dispose() {
	p.textures = make(map[string]*giu.Texture)
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("PaletteMapWidget_%s", p.id)
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
		textures: make(map[string]*giu.Texture),
	}

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
