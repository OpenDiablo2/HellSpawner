package hscommon

type EditorFocuser interface {
	FocusOn(editor EditorWindow)
}

type FocusController interface {
	Control(focuser EditorFocuser)
}
