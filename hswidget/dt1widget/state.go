package dt1widget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
)

type dt1Controls struct {
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
	*dt1Controls

	lastTileGroup int32

	tileGroups [][]*d2dt1.Tile
	textures   [][]map[string]*giu.Texture
}

// Dispose clears viewers state
func (is *widgetState) Dispose() {
	is.textures = nil
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
		dt1Controls: &dt1Controls{
			showGrid:  true,
			showFloor: true,
			showWall:  true,
		},
		tileGroups: p.groupTilesByIdentity(),
	}

	p.setState(state)
}
