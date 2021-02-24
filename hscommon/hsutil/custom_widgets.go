package hsutil

import (
	"fmt"

	"github.com/ianling/giu"
)

const (
	inputIntW = 30
)

// MakeImageButton is a hack for giu.ImageButton that creates image button
// as a giu.child
func MakeImageButton(id string, w, h int, t *giu.Texture, fn func()) giu.Layout {
	const (
		childIDSuffix = "child"
		padding       = 8 // pixels
	)

	// the image button
	btnW, btnH := float32(w), float32(h)
	button := giu.Layout{
		giu.ImageButton(t).Size(btnW, btnH).OnClick(fn),
	}

	// the container; needs to be padded to be larger than the button
	childW, childH := btnW+padding, btnH+padding
	childID := fmt.Sprintf("%s%s", id, childIDSuffix)
	con := giu.Child(childID).
		Border(false).
		Size(childW, childH).
		Layout(button).
		Flags(giu.WindowFlagsNoDecoration)

	return giu.Layout{con}
}

// SetByteToInt sets byte given to intager
// if intager > max possible byte size, sets to 255
func SetByteToInt(input int32, output *byte) {
	const (
		// nolint:gomnd // constant
		maxByteSize = byte(255)
	)

	if input > int32(maxByteSize) {
		*output = maxByteSize

		return
	}

	*output = byte(input)
}

// MakeInputInt creates input intager using POINTER given
// additionaly, for byte checks, if value smaller than 255
func MakeInputInt(id string, width int32, output interface{}) *giu.InputIntWidget {
	var input int32
	switch output.(type) {
	case *byte:
		input = int32(*output.(*byte))
	case *int:
		input = int32(*output.(*int))
	default:
		panic("invalid value type given")
	}

	return giu.InputInt(id, &input).Size(float32(width)).OnChange(func() {
		switch output.(type) {
		case *byte:
			SetByteToInt(input, &(*output.(*byte)))
		case *int:
			*output.(*int) = int(input)
		}
	})
}

// MakeCheckboxFromByte creates a checkbox using a byte as input/output
func MakeCheckboxFromByte(id string, value *byte) *giu.CheckboxWidget {
	v := (*value > 0)

	return giu.Checkbox(id, &v).OnChange(func() {
		if v {
			*value = 1
		} else {
			*value = 0
		}
	})
}
