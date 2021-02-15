package hsutil

import (
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"
)

// MakeImageButton is a hack for giu.ImageButton that manually sets the ImageButton ID by
// manually setting the ID through imgui
//nolint:unparam // width and height are always 15 at the time of writing, but may change
func MakeImageButton(id string, w, h int, t *giu.Texture, fn func()) giu.Layout {
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
