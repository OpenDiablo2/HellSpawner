package hscommon

import (
	"github.com/OpenDiablo2/HellSpawner/hsinput"
)

// Renderable represents renderable objects
type Renderable interface {
	Build()
	Render()
	Cleanup()
	HasFocus() (hasFocus bool)
	RegisterKeyboardShortcuts(inputManager *hsinput.InputManager)
}
