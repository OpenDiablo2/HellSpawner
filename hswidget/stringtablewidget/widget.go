package stringtablewidget

import (
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
	numEntries := len(state.keys)

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

	for keyIdx, key := range state.keys {
		// nolint:gomnd // first row is header
		rows[keyIdx+1] = giu.Row(
			giu.Label(key),
			giu.Label(p.dict[key]),
			giu.Line(
				giu.Button("delete##"+p.id+"deleteString"+strconv.Itoa(keyIdx)).Size(deleteW, deleteH).OnClick(func() {
				}),
			),
		)
	}

	giu.Layout{
		giu.Button("Add/Edit record##"+p.id+"addEditRecord").
			Size(addEditW, addEditH).OnClick(func() {
			state.mode = widgetModeAddEdit
		}),
		giu.Child("").Border(false).Layout(giu.Layout{
			giu.FastTable("").Border(true).Rows(rows),
		}),
	}.Build()
}

func (p *widget) buildAddEditLayout() {
	state := p.getState()

	giu.Layout{
		giu.Label("Key:"),
		giu.InputText("##"+p.id+"addEditKey", &state.key).OnChange(func() {
			str, found := p.dict[state.key]
			if found {
				state.value = str
			} else {
				state.value = ""
			}
		}),
		giu.Label("Value:"),
		giu.InputTextMultiline("##"+p.id+"addEditValue", &state.value),
		giu.Separator(),
		giu.Line(
			giu.Custom(func() {
				var btnStr string

				key := state.key

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
	}.Build()
}
