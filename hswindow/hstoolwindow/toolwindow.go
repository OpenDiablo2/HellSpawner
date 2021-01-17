package hstoolwindow

import "github.com/OpenDiablo2/HellSpawner/hswindow"

type ToolWindow struct {
	*hswindow.Window
}

func New(title string, x, y float32) *ToolWindow {
	return &ToolWindow{
		Window: hswindow.New(title, x, y),
	}
}
