package dccwidget

import (
	"fmt"
	image2 "image"
	"image/color"
	"log"
	"time"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

const miliseconds = 1000

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

const defaultTickTime = 100

type widgetState struct {
	controls struct {
		direction int32
		frame     int32
		scale     int32
	}

	isPlaying bool
	repeat    bool
	tickTime  int32
	playMode  animationPlayMode

	// cache - will not be saved
	textures []*giu.Texture

	leftRightDirection bool
	ticker             *time.Ticker
}

// Dispose cleans viewers state
func (s *widgetState) Dispose() {
	s.textures = nil
}

func (s *widgetState) Encode() []byte {
	sw := d2datautils.CreateStreamWriter()

	sw.PushInt32(s.controls.direction)
	sw.PushInt32(s.controls.frame)
	sw.PushInt32(s.controls.scale)

	if s.isPlaying {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	if s.repeat {
		sw.PushBytes(1)
	} else {
		sw.PushBytes(0)
	}

	sw.PushInt32(s.tickTime)
	sw.PushBytes(byte(s.playMode))

	return sw.GetBytes()
}

func (s *widgetState) Decode(data []byte) {
	var err error

	sr := d2datautils.CreateStreamReader(data)

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

	isPlaying, err := sr.ReadByte()
	if err != nil {
		log.Print(err)
		return
	}

	s.isPlaying = (isPlaying == 1)

	repeat, err := sr.ReadByte()
	if err != nil {
		log.Print(err)
		return
	}

	s.repeat = (repeat == 1)

	s.tickTime, err = sr.ReadInt32()
	if err != nil {
		log.Print(err)
		return
	}

	playMode, err := sr.ReadByte()
	if err != nil {
		log.Print(err)
		return
	}

	s.playMode = animationPlayMode(playMode)

	// update ticker
	s.ticker.Reset(time.Second * time.Duration(s.tickTime) / miliseconds)
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
	state := &widgetState{
		isPlaying: false,
		repeat:    false,
		tickTime:  defaultTickTime,
		playMode:  playModeForward,
	}

	state.ticker = time.NewTicker(time.Second * time.Duration(state.tickTime) / miliseconds)

	p.setState(state)

	go p.runPlayer(state)

	totalFrames := p.dcc.NumberOfDirections * p.dcc.FramesPerDirection
	images := make([]*image2.RGBA, totalFrames)

	for dirIdx := range p.dcc.Directions {
		fw := p.dcc.Directions[dirIdx].Box.Width
		fh := p.dcc.Directions[dirIdx].Box.Height

		for frameIdx := range p.dcc.Directions[dirIdx].Frames {
			absoluteFrameIdx := (dirIdx * p.dcc.FramesPerDirection) + frameIdx

			frame := p.dcc.Directions[dirIdx].Frames[frameIdx]
			pixels := frame.PixelData

			images[absoluteFrameIdx] = image2.NewRGBA(image2.Rect(0, 0, fw, fh))

			for y := 0; y < fh; y++ {
				for x := 0; x < fw; x++ {
					idx := x + (y * fw)
					if idx >= len(pixels) {
						continue
					}

					val := pixels[idx]

					RGBAColor := p.makeImagePixel(val)
					images[absoluteFrameIdx].Set(x, y, RGBAColor)
				}
			}
		}
	}

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			p.textureLoader.CreateTextureFromARGB(images[frameIndex], func(t *giu.Texture) {
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

func (p *widget) makeImagePixel(val byte) color.RGBA {
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

	RGBAColor := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: alpha,
	}

	return RGBAColor
}

// nolint:gocognit,gocyclo // will cut later
func (p *widget) runPlayer(state *widgetState) {
	for range state.ticker.C {
		if state.isPlaying {
			switch state.playMode {
			case playModeForward:
				if state.controls.frame < int32(p.dcc.FramesPerDirection-1) {
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
						state.controls.frame = int32(p.dcc.FramesPerDirection)
					} else {
						state.isPlaying = false
					}
				}
			case playModeLeftRight:
				if state.leftRightDirection {
					if fpd := int32(p.dcc.FramesPerDirection) - 1; state.controls.frame < fpd {
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
