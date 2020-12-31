package hsdc6editor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
)

func Create(path string, fullPath string, data []byte) (*DC6Editor, error) {
	dc6, err := d2dc6.Load(data)
	if err != nil {
		return nil, err
	}

	//numFrames := dc6.Directions * dc6.FramesPerDirection

	result := &DC6Editor{
		path:     path,
		fullPath: fullPath,
		dc6:      dc6,
		//decodedFrames: make([][]byte, numFrames),
		//textures:      make([]*g.Texture, numFrames),
	}

	return result, nil
}

type DC6Editor struct {
	hseditor.Editor
	path     string
	fullPath string
	dc6      *d2dc6.DC6
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

	g.WindowV(e.GetWindowTitle(), &e.Visible, g.WindowFlagsAlwaysAutoResize, 0, 0, 0, 0, g.Layout{
		hswidget.DC6Viewer(e.fullPath, e.dc6),
	})

}
