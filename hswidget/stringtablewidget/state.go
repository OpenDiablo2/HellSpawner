package stringtablewidget

import (
	"fmt"
	"log"
	"sort"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

type widgetMode int32

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

func (ws *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(int32(ws.mode))

	if ws.numOnly {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	sw.PushInt32(int32(len(ws.key)))
	sw.PushBytes([]byte(ws.key)...)

	sw.PushInt32(int32(len(ws.value)))
	sw.PushBytes([]byte(ws.value)...)

	if ws.addEditState.noName {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	if ws.addEditState.editable {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	sw.PushInt32(int32(len(ws.search)))
	sw.PushBytes([]byte(ws.search)...)

	return sw.GetBytes()
}

func (ws *widgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	mode, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	ws.mode = widgetMode(mode)

	numOnly, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	ws.numOnly = numOnly == 1

	l, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	editKey, err := sr.ReadBytes(int(l))
	if err != nil {
		log.Print(err)

		return
	}

	ws.addEditState.key = string(editKey)

	l, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	editValue, err := sr.ReadBytes(int(l))
	if err != nil {
		log.Print(err)

		return
	}

	ws.addEditState.value = string(editValue)

	noName, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	ws.noName = noName == 1

	editable, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	ws.editable = editable == 1

	l, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	search, err := sr.ReadBytes(int(l))
	if err != nil {
		log.Print(err)

		return
	}

	ws.search = string(search)
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
