package stringtablewidget

import (
	"fmt"
	"sort"

	"github.com/ianling/giu"
)

type widgetMode int

const (
	widgetModeViewer widgetMode = iota
	widgetModeAddEdit
)

type widgetState struct {
	mode    widgetMode
	keys    []string
	numOnly bool
	addEditState
	search string
}

func (ws *widgetState) Dispose() {
	ws.mode = widgetModeViewer
	ws.keys = make([]string, 0)
	ws.addEditState.Dispose()
	ws.search = ""
}

type addEditState struct {
	key   string
	value string
	// noName is true, when we're viewing only no-named indexes
	noName bool

	// if we used edit button by table entry,
	// we can't edit key value in edit layout
	editable bool
}

func (aes *addEditState) Dispose() {
	aes.key = ""
	aes.value = ""
	aes.noName = false
	aes.editable = false
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
