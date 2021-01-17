package hsdialog

import (
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"
)

type Dialog struct {
	*giu.PopupModalWidget
	title   string
	Visible bool
}

func New(title string) *Dialog {
	return &Dialog{
		PopupModalWidget: giu.PopupModal(title).Flags(imgui.WindowFlagsNoResize + imgui.WindowFlagsAlwaysAutoResize),
		title:            title,
	}
}

func (d *Dialog) ToggleVisibility() {
	d.Visible = !d.Visible
}

func (d *Dialog) Show() {
	d.Visible = true
}

func (d *Dialog) Render() {
	d.PopupModalWidget.Build()
}

func (d *Dialog) RegisterKeyboardShortcuts() {

}

func (d *Dialog) IsVisible() bool {
	return d.Visible
}

func (d *Dialog) Cleanup() {
	d.Visible = false
}
