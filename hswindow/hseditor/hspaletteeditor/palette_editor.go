// Package hspaletteeditor contains palette editor's data
package hspaletteeditor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegrideditorwidget"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// static check, to ensure, if palette editor implemented editoWindow
var _ hscommon.EditorWindow = &PaletteEditor{}

// PaletteEditor represents a palette editor
type PaletteEditor struct {
	*hseditor.Editor
	palette       d2interface.Palette
	textureLoader *hscommon.TextureLoader
}

// Create creates a new palette editor
func Create(_ *hsconfig.Config,
	tl *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	palette, err := d2dat.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dat palette: %w", err)
	}

	result := &PaletteEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		palette:       palette,
		textureLoader: tl,
	}

	return result, nil
}

// Build builds a palette editor
func (e *PaletteEditor) Build() {
	var col [256]palettegridwidget.PaletteColor
	for n, i := range e.palette.GetColors() {
		col[n] = palettegridwidget.PaletteColor(i)
	}

	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		palettegrideditorwidget.Create(e.textureLoader, e.GetID()+"_grid", &col),
	})
}

// UpdateMainMenuLayout updates a main menu layout to it contain palette editor's options
func (e *PaletteEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Editor").Layout(g.Layout{
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
func (e *PaletteEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates data to be saved
func (e *PaletteEditor) GenerateSaveData() []byte {
	palette, ok := e.palette.(*d2dat.DATPalette)
	if ok {
		data := palette.Marshal()

		return data
	}

	return make([]byte, 0)
}

// Save saves editor
func (e *PaletteEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides palette editor
func (e *PaletteEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
