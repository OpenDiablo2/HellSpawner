package animationwidget

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/ianling/giu"
	gim "github.com/ozankasikci/go-image-merge"
)

const maxAlpha = uint8(255)

const (
	dc6WidgetViewer widgetMode = iota
	dc6WidgetTiledView
)

type Dc6WidgetState struct {
	mode widgetMode
	viewerState
	tiledState
	WidgetState
}

func (w *Dc6WidgetState) Dispose() {
	w.viewerState.Dispose()
	w.mode = dc6WidgetViewer
	w.widgetDispose()
}

func (w *Dc6WidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	w.WidgetState.Encode(sw)

	sw.PushInt32(int32(w.mode))

	sw.PushInt32(w.controls.direction)
	sw.PushInt32(w.controls.frame)
	sw.PushInt32(w.controls.scale)

	sw.PushInt32(w.width)
	sw.PushInt32(w.height)

	return sw.GetBytes()
}

func (w *Dc6WidgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

	w.WidgetState.Decode(sr)

	mode, err := sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	w.mode = widgetMode(mode)

	w.controls.direction, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	w.controls.frame, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	w.controls.scale, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	w.width, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	w.height, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	// update ticker
	w.ticker.Reset(time.Second * time.Duration(w.tickTime) / miliseconds)
}

func (w *Dc6Widget) initState() {
	// Prevent multiple invocation to LoadImage.
	newState := &Dc6WidgetState{
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
		WidgetState: WidgetState{
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

func (w *Dc6Widget) getState() *Dc6WidgetState {
	var state *Dc6WidgetState

	s := giu.Context.GetState(w.widget.getStateID())

	if s != nil {
		state = s.(*Dc6WidgetState)
	} else {
		w.initState()
		state = w.getState()
	}

	return state
}

func (w *Dc6Widget) setState(s giu.Disposable) {
	giu.Context.SetState(w.widget.getStateID(), s)
}

func (w *Dc6Widget) runPlayer(state *Dc6WidgetState) {
	for range state.ticker.C {
		if !state.isPlaying {
			continue
		}

		numFrames := int32(w.dc6.FramesPerDirection - 1)
		isLastFrame := state.controls.frame == numFrames

		// update play direction
		switch state.playMode {
		case playModeForward:
			state.isForward = true
		case playModeBackward:
			state.isForward = false
		case playModePingPong:
			if isLastFrame || state.controls.frame == 0 {
				state.isForward = !state.isForward
			}
		}

		// now update the frame number
		if state.isForward {
			state.controls.frame++
		} else {
			state.controls.frame--
		}

		state.controls.frame = int32(hsutil.Wrap(int(state.controls.frame), int(w.dc6.FramesPerDirection)))

		// next, check for stopping/repeat
		isStoppingFrame := (state.controls.frame == 0) || (state.controls.frame == numFrames)

		if isStoppingFrame && !state.repeat {
			state.isPlaying = false
		}
	}
}

func (w *Dc6Widget) recalculateTiledViewWidth(state *Dc6WidgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.width = int32(w.dc6.FramesPerDirection) / state.height
	w.createImage(state)
}

func (w *Dc6Widget) recalculateTiledViewHeight(state *Dc6WidgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.height = int32(w.dc6.FramesPerDirection) / state.width
	w.createImage(state)
}

func (w *Dc6Widget) createImage(state *Dc6WidgetState) {
	firstFrame := state.controls.direction * int32(w.dc6.FramesPerDirection)

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


