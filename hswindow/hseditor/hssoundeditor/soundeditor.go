package hssoundeditor

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"

	g "github.com/ianling/giu"
)

type SoundEditor struct {
	*hseditor.Editor

	streamer beep.StreamSeekCloser
	control  *beep.Ctrl
	format   beep.Format
	file     string
}

func Create(pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	streamer, format, err := wav.Decode(bytes.NewReader(*data))

	if err != nil {
		log.Fatal(err)
	}

	control := &beep.Ctrl{
		Streamer: beep.Loop(-1, streamer),
		Paused:   false,
	}

	result := &SoundEditor{
		Editor:   hseditor.New(pathEntry, x, y, project),
		file:     filepath.Base(pathEntry.FullPath),
		streamer: streamer,
		control:  control,
		format:   format,
	}

	result.Path = pathEntry

	speaker.Play(result.control)

	return result, nil
}

func (s *SoundEditor) Build() {
	secondsCurrent := s.streamer.Position() / 22050
	secondsTotal := s.streamer.Len() / 22050

	s.IsOpen(&s.Visible).Flags(g.WindowFlagsNoResize).Size(300, 100).Layout(g.Layout{
		g.ProgressBar(float32(s.streamer.Position())/float32(s.streamer.Len())).Size(-1, 24).
			Overlay(fmt.Sprintf("%d:%02d / %d:%02d", secondsCurrent/60, secondsCurrent%60, secondsTotal/60, secondsTotal%60)),
		g.Separator(),
		g.Line(
			g.Button("Play").OnClick(s.play),
			g.Button("Stop").OnClick(s.stop),
		),
	})
}

func (s *SoundEditor) Cleanup() {
	speaker.Lock()
	s.control.Paused = true
	s.streamer.Close()

	if s.HasChanges(s) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			s.Path.FullPath).YesNo(); shouldSave {
			s.Save()
		}
	}

	s.Editor.Cleanup()
	speaker.Unlock()
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

func (s *SoundEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Sound Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			s.Cleanup()
		}),
	})

	*l = append(*l, m)
}

func (s *SoundEditor) GenerateSaveData() []byte {
	// TODO -- save real data for this editor
	data, _ := s.Path.GetFileBytes()

	return data
}

func (s *SoundEditor) Save() {
	s.Editor.Save(s)
}
