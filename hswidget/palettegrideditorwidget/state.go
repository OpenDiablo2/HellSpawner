package palettegrideditorwidget

import (
	"fmt"
	"image/color"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

type widgetMode int32

const (
	widgetModeGrid widgetMode = iota
	widgetModeEdit
)

// PaletteGridState represents palette grid's state
type widgetState struct {
	mode widgetMode
	editEntryState
}

// Dispose cleans palette grids state
func (ws *widgetState) Dispose() {
	ws.mode = widgetModeGrid
}

func (ws *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(int32(ws.mode))
	sw.PushInt32(int32(ws.idx))

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

	idx, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	ws.idx = int(idx)

	l, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s := make([]rune, int(l))

	for i := 0; i < int(l); i++ {
		r, err := sr.ReadByte()
		if err != nil {
			log.Print(err)

			return
		}

		s[i] = rune(r)
	}
}

type editEntryState struct {
	idx  int
	rgba color.RGBA // nolint:structcheck // bug in golangci-lint
}

func (ees *editEntryState) Dispose() {
	ees.idx = 0
}

func (p *PaletteGridEditorWidget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
}

func (p *PaletteGridEditorWidget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.setState(&widgetState{})
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *PaletteGridEditorWidget) initState() {
	state := &widgetState{
		mode: widgetModeGrid,
	}

	p.setState(state)
}

func (p *PaletteGridEditorWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
