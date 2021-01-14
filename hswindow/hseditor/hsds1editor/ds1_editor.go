package hsds1editor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

var _ hscommon.EditorWindow = &DS1Editor{}

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	ds1, err := d2ds1.LoadDS1(*data)
	if err != nil {
		return nil, err
	}

	result := &DS1Editor{
		ds1: ds1,
	}

	result.Path = pathEntry

	return result, nil
}

type DS1Editor struct {
	hseditor.Editor
	ds1 *d2ds1.DS1
}

func (e *DS1Editor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).
		IsOpen(&e.Visible).
		Flags(g.WindowFlagsAlwaysAutoResize).
		Pos(360, 30).
		Layout(g.Layout{
			hswidget.DS1Viewer(e.Path.GetUniqueId(), e.ds1),
			g.Custom(func() {
				e.Focused = imgui.IsWindowFocused(0)
			}),
		})
}

func (e *DS1Editor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("DS1 Editor").Layout(g.Layout{
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
