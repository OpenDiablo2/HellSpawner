package palettegridwidget

type PaletteColor interface {
	RGBA() uint32
	SetRGBA(uint32)
}
