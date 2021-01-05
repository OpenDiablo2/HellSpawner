package hssoundeditor

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
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

func Create(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error) {
	streamer, format, err := wav.Decode(bytes.NewBuffer(*data))

	if err != nil {
		log.Fatal(err)
	}

	control := &beep.Ctrl{
		Streamer: beep.Loop(-1, streamer),
		Paused:   true,
	}

	result := &SoundEditor{
		file:     filepath.Base(pathEntry.FullPath),
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

	g.Window(s.GetWindowTitle()).IsOpen(&s.Visible).Flags(g.WindowFlagsNoResize).Pos(50, 50).Size(300, 100).Layout(g.Layout{
		g.ProgressBar(float32(s.streamer.Position())/float32(s.streamer.Len())).Size(-1, 24).
			Overlay(fmt.Sprintf("%d:%02d / %d:%02d", secondsCurrent/60, secondsCurrent%60, secondsTotal/60, secondsTotal%60)),
		g.Separator(),
		g.Line(
			g.Button("Play").OnClick(s.play),
			g.Button("Stop").OnClick(s.stop),
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
