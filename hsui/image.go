package hsui

import (
	"image"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

func CreateImage(path string) (*Image, error) {
	imageData, err := ebitenutil.OpenFile(path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imageData)
	if err != nil {
		return nil, err
	}

	ebitenImage, err := ebiten.NewImageFromImage(img, ebiten.FilterDefault)
	if err != nil {
		return nil, err
	}

	hsImage := &Image{
		image: ebitenImage,
	}

	hsImage.Invalidate()

	return hsImage, nil
}

type Image struct {
	image   *ebiten.Image
	options *ebiten.DrawImageOptions
	fit     bool
}

func (i *Image) Render(screen *ebiten.Image, x, y, w, h int) {
	drawOptions := &ebiten.DrawImageOptions{}

	imgW, imgH := i.image.Size()
	sx, sy := float64(w)/float64(imgW), float64(h)/float64(imgH)

	if i.fit {
		drawOptions.GeoM.Scale(sx, sy)
	}

	drawOptions.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(i.image, drawOptions)
}

func (i *Image) Update() (dirty bool) {
	// nothing to do
	return false
}

func (i *Image) GetRequestedSize() (int, int) {
	return i.image.Size()
}

func (i *Image) Invalidate() {
	// nothing to do
}

func (i *Image) SetFit(b bool) {
	i.fit = b
}
