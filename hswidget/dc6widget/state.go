package dc6widget

import (
	"fmt"
	image2 "image"
	"image/color"
	"log"
	"time"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

const (
	miliseconds     = 1000
	defaultTickTime = 100
)

type animationPlayMode byte

const (
	playModeForward animationPlayMode = iota
	playModeBackword
	playModeLeftRight
)

func (a animationPlayMode) String() string {
	s := map[animationPlayMode]string{
		playModeForward:   "Forwards",
		playModeBackword:  "Backwords",
		playModeLeftRight: "Ping-Pong",
	}

	k, ok := s[a]
	if !ok {
		return "Unknown"
	}

	return k
}

type widgetMode int32

const (
	dc6WidgetViewer widgetMode = iota
)

// widgetState represents dc6 viewer's state
type widgetState struct {
	viewerState
	mode widgetMode

	isPlaying bool
	repeat    bool
	tickTime  int32
	playMode  animationPlayMode

	// cache - will not be saved
	textures []*giu.Texture

	leftRightDirection bool
	ticker             *time.Ticker
}

func (w *widgetState) Dispose() {
	w.viewerState.Dispose()
	w.mode = dc6WidgetViewer
	w.textures = nil
}

func (w *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(int32(w.mode))
	sw.PushInt32(w.controls.direction)
	sw.PushInt32(w.controls.frame)
	sw.PushInt32(w.controls.scale)

	return sw.GetBytes()
}

func (w *widgetState) Decode(data []byte) {
	sr := d2datautils.CreateStreamReader(data)

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
}

// nolint:structcheck // :-/ linter bug?! thes values are deffinitly used
type viewerState struct {
	controls struct {
		direction int32
		frame     int32
		scale     int32
	}

	lastFrame          int32
	lastDirection      int32
	framesPerDirection uint32
}

// Dispose disposes state
func (s *viewerState) Dispose() {
	// noop
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

func (p *widget) initState() {
	// Prevent multiple invocation to LoadImage.
	newState := &widgetState{
		mode: dc6WidgetViewer,
		viewerState: viewerState{
			lastFrame:          -1,
			lastDirection:      -1,
			framesPerDirection: p.dc6.FramesPerDirection,
		},

		isPlaying: false,
		repeat:    false,
		tickTime:  defaultTickTime,
		playMode:  playModeForward,
	}

	newState.ticker = time.NewTicker(time.Second * time.Duration(newState.tickTime) / miliseconds)

	go p.runPlayer(newState)

	p.setState(newState)

	totalFrames := int(p.dc6.Directions * p.dc6.FramesPerDirection)
	rgb := make([]*image2.RGBA, totalFrames)

	for frameIndex := 0; frameIndex < int(p.dc6.Directions*p.dc6.FramesPerDirection); frameIndex++ {
		rgb[frameIndex] = image2.NewRGBA(image2.Rect(0, 0, int(p.dc6.Frames[frameIndex].Width), int(p.dc6.Frames[frameIndex].Height)))
		decodedFrame := p.dc6.DecodeFrame(frameIndex)

		for y := 0; y < int(p.dc6.Frames[frameIndex].Height); y++ {
			for x := 0; x < int(p.dc6.Frames[frameIndex].Width); x++ {
				idx := x + (y * int(p.dc6.Frames[frameIndex].Width))
				val := decodedFrame[idx]

				alpha := maxAlpha

				if val == 0 {
					alpha = 0
				}

				var r, g, b uint8

				if p.palette != nil {
					col := p.palette[val]
					r, g, b = col.R(), col.G(), col.B()
				} else {
					r, g, b = val, val, val
				}

				rgb[frameIndex].Set(
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

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			p.textureLoader.CreateTextureFromARGB(rgb[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := p.getState()
		s.textures = textures
		p.setState(s)
	}()
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

// nolint:gocognit,gocyclo // will cut later
func (p *widget) runPlayer(state *widgetState) {
	for range state.ticker.C {
		if state.isPlaying {
			switch state.playMode {
			case playModeForward:
				if state.controls.frame < int32(p.dc6.FramesPerDirection-1) {
					state.controls.frame++
				} else {
					if state.repeat {
						state.controls.frame = 0
					} else {
						state.isPlaying = false
					}
				}
			case playModeBackword:
				if state.controls.frame > 0 {
					state.controls.frame--
				} else {
					if state.repeat {
						state.controls.frame = int32(p.dc6.FramesPerDirection)
					} else {
						state.isPlaying = false
					}
				}
			case playModeLeftRight:
				if state.leftRightDirection {
					if fpd := int32(p.dc6.FramesPerDirection) - 1; state.controls.frame < fpd {
						state.controls.frame++
						if state.controls.frame == fpd {
							state.leftRightDirection = false
						}
					}
				} else {
					if state.controls.frame > 0 {
						state.controls.frame--
						if state.controls.frame == 0 {
							state.leftRightDirection = true
							if !state.repeat {
								state.isPlaying = false
							}
						}
					}
				}
			}
		}
	}
}
