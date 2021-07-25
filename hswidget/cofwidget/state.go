package cofwidget

import (
	"fmt"

	"github.com/AllenDang/giu"

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
	Mode mode
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

// viewerState represents cof viewer's state
type viewerState struct {
	LayerIndex     int32
	DirectionIndex int32
	FrameIndex     int32
	layer          *d2cof.CofLayer
	confirmDialog  *hswidget.PopUpConfirmDialog
}

// Dispose clears viewer's layers
func (s *viewerState) Dispose() {
	s.layer = nil
}

type newLayerFields struct {
	LayerType   int32
	Shadow      byte
	Selectable  bool
	Transparent bool
	DrawEffect  int32
	WeaponClass int32
}

// Dispose disposes editor's state
func (s *newLayerFields) Dispose() {
	s.LayerType = 0
	s.DrawEffect = 0
	s.WeaponClass = 0
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
			state.viewerState.layer = &p.cof.CofLayers[state.viewerState.LayerIndex]
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
		Mode: modeViewer,
		viewerState: &viewerState{
			confirmDialog: &hswidget.PopUpConfirmDialog{},
		},
		newLayerFields: &newLayerFields{
			Selectable: true,
			DrawEffect: int32(d2enum.DrawEffectNone),
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
