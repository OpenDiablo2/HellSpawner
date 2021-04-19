package hscommon

import (
	"github.com/OpenDiablo2/HellSpawner/hsinput"
)

// Renderable represents renderable objects
type Renderable interface {
	Build()
	Render()
	Cleanup()
	RegisterKeyboardShortcuts(inputManager *hsinput.InputManager)
	IsVisible() bool
}
