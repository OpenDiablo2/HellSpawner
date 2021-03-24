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
// for example if state.slider1 = 30, state.selection = 0 (LightLevelVariation)
// and len p.pl2.InvColorVariations = 16, than
// (if we'd use map) we receive "index out of range" panic
func (p *widget) getPaletteIndices(state *widgetState) (indice *[256]byte) {
	switch state.selection {
	case transformLightLevelVariations:
		indice = &p.pl2.LightLevelVariations[state.slider1].Indices
	case transformInvColorVariations:
		indice = &p.pl2.InvColorVariations[state.slider1].Indices
	case transformSelectedUintShift:
		indice = &p.pl2.SelectedUintShift.Indices
	case transformAlphaBlend:
		indice = &p.pl2.AlphaBlend[state.slider2][state.slider1].Indices
	case transformAdditiveBlend:
		indice = &p.pl2.AdditiveBlend[state.slider1].Indices
	case transformMultiplicativeBlend:
		indice = &p.pl2.MultiplicativeBlend[state.slider1].Indices
	case transformHueVariations:
		indice = &p.pl2.HueVariations[state.slider1].Indices
	case transformRedTones:
		indice = &p.pl2.RedTones.Indices
	case transformGreenTones:
		indice = &p.pl2.GreenTones.Indices
	case transformBlueTones:
		indice = &p.pl2.BlueTones.Indices
	case transformUnknownVariations:
		indice = &p.pl2.UnknownVariations[state.slider1].Indices
	case transformMaxComponentBlend:
		indice = &p.pl2.MaxComponentBlend[state.slider1].Indices
	case transformDarkendColorShift:
		indice = &p.pl2.DarkendColorShift.Indices
	case transformTextColorShifts:
		indice = &p.pl2.TextColorShifts[state.slider1].Indices
	}

	return indice
}
