package hspalettemapeditor

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32) (hscommon.EditorWindow, error) {
	pl2, err := d2pl2.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteMapEditor{
		Editor: hseditor.New(pathEntry, x, y),
		pl2:    pl2,
	}

	result.Path = pathEntry

	return result, nil
}

type PaletteMapEditor struct {
	*hseditor.Editor
	pl2 *d2pl2.PL2
}

func (e *PaletteMapEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.PaletteMapViewer(e.Path.GetUniqueId(), e.pl2),
	})
}

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
