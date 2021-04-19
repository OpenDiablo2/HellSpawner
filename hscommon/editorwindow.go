package hscommon

import (
	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
)

// EditorWindow represents editor window
type EditorWindow interface {
	Renderable
	MainMenuUpdater

	// GetWindowTitle controls what the window title for this editor appears as
	GetWindowTitle() string
	// Show sets Visible to true
	Show()
	// IsVisible returns true if the editor has not been closed
	IsVisible() bool
	// SetVisible can be used to set Visible to false if the editor should be closed
	SetVisible(bool)
	// GetID returns a unique identifier for this editor window
	GetID() string
	// BringToFront brings this editor to the front of the application, giving it focus
	BringToFront()
	// State returns the current state of this editor, in a JSON-serializable struct
	State() hsstate.EditorState
	// Save writes any changes made in the editor to the file that is open in the editor.
	Save()

	Size(float32, float32) *giu.WindowWidget
}
