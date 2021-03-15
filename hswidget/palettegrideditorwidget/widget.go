package palettegrideditorwidget

import (
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

const (
	actionButtonW, actionButtonH = 250, 30
	inputIntW                    = 30
)

// PaletteGridEditorWidget represents a palette grid editor
type PaletteGridEditorWidget struct {
	id            string
	colors        *[256]palettegridwidget.PaletteColor
	textureLoader *hscommon.TextureLoader
	onChange      func()
}

// Create creates a new palette grid editor widget
func Create(state []byte,
	textureLoader *hscommon.TextureLoader,
	id string,
	colors *[256]palettegridwidget.PaletteColor) *PaletteGridEditorWidget {
	result := &PaletteGridEditorWidget{
		id:            id,
		colors:        colors,
		textureLoader: textureLoader,
		onChange:      nil,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
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
	grid := palettegridwidget.Create(p.textureLoader, p.id, p.colors).OnClick(func(idx int) {
		color := hsutil.Color(p.colors[idx].RGBA())
		state.r = color.R
		state.g = color.G
		state.b = color.B
		state.idx = idx

		state.mode = widgetModeEdit
	})

	switch state.mode {
	case widgetModeGrid:
		grid.Build()
	case widgetModeEdit:
		p.buildEditor(grid)
	}
}

func (p *PaletteGridEditorWidget) buildEditor(grid *palettegridwidget.PaletteGridWidget) {
	state := p.getState()

	giu.Layout{
		giu.Label("Edit Color: "),
		giu.Image(grid.GetColorTexture(state.idx)),
		giu.Separator(),
		p.makeRGBField("##"+p.id+"changeR", "R:", &state.r, grid),
		giu.Separator(),
		p.makeRGBField("##"+p.id+"changeG", "G:", &state.g, grid),
		giu.Separator(),
		p.makeRGBField("##"+p.id+"changeB", "B:", &state.b, grid),
		giu.Separator(),
		giu.Line(
			giu.Label("Hex: "),
			giu.InputText("##"+p.id+"editHex", &state.hex).OnChange(func() {
				r, g, b, err := Hex2RGB(state.hex)
				if err != nil {
					log.Print("error: ", err)
				}

				grid.UpdateColorTexture(state.idx)

				state.r, state.g, state.b = r, g, b
			}),
		),
		giu.Separator(),
		giu.Button("OK##"+p.id+"editColorOK").Size(actionButtonW, actionButtonH).OnClick(func() {
			if p.onChange != nil {
				p.onChange()
			}
			state.mode = widgetModeGrid
		}),
	}.Build()
}

func (p *PaletteGridEditorWidget) makeRGBField(id, label string, field *uint8, grid *palettegridwidget.PaletteGridWidget) giu.Layout {
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
					p.changeColor(state)
					grid.UpdateColorTexture(state.idx)
					if p.onChange != nil {
						p.onChange()
					}
					state.hex = RGB2Hex(state.r, state.g, state.b)
				},
			),
		),
		giu.SliderInt(id+"Slider", &f32, 0, 255).OnChange(func() {
			p.changeColor(state)
			grid.UpdateColorTexture(state.idx)
			if p.onChange != nil {
				p.onChange()
			}
			state.hex = RGB2Hex(state.r, state.g, state.b)
			hsutil.SetByteToInt(f32, field)
		}),
	}
}
