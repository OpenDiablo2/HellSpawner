package hscommon

import "github.com/AllenDang/giu"

type MainMenuUpdater interface {
	UpdateMainMenuLayout(layout *giu.Layout)
}
