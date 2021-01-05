package hsfonteditor

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

type FontEditor struct {
	hseditor.Editor

	pathEntry *hscommon.PathEntry
}

func (e *FontEditor) Cleanup() {

}

func (e *FontEditor) GetWindowTitle() string {
	return e.pathEntry.Name + "##" + e.GetId()
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	result := &FontEditor{
		pathEntry: pathEntry,
	}

	return result, nil
}

func (e *FontEditor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	g.Window(e.GetWindowTitle()).IsOpen(&e.Visible).Pos(50, 50).Size(400, 300).Layout(g.Layout{})
}
