package hscommon

// Renderable represents renderable objects
type Renderable interface {
	Build()
	Render()
	Cleanup()
	HasFocus() (hasFocus bool)
	RegisterKeyboardShortcuts()
}
