package hscommon

type Renderable interface {
	Build()
	Render()
	Cleanup()
	HasFocus() (hasFocus bool)
	RegisterKeyboardShortcuts()
}
