package palettegrideditorwidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

type PaletteGridEditorWidget struct {
	id            string
	colors        *[256]palettegridwidget.PaletteColor
	textureLoader *hscommon.TextureLoader
	onChange      func()
}

func Create(textureLoader *hscommon.TextureLoader, id string, colors *[256]palettegridwidget.PaletteColor) *PaletteGridEditorWidget {
	result := &PaletteGridEditorWidget{
		id:            id,
		colors:        colors,
		textureLoader: textureLoader,
		onChange:      nil,
	}

	return result
}

func (p *PaletteGridEditorWidget) OnChange(onChange func()) *PaletteGridEditorWidget {
	p.onChange = onChange
	return p
}

func (p *PaletteGridEditorWidget) Build() {
	grid := palettegridwidget.Create(p.textureLoader, p.id, p.colors)

	giu.Layout{
		grid,
	}.Build()
}
