package hscommon

const test = "huh?"

type EditorWindow interface {
	Renderable

	GetWindowTitle() string
	Show()
	IsVisible() bool
	GetId() string
	BringToFront()
}
