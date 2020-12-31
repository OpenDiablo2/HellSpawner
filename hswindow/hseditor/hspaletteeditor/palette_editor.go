package hspaletteeditor

import (
	"github.com/OpenDiablo2/HellSpawner/hswidget"

	"github.com/OpenDiablo2/giu/imgui"

	g "github.com/OpenDiablo2/giu"

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

	g.WindowV(
		e.GetWindowTitle(),
		&e.Visible,
		g.WindowFlagsAlwaysAutoResize,
		0, 0,
		0, 0,
		//float32(width+16), float32(height+40),
		g.Layout{
			hswidget.PaletteGrid(e.fullPath, e.palette.GetColors()),
		},
	)
}

func (e *PaletteEditor) Cleanup() {
	e.Visible = false
}
