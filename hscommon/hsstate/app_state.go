// Package hsstate contains structs that describe the state of the application.
// This allows us to save the state of the application to a file when exited,
// and then re-load the state when the application is opened again.
package hsstate

// AppState holds information related to the running state of HellSpawner.
type AppState struct {
	ProjectPath   string            `json:"projectPath"`
	EditorWindows []EditorState     `json:"editorWindows"`
	ToolWindows   []ToolWindowState `json:"toolWindows"`
}
