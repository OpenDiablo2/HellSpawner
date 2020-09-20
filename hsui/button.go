package hsui

import (
	"image/color"

	"github.com/OpenDiablo2/HellSpawner/hsutil"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
)

const buttonPaddingH = 4
const buttonPaddingV = 8

type Button struct {
	caption            string
	infoProvider       hsutil.InfoProvider
	textColor          color.Color
	disabledTextColor  color.Color
	fontBoundsX        int
	fontBoundsY        int
	hovered            bool
	canExecuteCallback bool
	enabled            bool
	reqWidth           int
	reqHeight          int
	onClick            func()
}

func CreateButton(infoProvider hsutil.InfoProvider, caption string, onClick func()) *Button {
	tc := infoProvider.GetAppConfig().Colors.Text
	dtc := infoProvider.GetAppConfig().Colors.DisabledText

	result := &Button{
		infoProvider:       infoProvider,
		caption:            caption,
		hovered:            false,
		canExecuteCallback: true,
		enabled:            true,
		textColor:          color.RGBA{R: tc[0], G: tc[1], B: tc[2], A: tc[3]},
		disabledTextColor:  color.RGBA{R: dtc[0], G: dtc[1], B: dtc[2], A: dtc[3]},
		onClick:            onClick,
	}

	result.Invalidate()

	return result
}

func (b *Button) Render(screen *ebiten.Image, x, y, width, height int) {
	primaryColor := b.infoProvider.GetAppConfig().Colors.Primary
	textColor := b.textColor

	if b.enabled {
		mouseX, mouseY := ebiten.CursorPosition()
		b.hovered = false
		if b.canExecuteCallback && mouseX >= hsutil.ScaleToDevice(x) && mouseX < hsutil.ScaleToDevice(x+width) &&
			mouseY >= hsutil.ScaleToDevice(y) && mouseY < hsutil.ScaleToDevice(y+height) {
			primaryColor = b.infoProvider.GetAppConfig().Colors.PrimaryHighlight
			b.hovered = true
		}
	} else {
		primaryColor = b.infoProvider.GetAppConfig().Colors.Disabled
		textColor = b.disabledTextColor
	}

	hsutil.DrawColoredRect(
		screen,
		hsutil.ScaleToDevice(x), hsutil.ScaleToDevice(y),
		hsutil.ScaleToDevice(width), hsutil.ScaleToDevice(height),
		primaryColor[0], primaryColor[1], primaryColor[2], primaryColor[3])

	font := b.infoProvider.GetNormalFont()
	heightDelta := int(float64(hsutil.ScaleToDevice(height)-b.fontBoundsY) * 0.30)
	offsetX := hsutil.ScaleToDevice(x+(width/2)) - (b.fontBoundsX / 2)
	offsetY := hsutil.ScaleToDevice(y) + b.fontBoundsY + heightDelta

	text.Draw(screen, b.caption, font, offsetX, offsetY, textColor)
}

func (b *Button) Update() {
	if !b.enabled {
		return
	}

	if b.canExecuteCallback {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if b.hovered {
				b.onClick()
			}
			b.canExecuteCallback = false
		}
	} else {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			b.canExecuteCallback = true
		}
	}

}

func (b *Button) GetRequestedSize() (int, int) {
	return b.reqWidth, b.reqHeight
}

func (b *Button) Invalidate() {
	font := b.infoProvider.GetNormalFont()
	b.fontBoundsX, b.fontBoundsY = hsutil.CalculateBounds(b.caption, font)
	b.reqWidth = hsutil.UnscaleFromDevice(b.fontBoundsX) + (buttonPaddingH * 2)
	b.reqHeight = hsutil.UnscaleFromDevice(b.fontBoundsY) + (buttonPaddingV * 2)
}

func (b *Button) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func (b *Button) IsEnabled() bool {
	return b.enabled
}
