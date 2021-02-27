package animdatawidget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
)

const (
	listW, listH       = 200, 400
	inputIntW          = 30
	backBtnW, backBtnH = 200, 30
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
	case widgetModeViewRecord:
		p.buildViewRecordLayout()
	}
}

func (p *widget) buildAnimationsList() {
	state := p.getState()

	list := make([]*giu.SelectableWidget, p.d2.GetRecordsCount())

	for idx, name := range state.mapKeys {
		currentIdx := idx
		list[idx] = giu.Selectable(name).OnClick(func() {
			state.mapIndex = int32(currentIdx)
			state.mode = widgetModeViewRecord
		})
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

func (p *widget) buildViewRecordLayout() {
	state := p.getState()

	name := state.mapKeys[state.mapIndex]
	records := p.d2.GetRecords(name)
	record := records[state.recordIdx]

	max := len(records) - 1

	fpd := int32(record.FramesPerDirection())
	speed := int32(record.Speed())

	giu.Layout{
		giu.Line(
			giu.ArrowButton("##"+p.id+"previousAnimation", giu.DirectionLeft).OnClick(func() {
				if state.mapIndex > 0 {
					state.mapIndex--
				}
			}),
			giu.Label(fmt.Sprintf("Animation name: %s", name)),
			giu.ArrowButton("##"+p.id+"nextAnimation", giu.DirectionRight).OnClick(func() {
				if int(state.mapIndex) < len(state.mapKeys)-1 {
					state.mapIndex++
				}
			}),
		),
		giu.Separator(),
		giu.Custom(func() {
			if max > 0 {
				giu.Layout{
					giu.SliderInt("record##"+p.id, &state.recordIdx, 0, int32(max)),
					giu.Separator(),
				}.Build()
			}
		}),
		giu.Line(
			giu.Label("Frames per direction: "),
			giu.InputInt("##"+p.id+"recordFramesPerDirection", &fpd).Size(inputIntW).OnChange(func() {
				// nolint:gocritic // just for editing in future
				// record.SetFramesPerDirection(fpd)
			}),
		),
		giu.Line(
			giu.Label("Speed: "),
			giu.InputInt("##"+p.id+"recordSpeed", &speed).Size(inputIntW).OnChange(func() {
				// nolint:gocritic // just for editing in future
				// record.SetSpeed(speed)
			}),
		),
		giu.Label(fmt.Sprintf("FPS: %v", record.FPS())),
		giu.Label(fmt.Sprintf("Frame duration: %v (miliseconds)", record.FrameDurationMS())),
		giu.Separator(),
		giu.Button("Select another record##"+p.id+"backToRecordSelection").Size(backBtnW, backBtnH).OnClick(func() {
			state.mode = widgetModeList
		}),
	}.Build()
}
