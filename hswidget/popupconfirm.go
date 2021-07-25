package hswidget

import (
	"log"

	"github.com/AllenDang/giu"
)

const (
	yesNoButtonW, yesNoButtonH = 40, 25
)

// PopUpConfirmDialog represents a pop up dialog
type PopUpConfirmDialog struct {
	header  string
	message string
	id      string
	yCB     func()
	nCB     func()
}

// NewPopUpConfirmDialog creates a new pop up dialog (with yes-no options)
func NewPopUpConfirmDialog(id, header, message string, yCB, nCB func()) *PopUpConfirmDialog {
	result := &PopUpConfirmDialog{
		header:  header,
		message: message,
		id:      id,
		yCB:     yCB,
		nCB:     nCB,
	}

	return result
}

// Build builds a pop up dialog
func (p *PopUpConfirmDialog) Build() {
	if p.header == "" {
		log.Print("Header is empty; please ensure, if you're building appropriate dialog")
	}

	open := true
	giu.Layout{
		giu.PopupModal(p.header + "##" + p.id).IsOpen(&open).Layout(giu.Layout{
			giu.Label(p.message),
			giu.Separator(),
			giu.Row(
				giu.Button("YES##"+p.id+"ConfirmDialog").Size(yesNoButtonW, yesNoButtonH).OnClick(p.yCB),
				giu.Button("NO##"+p.id+"confirmDialog").Size(yesNoButtonW, yesNoButtonH).OnClick(p.nCB),
			),
		}),
	}.Build()
}
