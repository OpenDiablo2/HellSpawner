package hscommon

type EditorWindow interface {
	Renderable

	GetWindowTitle() string
	Show()
	IsVisible() bool
	SetId(id string)
	GetId() string
	BringToFront()
}
