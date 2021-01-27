// Package hsstringtableeditor contains string tables editor's data
package hsstringtableeditor

import (
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 400, 300
)

// StringTableEditor represents a string table editor
type StringTableEditor struct {
	*hseditor.Editor
	// nolint:unused,structcheck // will be used
	header g.RowWidget
	rows   g.Rows
	dict   d2tbl.TextDictionary
}

// Create creates a new string table editor
func Create(_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dict, err := d2tbl.LoadTextDictionary(*data)
	if err != nil {
		return nil, err
	}

	result := &StringTableEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		dict:   dict,
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

// Build builds an editor
func (e *StringTableEditor) Build() {
	l := g.Layout{
		g.Child("").Border(false).Layout(g.Layout{
			g.FastTable("").Border(true).Rows(e.rows),
		}),
	}

	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Size(mainWindowW, mainWindowH).
		Layout(l)
}

// UpdateMainMenuLayout updates main menu layout to it contain editors options
func (e *StringTableEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("String Table Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Cleanup()
		}),
	})

	*l = append(*l, m)
}

// GenerateSaveData generates data to be saved
func (e *StringTableEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (e *StringTableEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *StringTableEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
