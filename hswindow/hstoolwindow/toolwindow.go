package hstoolwindow

import (
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
	"github.com/OpenDiablo2/HellSpawner/hswindow"
)

type ToolWindow struct {
	*hswindow.Window
	Type hsstate.ToolWindowType
}

func New(title string, toolWindowType hsstate.ToolWindowType, x, y float32) *ToolWindow {
	return &ToolWindow{
		Window: hswindow.New(title, x, y),
		Type:   toolWindowType,
	}
}

func (t *ToolWindow) State() hsstate.ToolWindowState {
	return hsstate.ToolWindowState{
		WindowState: t.Window.State(),
		Type:        t.Type,
	}
}
