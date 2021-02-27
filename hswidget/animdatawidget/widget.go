package animdatawidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
)

const (
	listW, listH = 200, 400
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

	switch state.mode {
	case widgetModeList:
		p.buildAnimationsList()
	}
}

func (p *widget) buildAnimationsList() {
	state := p.getState()

	list := make([]*giu.SelectableWidget, p.d2.GetRecordsCount())
	for idx, name := range state.mapKeys {
		list[idx] = giu.Selectable(name)
	}

	giu.Layout{
		giu.Child("##"+p.id+"keyList").Border(false).
			Size(listW, listH).
			Layout(giu.Layout{
				giu.Custom(func() {
					for _, i := range list {
						i.Build()
					}
				}),
			}),
	}.Build()
}
