package hsstate

// EditorState holds information about the state of an open editor
type EditorState struct {
	WindowState
	Path    []byte `json:"path"` // this gets exported as raw JSON to prevent an import loop
	Encoded []byte `json:"state"`
}
