package palettegridwidget

import (
	"sync"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	inputIntW = 30
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
					idx := y*gridWidth + x
					line = append(
						line,
						giu.ImageButton(state.texture[idx]).
							Size(cellSize, cellSize).OnClick(func() {
							color := p.colors[idx]
							state.idx = idx
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
	state := p.getState()

	giu.Layout{
		giu.Label("Edit Color: "),
		giu.Image(state.texture[state.idx]),
		p.makeRGBField("##"+p.id+"changeR", &state.r),
		p.makeRGBField("##"+p.id+"changeG", &state.g),
		p.makeRGBField("##"+p.id+"changeB", &state.b),
	}.Build()
}

func (p *widget) makeRGBField(id string, field *uint8) giu.Layout {
	state := p.getState()

	f32 := int32(*field)

	return giu.Layout{
		hsutil.MakeInputInt(
			id,
			inputIntW,
			field,
			func() {
				p.changeColor(
					state.r,
					state.g,
					state.b,
					state.idx,
				)
			},
		),
		giu.SliderInt(id+"Slider", &f32, 0, 255).OnChange(func() {
			// we need to lock, because sometimes crashes
			var mutex = &sync.Mutex{}
			mutex.Lock()
			p.changeColor(
				state.r,
				state.g,
				state.b,
				state.idx,
			)
			hsutil.SetByteToInt(f32, field)
			mutex.Unlock()
		}),
	}
}
