package dc6widget

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/ianling/giu"
	gim "github.com/ozankasikci/go-image-merge"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	miliseconds     = 1000
	defaultTickTime = 100
)

type animationPlayMode byte

const (
	playModeForward animationPlayMode = iota
	playModeBackword
	playModePingPong
)

func (a animationPlayMode) String() string {
	s := map[animationPlayMode]string{
		playModeForward:  "Forwards",
		playModeBackword: "Backwords",
		playModePingPong: "Ping-Pong",
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
	dc6WidgetTiledView
)

// widgetState represents dc6 viewer's state
type widgetState struct {
	viewerState
	tiledState
	Mode widgetMode

	IsPlaying bool
	Repeat    bool
	TickTime  int32
	PlayMode  animationPlayMode

	// cache - will not be saved
	rgb      []*image.RGBA
	textures []*giu.Texture

	IsForward bool
	ticker    *time.Ticker
}

func (w *widgetState) Dispose() {
	w.viewerState.Dispose()
	w.Mode = dc6WidgetViewer
	w.textures = nil
}

type viewerState struct {
	Controls struct {
		Direction int32
		Frame     int32
		Scale     int32
	}

	lastFrame          int32
	lastDirection      int32
	framesPerDirection uint32
}

// Dispose disposes state
func (s *viewerState) Dispose() {
	// noop
}

type tiledState struct {
	Width,
	Height int32
	tiled *giu.Texture
	Imgw,
	Imgh int
}

func (s *tiledState) Dispose() {
	s.Width, s.Height = 0, 0
	s.tiled = nil
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
		Mode: dc6WidgetViewer,
		viewerState: viewerState{
			lastFrame:          -1,
			lastDirection:      -1,
			framesPerDirection: p.dc6.FramesPerDirection,
		},
		tiledState: tiledState{
			Width:  int32(p.dc6.FramesPerDirection),
			Height: 1,
		},

		IsPlaying: false,
		Repeat:    false,
		TickTime:  defaultTickTime,
		PlayMode:  playModeForward,
	}

	newState.ticker = time.NewTicker(time.Second * time.Duration(newState.TickTime) / miliseconds)

	go p.runPlayer(newState)

	totalFrames := int(p.dc6.Directions * p.dc6.FramesPerDirection)
	newState.rgb = make([]*image.RGBA, totalFrames)

	for frameIndex := 0; frameIndex < int(p.dc6.Directions*p.dc6.FramesPerDirection); frameIndex++ {
		newState.rgb[frameIndex] = image.NewRGBA(image.Rect(0, 0, int(p.dc6.Frames[frameIndex].Width), int(p.dc6.Frames[frameIndex].Height)))
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

				newState.rgb[frameIndex].Set(
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

	p.setState(newState)

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			p.textureLoader.CreateTextureFromARGB(newState.rgb[frameIndex], func(t *giu.Texture) {
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

func (p *widget) runPlayer(state *widgetState) {
	for range state.ticker.C {
		if !state.IsPlaying {
			continue
		}

		numFrames := int32(p.dc6.FramesPerDirection - 1)
		isLastFrame := state.Controls.Frame == numFrames

		// update play direction
		switch state.PlayMode {
		case playModeForward:
			state.IsForward = true
		case playModeBackword:
			state.IsForward = false
		case playModePingPong:
			if isLastFrame || state.Controls.Frame == 0 {
				state.IsForward = !state.IsForward
			}
		}

		// now update the frame number
		if state.IsForward {
			state.Controls.Frame++
		} else {
			state.Controls.Frame--
		}

		state.Controls.Frame = int32(hsutil.Wrap(int(state.Controls.Frame), int(p.dc6.FramesPerDirection)))

		// next, check for stopping/repeat
		isStoppingFrame := (state.Controls.Frame == 0) || (state.Controls.Frame == numFrames)

		if isStoppingFrame && !state.Repeat {
			state.IsPlaying = false
		}
	}
}

func (p *widget) recalculateTiledViewWidth(state *widgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.Width = int32(p.dc6.FramesPerDirection) / state.Height
	p.createImage(state)
}

func (p *widget) recalculateTiledViewHeight(state *widgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.tiledState.Height = int32(p.dc6.FramesPerDirection) / state.Width
	p.createImage(state)
}

func (p *widget) createImage(state *widgetState) {
	firstFrame := state.Controls.Direction * int32(p.dc6.FramesPerDirection)

	grids := make([]*gim.Grid, 0)

	for j := int32(0); j < state.Height*state.Width; j++ {
		grids = append(grids, &gim.Grid{Image: image.Image(state.rgb[firstFrame+j])})
	}

	newimg, err := gim.New(grids, int(state.Width), int(state.Height)).Merge()
	if err != nil {
		log.Printf("merging image error: %v", err)
		return
	}

	p.textureLoader.CreateTextureFromARGB(newimg, func(t *giu.Texture) {
		state.tiled = t
	})

	state.Imgw, state.Imgh = newimg.Bounds().Dx(), newimg.Bounds().Dy()
}
