package hsutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var imgSquare *ebiten.Image

func DrawColoredRect(target *ebiten.Image, x, y, w, h int, c color.Color) {
	r, g, b, alpha := c.RGBA()
	max := float64(1 << 16)

	if imgSquare == nil {
		imgSquare = ebiten.NewImage(1, 1)
	}

	drawOptions := &ebiten.DrawImageOptions{}

	drawOptions.GeoM.Translate(float64(x)*(1/float64(w)), float64(y)*(1/float64(h)))
	drawOptions.GeoM.Scale(float64(w), float64(h))
	drawOptions.ColorM.Translate(float64(r)/max, float64(g)/max, float64(b)/max, float64(alpha)/max)

	target.DrawImage(imgSquare, drawOptions)
}

func ArrayToRGBA(tc []int) color.Color {
	return color.RGBA{R: uint8(tc[0]), G: uint8(tc[1]), B: uint8(tc[2]), A: uint8(tc[3])}
}
