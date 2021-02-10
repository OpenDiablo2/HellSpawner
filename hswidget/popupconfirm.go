package hswidget

import (
	"github.com/ianling/giu"
)

type PopUpConfirmDialog struct {
	header  string
	message string
	id      string
	yCB     func()
	nCB     func()
}

func NewPopUpConfirmDialog(id, header, message string, yCB func(), nCB func()) *PopUpConfirmDialog {
	result := &PopUpConfirmDialog{
		header:  header,
		message: message,
		id:      id,
		yCB:     yCB,
		nCB:     nCB,
	}

	return result
}

func (p *PopUpConfirmDialog) Build() {
	open := true
	giu.Layout{
		giu.PopupModal(p.header).IsOpen(&open).Layout(giu.Layout{
			giu.Label(p.message),
			giu.Separator(),
			giu.Line(
				giu.Button("YES##"+p.id+"ConfirmDialog").Size(40, 25).OnClick(p.yCB),
				giu.Button("NO##"+p.id+"confirmDialog").Size(40, 25).OnClick(p.nCB),
			),
		}),
	}.Build()
}
