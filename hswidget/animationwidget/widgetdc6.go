package animationwidget

import (
	"fmt"
	"log"
	"time"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

const (
	playPauseButtonSize = 15
	buttonW, buttonH    = 200, 30
	inputIntW           = 30
	comboW              = 125
)

type Dc6Widget struct {
	widget *Widget
	dc6    *d2dc6.DC6
}

func (w *Dc6Widget) exportGif(state *Dc6WidgetState) error {
	fpd := int32(w.dc6.FramesPerDirection)
	firstFrame := state.controls.direction * fpd
	images := state.images[firstFrame : firstFrame+fpd]

	err := hsutil.ExportToGif(images, state.tickTime)
	if err != nil {
		return fmt.Errorf("error creating gif file: %w", err)
	}

	return nil
}

func (w *Dc6Widget) makePlayerLayout(state *Dc6WidgetState) giu.Layout {
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

func (w *Dc6Widget) makeViewerLayout() giu.Layout {
	viewerState := w.getState()

	imageScale := uint32(viewerState.controls.scale)
	curFrameIndex := int(viewerState.controls.frame) + (int(viewerState.controls.direction) * int(w.dc6.FramesPerDirection))
	dirIdx := int(viewerState.controls.direction)

	textureIdx := dirIdx*int(w.dc6.FramesPerDirection) + int(viewerState.controls.frame)

	if imageScale < 1 {
		imageScale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	width := float32(w.dc6.Frames[curFrameIndex].Width * imageScale)
	height := float32(w.dc6.Frames[curFrameIndex].Height * imageScale)

	var widget *giu.ImageWidget
	if viewerState.textures == nil || len(viewerState.textures) <= int(viewerState.controls.frame) ||
		viewerState.textures[curFrameIndex] == nil {
		widget = giu.Image(nil).Size(width, height)
	} else {
		widget = giu.Image(viewerState.textures[textureIdx]).Size(width, height)
	}

	return giu.Layout{
		giu.Label(fmt.Sprintf(
			"Version: %v\t Flags: %b\t Encoding: %v\t",
			w.dc6.Version,
			int64(w.dc6.Flags),
			w.dc6.Encoding,
		)),
		giu.Label(fmt.Sprintf("Directions: %v\tFrames per Direction: %v", w.dc6.Directions, w.dc6.FramesPerDirection)),
		giu.Custom(func() {
			imgui.BeginGroup()
			if w.dc6.Directions > 1 {
				imgui.SliderInt("Direction", &viewerState.controls.direction, 0, int32(w.dc6.Directions-1))
			}

			if w.dc6.FramesPerDirection > 1 {
				imgui.SliderInt("Frames", &viewerState.controls.frame, 0, int32(w.dc6.FramesPerDirection-1))
			}

			const minScale, maxScale = 1, 8

			imgui.SliderInt("Scale", &viewerState.controls.scale, minScale, maxScale)

			imgui.EndGroup()
		}),
		giu.Separator(),
		w.makePlayerLayout(viewerState),
		giu.Separator(),
		widget,
		giu.Separator(),
		giu.Button("Tiled View##"+w.widget.id+"tiledViewButton").Size(buttonW, buttonH).OnClick(func() {
			viewerState.mode = dc6WidgetTiledView
			w.createImage(viewerState)
		}),
	}
}

func (w *Dc6Widget) makeTiledViewLayout(state *Dc6WidgetState) giu.Layout {
	return giu.Layout{
		giu.Row(
			giu.Label("Tiled view:"),
			giu.InputInt("Width##"+w.widget.id+"tiledWidth", &state.width).Size(inputIntW).OnChange(func() {
				w.recalculateTiledViewHeight(state)
			}),
			giu.InputInt("Height##"+w.widget.id+"tiledHeight", &state.height).Size(inputIntW).OnChange(func() {
				w.recalculateTiledViewWidth(state)
			}),
		),
		giu.Image(state.tiled).Size(float32(state.imgw), float32(state.imgh)),
		giu.Button("Back##"+w.widget.id+"tiledBack").Size(buttonW, buttonH).OnClick(func() {
			state.mode = dc6WidgetViewer
		}),
	}
}

// Build builds a widget
func (w *Dc6Widget) Build() {
	state := w.getState()

	switch state.mode {
	case dc6WidgetViewer:
		w.makeViewerLayout().Build()
	case dc6WidgetTiledView:
		w.makeTiledViewLayout(state).Build()
	}
}

func CreateDc6Widget(state []byte, palette *[256]d2interface.Color, textureLoader hscommon.TextureLoader, id string, dc6 *d2dc6.DC6) giu.Widget {
	widget := CreateWidget(palette, textureLoader, id)

	dc6Widget := &Dc6Widget{
		widget: widget,
		dc6:    dc6,
	}

	if giu.Context.GetState(dc6Widget.widget.getStateID()) == nil && state != nil {
		s := dc6Widget.getState()
		s.Decode(state)

		if s.mode == dc6WidgetTiledView {
			dc6Widget.createImage(s)
		}

		dc6Widget.setState(s)
	}

	return dc6Widget
}
