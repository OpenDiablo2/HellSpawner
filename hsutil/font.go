package hsutil

import (
	"golang.org/x/image/font"
)

func CalculateBounds(text string, font font.Face) (int, int) {
	width := 0
	height := 0

	for _, char := range text {
		glyphWidth, _ := font.GlyphAdvance(char)
		width += int(float64(glyphWidth) / 64)
		glyphHeight := font.Metrics().Height.Ceil()

		if glyphHeight > height {
			height = glyphHeight
		}
	}

	return width, height
}
