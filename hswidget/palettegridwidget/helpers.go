package palettegridwidget

import (
	"fmt"
	"strconv"
)

// Hex2RGB converts haxadecimal color into r, g, b
func Hex2RGB(hex string) (r, g, b uint8, err error) {
	const (
		base    = 16
		bitSize = 32
		mask    = 0xFF
		rOffset = 16
		gOffset = 8
	)

	values, err := strconv.ParseUint(hex, base, bitSize)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("error parsing uint: %w", err)
	}

	r = uint8(values >> rOffset)
	g = uint8((values >> gOffset) & mask)
	b = uint8(values & mask)

	return r, g, b, nil
}

func t2x(t int64) string {
	const base = 16
	result := strconv.FormatInt(t, base)

	if len(result) == 1 {
		result = "0" + result
	}

	return result
}

// RGB2Hex converts RGB into hexadecimal
func RGB2Hex(red, green, blue uint8) string {
	r := t2x(int64(red))
	g := t2x(int64(green))
	b := t2x(int64(blue))

	return r + g + b
}
