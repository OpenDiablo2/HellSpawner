package animdatawidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
)

type widget struct {
	id string
	d2 *d2animdata.AnimationData
}

// Create creates a new widget
func Create(id string, d2 *d2animdata.AnimationData) giu.Widget {
	result := &widget{
		id: id,
		d2: d2,
	}

	return result
}

// Build builds widget
func (p *widget) Build() {
	giu.Layout{
		giu.Label("an example animation data (D2) editor"),
	}.Build()
}
