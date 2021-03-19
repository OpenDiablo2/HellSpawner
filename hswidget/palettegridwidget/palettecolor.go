package palettegridwidget

// PaletteColor represents palette color
type PaletteColor interface {
	RGBA() uint32
	SetRGBA(uint32)
}
