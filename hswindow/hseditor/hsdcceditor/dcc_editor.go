package hsdcceditor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dcc, err := d2dcc.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &DCCEditor{
		dcc: dcc,
	}

	result.Path = pathEntry

	return result, nil
}

type DCCEditor struct {
	hseditor.Editor
	dcc *d2dcc.DCC
}

func (e *DCCEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.DCCViewer(e.Path.GetUniqueId(), e.dcc),
	})

}
