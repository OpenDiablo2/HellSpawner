package hstoolwindow

import "github.com/OpenDiablo2/HellSpawner/hswindow"

type ToolWindow struct {
	*hswindow.Window
}

func New(title string) *ToolWindow {
	return &ToolWindow{
		Window: hswindow.New(title),
	}
}
