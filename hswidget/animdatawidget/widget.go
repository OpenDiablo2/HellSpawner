package animdatawidget

import (
	"fmt"
	"strings"

	"github.com/OpenDiablo2/dialog"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
)

const (
	listW, listH                         = 200, 400
	inputIntW                            = 30
	actionBtnW, actionBtnH               = 200, 30
	saveCancelButtonW, saveCancelButtonH = 50, 30
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

	list := make([]giu.Widget, p.d2.GetRecordsCount())

	for idx, name := range state.mapKeys {
		currentIdx := idx
		list[idx] = giu.Selectable(name).OnClick(func() {
			state.mapIndex = int32(currentIdx)
			state.mode = widgetModeViewRecord
		})
	}

	giu.Layout{
		p.makeSearchLayout(),
		giu.Separator(),
		giu.Child("##"+p.id+"keyList").Border(false).
			Size(listW, listH).
			Layout(list),
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
				record.SetFramesPerDirection(uint32(fpd))
			}),
		),
		giu.Line(
			giu.Label("Speed: "),
			giu.InputInt("##"+p.id+"recordSpeed", &speed).Size(inputIntW).OnChange(func() {
				record.SetSpeed(uint16(speed))
			}),
		),
		giu.Label(fmt.Sprintf("FPS: %v", record.FPS())),
		giu.Label(fmt.Sprintf("Frame duration: %v (miliseconds)", record.FrameDurationMS())),
		giu.Separator(),
		giu.Button("Select another record##"+p.id+"backToRecordSelection").Size(actionBtnW, actionBtnH).OnClick(func() {
			state.mode = widgetModeList
		}),
		giu.Button("Add record##"+p.id+"addRecordBtn").Size(actionBtnW, actionBtnH).OnClick(func() {
			dialog.Message("available after merging https://github.com/OpenDiablo2/OpenDiablo2/pulls/1086").Info()
			// p.d2.PushRecord(name)
			// nolint:gomnd // list index
			// state.recordIdx = len(records)-1
		}),
	}.Build()
}

func (p *widget) makeSearchLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Label("Search or type new entry name:"),
		giu.InputText("##"+p.id+"newEntryName", &state.name).Size(listW).OnChange(func() {
			// formatting
			state.name = strings.ToUpper(state.name)
			state.name = strings.ReplaceAll(state.name, " ", "")
		}),
		giu.Custom(func() {
			if state.name == "" {
				return
			}

			found := (len(p.d2.GetRecords(state.name)) > 0)
			if found {
				giu.Line(
					giu.Button("View##"+p.id+"addEntryViewEntry").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
						p.viewRecord()
					}),
				).Build()

				return
			}

			giu.Line(
				giu.Button("Add##"+p.id+"addEntry").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
					dialog.Message("available after merging https://github.com/OpenDiablo2/OpenDiablo2/pulls/1086").Info()
					// p.d2.AddRecord(state.name)
					// p.viewRecord()
				}),
			).Build()
		}),
	}
}

func (p *widget) viewRecord() {
	state := p.getState()

	for n, i := range state.mapKeys {
		if i == state.name {
			state.mapIndex = int32(n)
		}
	}

	state.mode = widgetModeViewRecord
}
