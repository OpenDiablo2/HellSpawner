package hsutil

import (
	"github.com/ianling/giu"
)

// it looks like childWidget's size must be a bit greater then widget iside of it
const (
	stdMod = 8
)

// MakeImageButton is a hack for giu.ImageButton that creates image button
// as a giu.child
func MakeImageButton(id string, w, h int, t *giu.Texture, fn func()) giu.Layout {
	return giu.Layout{
		giu.Child(id+"child").Border(false).Size(float32(w+stdMod), float32(h+stdMod)).Layout(giu.Layout{
			giu.ImageButton(t).Size(float32(w), float32(h)).OnClick(fn),
		}).Flags(giu.WindowFlagsNoDecoration),
	}
}
