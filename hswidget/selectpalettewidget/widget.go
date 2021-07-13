package selectpalettewidget

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsprojectexplorer"
)

const (
	paletteSelectW, paletteSelectH = 400, 600
	actionButtonW, actionButtonH   = 200, 30
)

type selectPaletteState struct {
	mpqExplorer     *hsmpqexplorer.MPQExplorer
	projectExplorer *hsprojectexplorer.ProjectExplorer
}

func (s *selectPaletteState) Dispose() {
	s.mpqExplorer = nil
	s.projectExplorer = nil
}

// SelectPaletteWidget represents an pop-up MPQ explorer, when we're
// selectin DAT palette
type SelectPaletteWidget struct {
	isOpen   *bool
	id       string
	onSelect func(colors *[256]d2interface.Color)
	config   *hsconfig.Config
	project  *hsproject.Project
}

// NewSelectPaletteWidget creates a select palette widget
func NewSelectPaletteWidget(
	id string,
	project *hsproject.Project,
	config *hsconfig.Config,
) *SelectPaletteWidget {
	result := &SelectPaletteWidget{
		id:      id,
		config:  config,
		project: project,
	}

	return result
}

// OnSelect sets a callback for ppalette selection
func (p *SelectPaletteWidget) OnSelect(cb func(colors *[256]d2interface.Color)) *SelectPaletteWidget {
	p.onSelect = cb
	return p
}

func (p *SelectPaletteWidget) getState() *selectPaletteState {
	var state *selectPaletteState

	stateID := fmt.Sprintf("selectPalette_%s", p.id)
	s := giu.Context.GetState(stateID)

	if s != nil {
		state = s.(*selectPaletteState)
	} else {
		state = &selectPaletteState{}
		callback := func(path *hscommon.PathEntry) {
			bytes, bytesErr := path.GetFileBytes()
			if bytesErr != nil {
				log.Print(bytesErr)

				return
			}

			ft, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath), &bytes)
			if err != nil {
				log.Print(err)

				return
			}

			if ft == hsfiletypes.FileTypePalette {
				// load new palette:
				paletteData, err := path.GetFileBytes()
				if err != nil {
					log.Print(err)
				}

				palette, err := d2dat.Load(paletteData)
				if err != nil {
					log.Print(err)
				}

				colors := palette.GetColors()

				if p.onSelect != nil {
					p.onSelect(&colors)
				}

				*p.isOpen = false
			}
		}

		mpqExplorer, err := hsmpqexplorer.Create(callback, p.config, 0, 0)
		if err != nil {
			log.Print(err)
		}

		mpqExplorer.SetProject(p.project)

		state.mpqExplorer = mpqExplorer

		projectExplorer, err := hsprojectexplorer.Create(nil, callback, 0, 0)
		if err != nil {
			log.Print(err)
		}

		projectExplorer.SetProject(p.project)

		state.projectExplorer = projectExplorer
		giu.Context.SetState(stateID, state)
	}

	return state
}

// IsOpen sets pointer to isOpen variable - determinates if a widget is visible
func (p *SelectPaletteWidget) IsOpen(isOpen *bool) *SelectPaletteWidget {
	p.isOpen = isOpen
	return p
}

// Build builds a widget
func (p *SelectPaletteWidget) Build() {
	state := p.getState()
	giu.PopupModal("##" + p.id + "popUpSelectPalette").IsOpen(p.isOpen).Layout(giu.Layout{
		giu.Child("##"+p.id+"popUpSelectPaletteChildWidget").Size(paletteSelectW, paletteSelectH).Layout(giu.Layout{
			state.projectExplorer.GetProjectTreeNodes(),
			giu.Layout(state.mpqExplorer.GetMpqTreeNodes()),
			giu.Separator(),
			giu.Button("Don't use any palette##"+p.id+"selectPaletteDonotUseAny").
				Size(actionButtonW, actionButtonH).
				OnClick(func() {
					if p.onSelect != nil {
						p.onSelect(nil)
					}
					*p.isOpen = false
				}),
			giu.Button("Exit##"+p.id+"selectPaletteExit").
				Size(actionButtonW, actionButtonH).
				OnClick(func() {
					*p.isOpen = false
				}),
		}),
	}).Build()
}
