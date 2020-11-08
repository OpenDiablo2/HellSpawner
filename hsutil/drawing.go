package hsutil

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var imgSquare *ebiten.Image

func DrawColoredRect(target *ebiten.Image, x, y, w, h int, r, g, b, alpha uint8) {
	if imgSquare == nil {
		imgSquare = ebiten.NewImage(1, 1)
	}

	drawOptions := &ebiten.DrawImageOptions{}

	drawOptions.GeoM.Translate(float64(x)*(1/float64(w)), float64(y)*(1/float64(h)))
	drawOptions.GeoM.Scale(float64(w), float64(h))
	drawOptions.ColorM.Translate(float64(r)/255, float64(g)/255, float64(b)/255, float64(alpha)/255)

	target.DrawImage(imgSquare, drawOptions)
}

func ArrayToRGBA(tc []uint8) color.Color {
	return color.RGBA{R: tc[0], G: tc[1], B: tc[2], A: tc[3]}
}
