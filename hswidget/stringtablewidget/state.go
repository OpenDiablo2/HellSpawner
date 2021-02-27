package stringtablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"
)

type widgetState struct {
	keys []string
}

func (ws *widgetState) Dispose() {
	ws.keys = make([]string, 0)
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("StringTableWidget_%s", p.id)
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

	keys := make([]string, len(p.dict.Entries))
	n := 0
	for key := range p.dict.Entries {
		keys[n] = key
		n++
	}

	sort.Strings(keys)

	state.keys = keys

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
