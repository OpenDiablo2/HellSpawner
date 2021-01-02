package hspaletteeditor

import (
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

func Create(path string, fullPath string, data []byte) (*PaletteEditor, error) {
	palette, err := d2dat.Load(data)
	if err != nil {
		return nil, err
	}

	result := &PaletteEditor{
		path:     path,
		fullPath: fullPath,
		palette:  palette,
	}

	return result, nil
}

type PaletteEditor struct {
	hseditor.Editor
	palette  d2interface.Palette
	path     string
	fullPath string
}

func (e *PaletteEditor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *PaletteEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Flags(g.WindowFlagsAlwaysAutoResize).Pos(360, 30).Layout(g.Layout{
		hswidget.PaletteGrid(e.fullPath, e.palette.GetColors()),
	})
}

func (e *PaletteEditor) Cleanup() {
	e.Visible = false
}
