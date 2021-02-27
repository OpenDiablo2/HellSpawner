package cofwidget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

type mode int

const (
	modeViewer mode = iota
	modeAddLayer
	modeConfirm
)

type widgetState struct {
	*viewerState
	*newLayerFields
	mode
}

// Dispose clear widget's state
func (s *widgetState) Dispose() {
	s.viewerState.Dispose()
	s.newLayerFields.Dispose()
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
