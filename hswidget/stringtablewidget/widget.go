package stringtablewidget

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"
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

	columns := []string{"key", "value"}
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
		)
	}

	giu.Layout{
		giu.Child("").Border(false).Layout(giu.Layout{
			giu.FastTable("").Border(true).Rows(rows),
		}),
	}.Build()
}
