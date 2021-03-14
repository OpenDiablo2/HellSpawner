// Package hsdcceditor contains dcc editor's data
package hsdcceditor

import (
	"fmt"
	"log"
	"path/filepath"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/dccwidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
)

const (
	paletteSelectW, paletteSelectH = 400, 600
	actionButtonW, actionButtonH   = 200, 30
)

// static check, to ensure, if dc6 editor implemented editoWindow
var _ hscommon.EditorWindow = &DCCEditor{}

// DCCEditor represents a new dcc editor
type DCCEditor struct {
	*hseditor.Editor
	dcc           *d2dcc.DCC
	config        *hsconfig.Config
	selectPalette bool
	palette       *[256]d2interface.Color
	explorer      *hsmpqexplorer.MPQExplorer
}

// Create creates a new dcc editor
func Create(config *hsconfig.Config,
	_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dcc, err := d2dcc.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading dcc animation: %w", err)
	}

	result := &DCCEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		dcc:           dcc,
		config:        config,
		selectPalette: false,
	}

	return result, nil
}

// Build builds a dcc editor
func (e *DCCEditor) Build() {
	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if !e.selectPalette {
		e.Layout(g.Layout{
			hswidget.DCCViewer(e.palette, e.Path.GetUniqueID(), e.dcc),
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

// UpdateMainMenuLayout updates main menu to it contain editor's options
func (e *DCCEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DCC Editor").Layout(g.Layout{
		g.MenuItem("Change Palette").OnClick(func() {
			e.selectPalette = true
		}),
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
func (e *DCCEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates data to save
func (e *DCCEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves editor
func (e *DCCEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides editor
func (e *DCCEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
