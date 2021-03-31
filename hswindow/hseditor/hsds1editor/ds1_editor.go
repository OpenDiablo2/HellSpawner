// Package hsds1editor contains ds1 editor's data
package hsds1editor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/ds1widget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// static check if DS1Editor implemented hscommon.EditorWindow
var _ hscommon.EditorWindow = &DS1Editor{}

// DS1Editor represents ds1 editor
type DS1Editor struct {
	*hseditor.Editor
	ds1                 *d2ds1.DS1
	deleteButtonTexture *g.Texture
}

// Create creates a new ds1 editor
func Create(_ *hsconfig.Config,
	tl *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	ds1, err := d2ds1.Unmarshal(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading DS1 file")
	}

	result := &DS1Editor{
		Editor: hseditor.New(pathEntry, x, y, project),
		ds1:    ds1,
	}

	result.Path = pathEntry

	tl.CreateTextureFromFile(hsassets.DeleteIcon, func(texture *g.Texture) {
		result.deleteButtonTexture = texture
	})

	return result, nil
}

// Build builds an editor
func (e *DS1Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(g.Layout{
			ds1widget.Create(e.Path.GetUniqueID(), e.ds1, e.deleteButtonTexture),
		})
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *DS1Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DS1 Editor").Layout(g.Layout{
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

// RegisterKeyboardShortcuts adds a local shortcuts for this editor
func (e *DS1Editor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates data to be saved
func (e *DS1Editor) GenerateSaveData() []byte {
	data := e.ds1.Marshal()

	return data
}

// Save saves editors data
func (e *DS1Editor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides editor
func (e *DS1Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
