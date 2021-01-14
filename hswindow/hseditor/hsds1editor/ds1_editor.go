package hsds1editor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

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
		Layout(g.Layout{hswidget.DS1Viewer(e.Path.GetUniqueId(), e.ds1)})
}
