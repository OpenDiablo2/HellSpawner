// Package hscofeditor contains cof editor's data
package hscofeditor

import (
	"fmt"

	"github.com/OpenDiablo2/HellSpawner/hswidget/cofwidget"
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	upItemButtonPath     = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png"
	downItemButtonPath   = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png"
	leftArrowButtonPath  = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_left.png"
	rightArrowButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_right.png"
)

// static check, to ensure, if cof editor implemented editoWindow
var _ hscommon.EditorWindow = &COFEditor{}

// COFEditor represents a cof editor
type COFEditor struct {
	*hseditor.Editor
	cof           *d2cof.COF
	textureLoader *hscommon.TextureLoader
	textures      struct {
		up    *g.Texture
		down  *g.Texture
		right *g.Texture
		left  *g.Texture
	}
}

// Create creates a new cof editor
func Create(config *hsconfig.Config, tl *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	cof, err := d2cof.Unmarshal(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading cof file: %w", err)
	}

	result := &COFEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		cof:           cof,
		textureLoader: tl,
	}

	tl.CreateTextureFromFileAsync(upItemButtonPath, func(texture *g.Texture) {
		result.textures.up = texture
	})

	tl.CreateTextureFromFileAsync(downItemButtonPath, func(texture *g.Texture) {
		result.textures.down = texture
	})

	tl.CreateTextureFromFileAsync(leftArrowButtonPath, func(texture *g.Texture) {
		result.textures.left = texture
	})

	tl.CreateTextureFromFileAsync(rightArrowButtonPath, func(texture *g.Texture) {
		result.textures.right = texture
	})

	return result, nil
}

// Build builds a cof editor
func (e *COFEditor) Build() {
	uid := e.Path.GetUniqueID()
	cofWidget := cofwidget.Create(e.textures.up, e.textures.down, e.textures.right, e.textures.left, uid, e.cof)

	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)
	e.Layout(g.Layout{cofWidget})
}

// UpdateMainMenuLayout updates a main menu layout, to it contains COFViewer's settings
func (e *COFEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("COF Editor").Layout(g.Layout{
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
func (e *COFEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates data to be saved
func (e *COFEditor) GenerateSaveData() []byte {
	data := e.cof.Marshal()

	return data
}

// Save saves an editor
func (e *COFEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *COFEditor) Cleanup() {
	const strPrompt = "There are unsaved changes to %s, save before closing this editor?"

	if e.HasChanges(e) {
		if shouldSave := dialog.Message(strPrompt, e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
