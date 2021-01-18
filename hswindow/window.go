package hswindow

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
)

type Window struct {
	*giu.WindowWidget
	Visible bool
}

func New(title string, x, y float32) *Window {
	return &Window{
		WindowWidget: giu.Window(title).Pos(x, y),
	}
}

func (t *Window) State() hsstate.WindowState {
	x, y := t.CurrentPosition()

	return hsstate.WindowState{
		Visible: t.Visible,
		PosX:    x,
		PosY:    y,
	}
}

func (t *Window) ToggleVisibility() {
	t.Visible = !t.Visible
}

func (t *Window) Show() {
	t.Visible = true
}

func (t *Window) Build() {

}

func (t *Window) Render() {
	t.WindowWidget.Build()
}

func (t *Window) RegisterKeyboardShortcuts() {

}

func (t *Window) IsVisible() bool {
	return t.Visible
}

func (t *Window) SetVisible(visible bool) {
	t.Visible = visible
}

func (t *Window) Cleanup() {
	t.Visible = false
}
