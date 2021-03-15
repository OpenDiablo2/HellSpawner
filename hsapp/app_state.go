package hsapp

import (
	"encoding/json"
	"log"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"
)

// State creates a new app state
func (a *App) State() hsstate.AppState {
	appState := hsstate.AppState{
		ProjectPath:   a.project.GetProjectFilePath(),
		EditorWindows: []hsstate.EditorState{},
		ToolWindows:   []hsstate.ToolWindowState{},
	}

	for _, editor := range a.editors {
		appState.EditorWindows = append(appState.EditorWindows, editor.State())
	}

	appState.ToolWindows = append(appState.ToolWindows, a.mpqExplorer.State(), a.projectExplorer.State())

	return appState
}

// RestoreAppState restores an app state
func (a *App) RestoreAppState(state hsstate.AppState) {
	for _, toolState := range state.ToolWindows {
		var tool hscommon.ToolWindow

		switch toolState.Type {
		case hsstate.ToolWindowTypeConsole:
			tool = a.console
		case hsstate.ToolWindowTypeMPQExplorer:
			tool = a.mpqExplorer
		case hsstate.ToolWindowTypeProjectExplorer:
			tool = a.projectExplorer
		default:
			continue
		}

		tool.Pos(toolState.PosX, toolState.PosY)
		tool.SetVisible(toolState.Visible)
		tool.Size(toolState.Width, toolState.Height)
	}

	for _, editorState := range state.EditorWindows {
		editorState := editorState

		var path hscommon.PathEntry

		err := json.Unmarshal(editorState.Path, &path)
		if err != nil {
			log.Print("failed to restore editor: ", err)
			continue
		}

		go a.createEditor(&path, editorState.PosX, editorState.PosY, editorState.Width, editorState.Height)
	}
}
