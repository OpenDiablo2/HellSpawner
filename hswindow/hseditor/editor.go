package hseditor

import (
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow"
)

type Editor struct {
	hswindow.Window
	Path    *hscommon.PathEntry
	focuser hscommon.EditorFocuser

	ToFront bool
	Focused bool
}

func (e *Editor) Control(focuser hscommon.EditorFocuser) {
	e.focuser = focuser
}

func (e *Editor) IsVisible() bool {
	return e.Visible
}

func (e *Editor) IsFocused() bool {
	return e.Focused
}

func (e *Editor) GetId() string {
	return e.Path.GetUniqueId()
}

func (e *Editor) GetWindowTitle() string {
	return e.Path.Name + "##" + e.GetId()
}

func (e *Editor) BringToFront() {
	e.ToFront = true
}

func (e *Editor) Cleanup() {
	e.Visible = false
}
