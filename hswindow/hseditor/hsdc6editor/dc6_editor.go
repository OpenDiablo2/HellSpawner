package hsdc6editor

import (
	"path/filepath"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	dc6, err := d2dc6.Load(*data)
	if err != nil {
		return nil, err
	}

	//numFrames := dc6.Directions * dc6.FramesPerDirection

	result := &DC6Editor{
		path: filepath.Base(pathEntry.FullPath),
		dc6:  dc6,
		//decodedFrames: make([][]byte, numFrames),
		//textures:      make([]*g.Texture, numFrames),
	}

	return result, nil
}

type DC6Editor struct {
	hseditor.Editor
	path string
	dc6  *d2dc6.DC6
}

func (e *DC6Editor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *DC6Editor) Cleanup() {
	e.Visible = false
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
		hswidget.DC6Viewer(e.path, e.dc6),
	})

}
