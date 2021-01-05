package hsdcceditor

import (
	"path/filepath"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dcc, err := d2dcc.Load(*data)
	if err != nil {
		return nil, err
	}

	//numFrames := dcc.Directions * dcc.FramesPerDirection

	result := &DCCEditor{
		path: filepath.Base(pathEntry.FullPath),
		dcc:  dcc,
		//decodedFrames: make([][]byte, numFrames),
		//textures:      make([]*g.Texture, numFrames),
	}

	return result, nil
}

type DCCEditor struct {
	hseditor.Editor
	path     string
	fullPath string
	dcc      *d2dcc.DCC
}

func (e *DCCEditor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *DCCEditor) Cleanup() {
	e.Visible = false
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
		hswidget.DCCViewer(e.fullPath, e.dcc),
	})

}
