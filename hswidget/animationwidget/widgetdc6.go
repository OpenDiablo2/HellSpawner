package animationwidget

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"
	gim "github.com/ozankasikci/go-image-merge"
)

const (
	playPauseButtonSize = 15
	buttonW, buttonH    = 200, 30
	inputIntW           = 30
	comboW              = 125
)

type dc6Widget struct {
	*widget
	dc6 *d2dc6.DC6
}

func (w *dc6Widget) getDcImage() dcImage {
	return w.dc6
}

func (w *dc6Widget) makeViewerLayout() giu.Layout {
	viewerState := w.getState()

	imageScale := uint32(viewerState.scale)
	curFrameIndex := int(viewerState.frame) + (int(viewerState.direction) * int(w.dc6.FramesPerDirection))
	dirIdx := int(viewerState.direction)

	textureIdx := dirIdx*int(w.dc6.FramesPerDirection) + int(viewerState.frame)

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
	if viewerState.textures == nil || len(viewerState.textures) <= int(viewerState.frame) ||
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
				imgui.SliderInt("Direction", &viewerState.direction, 0, int32(w.dc6.Directions-1))
			}

			if w.dc6.FramesPerDirection > 1 {
				imgui.SliderInt("Frames", &viewerState.frame, 0, int32(w.dc6.FramesPerDirection-1))
			}

			const minScale, maxScale = 1, 8

			imgui.SliderInt("Scale", &viewerState.scale, minScale, maxScale)

			imgui.EndGroup()
		}),
		giu.Separator(),
		makePlayerLayout(w, viewerState),
		giu.Separator(),
		widget,
		giu.Separator(),
		giu.Button("Tiled View##"+w.widget.id+"tiledViewButton").Size(buttonW, buttonH).OnClick(func() {
			viewerState.mode = dc6WidgetTiledView
			w.createImage(viewerState)
		}),
	}
}

func (w *dc6Widget) makeTiledViewLayout(state *dc6WidgetState) giu.Layout {
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
func (w *dc6Widget) Build() {
	state := w.getState()

	switch state.mode {
	case dc6WidgetViewer:
		w.makeViewerLayout().Build()
	case dc6WidgetTiledView:
		w.makeTiledViewLayout(state).Build()
	}
}

func createDc6Widget(state []byte, widget *widget, dc6 *d2dc6.DC6) giu.Widget {
	dc6Widget := &dc6Widget{
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

func (w *dc6Widget) initState() {
	// Prevent multiple invocation to LoadImage.
	newState := &dc6WidgetState{
		mode: dc6WidgetViewer,
		viewerState: viewerState{
			lastFrame:          -1,
			lastDirection:      -1,
			framesPerDirection: w.dc6.FramesPerDirection,
		},
		tiledState: tiledState{
			width:  int32(w.dc6.FramesPerDirection),
			height: 1,
		},
		widgetState: widgetState{
			isPlaying: false,
			repeat:    false,
			tickTime:  defaultTickTime,
			playMode:  playModeForward,
		},
	}

	newState.ticker = time.NewTicker(time.Second * time.Duration(newState.tickTime) / miliseconds)

	go w.runPlayer(newState)

	totalFrames := int(w.dc6.Directions * w.dc6.FramesPerDirection)
	newState.images = make([]*image.RGBA, totalFrames)

	for frameIndex := 0; frameIndex < int(w.dc6.Directions*w.dc6.FramesPerDirection); frameIndex++ {
		newState.images[frameIndex] = image.NewRGBA(image.Rect(0, 0, int(w.dc6.Frames[frameIndex].Width), int(w.dc6.Frames[frameIndex].Height)))
		decodedFrame := w.dc6.DecodeFrame(frameIndex)

		for y := 0; y < int(w.dc6.Frames[frameIndex].Height); y++ {
			for x := 0; x < int(w.dc6.Frames[frameIndex].Width); x++ {
				idx := x + (y * int(w.dc6.Frames[frameIndex].Width))
				val := decodedFrame[idx]

				alpha := maxAlpha

				if val == 0 {
					alpha = 0
				}

				var r, g, b uint8

				if w.widget.palette != nil {
					col := w.widget.palette[val]
					r, g, b = col.R(), col.G(), col.B()
				} else {
					r, g, b = val, val, val
				}

				newState.images[frameIndex].Set(
					x, y,
					color.RGBA{
						R: r,
						G: g,
						B: b,
						A: alpha,
					},
				)
			}
		}
	}

	w.setState(newState)

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			w.widget.textureLoader.CreateTextureFromARGB(newState.images[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := w.getState()
		s.textures = textures
		w.setState(s)
	}()
}

func (w *dc6Widget) getState() *dc6WidgetState {
	var state *dc6WidgetState

	s := giu.Context.GetState(w.widget.getStateID())

	if s != nil {
		state = s.(*dc6WidgetState)
	} else {
		w.initState()
		state = w.getState()
	}

	return state
}

func (w *dc6Widget) setState(s giu.Disposable) {
	giu.Context.SetState(w.widget.getStateID(), s)
}

func (w *dc6Widget) runPlayer(state *dc6WidgetState) {
	for range state.ticker.C {
		if !state.isPlaying {
			continue
		}

		numFrames := int32(w.dc6.FramesPerDirection - 1)
		isLastFrame := state.frame == numFrames

		// update play direction
		switch state.playMode {
		case playModeForward:
			state.isForward = true
		case playModeBackward:
			state.isForward = false
		case playModePingPong:
			if isLastFrame || state.frame == 0 {
				state.isForward = !state.isForward
			}
		}

		// now update the frame number
		if state.isForward {
			state.frame++
		} else {
			state.frame--
		}

		state.frame = int32(hsutil.Wrap(int(state.frame), int(w.dc6.FramesPerDirection)))

		// next, check for stopping/repeat
		isStoppingFrame := (state.frame == 0) || (state.frame == numFrames)

		if isStoppingFrame && !state.repeat {
			state.isPlaying = false
		}
	}
}

func (w *dc6Widget) recalculateTiledViewWidth(state *dc6WidgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.width = int32(w.dc6.FramesPerDirection) / state.height
	w.createImage(state)
}

func (w *dc6Widget) recalculateTiledViewHeight(state *dc6WidgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.height = int32(w.dc6.FramesPerDirection) / state.width
	w.createImage(state)
}

func (w *dc6Widget) createImage(state *dc6WidgetState) {
	firstFrame := state.direction * int32(w.dc6.FramesPerDirection)

	grids := make([]*gim.Grid, 0)

	for j := int32(0); j < state.height*state.width; j++ {
		grids = append(grids, &gim.Grid{Image: image.Image(state.images[firstFrame+j])})
	}

	newimg, err := gim.New(grids, int(state.width), int(state.height)).Merge()
	if err != nil {
		log.Printf("merging image error: %v", err)
		return
	}

	w.widget.textureLoader.CreateTextureFromARGB(newimg, func(t *giu.Texture) {
		state.tiled = t
	})

	state.imgw, state.imgh = newimg.Bounds().Dx(), newimg.Bounds().Dy()
}
