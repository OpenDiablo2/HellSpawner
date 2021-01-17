package hscommon

import "github.com/ianling/giu"

type MainMenuUpdater interface {
	UpdateMainMenuLayout(layout *giu.Layout)
}
