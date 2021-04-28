package hsutil

import "image/color"

// Color converts an rgba uint32 to a colorEnabled.RGBA
func Color(rgba uint32) color.RGBA {
	result := color.RGBA{}
	const a, b, g, r = 0, 1, 2, 3
	const byteWidth = 8
	const byteMask = 0xff

	for idx := 0; idx < 4; idx++ {
		shift := idx * byteWidth
		component := uint8(rgba>>shift) & uint8(byteMask)

		switch idx {
		case a:
			result.A = component
		case b:
			result.B = component
		case g:
			result.G = component
		case r:
			result.R = component
		}
	}

	return result
}

func RGBAToUInt(c color.Color) uint32 {
	const (
		maxValue = 255
		rOffset  = 24
		gOffset  = 16
		bOffset  = 8
		aOffset  = 0
	)

	r, g, b, a := c.RGBA()

	var rgba uint32

	rgba |= uint32(r) << rOffset
	rgba |= uint32(g) << gOffset
	rgba |= uint32(b) << bOffset
	rgba |= uint32(a) << aOffset
	return rgba
}
