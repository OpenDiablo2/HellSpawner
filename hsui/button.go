package hsui

import (
	"image/color"
	"math"

	. "github.com/OpenDiablo2/HellSpawner/hsinput"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/OpenDiablo2/HellSpawner/hsutil"
)

const buttonPaddingH = 32
const buttonPaddingV = 16

type Button struct {
	caption            string
	infoProvider       hsutil.InfoProvider
	textColor          color.Color
	disabledTextColor  color.Color
	inputVector        *InputVector
	fontBoundsX        int
	fontBoundsY        int
	hovered            bool
	canExecuteCallback bool
	dirty              bool
	enabled            bool
	visible            bool
	reqWidth           int
	reqHeight          int
	onClick            func()
}

func CreateButton(infoProvider hsutil.InfoProvider, caption string, onClick func()) *Button {
	button := &Button{
		infoProvider:       infoProvider,
		caption:            caption,
		hovered:            false,
		visible:            true,
		canExecuteCallback: true,
		enabled:            true,
		dirty:              false,
		textColor:          hsutil.ArrayToRGBA(infoProvider.GetAppConfig().Colors.Text),
		disabledTextColor:  hsutil.ArrayToRGBA(infoProvider.GetAppConfig().Colors.DisabledText),
		onClick:            onClick,
	}

	button.inputVector = CreateInputVector().SetMouseButton(MouseButtonLeft)

	button.Invalidate()

	return button
}

func (b *Button) Render(screen *ebiten.Image, x, y, width, height int) {
	if width <= 0 || height <= 0 || !b.visible {
		return
	}

	primaryColor := hsutil.ArrayToRGBA(b.infoProvider.GetAppConfig().Colors.Primary)
	textColor := b.textColor

	if b.enabled {
		mouseX, mouseY := ebiten.CursorPosition()
		b.hovered = false
		if b.canExecuteCallback && mouseX >= hsutil.ScaleToDevice(x) && mouseX < hsutil.ScaleToDevice(x+width) &&
			mouseY >= hsutil.ScaleToDevice(y) && mouseY < hsutil.ScaleToDevice(y+height) {
			primaryColor = hsutil.ArrayToRGBA(b.infoProvider.GetAppConfig().Colors.PrimaryHighlight)
			b.hovered = true
		}
	} else {
		primaryColor = hsutil.ArrayToRGBA(b.infoProvider.GetAppConfig().Colors.Disabled)
		textColor = b.disabledTextColor
	}

	x = hsutil.ScaleToDevice(x)
	y = hsutil.ScaleToDevice(y)
	width = hsutil.ScaleToDevice(width)
	height = hsutil.ScaleToDevice(height)

	hsutil.DrawColoredRect(screen, x, y, width, height, primaryColor)

	font := b.infoProvider.GetNormalFont()
	heightDelta := int(math.Floor(float64(hsutil.ScaleToDevice(height)-b.fontBoundsY)*0.50)) - 3
	offsetX := hsutil.ScaleToDevice(x+(width/2)) - (b.fontBoundsX / 2)
	offsetY := hsutil.ScaleToDevice(y) + b.fontBoundsY + heightDelta

	text.Draw(screen, b.caption, font, offsetX, offsetY, textColor)
}

var lmbPressed bool

func (b *Button) Update() (dirty bool) {
	dirty = b.dirty

	if b.dirty {
		b.dirty = false
	}

	if !b.enabled || !b.visible {
		return dirty
	}

	lmbPressed = b.infoProvider.GetInputVector().Contains(b.inputVector)

	if b.canExecuteCallback && lmbPressed {
		if b.hovered {
			b.onClick()
		}

		b.canExecuteCallback = false
	} else if !lmbPressed {
		b.canExecuteCallback = true
	}

	return dirty
}

func (b *Button) GetRequestedSize() (int, int) {
	if !b.visible {
		return 0, 0
	}

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

func (b *Button) SetVisible(visible bool) {
	b.visible = visible
	b.dirty = true
}

func (b *Button) ToggleVisible() {
	b.SetVisible(!b.visible)
}

func (b *Button) GetVisible() bool {
	return b.visible
}
