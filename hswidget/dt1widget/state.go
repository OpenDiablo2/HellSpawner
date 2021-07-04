package dt1widget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
)

type controls struct {
	tileGroup   int32
	tileVariant int32
	showGrid    bool
	showFloor   bool
	showWall    bool
	subtileFlag int32
	scale       int32
}

// widgetState represents dt1 viewers state
type widgetState struct {
	*controls

	lastTileGroup int32

	tileGroups [][]*d2dt1.Tile
	textures   [][]map[string]*giu.Texture
}

// Dispose clears viewers state
func (s *widgetState) Dispose() {
	s.textures = nil
}

func (s *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(s.tileGroup)
	sw.PushInt32(s.tileVariant)

	if s.showGrid {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	if s.showFloor {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	if s.showWall {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	sw.PushInt32(s.subtileFlag)
	sw.PushInt32(s.scale)

	return sw.GetBytes()
}

func (s *widgetState) Decode(data []byte) {
	var err error

	sr := d2datautils.CreateStreamReader(data)

	s.tileGroup, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.tileVariant, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	showGrid, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s.showGrid = showGrid == 1

	showFloor, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s.showFloor = showFloor == 1

	showWall, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	s.showWall = showWall == 1

	s.subtileFlag, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.scale, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}
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
		p.makeTileTextures()
		state = p.getState()
	}

	return state
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) initState() {
	state := &widgetState{
		controls: &controls{
			showGrid:  true,
			showFloor: true,
			showWall:  true,
		},
		tileGroups: p.groupTilesByIdentity(),
	}

	p.setState(state)
}
