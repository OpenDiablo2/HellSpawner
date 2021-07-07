package palettemapwidget

func getPaletteTransformString() []string {
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

	return selections
}

// cannot use map or literal, because len of transforms isn't the same, so
// for example if state.Slider1 = 30, state.selection = 0 (LightLevelVariation)
// and len p.pl2.InvColorVariations = 16, than
// (if we'd use map) we receive "index out of range" panic
func (p *widget) getPaletteIndices(state *widgetState) (indice *[256]byte) {
	switch state.Selection {
	case transformLightLevelVariations:
		indice = &p.pl2.LightLevelVariations[state.Slider1].Indices
	case transformInvColorVariations:
		indice = &p.pl2.InvColorVariations[state.Slider1].Indices
	case transformSelectedUintShift:
		indice = &p.pl2.SelectedUintShift.Indices
	case transformAlphaBlend:
		indice = &p.pl2.AlphaBlend[state.Slider2][state.Slider1].Indices
	case transformAdditiveBlend:
		indice = &p.pl2.AdditiveBlend[state.Slider1].Indices
	case transformMultiplicativeBlend:
		indice = &p.pl2.MultiplicativeBlend[state.Slider1].Indices
	case transformHueVariations:
		indice = &p.pl2.HueVariations[state.Slider1].Indices
	case transformRedTones:
		indice = &p.pl2.RedTones.Indices
	case transformGreenTones:
		indice = &p.pl2.GreenTones.Indices
	case transformBlueTones:
		indice = &p.pl2.BlueTones.Indices
	case transformUnknownVariations:
		indice = &p.pl2.UnknownVariations[state.Slider1].Indices
	case transformMaxComponentBlend:
		indice = &p.pl2.MaxComponentBlend[state.Slider1].Indices
	case transformDarkendColorShift:
		indice = &p.pl2.DarkendColorShift.Indices
	case transformTextColorShifts:
		indice = &p.pl2.TextColorShifts[state.Slider1].Indices
	}

	return indice
}
