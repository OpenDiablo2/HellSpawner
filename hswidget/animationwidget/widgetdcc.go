package animationwidget

import (
	"fmt"
	"log"
	"time"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

const (
	imageW, imageH = 32, 32
)

type DccWidget struct {
	widget *Widget
	dcc    *d2dcc.DCC
}

func (p *DccWidget) exportGif(state *DccWidgetState) error {
	fpd := int32(p.dcc.FramesPerDirection)
	firstFrame := state.controls.direction * fpd
	images := state.images[firstFrame : firstFrame+fpd]

	err := hsutil.ExportToGif(images, state.tickTime)
	if err != nil {
		return fmt.Errorf("error creating gif file: %w", err)
	}

	return nil
}

func (w *DccWidget) makePlayerLayout(state *DccWidgetState) giu.Layout {
	playModeList := make([]string, 0)
	for i := playModeForward; i <= playModePingPong; i++ {
		playModeList = append(playModeList, i.String())
	}

	pm := int32(state.playMode)

	return giu.Layout{
		giu.Row(
			giu.Checkbox("Loop##"+w.widget.id+"PlayRepeat", &state.repeat),
			giu.Combo("##"+w.widget.id+"PlayModeList", playModeList[state.playMode], playModeList, &pm).OnChange(func() {
				state.playMode = animationPlayMode(pm)
			}).Size(comboW),
			giu.InputInt("Tick time##"+w.widget.id+"PlayTickTime", &state.tickTime).Size(inputIntW).OnChange(func() {
				state.ticker.Reset(time.Second * time.Duration(state.tickTime) / miliseconds)
			}),
			hswidget.PlayPauseButton("##"+w.widget.id+"PlayPauseAnimation", &state.isPlaying, w.widget.textureLoader).
				Size(playPauseButtonSize, playPauseButtonSize),
			giu.Button("Export GIF##"+w.widget.id+"exportGif").OnClick(func() {
				err := w.exportGif(state)
				if err != nil {
					dialog.Message(err.Error()).Error()
				}
			}),
		),
	}
}

// Build build a widget
func (w *DccWidget) Build() {
	viewerState := w.getState()

	imageScale := uint32(viewerState.controls.scale)
	dirIdx := int(viewerState.controls.direction)
	frameIdx := viewerState.controls.frame

	textureIdx := dirIdx*len(w.dcc.Directions[dirIdx].Frames) + int(frameIdx)

	if imageScale < 1 {
		imageScale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	var widget *giu.ImageWidget
	if viewerState.textures == nil || len(viewerState.textures) <= int(frameIdx) || viewerState.textures[frameIdx] == nil {
		widget = giu.Image(nil).Size(imageW, imageH)
	} else {
		bw := w.dcc.Directions[dirIdx].Box.Width
		bh := w.dcc.Directions[dirIdx].Box.Height
		w := float32(uint32(bw) * imageScale)
		h := float32(uint32(bh) * imageScale)
		widget = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
	}

	giu.Layout{
		giu.Row(
			giu.Label(fmt.Sprintf("Signature: %v", w.dcc.Signature)),
			giu.Label(fmt.Sprintf("Version: %v", w.dcc.Version)),
		),
		giu.Row(
			giu.Label(fmt.Sprintf("Directions: %v", w.dcc.NumberOfDirections)),
			giu.Label(fmt.Sprintf("Frames per Direction: %v", w.dcc.FramesPerDirection)),
		),
		giu.Custom(func() {
			imgui.BeginGroup()
			if w.dcc.NumberOfDirections > 1 {
				imgui.SliderInt("Direction", &viewerState.controls.direction, 0, int32(w.dcc.NumberOfDirections-1))
			}

			if w.dcc.FramesPerDirection > 1 {
				imgui.SliderInt("Frames", &viewerState.controls.frame, 0, int32(w.dcc.FramesPerDirection-1))
			}

			const minScale, maxScale = 1, 8

			imgui.SliderInt("Scale", &viewerState.controls.scale, minScale, maxScale)

			imgui.EndGroup()
		}),
		giu.Separator(),
		w.makePlayerLayout(viewerState),
		giu.Separator(),
		widget,
	}.Build()
}

func CreateDccWidget(tl hscommon.TextureLoader, state []byte, palette *[256]d2interface.Color, id string, dcc *d2dcc.DCC) giu.Widget {
	widget := &Widget{
		id: id,
		palette: palette,
		textureLoader: tl,
	}

	dccWidget := &DccWidget{
		widget: widget,
		dcc: dcc,
	}

	if giu.Context.GetState(dccWidget.widget.getStateID()) == nil && state != nil {
		s := dccWidget.getState()
		s.Decode(state)
		dccWidget.setState(s)
	}

	return dccWidget
}
