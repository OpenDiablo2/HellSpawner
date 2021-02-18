package fonttablewidget

import (
	"fmt"

	"github.com/ianling/giu"
)

type fontTableWidgetMode int

const (
	fontTableWidgetViewer fontTableWidgetMode = iota
	fontTableWidgetEditRune
)

type widgetState struct {
	mode fontTableWidgetMode
}

// Dispose cleans state
func (s *widgetState) Dispose() {
	// noop
}

func (p *FontTableWidget) getStateID() string {
	return fmt.Sprintf("COFWidget_%s", p.id)
}

func (p *FontTableWidget) getState() *widgetState {
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

func (p *FontTableWidget) initState() {
	p.setState(&widgetState{
		mode: fontTableWidgetViewer,
	})
}

func (p *FontTableWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
