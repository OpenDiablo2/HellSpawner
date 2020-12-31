package hspaletteeditor

import (
	"fmt"
	"image"
	"image/color"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

const (
	fmtTitle = "Palette Editor [%s]"
)

const (
	gridWidth  = 16
	gridHeight = 16
	cellSize   = 12
)

func Create(path string, data []byte) (*PaletteEditor, error) {
	palette, err := d2dat.Load(data)
	if err != nil {
		return nil, err
	}

	result := &PaletteEditor{
		path:    path,
		palette: palette,
	}

	return result, nil
}

type PaletteEditor struct {
	hseditor.Editor
	palette d2interface.Palette
	path    string
}

func (e *PaletteEditor) GetWindowTitle() string {
	return fmt.Sprintf(fmtTitle, e.path)
}

func (e *PaletteEditor) Render() {
	if !e.Visible {
		return
	}

	width := gridWidth * cellSize
	height := gridHeight * cellSize

	displayPalette := func() {
		canvas := g.GetCanvas()
		pos := g.GetCursorScreenPos()
		colors := e.palette.GetColors()

		for idx, c := range colors {
			x := (idx % gridWidth) * cellSize
			y := (idx / gridHeight) * cellSize

			tl, br := pos.Add(image.Pt(x, y)), pos.Add(image.Pt(x+cellSize, y+cellSize))

			theColor := color.RGBA{c.B(), c.G(), c.R(), c.A()}

			canvas.AddRectFilled(tl, br, theColor, 0, 0)
		}
	}

	g.WindowV(
		e.GetWindowTitle(),
		&e.Visible,
		g.WindowFlagsNoResize,
		0, 0,
		float32(width+16), float32(height+40),
		g.Layout{
			g.Custom(displayPalette),
		},
	)
}
