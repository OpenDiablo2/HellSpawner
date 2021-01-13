package hsapp

import (
	"fmt"
	"os"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	g "github.com/AllenDang/giu"
)

func (a *App) renderMainMenuBar() {
	projectOpened := a.project != nil

	g.MainMenuBar().Layout(g.Layout{
		g.Menu("File##MainMenuFile").Layout(g.Layout{
			g.Menu("New##MainMenuFileNew").Layout(g.Layout{
				g.MenuItem("Project...##MainMenuFileNewProject").OnClick(a.onNewProjectClicked),
			}),
			g.Menu("Open##MainMenuFileOpen").Layout(g.Layout{
				g.MenuItem("Project...##MainMenuFileOpenProject").OnClick(a.onOpenProjectClicked),
			}),
			g.Menu("Open Recent##MainMenuOpenRecent").Layout(g.Layout{
				g.Custom(func() {
					if len(a.config.RecentProjects) == 0 {
						g.MenuItem("No recent projects...##MainMenuOpenRecentItems").Build()
						return
					}
					for idx := range a.config.RecentProjects {
						projectName := a.config.RecentProjects[idx]
						g.MenuItem(fmt.Sprintf("%s##MainMenuOpenRecent_%d", projectName, idx)).OnClick(func() {
							a.loadProjectFromFile(projectName)
						}).Build()
					}
				}),
			}),
			g.Separator(),
			g.MenuItem("Preferences...##MainMenuFilePreferences").OnClick(a.onFilePreferencesClicked),
			g.Separator(),
			g.MenuItem("Exit##MainMenuFileExit").OnClick(func() { os.Exit(0) }),
		}),
		g.Menu("View##MainMenuView").Layout(a.buildViewMenu()),
		g.Menu("Project##MainMenuProject").Layout(g.Layout{
			g.MenuItem("Run in OpenDiablo2##MainMenuProjectRun").Enabled(projectOpened).OnClick(a.onProjectRunClicked),
			g.Separator(),
			g.MenuItem("Properties...##MainMenuProjectProperties").Enabled(projectOpened).OnClick(a.onProjectPropertiesClicked),
			g.Separator(),
			g.MenuItem("Export MPQ...##MainMenuProjectExport").Enabled(projectOpened).OnClick(a.onProjectExportMPQClicked),
		}),
		g.Menu("Help").Layout(g.Layout{
			g.MenuItem("About HellSpawner...##MainMenuHelpAbout").OnClick(a.onHelpAboutClicked),
		}),
	}).Build()
}

func (a *App) buildViewMenu() g.Layout {
	result := make([]g.Widget, 0)

	result = append(result, g.Menu("Tool Windows").Layout(g.Layout{
		g.MenuItem("Project Explorer").Selected(a.projectExplorer.Visible).Enabled(true).OnClick(a.toggleProjectExplorer),
		g.MenuItem("MPQ Explorer").Selected(a.mpqExplorer.Visible).Enabled(a.project != nil).OnClick(a.toggleMPQExplorer),
	}))

	if len(a.editors) == 0 {
		return result
	}

	result = append(result, g.Separator())

	for idx := range a.editors {
		i := idx
		result = append(result, g.MenuItem(a.editors[idx].GetWindowTitle()).OnClick(a.editors[i].BringToFront))
	}

	return result
}

func (a *App) onNewProjectClicked() {
	file, err := dialog.File().Filter("HellSpawner Project", "hsp").Save()
	if err != nil || len(file) == 0 {
		return
	}
	var project *hsproject.Project
	if project, err = hsproject.CreateNew(file); err != nil {
		return
	}
	a.loadProjectFromFile(project.GetProjectFilePath())
}

func (a *App) onOpenProjectClicked() {
	file, err := dialog.File().Filter("HellSpawner Project", "hsp").Load()
	if err != nil || len(file) == 0 {
		return
	}
	a.loadProjectFromFile(file)
}

func (a *App) onProjectPropertiesClicked() {
	a.projectPropertiesDialog.Show(a.project, a.config)
}

func (a *App) onFilePreferencesClicked() {
	a.preferencesDialog.Show(a.config)
}

func (a *App) onHelpAboutClicked() {
	a.aboutDialog.Show()
}

func (a *App) onProjectRunClicked() {

}

func (a *App) onProjectExportMPQClicked() {

}
