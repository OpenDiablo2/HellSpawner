// Package hsdc6editor represents a dc6 editor window
package hsdc6editor

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hswidget/dc6widget"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
)

const (
	paletteSelectW, paletteSelectH = 400, 600
	actionButtonW, actionButtonH   = 200, 30
)

// static check, to ensure, if dc6 editor implemented editoWindow
var _ hscommon.EditorWindow = &DC6Editor{}

// DC6Editor represents a dc6 editor
type DC6Editor struct {
	*hseditor.Editor
	dc6           *d2dc6.DC6
	textureLoader *hscommon.TextureLoader
	config        *hsconfig.Config
	selectPalette bool
	palette       *[256]d2interface.Color
	explorer      *hsmpqexplorer.MPQExplorer
}

// Create creates a new dc6 editor
func Create(config *hsconfig.Config,
	textureLoader *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
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
	}

	return result, nil
}

// Build builds a new dc6 editor
func (e *DC6Editor) Build() {
	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)

	if !e.selectPalette {
		e.Layout(g.Layout{
			dc6widget.Create(e.palette, e.textureLoader, e.Path.GetUniqueID(), e.dc6),
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
