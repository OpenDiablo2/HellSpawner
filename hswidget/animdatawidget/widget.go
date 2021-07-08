package animdatawidget

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ianling/giu"

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
		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding animation data widget state: %v", err)
		}

		result.setState(s)
	}

	return result
}

// Build builds widget
func (p *widget) Build() {
	state := p.getState()

	switch state.Mode {
	case widgetModeList:
		p.buildAnimationsList()
	case widgetModeViewRecord:
		p.buildViewRecordLayout()
	}
}

func (p *widget) buildAnimationsList() {
	state := p.getState()

	keys := make([]string, 0)

	if state.Name != "" {
		for _, key := range state.mapKeys {
			if strings.Contains(key, state.Name) {
				keys = append(keys, key)
			}
		}
	} else {
		keys = state.mapKeys
	}

	list := make([]giu.Widget, len(keys))

	const imageButtonSize = 13

	for idx, name := range keys {
		currentIdx := idx
		list[idx] = giu.Row(
			hswidget.MakeImageButton(
				"##"+p.id+"deleteEntry"+strconv.Itoa(currentIdx),
				imageButtonSize, imageButtonSize,
				state.deleteIcon,
				func() {
					p.deleteEntry(state.mapKeys[currentIdx])
				},
			),
			giu.Selectable(name).OnClick(func() {
				state.MapIndex = int32(currentIdx)
				state.Mode = widgetModeViewRecord
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

	name := state.mapKeys[state.MapIndex]
	records := p.d2.GetRecords(name)
	record := records[state.RecordIdx]

	max := len(records) - 1

	fpd := int32(record.FramesPerDirection())
	speed := int32(record.Speed())

	giu.Layout{
		giu.Row(
			giu.ArrowButton("##"+p.id+"previousAnimation", giu.DirectionLeft).OnClick(func() {
				state.RecordIdx = 0

				if state.MapIndex > 0 {
					state.MapIndex--
				}
			}),
			giu.Label(fmt.Sprintf("Animation name: %s", name)),
			giu.ArrowButton("##"+p.id+"nextAnimation", giu.DirectionRight).OnClick(func() {
				state.RecordIdx = 0

				if int(state.MapIndex) < len(state.mapKeys)-1 {
					state.MapIndex++
				}
			}),
		),
		giu.Separator(),
		giu.Custom(func() {
			if max > 0 {
				giu.Layout{
					giu.SliderInt("record##"+p.id, &state.RecordIdx, 0, int32(max)),
					giu.Separator(),
				}.Build()
			}
		}),
		giu.Row(
			giu.Label("Frames per direction: "),
			giu.InputInt("##"+p.id+"recordFramesPerDirection", &fpd).Size(inputIntW).OnChange(func() {
				record.SetFramesPerDirection(uint32(fpd))
			}),
		),
		giu.Row(
			giu.Label("Speed: "),
			giu.InputInt("##"+p.id+"recordSpeed", &speed).Size(inputIntW).OnChange(func() {
				record.SetSpeed(uint16(speed))
			}),
		),
		giu.Label(fmt.Sprintf("FPS: %v", record.FPS())),
		giu.Label(fmt.Sprintf("Frame duration: %v (miliseconds)", record.FrameDurationMS())),
		giu.Separator(),
		giu.Button("Back to entry preview##"+p.id+"backToRecordSelection").Size(actionBtnW, actionBtnH).OnClick(func() {
			state.Mode = widgetModeList
		}),
		giu.Button("Add record##"+p.id+"addRecordBtn").Size(actionBtnW, actionBtnH).OnClick(func() {
			p.d2.PushRecord(name)

			// no -1, because current records hasn't new field yet
			state.RecordIdx = int32(len(records))
		}),
		giu.Button("Delete record##"+p.id+"deleteRecordBtn").Size(actionBtnW, actionBtnH).OnClick(func() {
			if len(records) == 1 {
				state.RecordIdx = 0
				state.Mode = widgetModeList
				p.deleteEntry(name)

				return
			}
			if state.RecordIdx == int32(len(records)-1) {
				if state.RecordIdx > 0 {
					state.RecordIdx--
				} else {
					state.Mode = widgetModeList
				}
			}

			err := p.d2.DeleteRecord(name, int(state.RecordIdx))
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
		giu.InputText("##"+p.id+"newEntryName", &state.Name).Size(listW).OnChange(func() {
			// formatting
			state.Name = strings.ToUpper(state.Name)
			state.Name = strings.ReplaceAll(state.Name, " ", "")
		}),
		giu.Custom(func() {
			if state.Name == "" {
				return
			}

			found := (len(p.d2.GetRecords(state.Name)) > 0)
			if found {
				giu.Row(
					giu.Button("View##"+p.id+"addEntryViewEntry").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
						p.viewRecord()
					}),
				).Build()

				return
			}

			giu.Row(
				giu.Button("Add##"+p.id+"addEntry").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
					err := p.d2.AddEntry(state.Name)
					if err != nil {
						log.Print(err)
					}

					p.d2.PushRecord(state.Name)
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
		if i == state.Name {
			state.MapIndex = int32(n)
		}
	}

	state.Mode = widgetModeViewRecord
}

func (p *widget) deleteEntry(name string) {
	if err := p.d2.DeleteEntry(name); err != nil {
		log.Print(fmt.Errorf("deleting entry: %w", err))
	}

	p.reloadMapKeys()
}
