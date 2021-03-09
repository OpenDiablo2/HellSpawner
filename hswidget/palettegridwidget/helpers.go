package palettegridwidget

func (p *widget) changeColor(r, g, b uint8, idx int) {
	const (
		maxValue = 255
		rOffset  = 24
		gOffset  = 16
		bOffset  = 8
		aOffset  = 0
	)

	var rgba uint32
	rgba |= uint32(r) << rOffset
	rgba |= uint32(g) << gOffset
	rgba |= uint32(b) << bOffset
	rgba |= uint32(maxValue) << aOffset
	p.colors[idx].SetRGBA(rgba)
	p.loadTexture(idx)
}
