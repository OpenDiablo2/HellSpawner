package hspalettemapeditor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	pl2, err := d2pl2.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &PaletteMapEditor{
		pl2: pl2,
	}

	result.Path = pathEntry

	return result, nil
}

type PaletteMapEditor struct {
	hseditor.Editor
	pl2 *d2pl2.PL2
}

func (e *PaletteMapEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.PaletteMapViewer(e.Path.GetUniqueId(), e.pl2),
	})
}
