// Package hsdc6editor represents a dc6 editor window
package hsdc6editor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/dc6widget"
	"github.com/OpenDiablo2/HellSpawner/hswidget/selectpalettewidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// static check, to ensure, if dc6 editor implemented editoWindow
var _ hscommon.EditorWindow = &DC6Editor{}

// DC6Editor represents a dc6 editor
type DC6Editor struct {
	*hseditor.Editor
	dc6           *d2dc6.DC6
	textureLoader hscommon.TextureLoader
	config        *hsconfig.Config
	selectPalette bool
	palette       *[256]d2interface.Color
	state         []byte
}

// Create creates a new dc6 editor
func Create(config *hsconfig.Config,
	textureLoader hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dc6, err := d2dc6.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading DC6 animation: %w", err)
	}

	result := &DC6Editor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		dc6:           dc6,
		textureLoader: textureLoader,
		selectPalette: false,
		config:        config,
		state:         state,
	}

	return result, nil
}

// Build builds a new dc6 editor
func (e *DC6Editor) Build() {
	id := e.Path.GetUniqueID()

	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if e.selectPalette {
		selectPaletteWidget := selectpalettewidget.NewSelectPaletteWidget(
			id+"selectPalette",
			e.Project,
			e.config,
		).IsOpen(&e.selectPalette).OnSelect(
			func(palette *[256]d2interface.Color) {
				e.palette = palette
			},
		)

		e.Layout(selectPaletteWidget)
	}

	e.Layout(
		dc6widget.Create(e.state, e.palette, e.textureLoader, id, e.dc6),
	)
}

// UpdateMainMenuLayout updates main menu to it contain DC6's editor menu
func (e *DC6Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DC6 Editor").Layout(g.Layout{
		g.MenuItem("Change Palette").OnClick(func() {
			e.selectPalette = true
		}),
		g.Separator(),
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
func (e *DC6Editor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates save data
func (e *DC6Editor) GenerateSaveData() []byte {
	data := e.dc6.Marshal()

	return data
}

// Save saves editor's data
func (e *DC6Editor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides editor
func (e *DC6Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
