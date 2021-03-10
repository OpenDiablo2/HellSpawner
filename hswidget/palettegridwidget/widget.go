package palettegridwidget

import (
	"sync"

	"github.com/ianling/giu"

	//"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	inputIntW                    = 30
	actionButtonW, actionButtonH = 280, 30
)

const (
	gridWidth  = 16
	gridHeight = 16
	cellSize   = 12
)

type widget struct {
	id            string
	colors        *[256]PaletteColor
	textureLoader *hscommon.TextureLoader
}

// Create creates a new palette grid widget
func Create(tl *hscommon.TextureLoader, id string, colors *[256]PaletteColor) giu.Widget {
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
							color := hsutil.Color(p.colors[idx].RGBA())
							state.idx = idx
							state.r = color.R
							state.g = color.G
							state.b = color.B

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
		giu.Separator(),
		p.makeRGBField("##"+p.id+"changeR", "R:", &state.r),
		giu.Separator(),
		p.makeRGBField("##"+p.id+"changeG", "G:", &state.g),
		giu.Separator(),
		p.makeRGBField("##"+p.id+"changeB", "B:", &state.b),
		giu.Separator(),
		giu.Button("OK##"+p.id+"editColorOK").Size(actionButtonW, actionButtonH).OnClick(func() {
			state.mode = widgetModeGrid
		}),
	}.Build()
}

func (p *widget) makeRGBField(id, label string, field *uint8) giu.Layout {
	state := p.getState()

	f32 := int32(*field)

	return giu.Layout{
		giu.Line(
			giu.Label(label),
			hsutil.MakeInputInt(
				id,
				inputIntW,
				field,
				func() {
					p.changeColor(
						state,
					)
				},
			),
		),
		giu.SliderInt(id+"Slider", &f32, 0, 255).OnChange(func() {
			// we need to lock, because sometimes crashes
			var mutex = &sync.Mutex{}
			mutex.Lock()
			p.changeColor(
				state,
			)
			hsutil.SetByteToInt(f32, field)
			mutex.Unlock()
		}),
	}
}
