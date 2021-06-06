package palettegrideditorwidget

import (
	"fmt"
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
	sw.PushBytes(ws.r)
	sw.PushBytes(ws.g)
	sw.PushBytes(ws.b)
	sw.PushBytes(byte(len(ws.hex)))
	sw.PushBytes([]byte(ws.hex)...)

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

	ws.r, err = sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	ws.g, err = sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	ws.b, err = sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

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

	ws.hex = string(s)
}

type editEntryState struct {
	idx     int
	r, g, b uint8
	hex     string // nolint:structcheck // linter's bug
	texture *giu.Texture
}

func (ees *editEntryState) Dispose() {
	ees.idx = 0
	ees.r = 0
	ees.g = 0
	ees.b = 0
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

	p.updateEditedTexture()
}

func (p *PaletteGridEditorWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
