package hseditor

import "github.com/OpenDiablo2/HellSpawner/hswindow"

type Editor struct {
	hswindow.Window

	id      string
	ToFront bool
}

func (e *Editor) IsVisible() bool {
	return e.Visible
}

func (e *Editor) SetId(id string) {
	e.id = id
}

func (e *Editor) GetId() string {
	return e.id
}

func (e *Editor) BringToFront() {
	e.ToFront = true
}
