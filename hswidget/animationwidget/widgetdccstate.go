package animationwidget

import (
	"image"
	"image/color"
	"log"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/ianling/giu"
)

type DccWidgetState struct {
	controls *controlStructure
	WidgetState
}

func (s *DccWidgetState) Dispose() {
	s.widgetDispose()
}

func (s *DccWidgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	s.WidgetState.Encode(sw)

	sw.PushInt32(s.controls.direction)
	sw.PushInt32(s.controls.frame)
	sw.PushInt32(s.controls.scale)

	return sw.GetBytes()
}

func (s *DccWidgetState) Decode(data []byte) {
	var err error

	sr := d2datautils.CreateStreamReader(data)

	s.WidgetState.Decode(sr)

	s.controls.direction, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.controls.frame, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	s.controls.scale, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	// update ticker
	s.ticker.Reset(time.Second * time.Duration(s.tickTime) / miliseconds)
}

func (w *DccWidget) getState() *DccWidgetState {
	var state *DccWidgetState

	s := giu.Context.GetState(w.widget.getStateID())

	if s != nil {
		state = s.(*DccWidgetState)
	} else {
		w.initState()
		state = w.getState()
	}

	return state
}

func (w *DccWidget) setState(s giu.Disposable) {
	giu.Context.SetState(w.widget.getStateID(), s)
}

func (w *DccWidget) runPlayer(state *DccWidgetState) {
	for range state.ticker.C {
		if !state.isPlaying {
			continue
		}

		numFrames := int32(w.dcc.FramesPerDirection - 1)
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

		state.controls.frame = int32(hsutil.Wrap(int(state.controls.frame), w.dcc.FramesPerDirection))

		// next, check for stopping/repeat
		isStoppingFrame := (state.controls.frame == 0) || (state.controls.frame == numFrames)

		if isStoppingFrame && !state.repeat {
			state.isPlaying = false
		}
	}
}

func (w *DccWidget) makeImagePixel(val byte) color.RGBA {
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

	RGBAColor := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: alpha,
	}

	return RGBAColor
}

func (w *DccWidget) initState() {
	// Prevent multiple invocation to LoadImage.
	state := &DccWidgetState{
		WidgetState: WidgetState{
			isPlaying: false,
			repeat:    false,
			tickTime:  defaultTickTime,
			playMode:  playModeForward,
		},
	}

	state.ticker = time.NewTicker(time.Second * time.Duration(state.tickTime) / miliseconds)

	w.setState(state)

	go w.runPlayer(state)

	totalFrames := w.dcc.NumberOfDirections * w.dcc.FramesPerDirection
	state.images = make([]*image.RGBA, totalFrames)

	for dirIdx := range w.dcc.Directions {
		fw := w.dcc.Directions[dirIdx].Box.Width
		fh := w.dcc.Directions[dirIdx].Box.Height

		for frameIdx := range w.dcc.Directions[dirIdx].Frames {
			absoluteFrameIdx := (dirIdx * w.dcc.FramesPerDirection) + frameIdx

			frame := w.dcc.Directions[dirIdx].Frames[frameIdx]
			pixels := frame.PixelData

			state.images[absoluteFrameIdx] = image.NewRGBA(image.Rect(0, 0, fw, fh))

			for y := 0; y < fh; y++ {
				for x := 0; x < fw; x++ {
					idx := x + (y * fw)
					if idx >= len(pixels) {
						continue
					}

					val := pixels[idx]

					RGBAColor := w.makeImagePixel(val)
					state.images[absoluteFrameIdx].Set(x, y, RGBAColor)
				}
			}
		}
	}

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			w.widget.textureLoader.CreateTextureFromARGB(state.images[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := w.getState()
		s.textures = textures
		w.setState(s)
	}()
}
