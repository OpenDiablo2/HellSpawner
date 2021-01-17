package hspaletteeditor

import (
	"github.com/OpenDiablo2/HellSpawner/hscommon"
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

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32) (hscommon.EditorWindow, error) {
	palette, err := d2dat.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteEditor{
		Editor:  hseditor.New(pathEntry, x, y),
		palette: palette,
	}

	return result, nil
}

func (e *PaletteEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Pos(360, 30).Layout(g.Layout{
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
			e.Visible = false
		}),
	})

	*l = append(*l, m)
}
