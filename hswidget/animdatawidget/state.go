package animdatawidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"
)

type widgetMode int

const (
	widgetModeList widgetMode = iota
	widgetModeViewRecord
)

type widgetState struct {
	mode     widgetMode
	mapKeys  []string
	mapIndex int32
}

// Dispose clears widget's state
func (ws *widgetState) Dispose() {
	ws.mapKeys = make([]string, 0)
	ws.mapIndex = 0
	ws.mode = widgetModeList
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("AnimationDataWidget_%s", p.id)
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
	state := &widgetState{}

	state.mapKeys = p.d2.GetRecordNames()
	sort.Strings(state.mapKeys)

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
