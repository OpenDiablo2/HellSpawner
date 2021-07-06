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
	rgba |= uint32(state.rgba.R) << rOffset
	rgba |= uint32(state.rgba.G) << gOffset
	rgba |= uint32(state.rgba.B) << bOffset
	rgba |= uint32(maxValue) << aOffset
	(*p.colors)[state.idx].SetRGBA(rgba)
}
