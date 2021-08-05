package animationwidget

import (
	"fmt"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type widget struct {
	id            string
	palette       *[256]d2interface.Color
	textureLoader hscommon.TextureLoader
}

func (w *widget) getID() string {
	return w.id
}

func (w *widget) getTextureLoader() hscommon.TextureLoader {
	return w.textureLoader
}

func (w *widget) getStateID() string {
	return fmt.Sprintf("widget_%s", w.id)
}

func createWidget(palette *[256]d2interface.Color, textureLoader hscommon.TextureLoader, id string) *widget {
	widget := &widget{
		id:            id,
		palette:       palette,
		textureLoader: textureLoader,
	}

	return widget
}
