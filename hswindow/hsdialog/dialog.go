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
	result := &Dialog{
		PopupModalWidget: giu.PopupModal(title),
		title:            title,
	}

	return result
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

// IsOpen wrapps popupmodalwidget.isOpen
func (d *Dialog) IsOpen(isOpen *bool) *Dialog {
	d.PopupModalWidget.IsOpen(isOpen)
	return d
}

// Layout wrapps d.PopUpModalWidget.Layout
func (d *Dialog) Layout(widgets ...giu.Widget) *Dialog {
	d.PopupModalWidget.Layout(widgets...)
	return d
}

// Build builds the widget
func (d *Dialog) Build() {
	giu.OpenPopup(d.title)
	d.PopupModalWidget.Build()
}
