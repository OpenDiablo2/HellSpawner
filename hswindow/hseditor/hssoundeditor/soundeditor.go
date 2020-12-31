package hssoundeditor

import (
	"fmt"
	"log"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	g "github.com/OpenDiablo2/giu"
	"github.com/OpenDiablo2/giu/imgui"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

type SoundEditor struct {
	hseditor.Editor

	streamer beep.StreamSeekCloser
	control  *beep.Ctrl
	format   beep.Format
	file     string
}

func (s *SoundEditor) Cleanup() {
	speaker.Lock()
	s.control.Paused = true
	s.streamer.Close()
	speaker.Unlock()
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

func (s *SoundEditor) Render() {
	if !s.Visible {
		return
	}

	if s.ToFront {
		s.ToFront = false
		imgui.SetNextWindowFocus()
	}

	secondsCurrent := s.streamer.Position() / 22050
	secondsTotal := s.streamer.Len() / 22050

	g.WindowV(s.GetWindowTitle(), &s.Visible, g.WindowFlagsNoResize, 50, 50, 300, 100, g.Layout{
		g.ProgressBar(float32(s.streamer.Position())/float32(s.streamer.Len()), -1, 24,
			fmt.Sprintf("%d:%02d / %d:%02d", secondsCurrent/60, secondsCurrent%60, secondsTotal/60, secondsTotal%60)),
		g.Separator(),
		g.Line(
			g.Button("Play", s.play),
			g.Button("Stop", s.stop),
		),
	})
}

func (s *SoundEditor) GetWindowTitle() string {
	return s.file + "##" + s.GetId()
}

func (s *SoundEditor) play() {
	speaker.Lock()
	s.control.Paused = false
	speaker.Unlock()
}

func (s *SoundEditor) stop() {
	speaker.Lock()
	if s.control.Paused {
		if err := s.streamer.Seek(0); err != nil {
			log.Fatal(err)
		}
	}
	s.control.Paused = true
	speaker.Unlock()
}
