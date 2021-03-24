// Package hspalettemapeditor contains palette map editor's data
package hspalettemapeditor

import (
	"fmt"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettemapwidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// static check, to ensure, if palette map editor implemented editoWindow
var _ hscommon.EditorWindow = &PaletteMapEditor{}

// PaletteMapEditor represents a palette map editor
type PaletteMapEditor struct {
	*hseditor.Editor
	pl2           *d2pl2.PL2
	textureLoader *hscommon.TextureLoader
}

// Create creates a new palette map editor
func Create(_ *hsconfig.Config,
	textureLoader *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	pl2, err := d2pl2.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading PL2 file: %w", err)
	}

	result := &PaletteMapEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		pl2:           pl2,
		textureLoader: textureLoader,
	}

	result.Path = pathEntry

	return result, nil
}

// Build builds an editor
func (e *PaletteMapEditor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(g.Layout{
			palettemapwidget.Create(e.textureLoader, e.Path.GetUniqueID(), e.pl2),
		})
}

// UpdateMainMenuLayout updates a main menu layout to it contains editors options
func (e *PaletteMapEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Map Editor").Layout(g.Layout{
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
func (e *PaletteMapEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData creates data to be saved
func (e *PaletteMapEditor) GenerateSaveData() []byte {
	data := e.pl2.Marshal()

	return data
}

// Save saves an editor
func (e *PaletteMapEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *PaletteMapEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
