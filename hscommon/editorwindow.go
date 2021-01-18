package hscommon

type EditorWindow interface {
	Renderable
	MainMenuUpdater

	GetWindowTitle() string
	Show()
	IsVisible() bool
	GetId() string
	BringToFront()
}
