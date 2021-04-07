package ds1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"
)

func (p *widget) addFloor(idx int32) {
	p.ds1.InsertFloor(int(idx), &d2ds1.Layer{})
}

func (p *widget) deleteFloor(idx int32) {
	// state := p.getState()
	// here p.ds1.DeleteFloor(state.object)
	p.ds1.DeleteFloor(int(idx))
}

func (p *widget) addWall(idx int32) {
	p.ds1.InsertWall(int(idx), &d2ds1.Layer{})
}

func (p *widget) deleteWall(idx int32) {
	// state := p.getState()
	// here p.ds1.DeleteFloor(state.object)
	p.ds1.DeleteWall(int(idx))
}
