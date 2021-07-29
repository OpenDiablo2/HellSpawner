package animationwidget

import (
	"image"
	"log"
	"time"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
)

type animationPlayMode byte

const (
	playModeForward animationPlayMode = iota
	playModeBackward
	playModePingPong
)

const (
	miliseconds     = 1000
	defaultTickTime = 100
)

func (animation animationPlayMode) String() string {
	playModes := map[animationPlayMode]string{
		playModeForward:  "Forwards",
		playModeBackward: "Backwards",
		playModePingPong: "Ping-pong",
	}

	playMode, ok := playModes[animation]

	if !ok {
		return "Unknown"
	}

	return playMode
}

type WidgetState struct {
	isPlaying bool
	repeat    bool
	tickTime  int32
	playMode  animationPlayMode

	// Cache - will not be saved
	images   []*image.RGBA
	textures []*giu.Texture

	isForward bool // Determines a direction of animation
	ticker    *time.Ticker
}

func (s *WidgetState) getTickTime() int32 {
	return s.tickTime
}

func (s *WidgetState) getTick() *int32 {
	return &s.tickTime
}

func (s *WidgetState) getImages() []*image.RGBA {
	return s.images
}

func (s *WidgetState) getPlayMode() animationPlayMode {
	return s.playMode
}

func (s *WidgetState) setPlayMode(pm animationPlayMode) {
	s.playMode = pm
}

func (s *WidgetState) getRepeat() *bool {
	return &s.repeat
}

func (s *WidgetState) getPlaying() *bool {
	return &s.isPlaying
}

func (s *WidgetState) getTicker() *time.Ticker {
	return s.ticker
}

func (s *WidgetState) Encode(sw *d2datautils.StreamWriter) {
	sw.PushBytes(byte(hsutil.BoolToInt(s.isPlaying)))
	sw.PushBytes(byte(hsutil.BoolToInt(s.repeat)))

	sw.PushInt32(s.tickTime)
	sw.PushBytes(byte(s.playMode))
}

func (s *WidgetState) Decode(sr *d2datautils.StreamReader) {
	isPlaying, err := sr.ReadByte()
	if err != nil {
		log.Print(err)
		return
	}

	s.isPlaying = isPlaying == 1

	repeat, err := sr.ReadByte()
	if err != nil {
		log.Print(err)
		return
	}

	s.repeat = repeat == 1

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
}

func (s *WidgetState) widgetDispose() {
	s.textures = nil
}

type controlStructure struct {
	direction int32
	frame     int32
	scale     int32
}

type widgetMode int32

type viewerState struct {
	controls           *controlStructure
	lastFrame          int32
	lastDirection      int32
	framesPerDirection uint32
}

// Dispose dispose state
func (s *viewerState) Dispose() {
	// Noop
}

type tiledState struct {
	width,
	height int32
	tiled *giu.Texture
	imgw, // nolint:structcheck // linter's bug - it is used
	imgh int
}

func (s *tiledState) Dispose() {
	s.width, s.height = 0, 0
	s.tiled = nil
}
