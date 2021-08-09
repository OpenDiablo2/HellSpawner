package animationwidget

import (
	"log"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

type dccWidgetState struct {
	controls *controlStructure
	widgetState
	isForward bool // Determines a direction of animation
}

func (s *dccWidgetState) getDirection() int32 {
	return s.controls.direction
}

func (s *dccWidgetState) Dispose() {
	s.widgetDispose()
}

func (s *dccWidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	s.widgetState.encode(sw)

	sw.PushInt32(s.controls.direction)
	sw.PushInt32(s.controls.frame)
	sw.PushInt32(s.controls.scale)

	return sw.GetBytes()
}

func (s *dccWidgetState) Decode(data []byte) {
	var err error

	sr := d2datautils.CreateStreamReader(data)

	s.widgetState.decode(sr)

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

	// update ticker
	s.ticker.Reset(time.Second * time.Duration(s.tickTime) / miliseconds)
}
