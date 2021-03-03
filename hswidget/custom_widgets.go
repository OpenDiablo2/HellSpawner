package hswidget

import (
	"fmt"

	"github.com/ianling/giu"
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
// additionally, for byte checks, if value smaller than 255
func MakeInputInt(id string, width int32, output interface{}, optionalCB func()) *giu.InputIntWidget {
	var input int32
	switch o := output.(type) {
	case *byte:
		input = int32(*o)
	case *int:
		input = int32(*o)
	default:
		panic(fmt.Sprintf("MakeInputInt: invalid value type %T given", o))
	}

	return giu.InputInt(id, &input).Size(float32(width)).OnChange(func() {
		switch o := output.(type) {
		case *byte:
			SetByteToInt(input, o)
		case *int:
			*o = int(input)
		}

		if optionalCB != nil {
			optionalCB()
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
