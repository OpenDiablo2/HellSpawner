// Package hsdt1editor contains dt1 editor's data
package hsdt1editor

import (
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// DT1Editor represents a dt1 editor
type DT1Editor struct {
	*hseditor.Editor
	dt1       *d2dt1.DT1
	dt1Viewer *hswidget.DT1ViewerWidget
}

// Create creates new dt1 editor
func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dt1, err := d2dt1.LoadDT1(*data)
	if err != nil {
		return nil, err
	}

	result := &DT1Editor{
		Editor:    hseditor.New(pathEntry, x, y, project),
		dt1:       dt1,
		dt1Viewer: hswidget.DT1Viewer(pathEntry.GetUniqueID(), dt1),
	}

	return result, nil
}

// Build prepares the editor for rendering, but does not actually render it
func (e *DT1Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(g.Layout{
			e.dt1Viewer,
		})
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *DT1Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DT1 Editor").Layout(g.Layout{
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

// RegisterKeyboardShortcuts register a new keyboard shortcut
func (e *DT1Editor) RegisterKeyboardShortcuts() {
	// right arrow goes to the next tile group
	// TODO --- add local shocrcuts
	// nolint:gocritic // todo
	/*e.inputManager.RegisterShortcut(func() {
		e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() + 1)
	}, g.KeyRight, g.ModNone, false)*/

	// left arrow goes to the previous tile group
	// nolint:gocritic // todo
	/*e.inputManager.RegisterShortcut(func() {
		e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() - 1)
	}, g.KeyLeft, g.ModNone, false)*/
}

// GenerateSaveData generates data to be saved
func (e *DT1Editor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves editor
func (e *DT1Editor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides editor
func (e *DT1Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
