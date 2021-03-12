package palettemapwidget

import (
	"fmt"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegrideditorwidget"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

func (p *widget) makeGrid(key string, colors *[256]palettegridwidget.PaletteColor) {
	c := make([]palettegridwidget.PaletteColor, len(colors))
	for n := range colors {
		c[n] = colors[n]
	}

	state := p.getState()
	state.textures[key] = palettegridwidget.Create(p.textureLoader, p.id+key, &c).OnClick(func(idx int) {
		state.id = key
		state.idx = idx
		state.mode = widgetModeEditTransform
	})
}
func (p *widget) getColors(indices *[256]byte) *[256]palettegridwidget.PaletteColor {
	result := &[256]palettegridwidget.PaletteColor{}

	for idx := range indices {
		// nolint:gomnd // const
		if idx > 255 {
			break
		}

		result[idx] = palettegridwidget.PaletteColor(&p.pl2.BasePalette.Colors[indices[idx]])
	}

	return result
}

// single transform (256 palette indices)
// example: selected unit
func (p *widget) transformSingle(key string, transform *[256]byte) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	if tex, found := state.textures[key]; found {
		l = append(l, tex)
	} else {
		p.makeGrid(key, p.getColors(transform))
	}

	return l
}

// multiple transforms (n * 256 palette indices)
// light level variations, there's 32
func (p *widget) transformMulti(key string, transforms []d2pl2.PL2PaletteTransform) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	numSelections := int32(len(transforms))

	if state.slider1 >= numSelections {
		state.slider1 = numSelections - 1
		p.setState(state)
	}

	textureID := fmt.Sprintf("%s_%d", key, state.slider1)

	l = append(l, giu.SliderInt("##"+key+"_slider", &state.slider1, 0, numSelections-1))

	if tex, found := state.textures[textureID]; found {
		l = append(l, tex)
	} else {
		p.makeGrid(textureID, p.getColors(&transforms[state.slider1].Indices))
	}

	return l
}

// tranferMultiGroup - groups of multiple transforms (m * n * 256 palette indices)
// example: alpha blend, there's 3 alpha levels (25%, 50%, 75% ?), and each do a blend against all 256 colors
func (p *widget) transformMultiGroup(key string, groups ...[256]d2pl2.PL2PaletteTransform) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	numGroups := int32(len(groups))

	if state.slider2 >= numGroups {
		state.slider2 = numGroups - 1
		p.setState(state)
	}

	if numGroups > 1 {
		sliderKey := fmt.Sprintf("##%s_group", key)
		l = append(l, giu.SliderInt(sliderKey, &state.slider2, 0, numGroups-1))
	}

	groupIdx := state.slider2

	numSelections := int32(len(groups[groupIdx]) - 1)

	if state.slider1 >= numSelections {
		state.slider1 = numSelections - 1
		p.setState(state)
	}

	textureID := fmt.Sprintf("%s_%d_%d", key, state.slider2, state.slider1)

	l = append(l, giu.SliderInt("##"+key+"_slider", &state.slider1, 0, numSelections))

	if tex, found := state.textures[textureID]; found {
		l = append(l, tex)
	} else {
		col := p.getColors(&groups[groupIdx][state.slider1].Indices)
		p.makeGrid(textureID, col)
	}

	return l
}

func (p *widget) textColors(key string, colors []d2pl2.PL2Color24Bits) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	numSelections := int32(len(colors) - 1)

	if state.slider1 >= numSelections {
		state.slider1 = numSelections - 1
		p.setState(state)
	}

	textureID := fmt.Sprintf("%s_%d", key, state.slider1)
	if tex, found := state.textures[textureID]; found {
		l = append(l, tex)
	} else {
		c := make([]palettegridwidget.PaletteColor, len(p.pl2.TextColors))

		for n := range c {
			c[n] = palettegridwidget.PaletteColor(&p.pl2.TextColors[n])
		}

		grid := palettegrideditorwidget.Create(p.textureLoader, p.id+"transform24editColor", &c)

		state.textures[textureID] = grid
	}

	return l
}
