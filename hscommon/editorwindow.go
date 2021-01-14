package hscommon

type EditorWindow interface {
	Renderable
	MainMenuUpdater
	FocusController

	GetWindowTitle() string
	Show()
	IsVisible() bool
	IsFocused() bool
	GetId() string
	BringToFront()
}
