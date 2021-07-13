package animationwidget

import (
	"fmt"
	
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type Widget struct {
	id            string
	palette       *[256]d2interface.Color
	textureLoader hscommon.TextureLoader
}

func (w *Widget) getStateID() string {
	return fmt.Sprintf("widget_%s", w.id)
}

func CreateWidget(palette *[256]d2interface.Color, textureLoader hscommon.TextureLoader, id string) *Widget {
	widget := &Widget{
		id:            id,
		palette:       palette,
		textureLoader: textureLoader,
	}

	return widget
}