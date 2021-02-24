package hsutil

import (
	"fmt"

	"github.com/ianling/giu"
)

// MakeImageButton is a hack for giu.ImageButton that creates image button
// as a giu.child
func MakeImageButton(id string, w, h int, t *giu.Texture, fn func()) giu.Layout {
	const (
		childIdSuffix = "child"
		padding       = 8 // pixels
	)

	// the image button
	btnW, btnH := float32(w), float32(h)
	button := giu.Layout{
		giu.ImageButton(t).Size(btnW, btnH).OnClick(fn),
	}

	// the container; needs to be padded to be larger than the button
	childW, childH := btnW+padding, btnH+padding
	childID := fmt.Sprintf("%s%s", id, childIdSuffix)
	con := giu.Child(childID).
		Border(false).
		Size(float32(childW), float32(childH)).
		Layout(button).
		Flags(giu.WindowFlagsNoDecoration)

	return giu.Layout{con}
}
