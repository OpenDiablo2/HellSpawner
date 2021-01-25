package hseditor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
	"github.com/OpenDiablo2/HellSpawner/hswindow"
)

type Editor struct {
	*hswindow.Window
	Path    *hscommon.PathEntry
	Project *hsproject.Project
}

func New(path *hscommon.PathEntry, x, y float32, project *hsproject.Project) *Editor {
	return &Editor{
		Window:  hswindow.New(generateWindowTitle(path), x, y),
		Path:    path,
		Project: project,
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
	return e.Path.GetUniqueID()
}

func (e *Editor) Save(editor Saveable) {
	if e.Path.Source != hscommon.PathEntrySourceProject {
		// saving to MPQ not yet supported
		return
	}

	if editor, isSaveable := editor.(Saveable); isSaveable {
		saveData := editor.GenerateSaveData()
		if saveData == nil {
			return
		}

		existingFileData, err := e.Path.GetFileBytes()
		if err != nil {
			fmt.Println("failed to read file before saving: ", err)
			return
		}

		if bytes.Equal(saveData, existingFileData) {
			// nothing to save
			return
		}

		err = e.Path.WriteFile(saveData)
		if err != nil {
			fmt.Println("failed to save file: ", err)
			return
		}
	} else {
		return
	}
}

func (e *Editor) HasChanges(editor Saveable) bool {
	if e.Path.Source != hscommon.PathEntrySourceProject {
		// saving to MPQ not yet supported
		return false
	}

	if editor, isSaveable := editor.(Saveable); isSaveable {
		newData := editor.GenerateSaveData()
		if newData != nil {
			oldData, err := e.Path.GetFileBytes()
			if err == nil {
				return !bytes.Equal(oldData, newData)
			}
		}
	}

	// err on the side of caution; if any errors occurred, just say nothing has changed so no changes get saved
	return false
}

func (e *Editor) Cleanup() {
	e.Window.Cleanup()
}

func generateWindowTitle(path *hscommon.PathEntry) string {
	return path.Name + "##" + path.GetUniqueID()
}
