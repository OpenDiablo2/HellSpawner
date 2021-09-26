package animdatawidget

import (
	"fmt"
	"sort"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
)

type widgetMode int32

const (
	widgetModeList widgetMode = iota
	widgetModeViewRecord
)

type widgetState struct {
	Mode       widgetMode
	mapKeys    []string
	MapIndex   int32
	RecordIdx  int32
	deleteIcon *giu.Texture
	addEntryState
}

// Dispose clears widget's state
func (s *widgetState) Dispose() {
	s.Mode = widgetModeList
	s.mapKeys = make([]string, 0)
	s.MapIndex = 0
	s.RecordIdx = 0
	s.addEntryState.Dispose()
	s.deleteIcon = nil
}

type addEntryState struct {
	Name string
}

func (s *addEntryState) Dispose() {
	s.Name = ""
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

	p.textureLoader.CreateTextureFromFile(hsassets.DeleteIcon, func(texture *giu.Texture) {
		state.deleteIcon = texture
	})

	p.setState(state)

	p.reloadMapKeys()
}

func (p *widget) reloadMapKeys() {
	state := p.getState()
	state.mapKeys = p.d2.GetRecordNames()
	sort.Strings(state.mapKeys)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
