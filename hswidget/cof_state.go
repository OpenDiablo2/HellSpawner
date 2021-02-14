package hswidget

import "github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

type cofEditorMode int

const (
	cofEditorModeViewer cofEditorMode = iota
	cofEditorModeAddLayer
	cofEditorModeConfirm
)

// COFState represents cof editor's and viewer's state
type COFState struct {
	*viewerState
	*newLayerFields
	mode cofEditorMode
}

// Dispose clear widget's state
func (s *COFState) Dispose() {
	s.viewerState.Dispose()
	s.newLayerFields.Dispose()
}

// viewerState represents cof viewer's state
type viewerState struct {
	layerIndex     int32
	directionIndex int32
	frameIndex     int32
	layer          *d2cof.CofLayer
	confirmDialog  *PopUpConfirmDialog
}

// Dispose clears viewer's layers
func (s *viewerState) Dispose() {
	s.layer = nil
}

type newLayerFields struct {
	layerType   int32
	shadow      int32
	selectable  int32
	transparent int32
	drawEffect  int32
	weaponClass int32
}

// Dispose disposes editor's state
func (s *newLayerFields) Dispose() {
	s.layerType = 0
	s.drawEffect = 0
	s.weaponClass = 0
}
