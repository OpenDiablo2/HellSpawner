package palettemapwidget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegrideditorwidget"
	"github.com/OpenDiablo2/HellSpawner/hswidget/palettegridwidget"
)

const (
	leftLayoutSize       = 180
	bigImageW, bigImageH = 208, 208
)

type widget struct {
	id            string
	pl2           *d2pl2.PL2
	textureLoader *hscommon.TextureLoader
	baseColors    *[256]palettegridwidget.PaletteColor
}

// Create creates a new palette map viewer's widget
func Create(textureLoader *hscommon.TextureLoader, id string, pl2 *d2pl2.PL2) giu.Widget {
	result := &widget{
		id:            id,
		pl2:           pl2,
		textureLoader: textureLoader,
	}

	return result
}

// Build builds a new widget
func (p *widget) Build() {
	// nolint:ifshort // state should be a global variable here
	state := p.getState()

	switch state.mode {
	case widgetModeView:
		p.buildViewer(state)
	case widgetModeEditTransform:
		p.buildEditor(state)
	}
}

func (p *widget) buildViewer(state *widgetState) {
	selections := []string{
		"Light Level Variations",
		"InvColor Variations",
		"Selected Unit Shift",
		"Alpha Blend",
		"Additive Blend",
		"Multiplicative Blend",
		"Hue Variations",
		"Red Tones",
		"Green Tones",
		"Blue Tones",
		"Unknown Variations",
		"MaxComponent Blend",
		"Darkened Color Shift",
		"Text Colors",
		"Text ColorShifts",
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	var baseColors = make([]palettegridwidget.PaletteColor, 256)
	colors := &p.pl2.BasePalette.Colors
	for n := range colors {
		baseColors[n] = palettegridwidget.PaletteColor(&p.pl2.BasePalette.Colors[n])
	}

	left := giu.Layout{
		giu.Label("Base Palette"),
		palettegrideditorwidget.Create(p.textureLoader, p.id+"basePalette", &baseColors).OnChange(func() {
			//state.textures = make(map[string]*palettegridwidget.PaletteGridWidget)
			state.textures = make(map[string]giu.Widget)
		}),
	}

	right := giu.Layout{
		giu.Label("Palette Map"),
		giu.Layout{
			giu.Combo("", selections[state.selection], selections, &state.selection).Size(leftLayoutSize),
			p.getTransformViewLayout(state.selection),
		},
	}

	// nolint:gomnd // constant
	w1, h1 := float32(256+32), float32(256+48)
	w2, h2 := w1, h1

	// nolint:gomnd // special case for alpha blend
	if state.selection == 3 {
		h2 += 32
	}

	layout := giu.Layout{
		giu.Child("left").Size(w1, h1).Layout(left),
		giu.Child("right").Size(w2, h2).Layout(right),
	}

	layout.Build()
}

func (p *widget) buildEditor(state *widgetState) {
	var grid giu.Widget
	var colors = make([]palettegridwidget.PaletteColor, len(p.pl2.BasePalette.Colors))
	for n := range colors {
		colors[n] = palettegridwidget.PaletteColor(&p.pl2.BasePalette.Colors[n])
	}

	indices := []*[256]uint8{
		&p.pl2.LightLevelVariations[state.slider1].Indices,
		&p.pl2.InvColorVariations[state.slider1].Indices,
		&p.pl2.SelectedUintShift.Indices,
		&p.pl2.AlphaBlend[state.slider2][state.slider1].Indices,
		&p.pl2.AdditiveBlend[state.slider1].Indices,
		&p.pl2.MultiplicativeBlend[state.slider1].Indices,
		&p.pl2.HueVariations[state.slider1].Indices,
		&p.pl2.RedTones.Indices,
		&p.pl2.GreenTones.Indices,
		&p.pl2.BlueTones.Indices,
		&p.pl2.UnknownVariations[state.slider1].Indices,
		&p.pl2.MaxComponentBlend[state.slider1].Indices,
		&p.pl2.DarkendColorShift.Indices,
		nil,
		&p.pl2.TextColorShifts[state.slider1].Indices,
	}

	cols := indices[state.selection]

	grid = palettegridwidget.Create(p.textureLoader, p.id+"transformEdit", &colors).OnClick(func(idx int) {
		cols[state.idx] = byte(idx)
		state.mode = widgetModeView
	})

	giu.Layout{grid}.Build()
}

func (p *widget) getTransformViewLayout(transformIdx int32) giu.Layout {
	buildLayout := []func() giu.Layout{
		func() giu.Layout {
			return p.transformMulti("LightLevelVariations", p.pl2.LightLevelVariations[:])
		},
		func() giu.Layout {
			return p.transformMulti("InvColorVariations", p.pl2.InvColorVariations[:])
		},
		func() giu.Layout {
			return p.transformSingle("SelectedUintShift", &p.pl2.SelectedUintShift.Indices)
		},
		func() giu.Layout {
			return p.transformMultiGroup("AlphaBlend", p.pl2.AlphaBlend[:]...)
		},
		func() giu.Layout {
			return p.transformMulti("AdditiveBlend", p.pl2.AdditiveBlend[:])
		},
		func() giu.Layout {
			return p.transformMulti("MultiplicativeBlend", p.pl2.MultiplicativeBlend[:])
		},
		func() giu.Layout {
			return p.transformMulti("HueVariations", p.pl2.HueVariations[:])
		},
		func() giu.Layout {
			return p.transformSingle("RedTones", &p.pl2.RedTones.Indices)
		},
		func() giu.Layout {
			return p.transformSingle("GreenTones", &p.pl2.GreenTones.Indices)
		},
		func() giu.Layout {
			return p.transformSingle("BlueTones", &p.pl2.BlueTones.Indices)
		},
		func() giu.Layout {
			return p.transformMulti("UnknownVariations", p.pl2.UnknownVariations[:])
		},
		func() giu.Layout {
			return p.transformMulti("MaxComponentBlend", p.pl2.MaxComponentBlend[:])
		},
		func() giu.Layout {
			return p.transformSingle("DarkendColorShift", &p.pl2.DarkendColorShift.Indices)
		},
		func() giu.Layout {
			return p.textColors("TextColors", p.pl2.TextColors[:])
		},
		func() giu.Layout {
			return p.transformMulti("TextColorShifts", p.pl2.TextColorShifts[:])
		},
	}

	return buildLayout[transformIdx]()
}

func (p *widget) makeTexture(key string, colors *[256]palettegridwidget.PaletteColor) {
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
		p.makeTexture(key, p.getColors(transform))
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
		p.makeTexture(textureID, p.getColors(&transforms[state.slider1].Indices))
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
		p.makeTexture(textureID, col)
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

func (p *widget) paletteView() giu.Layout {
	baseTransform := [256]byte{}

	for idx := range baseTransform {
		baseTransform[idx] = byte(idx)
	}

	return p.transformSingle("base_palette", &baseTransform)
}

type colorFace struct {
	r, g, b uint8
}

func (c colorFace) R() uint8 {
	return c.r
}

func (c colorFace) G() uint8 {
	return c.g
}

func (c colorFace) B() uint8 {
	return c.b
}

func (c colorFace) A() uint8 {
	// nolint:gomnd // full-opacity (RGBA a=255)
	return 0xff
}

func (c colorFace) RGBA() uint32 {
	// nolint:gomnd // constants
	return uint32(c.r)<<24 | uint32(c.g)<<16 | uint32(c.b)<<8 | uint32(0xff)
}

// nolint:gomnd // constants
func (c colorFace) SetRGBA(u uint32) {
	c.r = byte((u >> 24) & 0xff)
	c.g = byte((u >> 16) & 0xff)
	c.b = byte((u >> 8) & 0xff)
}

func (c colorFace) BGRA() uint32 {
	// nolint:gomnd // constants
	return uint32(c.b)<<8 | uint32(c.g)<<16 | uint32(c.r)<<24 | uint32(0xff)
}

// nolint:gomnd // constants
func (c colorFace) SetBGRA(u uint32) {
	c.b = byte((u >> 24) & 0xff)
	c.g = byte((u >> 16) & 0xff)
	c.r = byte((u >> 8) & 0xff)
}
