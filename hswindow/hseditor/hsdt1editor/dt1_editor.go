// Package hsdt1editor contains dt1 editor's data
package hsdt1editor

import (
	"fmt"
	"log"
	"path/filepath"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/dt1widget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
)

const (
	paletteSelectW, paletteSelectH = 400, 600
	actionButtonW, actionButtonH   = 200, 30
)

// static check, to ensure, if dt1 editor implemented editoWindow
var _ hscommon.EditorWindow = &DT1Editor{}

// DT1Editor represents a dt1 editor
type DT1Editor struct {
	*hseditor.Editor
	dt1           *d2dt1.DT1
	dt1Viewer     *hswidget.DT1ViewerWidget
	config        *hsconfig.Config
	selectPalette bool
	palette       *[256]d2interface.Color
	explorer      *hsmpqexplorer.MPQExplorer
	textureLoader *hscommon.TextureLoader
}

// Create creates new dt1 editor
func Create(config *hsconfig.Config,
	textureLoader *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
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
	}

	return result, nil
}

// Build prepares the editor for rendering, but does not actually render it
func (e *DT1Editor) Build() {
	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if !e.selectPalette {
		dt1Viewer := hswidget.DT1Viewer(e.palette, e.textureLoader, e.Path.GetUniqueID(), e.dt1)
		e.Layout(g.Layout{
			dt1Viewer,
		})

		return
	}

	// create mpq explorer if doesn't exist for now
	if e.explorer == nil {
		mpqExplorer, err := hsmpqexplorer.Create(
			func(path *hscommon.PathEntry) {
				bytes, bytesErr := path.GetFileBytes()
				if bytesErr != nil {
					log.Print(bytesErr)

					return
				}

				ft, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath), &bytes)
				if err != nil {
					log.Print(err)

					return
				}

				if ft == hsfiletypes.FileTypePalette {
					// load new palette:
					paletteData, err := path.GetFileBytes()
					if err != nil {
						log.Print(err)
					}

					palette, err := d2dat.Load(paletteData)
					if err != nil {
						log.Print(err)
					}

					colors := palette.GetColors()

					e.palette = &colors

					e.selectPalette = false
				}
			},
			e.config,
			0, 0,
		)

		mpqExplorer.SetProject(e.Project)

		if err != nil {
			log.Print(err)

			return
		}

		mpqExplorer.Visible = e.Visible

		e.explorer = mpqExplorer
	}

	e.Layout(g.Layout{
		g.PopupModal("something").IsOpen(&e.Visible).Layout(g.Layout{
			g.Child("somethingChild").Size(paletteSelectW, paletteSelectH).Layout(g.Layout{
				e.explorer.Layout(),
				g.Separator(),
				g.Button("Exit##"+e.Path.GetUniqueID()+"selectPaletteExit").
					Size(actionButtonW, actionButtonH).
					OnClick(func() {
						e.selectPalette = false
					}),
			}),
		}),
	})
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

// RegisterKeyboardShortcuts register a new keyboard shortcut
// nolint:wsl // I can't put this ccommented out code anywhere else
func (e *DT1Editor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)

	// nolint:gocritic // we may want to use this code
	/*
		// right arrow goes to the next tile group
		inputManager.RegisterShortcut(func() {
			e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() + 1)
		}, g.KeyRight, g.ModNone, false)

		// left arrow goes to the previous tile group
		inputManager.RegisterShortcut(func() {
			e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() - 1)
		}, g.KeyLeft, g.ModNone, false)
	*/
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
