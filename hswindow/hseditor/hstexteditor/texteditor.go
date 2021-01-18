package hstexteditor

import (
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type TextEditor struct {
	*hseditor.Editor

	text      string
	tableView bool
	tableRows g.Rows
	columns   int
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32) (hscommon.EditorWindow, error) {
	result := &TextEditor{
		Editor: hseditor.New(pathEntry, x, y),
		text:   string(*data),
	}

	result.Path = pathEntry

	lines := strings.Split(result.text, "\n")
	firstLine := lines[0]
	result.tableView = strings.Count(firstLine, "\t") > 0

	if !result.tableView {
		return result, nil
	}

	result.tableRows = make([]*g.RowWidget, len(lines))

	columns := strings.Split(firstLine, "\t")
	result.columns = len(columns)
	columnWidgets := make([]g.Widget, len(columns))
	for idx := range columns {
		columnWidgets[idx] = g.Label(columns[idx])
	}
	result.tableRows[0] = g.Row(columnWidgets...)

	for lineIdx := range lines[1:] {
		columns := strings.Split(lines[lineIdx+1], "\t")
		columnWidgets := make([]g.Widget, len(columns))
		for idx := range columns {
			columnWidgets[idx] = g.Label(columns[idx])
		}
		result.tableRows[lineIdx+1] = g.Row(columnWidgets...)
	}

	return result, nil
}

func (e *TextEditor) Build() {
	if !e.tableView {
		e.IsOpen(&e.Visible).Size(400, 300).Layout(g.Layout{
			g.InputTextMultiline("", &e.text).Size(-1, -1).Flags(g.InputTextFlagsAllowTabInput),
		})
	} else {
		e.IsOpen(&e.Visible).Flags(g.WindowFlagsHorizontalScrollbar).Pos(50, 50).Size(400, 300).Layout(g.Layout{
			g.Child("").Border(false).Size(float32(e.columns*80), 0).Layout(g.Layout{
				g.FastTable("").Border(true).Rows(e.tableRows),
			}),
		})
	}
}

func (e *TextEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Text Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Visible = false
		}),
	})

	*l = append(*l, m)
}
