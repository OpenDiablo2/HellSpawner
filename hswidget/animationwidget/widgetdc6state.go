package animationwidget

import (
	"log"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

const maxAlpha = uint8(255)

const (
	dc6WidgetViewer widgetMode = iota
	dc6WidgetTiledView
)

type Dc6WidgetState struct {
	mode widgetMode
	viewerState
	tiledState
	WidgetState
}

func (s *Dc6WidgetState) getDirection() int32{
	return s.controls.direction
}

func (s *Dc6WidgetState) Dispose() {
	s.viewerState.Dispose()
	s.mode = dc6WidgetViewer
	s.widgetDispose()
}

func (s *Dc6WidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	s.WidgetState.Encode(sw)

	sw.PushInt32(int32(s.mode))

	sw.PushInt32(s.controls.direction)
	sw.PushInt32(s.controls.frame)
	sw.PushInt32(s.controls.scale)

	sw.PushInt32(s.width)
	sw.PushInt32(s.height)

	return sw.GetBytes()
}

func (s *Dc6WidgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	s.WidgetState.Decode(sr)

	mode, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.mode = widgetMode(mode)

	s.controls.direction, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.controls.frame, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.controls.scale, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.width, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.height, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	// update ticker
	s.ticker.Reset(time.Second * time.Duration(s.tickTime) / miliseconds)
}
