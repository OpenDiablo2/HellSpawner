package palettegridwidget

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
