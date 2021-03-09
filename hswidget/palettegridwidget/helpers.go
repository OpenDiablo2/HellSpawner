package palettegridwidget

func (p *widget) changeColor(r, g, b uint8, idx int) {
	var rgba uint32
	rgba |= uint32(r) << 24
	rgba |= uint32(g) << 16
	rgba |= uint32(b) << 8
	rgba |= uint32(255) << 0
	p.colors[idx].SetRGBA(rgba)
	p.loadTexture(idx)
}
