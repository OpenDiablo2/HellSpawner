package hscommon

type EditorWindow interface {
	Renderable

	GetWindowTitle() string
	Show()
}
