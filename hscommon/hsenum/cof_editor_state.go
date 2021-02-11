package hsenum

type COFEditorState int

const (
	COFEditorStateViewer COFEditorState = iota
	COFEditorStateAddLayer
	COFEditorStateAddDirection
	COFEditorStateConfirm
)
