package hswidget

import (
	"fmt"
	"image"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

// PaletteMapViewerState creates a new palette map viewer's state
type PaletteMapViewerState struct {
	selection int32
	slider1   int32
	slider2   int32
	textures  map[string]*giu.Texture
}

// Dispose cleans viewer's state
func (p *PaletteMapViewerState) Dispose() {
	p.textures = make(map[string]*giu.Texture)
}

// PaletteMapViewerWidget represents a palette map viewer's widget
type PaletteMapViewerWidget struct {
	id  string
	pl2 *d2pl2.PL2
}

// PaletteMapViewer creates a new palette map viewer's widget
func PaletteMapViewer(id string, pl2 *d2pl2.PL2) *PaletteMapViewerWidget {
	result := &PaletteMapViewerWidget{
		id:  id,
		pl2: pl2,
	}

	return result
}

func (p *PaletteMapViewerWidget) getStateID() string {
	return fmt.Sprintf("PaletteMapViewerWidget_%s", p.id)
}

func (p *PaletteMapViewerWidget) initState() {
	state := &PaletteMapViewerState{
		textures: make(map[string]*giu.Texture),
	}

	p.setState(state)
}

func (p *PaletteMapViewerWidget) getState() *PaletteMapViewerState {
	var state *PaletteMapViewerState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*PaletteMapViewerState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *PaletteMapViewerWidget) setState(s *PaletteMapViewerState) {
	giu.Context.SetState(p.getStateID(), s)
}

// Build builds a new widget
func (p *PaletteMapViewerWidget) Build() {
	state := p.getState()

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

	left := giu.Layout{
		giu.Label("Base Palette"),
		p.paletteView(),
	}

	right := giu.Layout{
		giu.Label("Palette Map"),
		giu.Layout{
			giu.Combo("", selections[state.selection], selections, &state.selection).Size(180),
			p.getTransformViewLayout(state.selection),
		},
	}

	w1, h1 := float32(256+32), float32(256+48)
	w2, h2 := w1, h1

	// special case for alpha blend
	if state.selection == 3 {
		h2 += 32
	}

	layout := giu.Layout{
		giu.Child("left").Size(w1, h1).Layout(left),
		giu.Child("right").Size(w2, h2).Layout(right),
	}

	layout.Build()
}

func (p *PaletteMapViewerWidget) getTransformViewLayout(transformIdx int32) giu.Layout {
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

func (p *PaletteMapViewerWidget) makeTexture(key string, colors *[256]d2interface.Color) {
	pix := make([]byte, 256*4)

	img := &image.RGBA{
		Rect: image.Rectangle{
			image.Point{0, 0},
			image.Point{16, 16},
		},
	}

	for idx := range colors {
		pix[idx*4+0] = colors[idx].R()
		pix[idx*4+1] = colors[idx].G()
		pix[idx*4+2] = colors[idx].B()
		pix[idx*4+3] = 0xFF
	}

	img.Pix = pix

	makeTexture := func(tex *giu.Texture) {
		state := p.getState()
		state.textures[key] = tex
	}

	hscommon.CreateTextureFromARGB(img, makeTexture)
}

func (p *PaletteMapViewerWidget) getColors(indices *[256]byte) *[256]d2interface.Color {
	result := &[256]d2interface.Color{}

	for idx := range indices {
		if idx > 255 {
			break
		}

		c := p.pl2.BasePalette.Colors[indices[idx]]

		cface := &colorFace{
			r: c.R,
			g: c.G,
			b: c.B,
		}

		result[idx] = cface
	}

	return result
}

// single transform (256 palette indices)
// example: selected unit
func (p *PaletteMapViewerWidget) transformSingle(key string, transform *[256]byte) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	if tex, found := state.textures[key]; found {
		l = append(l, giu.Image(tex).Size(255, 255))
	} else {
		p.makeTexture(key, p.getColors(transform))
	}

	return l
}

// multiple transforms (n * 256 palette indices)
// light level variations, there's 32
func (p *PaletteMapViewerWidget) transformMulti(key string, transforms []d2pl2.PL2PaletteTransform) giu.Layout {
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
		l = append(l, giu.Image(tex).Size(208, 208))
	} else {
		p.makeTexture(textureID, p.getColors(&transforms[state.slider1].Indices))
	}

	return l
}

// tranferMultiGroup - groups of multiple transforms (m * n * 256 palette indices)
// example: alpha blend, there's 3 alpha levels (25%, 50%, 75% ?), and each do a blend against all 256 colors
func (p *PaletteMapViewerWidget) transformMultiGroup(key string, groups ...[256]d2pl2.PL2PaletteTransform) giu.Layout {
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
		l = append(l, giu.Image(tex).Size(208, 208))
	} else {
		col := p.getColors(&groups[groupIdx][state.slider1].Indices)
		p.makeTexture(textureID, col)
	}

	return l
}

func (p *PaletteMapViewerWidget) textColors(key string, colors []d2pl2.PL2Color24Bits) giu.Layout {
	state := p.getState()

	l := giu.Layout{}

	numSelections := int32(len(colors) - 1)

	if state.slider1 >= numSelections {
		state.slider1 = numSelections - 1
		p.setState(state)
	}

	textureID := fmt.Sprintf("%s_%d", key, state.slider1)
	if tex, found := state.textures[textureID]; found {
		l = append(l, giu.Image(tex).Size(float32(len(colors)*16), 16))
	} else {
		colorFaces := make([]d2interface.Color, len(colors))

		for idx := range colors {
			cface := &colorFace{
				r: colors[idx].R,
				g: colors[idx].G,
				b: colors[idx].B,
			}

			colorFaces[idx] = cface
		}

		pix := make([]byte, len(colors)*4)

		img := &image.RGBA{
			Rect: image.Rectangle{
				image.Point{0, 0},
				image.Point{len(colors), 1},
			},
		}

		for idx := range colors {
			pix[idx*4+0] = colorFaces[idx].R()
			pix[idx*4+1] = colorFaces[idx].G()
			pix[idx*4+2] = colorFaces[idx].B()
			pix[idx*4+3] = 0xFF
		}

		img.Pix = pix

		makeTexture := func(tex *giu.Texture) {
			state := p.getState()
			state.textures[textureID] = tex
		}

		hscommon.CreateTextureFromARGB(img, makeTexture)
	}

	return l
}

func (p *PaletteMapViewerWidget) paletteView() giu.Layout {
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
	return 0xff
}

func (c colorFace) RGBA() uint32 {
	return uint32(c.r)<<24 | uint32(c.g)<<16 | uint32(c.b)<<8 | uint32(0xff)
}

func (c colorFace) SetRGBA(u uint32) {
	c.r = byte((u >> 24) & 0xff)
	c.g = byte((u >> 16) & 0xff)
	c.b = byte((u >> 8) & 0xff)
}

func (c colorFace) BGRA() uint32 {
	return uint32(c.b)<<8 | uint32(c.g)<<16 | uint32(c.r)<<24 | uint32(0xff)
}

func (c colorFace) SetBGRA(u uint32) {
	c.b = byte((u >> 24) & 0xff)
	c.g = byte((u >> 16) & 0xff)
	c.r = byte((u >> 8) & 0xff)
}
