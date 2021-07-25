package hscommon

import "github.com/AllenDang/giu"

// MainMenuUpdater represents main menu updater
type MainMenuUpdater interface {
	// UpdateMainMenuLayout receives a pointer to the current layout of the menu bar at the top of the application,
	// allowing a struct implementing this interface to alter the menu bar.
	// This is generally used for adding a menu to the bar specific to the struct implementing this method, with options
	// that would be useful for that struct.
	UpdateMainMenuLayout(layout *giu.Layout)
}
