package hssoundeditor

import (
	"fmt"
	"log"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

const sampleRate = 44100

type SoundEditor struct {
	hseditor.Editor

	streamer beep.StreamSeekCloser
	control  *beep.Ctrl
	format   beep.Format
	file     string
}

func Create(file string, audioStream d2interface.DataStream) (*SoundEditor, error) {
	streamer, format, err := wav.Decode(audioStream)

	if err != nil {
		log.Fatal(err)
	}

	control := &beep.Ctrl{
		Streamer: beep.Loop(-1, streamer),
		Paused:   true,
	}

	result := &SoundEditor{
		file:     file,
		streamer: streamer,
		control:  control,
		format:   format,
	}

	speaker.Play(result.control)

	return result, nil
}

func (s SoundEditor) Render() {
	secondsCurrent := s.streamer.Position() / 22050
	secondsTotal := s.streamer.Len() / 22050

	g.WindowV(s.GetWindowTitle(), nil, g.WindowFlagsNoResize, 50, 50, 300, 100, g.Layout{
		g.ProgressBar(float32(s.streamer.Position())/float32(s.streamer.Len()), -1, 24,
			fmt.Sprintf("%d:%02d / %d:%02d", secondsCurrent/60, secondsCurrent%60, secondsTotal/60, secondsTotal%60)),
		g.Separator(),
		g.Line(
			g.Button("Play", s.play),
			g.Button("Stop", s.stop),
		),
	})
}

func (s SoundEditor) GetWindowTitle() string {
	return "Sound Editor [" + s.file + "]"
}

func (s SoundEditor) play() {
	speaker.Lock()
	s.control.Paused = false
	speaker.Unlock()
}

func (s SoundEditor) stop() {
	speaker.Lock()
	if s.control.Paused {
		if err := s.streamer.Seek(0); err != nil {
			log.Fatal(err)
		}
	}
	s.control.Paused = true
	speaker.Unlock()
}
