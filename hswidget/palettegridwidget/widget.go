package palettegridwidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

const (
	gridWidth  = 16
	gridHeight = 16
	cellSize   = 12
)

// PaletteGridWidget represents a palette grid
type PaletteGridWidget struct {
	id            string
	colors        *[]PaletteColor
	textureLoader hscommon.TextureLoader
	onClick       func(idx int)
}

// Create creates a new palette grid widget
func Create(tl hscommon.TextureLoader, id string, colors *[]PaletteColor) *PaletteGridWidget {
	result := &PaletteGridWidget{
		id:            id,
		colors:        colors,
		textureLoader: tl,
		onClick:       nil,
	}

	return result
}

// OnClick sets onClick callback
func (p *PaletteGridWidget) OnClick(onClick func(idx int)) *PaletteGridWidget {
	p.onClick = onClick
	return p
}

// GetColorTexture returns selected color texture
func (p *PaletteGridWidget) GetColorTexture(idx int) *giu.Texture {
	state := p.getState()
	return state.texture[idx]
}

// UpdateColorTexture updates specified texture
func (p *PaletteGridWidget) UpdateColorTexture(idx int) {
	p.loadTexture(idx)
}

// Build build a new widget
func (p *PaletteGridWidget) Build() {
	state := p.getState()

	giu.Layout{
		giu.Custom(func() {
			var grid giu.Layout = make([]giu.Widget, 0)

			for y := 0; y < gridHeight; y++ {
				line := make([]giu.Widget, 0)

				for x := 0; x < gridWidth; x++ {
					idx := y*gridWidth + x
					line = append(
						line,
						giu.ImageButton(state.texture[idx]).
							Size(cellSize, cellSize).OnClick(func() {
							if p.onClick != nil {
								p.onClick(idx)
								p.loadTexture(idx)
							}
						}),
					)
				}

				grid = append(grid, giu.Line(line...))
			}

			grid.Build()
		}),
	}.Build()
}
