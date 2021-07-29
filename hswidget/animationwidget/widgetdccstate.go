package animationwidget

import (
	"log"
	"time"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

type DccWidgetState struct {
	controls *controlStructure
	WidgetState
}

func (s *DccWidgetState) getDirection() int32{
	return s.controls.direction
}

func (s *DccWidgetState) Dispose() {
	s.widgetDispose()
}

func (s *DccWidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	s.WidgetState.Encode(sw)

	sw.PushInt32(s.controls.direction)
	sw.PushInt32(s.controls.frame)
	sw.PushInt32(s.controls.scale)

	return sw.GetBytes()
}

func (s *DccWidgetState) Decode(data []byte) {
	var err error

	sr := d2datautils.CreateStreamReader(data)

	s.WidgetState.Decode(sr)

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
