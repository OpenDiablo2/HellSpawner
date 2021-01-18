package hscommon

import "github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"

type EditorWindow interface {
	Renderable
	MainMenuUpdater

	GetWindowTitle() string
	Show()
	IsVisible() bool
	SetVisible(bool)
	GetId() string
	BringToFront()
	State() hsstate.EditorState
}
