package palettegrideditorwidget

func (p *PaletteGridEditorWidget) changeColor(state *widgetState) {
	const (
		maxValue = 255
		rOffset  = 24
		gOffset  = 16
		bOffset  = 8
		aOffset  = 0
	)

	var rgba uint32
	rgba |= uint32(state.RGBA.R) << rOffset
	rgba |= uint32(state.RGBA.G) << gOffset
	rgba |= uint32(state.RGBA.B) << bOffset
	rgba |= uint32(maxValue) << aOffset
	(*p.colors)[state.Idx].SetRGBA(rgba)
}
