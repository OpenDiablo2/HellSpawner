package hsenum

// COFEditorState represents cof editor's states
type COFEditorState int

// cof editor's states
const (
	COFEditorStateViewer COFEditorState = iota
	COFEditorStateAddLayer
	COFEditorStateAddDirection
	COFEditorStateConfirm
)
