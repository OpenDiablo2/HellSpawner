package animdatawidget

import (
	"fmt"
	"log"
	"sort"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

type widgetMode int32

const (
	widgetModeList widgetMode = iota
	widgetModeViewRecord
)

type AnimationDataWidgetState struct {
	mode      widgetMode
	mapKeys   []string
	mapIndex  int32
	recordIdx int32
	addEntryState
}

// Dispose clears widget's state
func (ws *AnimationDataWidgetState) Dispose() {
	ws.mode = widgetModeList
	ws.mapKeys = make([]string, 0)
	ws.mapIndex = 0
	ws.recordIdx = 0
	ws.addEntryState.Dispose()
}

type addEntryState struct {
	name string
}

func (aes *addEntryState) Dispose() {
	aes.name = ""
}

func (ws *AnimationDataWidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(int32(ws.mode))
	sw.PushInt32(ws.mapIndex)
	sw.PushInt32(ws.recordIdx)

	return sw.GetBytes()
}

func (ws *AnimationDataWidgetState) Decode(data []byte) {
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

func (p *widget) getState() *AnimationDataWidgetState {
	var state *AnimationDataWidgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*AnimationDataWidgetState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *widget) initState() {
	state := &widgetState{}
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
