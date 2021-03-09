package palettegridwidget

import (
	image2 "image"
	"image/color"

	"github.com/ianling/giu"
)

func (p *widget) loadTexture(i int) {
	state := p.getState()

	rgb := image2.NewRGBA(image2.Rect(0, 0, cellSize, cellSize))

	for y := 0; y < cellSize; y++ {
		for x := 0; x < cellSize; x++ {
			col := p.colors[i]

			// nolint:gomnd // opacity
			rgb.Set(x, y, color.RGBA{R: col.R(), G: col.G(), B: col.B(), A: 255})
		}
	}

	go p.textureLoader.CreateTextureFromARGB(rgb, func(texture *giu.Texture) {
		state.texture[i] = texture
		giu.Context.SetState(p.getStateID(), state)
	})
}

func (p *widget) reloadTextures() {
	for x := 0; x < 256; x++ {
		p.loadTexture(x)
	}
}
