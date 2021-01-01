package hscofeditor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
)

func Create(path string, fullPath string, data []byte) (*COFEditor, error) {
	cof, err := d2cof.Load(data)
	if err != nil {
		return nil, err
	}

	result := &COFEditor{
		path:     path,
		fullPath: fullPath,
		cof:      cof,
	}

	return result, nil
}

type COFEditor struct {
	hseditor.Editor
	path     string
	fullPath string
	cof      *d2cof.COF
}

func (e *COFEditor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *COFEditor) Cleanup() {
	e.Visible = false
}

func (e *COFEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.WindowV(e.GetWindowTitle(), &e.Visible, g.WindowFlagsAlwaysAutoResize, 0, 0, 0, 0, g.Layout{
		hswidget.COFViewer(e.fullPath, e.cof),
	})
}
