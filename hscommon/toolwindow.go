package hscommon

import (
	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
)

// ToolWindow represents tool windows
type ToolWindow interface {
	Renderable

	Show()
	IsVisible() bool
	SetVisible(bool)
	BringToFront()
	State() hsstate.ToolWindowState
	Pos(x, y float32) *giu.WindowWidget
	Size(float32, float32) *giu.WindowWidget
	CurrentSize() (float32, float32)
}
