package stringtablewidget

import (
	"sort"
	"strconv"
	"strings"

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
func Create(id string, dict d2tbl.TextDictionary) giu.Widget {
	result := &widget{
		id:   id,
		dict: dict,
	}

	return result
}

func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case widgetModeViewer:
		p.buildTableLayout()
	case widgetModeAddEdit:
		p.buildAddEditLayout()
	}
}

func (p *widget) buildTableLayout() {
	state := p.getState()

	keys := make([]string, 0)

	switch {
	case state.numOnly:
		for _, key := range state.keys {
			if key[0] == '#' {
				keys = append(keys, key)
			} else {
				// labels are sorted, so no-name (starting from # are on top)
				break
			}
		}
	case state.search != "":
		for _, key := range state.keys {
			k := strings.ToLower(key)
			v := strings.ToLower(p.dict[key])

			switch {
			case strings.Contains(k, state.search),
				strings.Contains(v, state.search):
				keys = append(keys, key)
			}
		}
	default:
		keys = state.keys
	}

	rows := make([]*giu.RowWidget, len(keys)+1)

	columns := []string{"key", "value", "action"}
	columnWidgets := make([]giu.Widget, len(columns))

	for idx := range columns {
		columnWidgets[idx] = giu.Label(columns[idx])
	}

	rows[0] = giu.Row(columnWidgets...)

	for keyIdx, key := range keys {
		currentKey := key
		// nolint:gomnd // first row is header
		rows[keyIdx+1] = giu.Row(
			giu.Label(currentKey),
			giu.Label(p.dict[currentKey]),
			giu.Line(
				giu.Button("delete##"+p.id+"deleteString"+strconv.Itoa(keyIdx)).Size(deleteW, deleteH).OnClick(func() {
					delete(p.dict, currentKey)
					p.reloadMapValues()
				}),
				giu.Button("edit##"+p.id+"editButton"+strconv.Itoa(keyIdx)).Size(deleteW, deleteH).OnClick(func() {
					state.key = currentKey
					state.editable = false
					p.updateValueText()
					state.mode = widgetModeAddEdit
				}),
			),
		)
	}

	giu.Layout{
		giu.Button("Add/Edit record##"+p.id+"addEditRecord").
			Size(addEditW, addEditH).OnClick(func() {
			state.editable = true
			state.mode = widgetModeAddEdit
		}),
		giu.Separator(),
		giu.Line(
			giu.Checkbox("only no-named (starting from #) labels##"+p.id+"numOnly", &state.numOnly),
		),
		giu.Custom(func() {
			if !state.numOnly {
				giu.Line(
					giu.Label("Search:"),
					giu.InputText("##"+p.id+"search", &state.search).OnChange(func() {
						p.formatKey(&state.search)
					}),
				).Build()
			}
		}),
		giu.Separator(),
		giu.Child("").Border(false).Layout(giu.Layout{
			giu.FastTable("").Border(true).Rows(rows),
		}),
	}.Build()
}

func (p *widget) buildAddEditLayout() {
	state := p.getState()

	giu.Layout{
		giu.Label("Key:"),
		giu.Custom(func() {
			checkbox := giu.Checkbox("no-name##"+p.id+"addEditNoName", &state.noName).OnChange(func() {
				if state.noName {
					firstFreeNoName := p.calculateFirstFreeNoName()
					state.key = "#" + strconv.Itoa(firstFreeNoName)
					p.updateValueText()
				}
			})

			if state.editable {
				giu.Line(
					p.makeKeyField("##"+p.id+"addEditKey"),
					checkbox,
				).Build()
			} else {
				giu.Label(state.key).Build()
			}
		}),
		giu.Label("Value:"),
		giu.InputTextMultiline("##"+p.id+"addEditValue", &state.value),
		giu.Separator(),
		giu.Line(
			giu.Custom(func() {
				var btnStr string

				key := state.key
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
						p.dict[key] = state.value
						p.reloadMapValues()
						state.mode = widgetModeViewer
					}).
					Build()
			}),
			giu.Button("cancel##"+p.id+"addEditCancel").
				Size(actionButtonW, actionButtonH).OnClick(func() {
				state.mode = widgetModeViewer
			}),
		),
		giu.Separator(),
		giu.Label("Tip: enter existing key in key field to edit it"),
		giu.Label("Tip: you don't have to enter key; you can just select \"no-name\""),
	}.Build()
}

func (p *widget) makeKeyField(id string) giu.Widget {
	state := p.getState()

	return giu.InputText(id, &state.key).OnChange(func() {
		p.formatKey(&state.key)
		p.updateValueText()
	})
}

func (p *widget) calculateFirstFreeNoName() (firstFreeNoName int) {
	state := p.getState()

	ints := make([]int, 0)

	for _, key := range state.keys {
		if key[0] == '#' {
			idx, err := strconv.Atoi(key[1:])
			if err != nil {
				continue
			}

			ints = append(ints, idx)
		}
	}

	sort.Ints(ints)

	for n, i := range ints {
		if n != i {
			firstFreeNoName = n
			break
		}
	}

	return
}

func (p *widget) updateValueText() {
	state := p.getState()

	str, found := p.dict[state.key]
	if found {
		state.value = str
	} else {
		state.value = ""
	}
}

func (p *widget) formatKey(s *string) {
	*s = strings.ReplaceAll(*s, " ", "_")
}
