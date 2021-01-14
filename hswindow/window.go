package hswindow

import "github.com/AllenDang/giu"

type Window struct {
	Widget  *giu.WindowWidget
	Visible bool
}

func (t *Window) ToggleVisibility() {
	t.Visible = !t.Visible
}

func (t *Window) Show() {
	t.Visible = true
}
