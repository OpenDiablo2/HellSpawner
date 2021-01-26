// Pakcage hspalettemapeditor contains palette map editor's data
package hspalettemapeditor

import (
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

// PaletteMapEditor represents a palette map editor
type PaletteMapEditor struct {
	*hseditor.Editor
	pl2 *d2pl2.PL2
}

// Create creates a new palette map editor
func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	pl2, err := d2pl2.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteMapEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		pl2:    pl2,
	}

	result.Path = pathEntry

	return result, nil
}

// Build builds an editor
func (e *PaletteMapEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.PaletteMapViewer(e.Path.GetUniqueID(), e.pl2),
	})
}

// UpdateMainMenuLayout updates a main menu layout to it contains editors options
func (e *PaletteMapEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Map Editor").Layout(g.Layout{
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

// GenerateSaveData creates data to be saved
func (e *PaletteMapEditor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := e.Path.GetFileBytes()

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
