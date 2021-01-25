package hscofeditor

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	cof, err := d2cof.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &COFEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		cof:    cof,
	}

	return result, nil
}

type COFEditor struct {
	*hseditor.Editor
	cof *d2cof.COF
}

func (e *COFEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.COFViewer(e.Path.GetUniqueId(), e.cof),
	})
}

func (e *COFEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("COF Editor").Layout(g.Layout{
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

func (e *COFEditor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := e.Path.GetFileBytes()

	return data
}

func (e *COFEditor) Save() {
	e.Editor.Save(e)
}

func (e *COFEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
