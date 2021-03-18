// Package hsanimdataeditor contains D2 editor's data
package hsanimdataeditor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget/animdatawidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	delItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_delete.png"
)

// static check, to ensure, if D2 editor implemented editoWindow
var _ hscommon.EditorWindow = &AnimationDataEditor{}

// AnimationDataEditor represents a cof editor
type AnimationDataEditor struct {
	*hseditor.Editor
	d2  *d2animdata.AnimationData
	del *g.Texture
}

// Create creates a new cof editor
func Create(_ *hsconfig.Config,
	tl *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	d2, err := d2animdata.Load(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading animation data file: %w", err)
	}

	result := &AnimationDataEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		d2:     d2,
	}

	tl.CreateTextureFromFileAsync(delItemButtonPath, func(texture *g.Texture) {
		result.del = texture
	})

	return result, nil
}

// Build builds a D2 editor
func (e *AnimationDataEditor) Build() {
	uid := e.Path.GetUniqueID()
	animDataWidget := animdatawidget.Create(e.del, uid, e.d2)

	e.IsOpen(&e.Visible)
	e.Flags(g.WindowFlagsAlwaysAutoResize)
	e.Layout(g.Layout{animDataWidget})
}

// UpdateMainMenuLayout updates a main menu layout, to it contains anim data viewer's settings
func (e *AnimationDataEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Animation Data Editor").Layout(g.Layout{
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
func (e *AnimationDataEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
}

// GenerateSaveData generates data to be saved
func (e *AnimationDataEditor) GenerateSaveData() []byte {
	data := e.d2.Marshal()

	return data
}

// Save saves an editor
func (e *AnimationDataEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *AnimationDataEditor) Cleanup() {
	const strPrompt = "There are unsaved changes to %s, save before closing this editor?"

	if e.HasChanges(e) {
		if shouldSave := dialog.Message(strPrompt, e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
