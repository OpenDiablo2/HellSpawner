package fonttablewidget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
)

type widgetMode int32

const (
	modeViewer widgetMode = iota
	modeEditRune
	modeAddItem
)

type widgetState struct {
	mode                widgetMode
	editRuneState       editRuneState
	addItemState        addItemState
	deleteButtonTexture *giu.Texture
}

// Dispose cleans state
func (s *widgetState) Dispose() {
	s.editRuneState.Dispose()
	s.addItemState.Dispose()
}

func (s *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(int32(s.mode))
	sw.PushInt32(s.editRuneState.editedRune)
	sw.PushInt16(int16(s.editRuneState.runeBefore))
	sw.PushInt32(s.addItemState.newRune)
	sw.PushInt32(s.addItemState.width)
	sw.PushInt32(s.addItemState.height)

	return sw.GetBytes()
}

func (s *widgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	mode, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.mode = widgetMode(mode)

	s.editRuneState.editedRune, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	runeBefore, err := sr.ReadInt16()
	if err != nil {
		log.Print(err)

		return
	}

	s.editRuneState.runeBefore = rune(runeBefore)

	s.addItemState.newRune, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.addItemState.width, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.addItemState.height, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}
}

type editRuneState struct {
	editedRune int32
	runeBefore rune
}

// Dispose disposes a rune state
func (e *editRuneState) Dispose() {
	e.editedRune = rune(0)
	e.runeBefore = rune(0)
}

type addItemState struct {
	newRune,
	width,
	height int32
}

func (s *addItemState) Dispose() {
	s.newRune = rune(0)
	s.height = 0
	s.width = 0
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
	state := &widgetState{
		mode: modeViewer,
	}

	p.textureLoader.CreateTextureFromFile(hsassets.DeleteIcon, func(texture *giu.Texture) {
		state.deleteButtonTexture = texture
	})

	p.setState(state)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
