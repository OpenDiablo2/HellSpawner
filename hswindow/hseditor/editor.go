package hseditor

import (
	"encoding/json"
	"log"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
	"github.com/OpenDiablo2/HellSpawner/hswindow"
)

type Editor struct {
	*hswindow.Window
	Path *hscommon.PathEntry
}

func New(path *hscommon.PathEntry, x, y float32) *Editor {
	return &Editor{
		Window: hswindow.New(generateWindowTitle(path), x, y),
		Path:   path,
	}
}

func (e *Editor) State() hsstate.EditorState {
	path, err := json.Marshal(e.Path)
	if err != nil {
		log.Print("failed to marshal editor path to JSON: ", err)
	}

	return hsstate.EditorState{
		WindowState: e.Window.State(),
		Path:        path,
	}
}

func (e *Editor) GetWindowTitle() string {
	return generateWindowTitle(e.Path)
}

func (e *Editor) GetId() string {
	return e.Path.GetUniqueId()
}

func generateWindowTitle(path *hscommon.PathEntry) string {
	return path.Name + "##" + path.GetUniqueId()
}
