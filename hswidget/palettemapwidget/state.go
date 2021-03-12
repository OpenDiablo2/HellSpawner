package palettemapwidget

import (
	"fmt"

	"github.com/ianling/giu"
	//"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

type widgetMode int

const (
	widgetModeView widgetMode = iota
	widgetModeEditTransform
)

type widgetState struct {
	mode      widgetMode
	selection int32
	slider1   int32
	slider2   int32
	//textures  map[string]*palettegridwidget.PaletteGridWidget
	textures map[string]giu.Widget
	editTransformState
}

// Dispose cleans viewer's state
func (p *widgetState) Dispose() {
	//p.textures = make(map[string]*palettegridwidget.PaletteGridWidget)
	p.textures = make(map[string]giu.Widget)
	p.editTransformState.Dispose()
}

type editTransformState struct {
	id  string
	idx int
}

func (p *editTransformState) Dispose() {
	p.id = ""
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
		mode: widgetModeView,
		//textures: make(map[string]*palettegridwidget.PaletteGridWidget),
		textures: make(map[string]giu.Widget),
	}

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
