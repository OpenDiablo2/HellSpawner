// Package hsstringtableeditor contains string tables editor's data
package hsstringtableeditor

import (
	"fmt"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswidget/stringtablewidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 600, 500
)

// static check, to ensure, if string table editor implemented editoWindow
var _ hscommon.EditorWindow = &StringTableEditor{}

// StringTableEditor represents a string table editor
type StringTableEditor struct {
	*hseditor.Editor
	dict  d2tbl.TextDictionary
	state []byte
}

// Create creates a new string table editor
func Create(_ *hsconfig.Config,
	_ hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dict, err := d2tbl.LoadTextDictionary(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading string table: %w", err)
	}

	result := &StringTableEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		dict:   dict,
		state:  state,
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	result.Path = pathEntry

	return result, nil
}

// Build builds an editor
func (e *StringTableEditor) Build() {
	l := stringtablewidget.Create(e.state, e.Path.GetUniqueID(), e.dict)

	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(g.Layout{l})
}

// UpdateMainMenuLayout updates main menu layout to it contain editors options
func (e *StringTableEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("String Table Editor").Layout(g.Layout{
		g.MenuItem("Save\t\t\t\tCtrl+Shift+S").OnClick(e.Save),
		g.Separator(),
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
	data := e.dict.Marshal()

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
