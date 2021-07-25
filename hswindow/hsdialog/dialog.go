// Package hsdialog contains project's dialogs
package hsdialog

import (
	"github.com/AllenDang/giu"
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

// RegisterKeyboardShortcuts registers a shortcuts for popup dialog
func (d *Dialog) RegisterKeyboardShortcuts(_ ...giu.WindowShortcut) {
	// https://github.com/OpenDiablo2/HellSpawner/issues/327
}

// KeyboardShortcuts returns a list of shortcuts
func (d *Dialog) KeyboardShortcuts() []giu.WindowShortcut {
	return []giu.WindowShortcut{}
}

// IsVisible returns true if dialog is visible
func (d *Dialog) IsVisible() bool {
	return d.Visible
}

// Cleanup hides dialog
func (d *Dialog) Cleanup() {
	d.Visible = false
}
