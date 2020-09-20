package hsui

import "github.com/hajimehoshi/ebiten"

// Widget is an abstract representation of a UI widget.
type Widget interface {
	// Render renders the widget.
	Render(screen *ebiten.Image, x, y, width, height int)

	// Update updates the widget.
	Update()

	// GetRequestedSize returns the size the widget wants to be.
	GetRequestedSize() (int, int)

	// Invalidate causes the widget to recalculate itself and invalid all of its children.
	Invalidate()
}
