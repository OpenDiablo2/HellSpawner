package animdatawidget

import (
	"fmt"
	"log"
	"sort"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
)

type widgetMode int32

const (
	widgetModeList widgetMode = iota
	widgetModeViewRecord
)

type widgetState struct {
	mode       widgetMode
	mapKeys    []string
	mapIndex   int32
	recordIdx  int32
	deleteIcon *giu.Texture
	addEntryState
}

// Dispose clears widget's state
func (ws *widgetState) Dispose() {
	ws.mode = widgetModeList
	ws.mapKeys = make([]string, 0)
	ws.mapIndex = 0
	ws.recordIdx = 0
	ws.addEntryState.Dispose()
	ws.deleteIcon = nil
}

type addEntryState struct {
	name string
}

func (aes *addEntryState) Dispose() {
	aes.name = ""
}

// Encode encodes state into byte slice to save it
func (ws *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(int32(ws.mode))
	sw.PushInt32(ws.mapIndex)
	sw.PushInt32(ws.recordIdx)

	return sw.GetBytes()
}

// Decode decodes byte slice into widget state
func (ws *widgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	mode, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	mapIndex, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	recordIdx, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	ws.mode = widgetMode(mode)
	ws.mapIndex = mapIndex
	ws.recordIdx = recordIdx
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
