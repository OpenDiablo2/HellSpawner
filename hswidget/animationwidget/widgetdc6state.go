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

type dc6WidgetState struct {
	mode widgetMode
	viewerState
	tiledState
	widgetState
}

func (s *dc6WidgetState) getDirection() int32 {
	return s.direction
}

func (s *dc6WidgetState) Dispose() {
	s.viewerState.Dispose()
	s.mode = dc6WidgetViewer
	s.widgetDispose()
}

func (s *dc6WidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	s.widgetState.encode(sw)

	sw.PushInt32(int32(s.mode))

	sw.PushInt32(s.direction)
	sw.PushInt32(s.frame)
	sw.PushInt32(s.scale)

	sw.PushInt32(s.width)
	sw.PushInt32(s.height)

	return sw.GetBytes()
}

func (s *dc6WidgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	s.widgetState.decode(sr)

	mode, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.mode = widgetMode(mode)

	s.direction, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.frame, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.scale, err = sr.ReadInt32()
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
