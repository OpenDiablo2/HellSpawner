package hsdc6editor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dc6, err := d2dc6.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &DC6Editor{
		dc6: dc6,
	}

	result.Path = pathEntry

	return result, nil
}

type DC6Editor struct {
	hseditor.Editor
	dc6 *d2dc6.DC6
}

func (e *DC6Editor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.DC6Viewer(e.Path.GetUniqueId(), e.dc6),
	})

}
