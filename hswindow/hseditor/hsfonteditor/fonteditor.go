package hsfonteditor

import (
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type FontEditor struct {
	*hseditor.Editor
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32) (hscommon.EditorWindow, error) {
	result := &FontEditor{
		Editor: hseditor.New(pathEntry, x, y),
	}

	return result, nil
}

func (e *FontEditor) Build() {
	e.IsOpen(&e.Visible).Size(400, 300).Layout(g.Layout{})
}

func (e *FontEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Visible = false
		}),
	})

	*l = append(*l, m)
}
