package hstexteditor

import (
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type TextEditor struct {
	hseditor.Editor

	fontFixed imgui.Font
	file      string
	text      string
	tableView bool
	tableRows g.Rows
	columns   int
}

func (e *TextEditor) Cleanup() {

}

func (e *TextEditor) GetWindowTitle() string {
	return e.file + "##" + e.GetId()
}

func Create(file, text string, fontFixed imgui.Font) (*TextEditor, error) {
	result := &TextEditor{
		file:      file,
		text:      text,
		fontFixed: fontFixed,
	}

	lines := strings.Split(text, "\n")
	firstLine := lines[0]
	result.tableView = strings.Count(firstLine, "\t") > 2

	if !result.tableView {
		return result, nil
	}

	result.tableRows = make([]*g.RowWidget, len(lines))

	columns := strings.Split(firstLine, "\t")
	result.columns = len(columns)
	columnWidgets := make([]g.Widget, len(columns))
	for idx := range columns {
		columnWidgets[idx] = g.LabelV(columns[idx], false, nil, &result.fontFixed)
	}
	result.tableRows[0] = g.Row(columnWidgets...)

	for lineIdx := range lines[1:] {
		columns := strings.Split(lines[lineIdx+1], "\t")
		columnWidgets := make([]g.Widget, len(columns))
		for idx := range columns {
			columnWidgets[idx] = g.LabelV(columns[idx], false, nil, &result.fontFixed)
		}
		result.tableRows[lineIdx+1] = g.Row(columnWidgets...)
	}

	return result, nil
}

func (e *TextEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	if !e.tableView {
		g.WindowV(e.GetWindowTitle(), &e.Visible, g.WindowFlagsNone, 50, 50, 400, 300, g.Layout{
			g.InputTextMultiline("", &e.text, -1, -1, g.InputTextFlagsAllowTabInput, nil, func() {
				// On Change Event
			}),
		})

		return
	}

	g.WindowV(e.GetWindowTitle(), &e.Visible, g.WindowFlagsHorizontalScrollbar, 0, 0, 400, 300, g.Layout{
		g.Child("", false, float32(e.columns*80), 0,
			g.WindowFlagsNone, g.Layout{
				g.FastTable("", true, e.tableRows),
			}),
	})
}
