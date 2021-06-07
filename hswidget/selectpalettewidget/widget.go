package selectpalettewidget

import (
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

// SelectPaletteWidget represents an pop-up MPQ explorer, when we're
// selectin DAT palette
type SelectPaletteWidget struct {
	mpqExplorer     *hsmpqexplorer.MPQExplorer
	projectExplorer *hsprojectexplorer.ProjectExplorer
	isOpen          *bool
	id              string
	onSelect        func(colors *[256]d2interface.Color)
}

// NewSelectPaletteWidget creates a select palette widget
func NewSelectPaletteWidget(
	id string,
	project *hsproject.Project,
	config *hsconfig.Config,
) *SelectPaletteWidget {
	result := &SelectPaletteWidget{
		id: id,
	}

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

			if result.onSelect != nil {
				result.onSelect(&colors)
			}

			*result.isOpen = false
		}
	}

	mpqExplorer, err := hsmpqexplorer.Create(callback, config, 0, 0)
	if err != nil {
		log.Print(err)
	}

	mpqExplorer.SetProject(project)

	result.mpqExplorer = mpqExplorer

	projectExplorer, err := hsprojectexplorer.Create(nil, callback, 0, 0)
	if err != nil {
		log.Print(err)
	}

	projectExplorer.SetProject(project)

	result.projectExplorer = projectExplorer

	return result
}

// OnSelect sets a callback for ppalette selection
func (p *SelectPaletteWidget) OnSelect(cb func(colors *[256]d2interface.Color)) *SelectPaletteWidget {
	p.onSelect = cb
	return p
}

// IsOpen sets pointer to isOpen variable - determinates if a widget is visible
func (p *SelectPaletteWidget) IsOpen(isOpen *bool) *SelectPaletteWidget {
	p.isOpen = isOpen
	return p
}

// Build builds a widget
func (p *SelectPaletteWidget) Build() {
	giu.PopupModal("##" + p.id + "popUpSelectPalette").IsOpen(p.isOpen).Layout(giu.Layout{
		giu.Child("##"+p.id+"popUpSelectPaletteChildWidget").Size(paletteSelectW, paletteSelectH).Layout(giu.Layout{
			p.projectExplorer.GetProjectTreeNodes(),
			giu.Layout(p.mpqExplorer.GetMpqTreeNodes()),
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
