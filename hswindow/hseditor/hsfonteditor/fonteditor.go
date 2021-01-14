package hsfonteditor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type FontEditor struct {
	hseditor.Editor
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	result := &FontEditor{}

	result.Path = pathEntry

	return result, nil
}

func (e *FontEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Pos(50, 50).Size(400, 300).Layout(g.Layout{
		g.Custom(func() {
			e.Focused = imgui.IsWindowFocused(0)
		}),
	})
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
