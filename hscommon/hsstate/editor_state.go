package hsstate

// EditorState holds information about all of the open editors.
type EditorState struct {
	WindowState
	Path []byte `json:"path"` // this gets exported as raw JSON to prevent an import loop
}
