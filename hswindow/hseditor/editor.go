package hseditor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
	"github.com/OpenDiablo2/HellSpawner/hswindow"
)

// Editor represents an editor
type Editor struct {
	*hswindow.Window
	Path    *hscommon.PathEntry
	Project *hsproject.Project
}

// New creates a new editor
func New(path *hscommon.PathEntry, x, y float32, project *hsproject.Project) *Editor {
	return &Editor{
		Window:  hswindow.New(generateWindowTitle(path), x, y),
		Path:    path,
		Project: project,
	}
}

// State returns editors state
func (e *Editor) State() hsstate.EditorState {
	path, err := json.Marshal(e.Path)
	if err != nil {
		log.Print("failed to marshal editor path to JSON: ", err)
	}

	result := hsstate.EditorState{
		WindowState: e.Window.State(),
		Path:        path,
		Encoded:     e.EncodeState(),
	}

	return result
}

// GetWindowTitle returns window title
func (e *Editor) GetWindowTitle() string {
	return generateWindowTitle(e.Path)
}

// GetID returns editors ID
func (e *Editor) GetID() string {
	return e.Path.GetUniqueID()
}

// Save saves an editor
func (e *Editor) Save(editor Saveable) {
	if e.Path.Source != hscommon.PathEntrySourceProject {
		// saving to MPQ not yet supported
		return
	}

	if _, isSaveable := editor.(Saveable); isSaveable {
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

// HasChanges returns true if editor has changed data
func (e *Editor) HasChanges(editor Saveable) bool {
	if e.Path.Source != hscommon.PathEntrySourceProject {
		// saving to MPQ not yet supported
		return false
	}

	if _, isSaveable := editor.(Saveable); isSaveable {
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

// Cleanup cides an editor
func (e *Editor) Cleanup() {
	e.Window.Cleanup()
}

func generateWindowTitle(path *hscommon.PathEntry) string {
	return path.Name + "##" + path.GetUniqueID()
}

// EncodeState returns widget's state (unique for each editor type) in byte slice format
func (e *Editor) EncodeState() []byte {
	id := fmt.Sprintf("widget_%s", e.Path.GetUniqueID())

	if s := giu.Context.GetState(id); s != nil {
		state, ok := s.(interface {
			Dispose()
			Encode() []byte
		})
		if !ok {
			log.Printf("editor on path %s doesn't support saving state", e.Path.GetUniqueID())
			return nil
		}

		return state.Encode()
	}

	return nil
}
