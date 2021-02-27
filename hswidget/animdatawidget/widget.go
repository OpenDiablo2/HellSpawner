package animdatawidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
)

type widget struct {
	id string
	d2 *d2animdata.AnimationData
}

func Create(id string, d2 *d2animdata.AnimationData) giu.Widget {
	result := &widget{
		id: id,
		d2: d2,
	}

	return result
}

func (p *widget) Build() {
	giu.Layout(giu.Layout{
		giu.Label("an example animation data (D2) editor"),
	}).Build()
}
