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
	mode      widgetMode
	mapKeys   []string
	mapIndex  int32
	recordIdx int32
}

// Dispose clears widget's state
func (ws *widgetState) Dispose() {
	ws.mode = widgetModeList
	ws.mapKeys = make([]string, 0)
	ws.mapIndex = 0
	ws.recordIdx = 0
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
