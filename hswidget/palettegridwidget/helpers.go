package palettegridwidget

import (
	"strconv"
)

func (p *widget) changeColor(state *widgetState) {
	const (
		maxValue = 255
		rOffset  = 24
		gOffset  = 16
		bOffset  = 8
		aOffset  = 0
	)

	var rgba uint32
	rgba |= uint32(state.r) << rOffset
	rgba |= uint32(state.g) << gOffset
	rgba |= uint32(state.b) << bOffset
	rgba |= uint32(maxValue) << aOffset
	p.colors[state.idx].SetRGBA(rgba)
	p.loadTexture(state.idx)
}

func Hex2RGB(hex string) (r, g, b uint8, err error) {
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return 0, 0, 0, err
	}

	r = uint8(values >> 16)
	g = uint8((values >> 8) & 0xFF)
	b = uint8(values & 0xFF)

	return r, g, b, nil
}

func t2x(t int64) string {
	result := strconv.FormatInt(t, 16)
	if len(result) == 1 {
		result = "0" + result
	}
	return result
}

func RGB2Hex(red, green, blue uint8) string {
	r := t2x(int64(red))
	g := t2x(int64(green))
	b := t2x(int64(blue))
	return string(r + g + b)
}
