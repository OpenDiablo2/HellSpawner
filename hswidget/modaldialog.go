package hswidget

import (
	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
)

func ModalDialog(id string, visible *bool, layout g.Layout) {
	if *visible {
		imgui.OpenPopup(id)
	}

	if imgui.BeginPopupModalV(id, visible, imgui.WindowFlagsNoResize|imgui.WindowFlagsAlwaysAutoResize) {
		layout.Build()
		imgui.EndPopup()
	}
}
