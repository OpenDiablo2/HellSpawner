package stringtablewidget

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"
)

const (
	deleteW, deleteH             = 50, 25
	addEditW, addEditH           = 200, 30
	actionButtonW, actionButtonH = 100, 30
)

type widget struct {
	id   string
	dict d2tbl.TextDictionary
}

// Create creates a new string table editor widget
func Create(state []byte, id string, dict d2tbl.TextDictionary) giu.Widget {
	result := &widget{
		id:   id,
		dict: dict,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding string table editor state: %v", err)
		}

		result.setState(s)
	}

	return result
}

func (p *widget) Build() {
	state := p.getState()

	switch state.Mode {
	case widgetModeViewer:
		p.buildTableLayout()
	case widgetModeAddEdit:
		p.buildAddEditLayout()
	}
}

func (p *widget) buildTableLayout() {
	state := p.getState()

	keys := p.generateTableKeys()

	rows := make([]*giu.TableRowWidget, len(keys)+1)

	columns := []string{"key", "value", "action"}
	columnWidgets := make([]giu.Widget, len(columns))

	for idx := range columns {
		columnWidgets[idx] = giu.Label(columns[idx])
	}

	rows[0] = giu.TableRow(columnWidgets...)

	for keyIdx, key := range keys {
		// first row is header
		rows[keyIdx+1] = p.makeTableRow(key)
	}

	giu.Layout{
		giu.Button("Add/Edit record##"+p.id+"addEditRecord").
			Size(addEditW, addEditH).OnClick(func() {
			state.Editable = true
			state.Mode = widgetModeAddEdit
		}),
		giu.Separator(),
		p.makeSearchSection(),
		giu.Separator(),
		giu.Custom(func() {
			if len(keys) == 0 {
				giu.Label("Nothing to display.").Build()

				return
			}
			giu.Layout{
				giu.Child("##" + p.id + "tableArea").Border(false).Layout(giu.Layout{
					giu.Table("##" + p.id + "table").FastMode(true).Rows(rows...),
				}),
			}.Build()
		}),
	}.Build()
}

func (p *widget) makeTableRow(key string) *giu.TableRowWidget {
	state := p.getState()

	return giu.TableRow(
		giu.Label(key),
		giu.Label(p.dict[key]),
		giu.Row(
			giu.Button("delete##"+p.id+"deleteString"+key).Size(deleteW, deleteH).OnClick(func() {
				delete(p.dict, key)
				p.reloadMapValues()
			}),
			giu.Button("edit##"+p.id+"editButton"+key).Size(deleteW, deleteH).OnClick(func() {
				state.Key = key
				state.Editable = false
				p.updateValueText()
				state.Mode = widgetModeAddEdit
			}),
		),
	)
}

func (p *widget) makeSearchSection() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.Checkbox("only no-named (starting from #) labels##"+p.id+"numOnly", &state.NumOnly),
		giu.Custom(func() {
			if !state.NumOnly {
				giu.Row(
					giu.Label("Search:"),
					giu.InputText("##"+p.id+"search", &state.Search),
				).Build()
			}
		}),
	}
}

func (p *widget) buildAddEditLayout() {
	state := p.getState()

	giu.Layout{
		giu.Label("Key:"),
		giu.Custom(func() {
			checkbox := giu.Checkbox("no-name##"+p.id+"addEditNoName", &state.NoName).OnChange(func() {
				if state.NoName {
					firstFreeNoName := p.calculateFirstFreeNoName()
					state.Key = "#" + strconv.Itoa(firstFreeNoName)
					p.updateValueText()
				}
			})

			if state.Editable {
				giu.Row(
					p.makeKeyField("##"+p.id+"addEditKey"),
					checkbox,
				).Build()
			} else {
				giu.Label(state.Key).Build()
			}
		}),
		giu.Label("Value:"),
		giu.InputTextMultiline("##"+p.id+"addEditValue", &state.Value),
		giu.Separator(),
		giu.Row(
			giu.Custom(func() {
				var btnStr string

				key := state.Key
				if key == "" {
					return
				}

				_, found := p.dict[key]
				if found {
					btnStr = "Edit"
				} else {
					btnStr = "Add"
				}

				giu.Button(btnStr+"##"+p.id+"addEditAcceptButton").
					Size(actionButtonW, actionButtonH).
					OnClick(func() {
						p.dict[key] = state.Value
						p.reloadMapValues()
						state.Mode = widgetModeViewer
					}).
					Build()
			}),
			giu.Button("cancel##"+p.id+"addEditCancel").
				Size(actionButtonW, actionButtonH).OnClick(func() {
				state.Mode = widgetModeViewer
			}),
		),
		giu.Separator(),
		giu.Label("Tip: enter existing key in key field to edit it"),
		giu.Label("Tip: you don't have to enter key; you can just select \"no-name\""),
	}.Build()
}

func (p *widget) makeKeyField(id string) giu.Widget {
	state := p.getState()

	return giu.InputText(id, &state.Key).OnChange(func() {
		p.formatKey(&state.Key)
		p.updateValueText()
	})
}
