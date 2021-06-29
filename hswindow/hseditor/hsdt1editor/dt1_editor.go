// Package hsdt1editor contains dt1 editor's data
package hsdt1editor

import (
	"fmt"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswidget/dt1widget"
	"github.com/OpenDiablo2/HellSpawner/hswidget/selectpalettewidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// static check, to ensure, if dt1 editor implemented editoWindow
var _ hscommon.EditorWindow = &DT1Editor{}

// DT1Editor represents a dt1 editor
type DT1Editor struct {
	*hseditor.Editor
	dt1                 *d2dt1.DT1
	textureLoader       hscommon.TextureLoader
	config              *hsconfig.Config
	selectPalette       bool
	palette             *[256]d2interface.Color
	selectPaletteWidget g.Widget
	state               []byte
}

// Create creates new dt1 editor
func Create(config *hsconfig.Config,
	textureLoader hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dt1, err := d2dt1.LoadDT1(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dt1 file: %w", err)
	}

	result := &DT1Editor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		dt1:           dt1,
		config:        config,
		selectPalette: false,
		textureLoader: textureLoader,
		state:         state,
	}

	return result, nil
}

// Build prepares the editor for rendering, but does not actually render it
func (e *DT1Editor) Build() {
	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if !e.selectPalette {
		dt1Viewer := dt1widget.Create(e.state, e.palette, e.textureLoader, e.Path.GetUniqueID(), e.dt1)
		e.Layout(g.Layout{
			dt1Viewer,
		})

		return
	}

	// create mpq explorer if doesn't exist for now
	if e.selectPaletteWidget == nil {
		e.selectPaletteWidget = selectpalettewidget.NewSelectPaletteWidget(
			e.Path.GetUniqueID(),
			e.Project,
			e.config,
			func(colors *[256]d2interface.Color) {
				e.palette = colors
			},
			func() {
				e.selectPalette = false
			},
		)
	}

	e.Layout(g.Layout{e.selectPaletteWidget})
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *DT1Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DT1 Editor").Layout(g.Layout{
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

// KeyboardShortcuts register a new keyboard shortcut
func (e *DT1Editor) KeyboardShortcuts() []g.WindowShortcut {
	// nolint:gocritic // we may want to use this code
	return []g.WindowShortcut{
		/*
			// right arrow goes to the next tile group
			giu.WindowShortcut{
				Callback: func() {
					e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() + 1)
				},
				Key:      g.KeyRight,
				Modifier: g.ModNone,
			},

			// left arrow goes to the previous tile group
			giu.WindowShortcut{
				Callback: func() {
					e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() - 1)
				},
				Key:      g.KeyLeft,
				Modifier: g.ModNone,
			},
		*/
	}
}

// GenerateSaveData generates data to be saved
func (e *DT1Editor) GenerateSaveData() []byte {
	data := e.dt1.Marshal()

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
