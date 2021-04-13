package dc6widget

import (
	"fmt"
	"image"
	image2 "image"
	"image/color"
	"log"
	"time"

	"github.com/ianling/giu"
	gim "github.com/ozankasikci/go-image-merge"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"

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
	dc6WidgetMerge
)

// widgetState represents dc6 viewer's state
type widgetState struct {
	viewerState
	mergeState
	mode widgetMode

	isPlaying bool
	repeat    bool
	tickTime  int32
	playMode  animationPlayMode

	// cache - will not be saved
	rgb      []*image2.RGBA
	textures []*giu.Texture

	isForward bool
	ticker    *time.Ticker
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

	sw.PushBytes(byte(hsutil.BoolToInt(w.isPlaying)))
	sw.PushBytes(byte(hsutil.BoolToInt(w.repeat)))
	sw.PushInt32(w.tickTime)
	sw.PushBytes(byte(w.playMode))

	sw.PushInt32(w.width)
	sw.PushInt32(w.height)

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

	isPlaying, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	w.isPlaying = isPlaying == 1

	repeat, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	w.repeat = repeat == 1

	w.tickTime, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)

		return
	}

	playMode, err := sr.ReadByte()
	if err != nil {
		log.Print(err)

		return
	}

	w.playMode = animationPlayMode(playMode)

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

type mergeState struct {
	width,
	height int32
	merged *giu.Texture
	imgw,  // nolint:structcheck // linter's bug - it is used
	imgh int
}

func (s *mergeState) Dispose() {
	s.width, s.height = 0, 0
	s.merged = nil
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
		mergeState: mergeState{
			width:  int32(p.dc6.FramesPerDirection),
			height: 1,
		},

		isPlaying: false,
		repeat:    false,
		tickTime:  defaultTickTime,
		playMode:  playModeForward,
	}

	newState.ticker = time.NewTicker(time.Second * time.Duration(newState.tickTime) / miliseconds)

	go p.runPlayer(newState)

	totalFrames := int(p.dc6.Directions * p.dc6.FramesPerDirection)
	newState.rgb = make([]*image2.RGBA, totalFrames)

	for frameIndex := 0; frameIndex < int(p.dc6.Directions*p.dc6.FramesPerDirection); frameIndex++ {
		newState.rgb[frameIndex] = image2.NewRGBA(image2.Rect(0, 0, int(p.dc6.Frames[frameIndex].Width), int(p.dc6.Frames[frameIndex].Height)))
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
		if !state.isPlaying {
			continue
		}

		numFrames := int32(p.dc6.FramesPerDirection - 1)
		isLastFrame := state.controls.frame == numFrames

		// update play direction
		switch state.playMode {
		case playModeForward:
			state.isForward = true
		case playModeBackword:
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

		state.controls.frame = int32(hsutil.Wrap(int(state.controls.frame), int(p.dc6.FramesPerDirection)))

		// next, check for stopping/repeat
		isStoppingFrame := (state.controls.frame == 0) || (state.controls.frame == numFrames)

		if isStoppingFrame && !state.repeat {
			state.isPlaying = false
		}
	}
}

func (p *widget) recalculateMergeWidth(state *widgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.width = int32(p.dc6.FramesPerDirection) / state.height
	p.createImage(state)
}

func (p *widget) recalculateMergeHeight(state *widgetState) {
	// the area of our rectangle must be less or equal than FramesPerDirection
	state.height = int32(p.dc6.FramesPerDirection) / state.width
	p.createImage(state)
}

func (p *widget) createImage(state *widgetState) {
	firstFrame := state.controls.direction * int32(p.dc6.FramesPerDirection)

	grids := make([]*gim.Grid, 0)
	for j := int32(0); j < state.height*state.width; j++ {
		// https://github.com/ozankasikci/go-image-merge/pull/9
		img := image.Image(state.rgb[firstFrame+j])
		grids = append(grids, &gim.Grid{Image: &img})
	}

	newimg, err := gim.New(grids, int(state.width), int(state.height)).Merge()
	if err != nil {
		log.Print(err)
	}

	p.textureLoader.CreateTextureFromARGB(newimg, func(t *giu.Texture) {
		state.merged = t
	})

	state.imgw, state.imgh = newimg.Bounds().Dx(), newimg.Bounds().Dy()
}
