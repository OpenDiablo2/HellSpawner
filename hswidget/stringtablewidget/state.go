package stringtablewidget

import (
	"fmt"
	"sort"

	"github.com/AllenDang/giu"
)

type widgetMode int32

const (
	widgetModeViewer widgetMode = iota
	widgetModeAddEdit
)

type widgetState struct {
	Mode    widgetMode
	keys    []string
	NumOnly bool
	addEditState
	Search string
}

func (ws *widgetState) Dispose() {
	ws.Mode = widgetModeViewer
	ws.keys = make([]string, 0)
	ws.addEditState.Dispose()
	ws.Search = ""
}

type addEditState struct {
	Key   string
	Value string
	// NoName is true, when we're viewing only no-named indexes
	NoName bool

	// if we used edit button by table entry,
	// we can't edit key value in edit layout
	Editable bool
}

func (aes *addEditState) Dispose() {
	aes.Key = ""
	aes.Value = ""
	aes.NoName = false
	aes.Editable = false
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
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

	p.setState(state)

	p.reloadMapValues()
}

func (p *widget) reloadMapValues() {
	state := p.getState()

	keys := make([]string, len(p.dict))

	n := 0

	for key := range p.dict {
		keys[n] = key
		n++
	}

	sort.Strings(keys)

	state.keys = keys
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
