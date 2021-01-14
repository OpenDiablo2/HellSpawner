package hsstringtableeditor

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

//nolint:structcheck,unused // will be used
type StringTableEditor struct {
	hseditor.Editor
	header g.RowWidget
	rows   g.Rows
	dict   d2tbl.TextDictionary
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dict, err := d2tbl.LoadTextDictionary(*data)
	if err != nil {
		return nil, err
	}

	result := &StringTableEditor{
		dict: dict,
	}

	result.Path = pathEntry

	numEntries := len(result.dict)

	if !(numEntries > 0) {
		return result, nil
	}

	result.rows = make([]*g.RowWidget, numEntries+1)

	columns := []string{"key", "value"}
	columnWidgets := make([]g.Widget, len(columns))

	for idx := range columns {
		columnWidgets[idx] = g.Label(columns[idx])
	}

	result.rows[0] = g.Row(columnWidgets...)

	keyIdx := 0
	for key := range result.dict {
		result.rows[keyIdx+1] = g.Row(
			g.Label(key),
			g.Label(result.dict[key]),
		)

		keyIdx++
	}

	return result, nil
}

func (e *StringTableEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	l := g.Layout{
		g.Child("").Border(false).Layout(g.Layout{
			g.FastTable("").Border(true).Rows(e.rows),
		}),
		g.Custom(func() {
			e.Focused = imgui.IsWindowFocused(0)
		}),
	}

	g.Window(e.GetWindowTitle()).
		IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Pos(50, 50).
		Size(400, 300).
		Layout(l)
}

func (e *StringTableEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("String Table Editor").Layout(g.Layout{
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
