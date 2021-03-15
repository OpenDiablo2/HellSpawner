package cofwidget

import (
	"log"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

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

	s.selectable = (selectable == 1)

	transparent, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s.transparent = (transparent == 1)

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
