package hsdt1editor

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dt1, err := d2dt1.LoadDT1(*data)
	if err != nil {
		return nil, err
	}

	result := &DT1Editor{
		Editor:    hseditor.New(pathEntry, x, y, project),
		dt1:       dt1,
		dt1Viewer: hswidget.DT1Viewer(pathEntry.GetUniqueID(), dt1),
	}

	return result, nil
}

type DT1Editor struct {
	*hseditor.Editor
	dt1       *d2dt1.DT1
	dt1Viewer *hswidget.DT1ViewerWidget
}

// Build prepares the editor for rendering, but does not actually render it
func (e *DT1Editor) Build() {
	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Layout(g.Layout{
			e.dt1Viewer,
		})
}

func (e *DT1Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DT1 Editor").Layout(g.Layout{
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

func (e *DT1Editor) RegisterKeyboardShortcuts() {
	// right arrow goes to the next tile group
	hsinput.RegisterShortcut(func() {
		e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() + 1)
	}, g.KeyRight, g.ModNone, false)
	// left arrow goes to the previous tile group
	hsinput.RegisterShortcut(func() {
		e.dt1Viewer.SetTileGroup(e.dt1Viewer.TileGroup() - 1)
	}, g.KeyLeft, g.ModNone, false)
}

func (e *DT1Editor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := e.Path.GetFileBytes()

	return data
}

func (e *DT1Editor) Save() {
	e.Editor.Save(e)
}

func (e *DT1Editor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
