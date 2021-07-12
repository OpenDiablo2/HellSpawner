package hswindow

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
)

// Window represents project's window
type Window struct {
	*giu.WindowWidget
	Visible bool
}

// New creates new window
func New(title string, x, y float32) *Window {
	return &Window{
		WindowWidget: giu.Window(title).Pos(x, y),
	}
}

// State returns window's state
func (t *Window) State() hsstate.WindowState {
	x, y := t.CurrentPosition()
	w, h := t.CurrentSize()

	return hsstate.WindowState{
		Visible: t.Visible,
		PosX:    x,
		PosY:    y,
		Width:   w,
		Height:  h,
	}
}

// ToggleVisibility toggles visibility
func (t *Window) ToggleVisibility() {
	t.Visible = !t.Visible
}

// Show turns visibility to true
func (t *Window) Show() {
	t.Visible = true
}

// Build builds window
func (t *Window) Build() {
	// noop
}

// RegisterKeyboardShortcuts sets a local shortcuts list
func (t *Window) RegisterKeyboardShortcuts(s ...giu.WindowShortcut) {
	t.WindowWidget.RegisterKeyboardShortcuts(s...)
}

// KeyboardShortcuts returns a list of local keyboard shortcuts
func (t *Window) KeyboardShortcuts() []giu.WindowShortcut {
	return []giu.WindowShortcut{}
}

// IsVisible returns true if window is visible
func (t *Window) IsVisible() bool {
	return t.Visible
}

// SetVisible sets window's visibility
func (t *Window) SetVisible(visible bool) {
	t.Visible = visible
}

// Cleanup hides window
func (t *Window) Cleanup() {
	t.Visible = false
}
