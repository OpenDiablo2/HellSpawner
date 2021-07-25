package dt1widget

import (
	"fmt"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
)

type controls struct {
	TileGroup   int32
	TileVariant int32
	ShowGrid    bool
	ShowFloor   bool
	ShowWall    bool
	SubtileFlag int32
	Scale       int32
}

// widgetState represents dt1 viewers state
type widgetState struct {
	*controls

	LastTileGroup int32

	tileGroups [][]*d2dt1.Tile
	textures   [][]map[string]*giu.Texture
}

// Dispose clears viewers state
func (s *widgetState) Dispose() {
	s.textures = nil
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
			ShowGrid:  true,
			ShowFloor: true,
			ShowWall:  true,
		},
		tileGroups: p.groupTilesByIdentity(),
	}

	p.setState(state)
}
