// Package hscofeditor contains cof editor's data
package hscofeditor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"

	"github.com/OpenDiablo2/HellSpawner/hswidget/cofwidget"
)

// static check, to ensure, if cof editor implemented editoWindow
var _ hscommon.EditorWindow = &COFEditor{}

// COFEditor represents a cof editor
type COFEditor struct {
	*hseditor.Editor
	cof           *d2cof.COF
	textureLoader hscommon.TextureLoader
	state         []byte
}

// Create creates a new cof editor
func Create(_ *hsconfig.Config,
	tl hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	state []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	cof, err := d2cof.Unmarshal(*data)
	if err != nil {
		return nil, fmt.Errorf("error loading cof file: %w", err)
	}

	result := &COFEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		cof:           cof,
		textureLoader: tl,
		state:         state,
	}

	return result, nil
}

// Build builds a cof editor
func (e *COFEditor) Build() {
	uid := e.Path.GetUniqueID()
	cofWidget := cofwidget.Create(e.state, e.textureLoader, uid, e.cof)

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
