package palettegrideditorwidget

import (
	"encoding/json"
	"log"
	"math"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

const (
	actionButtonW, actionButtonH = 250, 30
)

// PaletteGridEditorWidget represents a palette grid editor
type PaletteGridEditorWidget struct {
	id            string
	colors        *[]palettegridwidget.PaletteColor
	textureLoader hscommon.TextureLoader
	onChange      func()
}

// Create creates a new palette grid editor widget
func Create(state []byte,
	textureLoader hscommon.TextureLoader,
	id string,
	colors *[]palettegridwidget.PaletteColor) *PaletteGridEditorWidget {
	result := &PaletteGridEditorWidget{
		id:            id,
		colors:        colors,
		textureLoader: textureLoader,
		onChange:      nil,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()

		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error loading palette grid editor state: %v", err)
		}

		result.setState(s)
	}

	return result
}

// OnChange sets on change callback
// this callback is ran, when editor's slider or field gets change
func (p *PaletteGridEditorWidget) OnChange(onChange func()) *PaletteGridEditorWidget {
	p.onChange = onChange
	return p
}

// Build Builds a widget
func (p *PaletteGridEditorWidget) Build() {
	state := p.getState()

	colors := make([]palettegridwidget.PaletteColor, len(*p.colors))
	for n := range *(p.colors) {
		colors[n] = (*p.colors)[n]
	}

	grid := palettegridwidget.Create(p.textureLoader, p.id, &colors).OnClick(func(idx int) {
		color := hsutil.Color((*p.colors)[idx].RGBA())
		state.RGBA = color
		state.Idx = idx

		state.Mode = widgetModeEdit
	})

	grid.Build()

	if state.Mode == widgetModeEdit {
		p.buildEditor(grid)
	}
}

func (p *PaletteGridEditorWidget) buildEditor(grid *palettegridwidget.PaletteGridWidget) {
	state := p.getState()

	isOpen := state.Mode == widgetModeEdit
	onChange := func() {
		p.changeColor(state)
		grid.UpdateImage()

		if p.onChange != nil {
			p.onChange()
		}
	}

	giu.Layout{
		giu.PopupModal("Edit color").IsOpen(&isOpen).Layout(
			giu.ColorEdit("##edit color", &state.RGBA).Flags(giu.ColorEditFlagsNoAlpha),
			giu.Separator(),
			giu.Button("OK##"+p.id+"editColorOK").Size(actionButtonW, actionButtonH).OnClick(func() {
				onChange()
				state.Mode = widgetModeGrid
			}),
		),
		// handle clicking on "X" button of popup
		giu.Custom(func() {
			if !isOpen {
				onChange()
				state.mode = widgetModeGrid
			}
		}),
	}.Build()
}

func (p *PaletteGridEditorWidget) makeRGBField(id, label string, field *uint8, grid *palettegridwidget.PaletteGridWidget) giu.Layout {
	state := p.getState()

	f32 := int32(*field)

	return giu.Layout{
		giu.Row(
			giu.Label(label),
			hswidget.MakeInputInt(
				id,
				inputIntW,
				field,
				func() {
					p.changeColor(state)
					grid.UpdateColorTexture(state.Idx)
					if p.onChange != nil {
						p.onChange()
					}
					state.Hex = hsutil.RGB2Hex(state.R, state.G, state.B)
				},
			),
		),
		giu.SliderInt(id+"Slider", &f32, 0, math.MaxUint8).OnChange(func() {
			p.changeColor(state)
			grid.UpdateColorTexture(state.Idx)
			if p.onChange != nil {
				p.onChange()
			}
			state.Hex = hsutil.RGB2Hex(state.R, state.G, state.B)
			hswidget.SetByteToInt(f32, field)
		}),
	}
}
