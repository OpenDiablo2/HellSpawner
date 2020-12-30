package hswindow

type Window struct {
	Visible bool
}

func (t *Window) ToggleVisibility() {
	t.Visible = !t.Visible
}

func (t *Window) Show() {
	t.Visible = true
}
