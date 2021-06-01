// Package hsdialog contains project's dialogs
package hsdialog

import (
	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hsinput"
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
		PopupModalWidget: giu.PopupModal(title),
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
