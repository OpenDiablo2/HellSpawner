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
	state := p.getState()
	list := make([]*giu.SelectableWidget, p.d2.GetRecordsCount())
	for n, index := range state.mapKeys {
		list[n] = giu.Selectable(index)
	}

	giu.Layout{
		giu.Child("##"+p.id+"keyList").Border(false).
			Size(200, 200).
			Layout(giu.Layout{
				giu.Custom(func() {
					for _, i := range list {
						i.Build()
					}
				}),
			}),
	}.Build()
}
