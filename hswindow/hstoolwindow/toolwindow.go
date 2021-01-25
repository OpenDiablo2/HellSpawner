package hstoolwindow

import (
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
	"github.com/OpenDiablo2/HellSpawner/hswindow"
)

// ToolWindow represents a tool window
type ToolWindow struct {
	*hswindow.Window
	Type hsstate.ToolWindowType
}

// New creates a new tool window
func New(title string, toolWindowType hsstate.ToolWindowType, x, y float32) *ToolWindow {
	return &ToolWindow{
		Window: hswindow.New(title, x, y),
		Type:   toolWindowType,
	}
}

// State returns state of tool window
func (t *ToolWindow) State() hsstate.ToolWindowState {
	return hsstate.ToolWindowState{
		WindowState: t.Window.State(),
		Type:        t.Type,
	}
}
