package hspreferencesdialog

import (
	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
	"github.com/OpenDiablo2/dialog"
)

type PreferencesDialog struct {
	hsdialog.Dialog

	config          hsconfig.Config
	onConfigChanged func(config hsconfig.Config)
}

func Create(onConfigChanged func(config hsconfig.Config)) *PreferencesDialog {
	result := &PreferencesDialog{
		onConfigChanged: onConfigChanged,
	}
	result.Visible = false

	return result
}

func (p *PreferencesDialog) Show(config *hsconfig.Config) {
	p.Dialog.Show()

	p.config = *config
}

func (p *PreferencesDialog) Render() {
	hswidget.ModalDialog("Preferences##AppPreferences", &p.Visible, g.Layout{
		g.Label("Auxiliary MPQ Path"),
		g.Line(
			g.ButtonV("...##AppPreferencesAuxMPQPathBrowse", 30, 0, p.onBrowseAuxMpqPathClicked),
			g.InputTextV("##AppPreferencesAuxMPQPath", 250, &p.config.AuxiliaryMpqPath, g.InputTextFlagsReadOnly, nil, nil),
		),
		g.Separator(),
		g.Line(
			g.Button("Save##AppPreferencesSave", p.onSaveClicked),
			g.Button("Cancel##AppPreferencesCancel", p.onCancelClicked),
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

func (p *PreferencesDialog) onSaveClicked() {
	p.onConfigChanged(p.config)
	p.Visible = false
}

func (p *PreferencesDialog) onCancelClicked() {
	p.Visible = false
}
