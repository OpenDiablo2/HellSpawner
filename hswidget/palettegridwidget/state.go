package palettegridwidget

import (
	"fmt"
	"image"
	"image/color"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

// PaletteGridState represents palette grid's state
type widgetState struct {
	rgba *giu.Texture
}

// Dispose cleans palette grids state
func (s *widgetState) Dispose() {
	s.rgba = nil
}

func (p *PaletteGridWidget) getStateID() string {
	return fmt.Sprintf("PaletteGridWidget_%s", p.id)
}

func (p *PaletteGridWidget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.setState(&widgetState{})
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *PaletteGridWidget) initState() {
	state := &widgetState{}
	p.setState(state)

	p.rebuildImage()
}

func (p *PaletteGridWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *PaletteGridWidget) rebuildImage() {
	rgb := image.NewRGBA(image.Rect(0, 0, gridWidth*cellSize, gridHeight*cellSize))

	for y := 0; y < gridHeight*cellSize; y++ {
		if y%cellSize == 0 {
			continue
		}

		for x := 0; x < gridWidth*cellSize; x++ {
			if x%cellSize == 0 {
				continue
			}

			idx := (x / cellSize) + ((y / cellSize) * gridWidth)

			c := (*p.colors)[idx]
			col := hsutil.Color(c.RGBA())

			// nolint:gomnd // const
			rgb.Set(x, y, color.RGBA{R: col.R, G: col.G, B: col.B, A: 255})
		}
	}

	go func() {
		p.textureLoader.CreateTextureFromARGB(rgb, func(texture *giu.Texture) {
			p.setState(&widgetState{rgba: texture})
		})
	}()
}
