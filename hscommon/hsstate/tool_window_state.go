package hsstate

// ToolWindowType represents type of tool window
type ToolWindowType string

// ToolWindows types
const (
	ToolWindowTypeMPQExplorer     = ToolWindowType("MPQ Explorer")
	ToolWindowTypeProjectExplorer = ToolWindowType("Project Explorer")
	ToolWindowTypeConsole         = ToolWindowType("Console")
)

// ToolWindowState holds information about tool windows (e.g. MPQ Explorer)
type ToolWindowState struct {
	Type ToolWindowType `json:"type"`
	WindowState
}
