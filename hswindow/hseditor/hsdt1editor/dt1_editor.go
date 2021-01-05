package hsdt1editor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dt1, err := d2dt1.LoadDT1(*data)
	if err != nil {
		return nil, err
	}

	result := &DT1Editor{
		path:     pathEntry.Name,
		fullPath: pathEntry.FullPath,
		dt1:      dt1,
	}

	return result, nil
}

type DT1Editor struct {
	hseditor.Editor
	path     string
	fullPath string
	dt1      *d2dt1.DT1
}

func (e *DT1Editor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *DT1Editor) Cleanup() {
	e.Visible = false
}

func (e *DT1Editor) Render() {
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
		Layout(g.Layout{hswidget.DT1Viewer(e.fullPath, e.dt1)})
}
