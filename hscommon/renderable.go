package hscommon

import "github.com/ianling/giu"

// Renderable represents renderable objects
type Renderable interface {
	Build()
	Cleanup()
	// KeyboardShortcuts returns a list of keyboard shortcuts
	KeyboardShortcuts() []giu.WindowShortcut
	IsVisible() bool
	// RegisterKeyboardShortcuts wraps giu.RegisterKeyboardShortcuts
	RegisterKeyboardShortcuts(...giu.WindowShortcut)
}
