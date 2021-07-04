// Package hssoundeditor represents a soundEditor's window
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

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"

	g "github.com/ianling/giu"
)

const (
	mainWindowW, mainWindowH  = 300, 70
	progressIndicatorModifier = 60
	progressTimeModifier      = 22050
	btnSize                   = 20
)

// static check, to ensure, if sound editor implemented editoWindow
var _ hscommon.EditorWindow = &SoundEditor{}

// SoundEditor represents a sound editor
type SoundEditor struct {
	*hseditor.Editor

	streamer      beep.StreamSeekCloser
	control       *beep.Ctrl
	format        beep.Format
	file          string
	textureLoader hscommon.TextureLoader
}

// Create creates a new sound editor
func Create(_ *hsconfig.Config,
	tl hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	streamer, format, err := wav.Decode(bytes.NewReader(*data))
	if err != nil {
		return nil, fmt.Errorf("wav decode error: %w", err)
	}

	control := &beep.Ctrl{
		Streamer: beep.Loop(-1, streamer),
		Paused:   false,
	}

	result := &SoundEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		file:          filepath.Base(pathEntry.FullPath),
		streamer:      streamer,
		control:       control,
		format:        format,
		textureLoader: tl,
	}

	result.Path = pathEntry

	speaker.Play(result.control)

	return result, nil
}

// Build builds a sound editor
func (s *SoundEditor) Build() {
	isPlaying := !s.control.Paused

	secondsCurrent := s.streamer.Position() / progressTimeModifier
	secondsTotal := s.streamer.Len() / progressTimeModifier

	const progressBarHeight = 24 // px

	progress := float32(s.streamer.Position()) / float32(s.streamer.Len())

	s.IsOpen(&s.Visible).
		Flags(g.WindowFlagsNoResize).
		Size(mainWindowW, mainWindowH).
		Layout(g.Layout{
			g.Row(
				hswidget.PlayPauseButton("##"+s.Path.GetUniqueID()+"playPause", &isPlaying, s.textureLoader).
					OnPlayClicked(s.play).OnPauseClicked(s.stop).Size(btnSize, btnSize),
				g.ProgressBar(progress).Size(-1, progressBarHeight).
					Overlay(fmt.Sprintf("%d:%02d / %d:%02d",
						secondsCurrent/progressIndicatorModifier,
						secondsCurrent%progressIndicatorModifier,
						secondsTotal/progressIndicatorModifier,
						secondsTotal%progressIndicatorModifier,
					)),
			),
		})
}

// Cleanup closes an editor
func (s *SoundEditor) Cleanup() {
	speaker.Lock()
	s.control.Paused = true

	if err := s.streamer.Close(); err != nil {
		log.Print(err)
	}

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
			log.Print(err)
			return
		}
	}

	s.control.Paused = true

	speaker.Unlock()
}

// UpdateMainMenuLayout updates mainMenu's layout to it contain soundEditor's options
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

// GenerateSaveData generates data to be saved
func (s *SoundEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := s.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (s *SoundEditor) Save() {
	s.Editor.Save(s)
}
