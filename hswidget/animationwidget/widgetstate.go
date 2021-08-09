package animationwidget

import (
	"image"
	"log"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2datautils"
	"github.com/ianling/giu"
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

type widgetState struct {
	isPlaying bool
	repeat    bool
	tickTime  int32
	playMode  animationPlayMode

	// Cache - will not be saved
	images   []*image.RGBA
	textures []*giu.Texture

	ticker *time.Ticker
}

func (s *widgetState) getTickTime() int32 {
	return s.tickTime
}

func (s *widgetState) getTick() *int32 {
	return &s.tickTime
}

func (s *widgetState) getImages() []*image.RGBA {
	return s.images
}

func (s *widgetState) getPlayMode() animationPlayMode {
	return s.playMode
}

func (s *widgetState) setPlayMode(pm animationPlayMode) {
	s.playMode = pm
}

func (s *widgetState) getRepeat() *bool {
	return &s.repeat
}

func (s *widgetState) getPlayingPointer() *bool {
	return &s.isPlaying
}

func (s *widgetState) getTicker() *time.Ticker {
	return s.ticker
}

func (s *widgetState) encode(sw *d2datautils.StreamWriter) {
	sw.PushBytes(byte(hsutil.BoolToInt(s.isPlaying)))
	sw.PushBytes(byte(hsutil.BoolToInt(s.repeat)))

	sw.PushInt32(s.tickTime)
	sw.PushBytes(byte(s.playMode))
}

func (s *widgetState) decode(sr *d2datautils.StreamReader) {
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

func (s *widgetState) widgetDispose() {
	s.textures = nil
}

type controlStructure struct {
	direction int32
	frame     int32
	scale     int32
}

type widgetMode int32

type viewerState struct {
	*controlStructure
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
