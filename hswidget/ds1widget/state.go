package ds1widget

import (
	"fmt"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

type widgetMode int32

const (
	widgetModeViewer widgetMode = iota
	widgetModeAddFile
	widgetModeAddObject
	widgetModeAddPath
	widgetModeConfirm
)

type ds1Controls struct {
	TileX, TileY int32
	Object       int32
	Subgroup     int32
	Tile         struct {
		Floor, Wall, Shadow, Sub int32
	}
	noObjectsImageTexture *giu.Texture
}

// ds1AddObjectState represents state of new Object
type ds1AddObjectState struct {
	ObjType  int32
	ObjID    int32
	ObjX     int32
	ObjY     int32
	ObjFlags int32
}

// Dispose clears state
func (t *ds1AddObjectState) Dispose() {
	// noop
}

// ds1AddPathState contains data about new path
type ds1AddPathState struct {
	PathAction int32
	PathX      int32
	PathY      int32
}

// Dispose clears state
func (t *ds1AddPathState) Dispose() {
	// noop
}

// widgetState represents ds1 viewers state
type widgetState struct {
	*ds1Controls
	Mode           widgetMode
	confirmDialog  *hswidget.PopUpConfirmDialog
	NewFilePath    string
	addObjectState ds1AddObjectState
	addPathState   ds1AddPathState
}

// Dispose clears viewers state
func (is *widgetState) Dispose() {
	is.addObjectState.Dispose()
	is.addPathState.Dispose()
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
}

func (p *widget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *widget) initState() {
	state := &widgetState{
		ds1Controls: &ds1Controls{},
	}

	p.textureLoader.CreateTextureFromFile(hsassets.ImageShrug, func(t *giu.Texture) {
		state.ds1Controls.noObjectsImageTexture = t
	})

	p.setState(state)
}
