package hspreferencesdialog

import (
	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

type PreferencesDialog struct {
	*hsdialog.Dialog

	config          hsconfig.Config
	onConfigChanged func(config hsconfig.Config)
}

func Create(onConfigChanged func(config hsconfig.Config)) *PreferencesDialog {
	result := &PreferencesDialog{
		Dialog:          hsdialog.New("Preferences"),
		onConfigChanged: onConfigChanged,
	}
	result.Visible = false

	return result
}

func (p *PreferencesDialog) Show(config *hsconfig.Config) {
	p.Dialog.Show()

	p.config = *config
}

func (p *PreferencesDialog) Build() {
	p.IsOpen(&p.Visible).Layout(g.Layout{
		g.Child("PreferencesLayout").Size(300, 150).Layout(g.Layout{
			g.Label("Auxiliary MPQ Path"),
			g.Line(
				g.Button("...##AppPreferencesAuxMPQPathBrowse").Size(30, 0).OnClick(p.onBrowseAuxMpqPathClicked),
				g.InputText("##AppPreferencesAuxMPQPath", &p.config.AuxiliaryMpqPath).Size(-1).Flags(g.InputTextFlagsReadOnly),
			),
			g.Label("External MPQ listfile Path"),
			g.Line(
				g.Button("...##AppPreferencesListfilePathBrowse").Size(30, 0).OnClick(p.onBrowseExternalListfileClicked),
				g.InputText("##AppPreferencesListfilePath", &p.config.ExternalListfile).Size(-1).Flags(g.InputTextFlagsReadOnly),
			),
			g.Checkbox("Open most recent project on start-up", &p.config.OpenMostRecentOnStartup),
		}),
		g.Line(
			g.Button("Save##AppPreferencesSave").OnClick(p.onSaveClicked),
			g.Button("Cancel##AppPreferencesCancel").OnClick(p.onCancelClicked),
		),
	})
}

func (p *PreferencesDialog) onBrowseAuxMpqPathClicked() {
	path, err := dialog.Directory().Browse()
	if err != nil || len(path) == 0 {
		return
	}
	p.config.AuxiliaryMpqPath = path
}

func (p *PreferencesDialog) onBrowseExternalListfileClicked() {
	path := dialog.File()
	path.Filter("Text file", "txt")
	filePath, err := path.Load()

	if err != nil || len(filePath) == 0 {
		return
	}
	p.config.ExternalListfile = filePath
}

func (p *PreferencesDialog) onSaveClicked() {
	p.onConfigChanged(p.config)
	p.Visible = false
}

func (p *PreferencesDialog) onCancelClicked() {
	p.Visible = false
}
