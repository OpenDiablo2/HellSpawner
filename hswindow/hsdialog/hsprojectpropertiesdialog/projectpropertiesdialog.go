package hsprojectpropertiesdialog

import (
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

type ProjectPropertiesDialog struct {
	hsdialog.Dialog

	project                    hsproject.Project
	onProjectPropertiesChanged func(project hsproject.Project)
}

func Create(onProjectPropertiesChanged func(project hsproject.Project)) *ProjectPropertiesDialog {
	result := &ProjectPropertiesDialog{
		onProjectPropertiesChanged: onProjectPropertiesChanged,
	}

	return result
}

func (p *ProjectPropertiesDialog) Show(project *hsproject.Project) {
	p.Dialog.Show()

	p.project = *project
}

func (p *ProjectPropertiesDialog) Render() {
	if !p.Visible {
		return
	}

	canSave := len(strings.TrimSpace(p.project.ProjectName)) > 0

	imgui.SetNextWindowFocus()
	g.WindowV("Project Properties##ProjectPropertiesDialog", &p.Visible, g.WindowFlagsNoResize|g.WindowFlagsAlwaysAutoResize,
		200, 50, 0, 0, g.Layout{
			g.Label("Project Name:"),
			g.InputText("##ProjectPropertiesDialogProjectName", 250, &p.project.ProjectName),
			g.Label("Description:"),
			g.InputText("##ProjectPropertiesDialogDescription", 250, &p.project.Description),
			g.Label("Author:"),
			g.InputText("##ProjectPropertiesDialogAuthor", 250, &p.project.Author),
			g.Separator(),
			g.Line(
				g.Custom(func() {
					if !canSave {
						imgui.PushStyleVarFloat(imgui.StyleVarAlpha, 0.5)
					}
				}),
				g.Button("Save##ProjectPropertiesDialogSave", p.onSaveClicked),
				g.Custom(func() {
					if !canSave {
						imgui.PopStyleVar()
					}
				}),
				g.Button("Cancel##ProjectPropertiesDialogCancel", p.onCancelClicked),
			),
		},
	)
}

func (p *ProjectPropertiesDialog) onSaveClicked() {
	if len(strings.TrimSpace(p.project.ProjectName)) <= 0 {
		return
	}

	p.onProjectPropertiesChanged(p.project)
	p.Visible = false
}

func (p *ProjectPropertiesDialog) onCancelClicked() {
	p.Visible = false
}
