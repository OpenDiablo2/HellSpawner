package stringtablewidget

import (
	"sort"
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

	var numEntries int

	if !state.numOnly {
		numEntries = len(state.keys)
	} else {
		for _, key := range state.keys {
			if key[0] == '#' {
				numEntries++
			} else {
				// labels are sorted, so no-name (starting from # are on top)
				break
			}
		}
	}

	// wprobably will remove
	if !(numEntries > 0) {
		giu.Layout{}.Build()
	}

	rows := make([]*giu.RowWidget, numEntries+1)

	columns := []string{"key", "value", "action"}
	columnWidgets := make([]giu.Widget, len(columns))

	for idx := range columns {
		columnWidgets[idx] = giu.Label(columns[idx])
	}

	rows[0] = giu.Row(columnWidgets...)

	for keyIdx := 0; keyIdx < numEntries; keyIdx++ {
		key := state.keys[keyIdx]
		// nolint:gomnd // first row is header
		rows[keyIdx+1] = giu.Row(
			giu.Label(key),
			giu.Label(p.dict[key]),
			giu.Line(
				giu.Button("delete##"+p.id+"deleteString"+strconv.Itoa(keyIdx)).Size(deleteW, deleteH).OnClick(func() {
					delete(p.dict, key)
					p.reloadMapValues()
				}),
				giu.Button("edit##"+p.id+"editButton"+strconv.Itoa(keyIdx)).Size(deleteW, deleteH).OnClick(func() {
					state.key = key
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
			if state.editable {
				giu.Line(
					giu.InputText("##"+p.id+"addEditKey", &state.key).OnChange(func() {
						p.updateValueText()
					}),
					giu.Checkbox("no-name##"+p.id+"addEditNoName", &state.noName).OnChange(func() {
						if state.noName {
							firstFreeNoName := p.calculateFirstFreeNoName()
							state.key = "#" + strconv.Itoa(firstFreeNoName)
							p.updateValueText()
						}
					}),
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
		giu.Label("Tip: enter existing key in value vield to edit it"),
	}.Build()
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
