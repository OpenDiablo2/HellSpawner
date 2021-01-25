package hspaletteeditor

import (
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type PaletteEditor struct {
	*hseditor.Editor
	palette d2interface.Palette
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	palette, err := d2dat.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteEditor{
		Editor:  hseditor.New(pathEntry, x, y, project),
		palette: palette,
	}

	return result, nil
}

func (e *PaletteEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.PaletteGrid(e.GetId()+"_grid", e.palette.GetColors()),
	})
}

func (e *PaletteEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Palette Editor").Layout(g.Layout{
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

func (e *PaletteEditor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := e.Path.GetFileBytes()

	return data
}

func (e *PaletteEditor) Save() {
	e.Editor.Save(e)
}

func (e *PaletteEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
