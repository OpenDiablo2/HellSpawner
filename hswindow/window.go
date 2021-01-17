package hswindow

import "github.com/ianling/giu"

type Window struct {
	*giu.WindowWidget
	Visible bool
}

func New(title string, x, y float32) *Window {
	return &Window{
		WindowWidget: giu.Window(title).Pos(x, y),
	}
}

func (t *Window) ToggleVisibility() {
	t.Visible = !t.Visible
}

func (t *Window) Show() {
	t.Visible = true
}

func (t *Window) Render() {
	t.WindowWidget.Build()
}

func (t *Window) RegisterKeyboardShortcuts() {

}

func (t *Window) IsVisible() bool {
	return t.Visible
}

func (t *Window) Cleanup() {
	t.Visible = false
}
