package hscofeditor

import (
	"path/filepath"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
)

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	cof, err := d2cof.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &COFEditor{
		path: filepath.Base(pathEntry.FullPath),
		cof:  cof,
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

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Layout(g.Layout{
		hswidget.COFViewer(e.fullPath, e.cof),
	})
}
