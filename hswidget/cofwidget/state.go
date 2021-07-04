package cofwidget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

type mode int32

const (
	modeViewer mode = iota
	modeAddLayer
	modeConfirm
)

type widgetState struct {
	*viewerState
	*newLayerFields
	mode
	textures
}

type textures struct {
	up    *giu.Texture
	down  *giu.Texture
	left  *giu.Texture
	right *giu.Texture
}

// Dispose clear widget's state
func (s *widgetState) Dispose() {
	s.viewerState.Dispose()
	s.newLayerFields.Dispose()
}

// Encode encodes state into byte slice to save it
func (s *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	if s.mode == modeConfirm {
		s.mode = modeViewer
	}

	sw.PushInt32(int32(s.mode))
	sw.PushInt32(s.layerIndex)
	sw.PushInt32(s.frameIndex)
	sw.PushInt32(s.directionIndex)
	sw.PushInt32(s.layerType)
	sw.PushBytes(s.shadow)

	if s.selectable {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	if s.transparent {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	sw.PushInt32(s.drawEffect)
	sw.PushInt32(s.weaponClass)

	return sw.GetBytes()
}

// Decode decodes byt slice into state
func (s *widgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	m, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.mode = mode(m)

	s.layerIndex, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.frameIndex, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.directionIndex, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.layerType, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.shadow, err = sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	selectable, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s.selectable = selectable == 1

	transparent, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s.transparent = transparent == 1

	s.drawEffect, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.weaponClass, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}
}

// viewerState represents cof viewer's state
type viewerState struct {
	layerIndex     int32
	directionIndex int32
	frameIndex     int32
	layer          *d2cof.CofLayer
	confirmDialog  *hswidget.PopUpConfirmDialog
}

// Dispose clears viewer's layers
func (s *viewerState) Dispose() {
	s.layer = nil
}

type newLayerFields struct {
	layerType   int32
	shadow      byte
	selectable  bool
	transparent bool
	drawEffect  int32
	weaponClass int32
}

// Dispose disposes editor's state
func (s *newLayerFields) Dispose() {
	s.layerType = 0
	s.drawEffect = 0
	s.weaponClass = 0
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
}

func (p *widget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
		if len(p.cof.CofLayers) > 0 {
			state.viewerState.layer = &p.cof.CofLayers[state.viewerState.layerIndex]
		}
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) initState() {
	state := &widgetState{
		mode: modeViewer,
		viewerState: &viewerState{
			confirmDialog: &hswidget.PopUpConfirmDialog{},
		},
		newLayerFields: &newLayerFields{
			selectable: true,
			drawEffect: int32(d2enum.DrawEffectNone),
		},
	}

	if len(p.cof.CofLayers) > 0 {
		state.viewerState.layer = &p.cof.CofLayers[0]
	}

	p.textureLoader.CreateTextureFromFile(hsassets.UpArrowIcon, func(texture *giu.Texture) {
		state.textures.up = texture
	})

	p.textureLoader.CreateTextureFromFile(hsassets.DownArrowIcon, func(texture *giu.Texture) {
		state.textures.down = texture
	})

	p.textureLoader.CreateTextureFromFile(hsassets.LeftArrowIcon, func(texture *giu.Texture) {
		state.textures.left = texture
	})

	p.textureLoader.CreateTextureFromFile(hsassets.RightArrowIcon, func(texture *giu.Texture) {
		state.textures.right = texture
	})

	p.setState(state)
}
