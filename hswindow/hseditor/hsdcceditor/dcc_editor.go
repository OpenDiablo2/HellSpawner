package hsdcceditor

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	dcc, err := d2dcc.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &DCCEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		dcc:    dcc,
	}

	return result, nil
}

type DCCEditor struct {
	*hseditor.Editor
	dcc *d2dcc.DCC
}

func (e *DCCEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.DCCViewer(e.Path.GetUniqueId(), e.dcc),
	})
}

func (e *DCCEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DCC Editor").Layout(g.Layout{
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

func (e *DCCEditor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := e.Path.GetFileBytes()

	return data
}

func (e *DCCEditor) Save() {
	e.Editor.Save(e)
}

func (e *DCCEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
