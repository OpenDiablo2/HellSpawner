package cofwidget

import (
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
)

//nolint:unparam // width and height are always 15 at the time of writing, but may change
func makeImageButton(id string, w, h int, t *giu.Texture, fn func()) giu.Layout {
	return giu.Layout{
		giu.ImageButton(t).Size(float32(w), float32(h)).OnClick(fn),
		giu.Custom(func() {
			// make this button unique across all editor instances
			// at the time of writing, ImageButton uses the texture ID as the button ID
			// so it wont be unique across multiple instances if we use the same texture...
			// we need to step over giu and manually tell imgui to pop the last ID and
			// push the desired one onto the stack
			imgui.PopID()
			imgui.PushID(id)
		}),
	}
}

// this likely needs to be a method of d2cof.COF
func speedToFPS(speed int) float64 {
	const (
		baseFPS      = 25
		speedDivisor = 256
	)

	fps := baseFPS * (float64(speed) / speedDivisor)
	if fps == 0 {
		fps = baseFPS
	}

	return fps
}

// this should also probably be a method of COF
func calculateDuration(cof *d2cof.COF) float64 {
	const (
		milliseconds = 1000
	)

	frameDelay := milliseconds / speedToFPS(cof.Speed)

	return float64(cof.FramesPerDirection) * frameDelay
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
