package palettegridwidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

const (
	gridWidth  = 16
	gridHeight = 16
	cellSize   = 12
)

type widget struct {
	id            string
	colors        *[256]d2interface.Color
	textureLoader *hscommon.TextureLoader
}

// Create creates a new palette grid widget
func Create(tl *hscommon.TextureLoader, id string, colors *[256]d2interface.Color) giu.Widget {
	result := &widget{
		id:            id,
		colors:        colors,
		textureLoader: tl,
	}

	return result
}

// Build build a new widget
func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case widgetModeGrid:
		p.buildGrid()
	case widgetModeEdit:
		p.buildEditor()
	}
}

func (p *widget) buildGrid() {
	state := p.getState()

	giu.Layout{
		giu.Custom(func() {
			var grid giu.Layout = make([]giu.Widget, 0)

			for y := 0; y < gridHeight; y++ {
				line := make([]giu.Widget, 0)

				for x := 0; x < gridWidth; x++ {
					currentX := x
					line = append(
						line,
						giu.ImageButton(state.texture[y*gridWidth+x]).
							Size(cellSize, cellSize).OnClick(func() {
							color := p.colors[currentX]
							state.idx = currentX
							state.r = color.R()
							state.g = color.G()
							state.b = color.B()

							state.mode = widgetModeEdit
						}),
					)
				}

				grid = append(grid, giu.Line(line...))
			}

			grid.Build()
		}),
	}.Build()
}

func (p *widget) buildEditor() {
	giu.Label("edit color").Build()
}
