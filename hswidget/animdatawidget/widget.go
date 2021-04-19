package animdatawidget

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

const (
	listW, listH                         = 200, 400
	inputIntW                            = 30
	actionBtnW, actionBtnH               = 200, 30
	saveCancelButtonW, saveCancelButtonH = 50, 30
)

type widget struct {
	id            string
	d2            *d2animdata.AnimationData
	textureLoader hscommon.TextureLoader
}

// Create creates a new widget
func Create(textureLoader hscommon.TextureLoader, state []byte, id string, d2 *d2animdata.AnimationData) giu.Widget {
	result := &widget{
		id:            id,
		d2:            d2,
		textureLoader: textureLoader,
	}

	if state != nil && giu.Context.GetState(result.getStateID()) == nil {
		s := result.getState()
		s.Decode(state)
		result.setState(s)
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

	keys := make([]string, 0)

	if state.name != "" {
		for _, key := range state.mapKeys {
			if strings.Contains(key, state.name) {
				keys = append(keys, key)
			}
		}
	} else {
		keys = state.mapKeys
	}

	list := make([]giu.Widget, len(keys))

	for idx, name := range keys {
		currentIdx := idx
		list[idx] = giu.Line(
			hswidget.MakeImageButton(
				"##"+p.id+"deleteEntry"+strconv.Itoa(currentIdx),
				13, 13,
				state.deleteIcon,
				func() {
					p.deleteEntry(state.mapKeys[currentIdx])
				},
			),
			giu.Selectable(name).OnClick(func() {
				state.mapIndex = int32(currentIdx)
				state.mode = widgetModeViewRecord
			}),
		)
	}

	giu.Layout{
		p.makeSearchLayout(),
		giu.Separator(),
		giu.Child("##"+p.id+"keyList").Border(false).
			Size(listW, listH).
			Layout(giu.Layout{
				giu.Custom(func() {
					if len(list) > 0 {
						giu.Layout(list).Build()

						return
					}

					giu.Label("Nothing matches...").Build()
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
				state.recordIdx = 0

				if state.mapIndex > 0 {
					state.mapIndex--
				}
			}),
			giu.Label(fmt.Sprintf("Animation name: %s", name)),
			giu.ArrowButton("##"+p.id+"nextAnimation", giu.DirectionRight).OnClick(func() {
				state.recordIdx = 0

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
		giu.Button("Back to entry preview##"+p.id+"backToRecordSelection").Size(actionBtnW, actionBtnH).OnClick(func() {
			state.mode = widgetModeList
		}),
		giu.Button("Add record##"+p.id+"addRecordBtn").Size(actionBtnW, actionBtnH).OnClick(func() {
			p.d2.PushRecord(name)

			// no -1, because current records hasn't new field yet
			state.recordIdx = int32(len(records))
		}),
		giu.Button("Delete record##"+p.id+"deleteRecordBtn").Size(actionBtnW, actionBtnH).OnClick(func() {
			if len(records) == 1 {
				state.recordIdx = 0
				state.mode = widgetModeList
				p.deleteEntry(name)

				return
			}
			if state.recordIdx == int32(len(records)-1) {
				if state.recordIdx > 0 {
					state.recordIdx--
				} else {
					state.mode = widgetModeList
				}
			}

			err := p.d2.DeleteRecord(name, int(state.recordIdx))
			if err != nil {
				log.Print(err)
			}
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
					err := p.d2.AddEntry(state.name)
					if err != nil {
						log.Print(err)
					}

					p.d2.PushRecord(state.name)
					p.reloadMapKeys()
					p.viewRecord()
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

func (p *widget) deleteEntry(name string) {
	if err := p.d2.DeleteEntry(name); err != nil {
		log.Print(fmt.Errorf("deleting entry: %w", err))
	}

	p.reloadMapKeys()
}
