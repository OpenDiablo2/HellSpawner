// Package hsdialog contains project's dialogs
package hsdialog

import (
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"
)

// Dialog represents HellSpawner's dialog
type Dialog struct {
	*giu.PopupModalWidget
	title   string
	Visible bool
}

// New creates a new dialog
func New(title string) *Dialog {
	return &Dialog{
		PopupModalWidget: giu.PopupModal(title).Flags(imgui.WindowFlagsNoResize + imgui.WindowFlagsAlwaysAutoResize),
		title:            title,
	}
}

// ToggleVisibility toggles dialog's visibility
func (d *Dialog) ToggleVisibility() {
	d.Visible = !d.Visible
}

// Show shows dialog
func (d *Dialog) Show() {
	d.Visible = true
}

// Render renders dialog
func (d *Dialog) Render() {
	d.PopupModalWidget.Build()
}

// RegisterKeyboardShortcuts registers a new shortcut
func (d *Dialog) RegisterKeyboardShortcuts(_ *hsinput.InputManager) {
	// noop
}

// IsVisible returns true if dialog is visible
func (d *Dialog) IsVisible() bool {
	return d.Visible
}

// Cleanup hides dialog
func (d *Dialog) Cleanup() {
	d.Visible = false
}
