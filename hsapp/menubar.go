package hsapp

import (
	"fmt"
	"log"
	"os"

	g "github.com/ianling/giu"
	"github.com/pkg/browser"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
)

const (
	githubURL            = "https://github.com/OpenDiablo2/HellSpawner"
	discordInvitationURL = "https://discord.gg/pRy8tdc"
	twitchURL            = "https://www.twitch.tv/essial/"
	supportURL           = "https://www.patreon.com/bePatron?u=37261055"
)

func (a *App) renderMainMenuBar() {
	var runAbyssEngineLabel string

	projectOpened := a.project != nil
	enginePathSet := len(a.config.AbyssEnginePath) > 0

	if a.abyssWrapper.IsRunning() {
		runAbyssEngineLabel = "Stop Abyss Engine"
	} else {
		runAbyssEngineLabel = "Run in Abyss Engine"
	}

	menuLayout := g.Layout{
		g.Menu("File##MainMenuFile").Layout(g.Layout{
			g.Menu("New##MainMenuFileNew").Layout(g.Layout{
				g.MenuItem("Project...\t\tCtrl+Shift+N##MainMenuFileNewProject").OnClick(a.onNewProjectClicked),
			}),
			g.Menu("Open##MainMenuFileOpen").Layout(g.Layout{
				g.MenuItem("Project...\t\tCtrl+O##MainMenuFileOpenProject").OnClick(a.onOpenProjectClicked),
			}),
			g.MenuItem("Save\t\t\t\t\t\tCtrl+S##MainMenuFileSaveProject").OnClick(a.Save),
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
			g.MenuItem("Preferences...\t\tAlt+P##MainMenuFilePreferences").OnClick(a.onFilePreferencesClicked),
			g.Separator(),
			g.MenuItem("Exit\t\t\t\t\t\t  Alt+Q##MainMenuFileExit").OnClick(func() {
				a.Quit()
				os.Exit(0)
			}),
		}),
		g.Menu("View##MainMenuView").Layout(a.buildViewMenu()),
		g.Menu("Project##MainMenuProject").Layout(g.Layout{
			g.MenuItem(runAbyssEngineLabel + "##MainMenuProjectRun").
				Enabled(projectOpened && enginePathSet).
				OnClick(a.onProjectRunClicked),
			g.Separator(),
			g.MenuItem("Properties...##MainMenuProjectProperties").
				Enabled(projectOpened).
				OnClick(a.onProjectPropertiesClicked),
			g.Separator(),
			g.MenuItem("Export MPQ...##MainMenuProjectExport").
				Enabled(projectOpened).
				OnClick(a.onProjectExportMPQClicked),
		}),
		g.Menu("Help").Layout(g.Layout{
			g.MenuItem("About HellSpawner...\tF1##MainMenuHelpAbout").OnClick(a.onHelpAboutClicked),
			g.Separator(),
			g.MenuItem("GitHub repository").OnClick(func() {
				if err := browser.OpenURL(githubURL); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Join Discord server").OnClick(func() {
				if err := browser.OpenURL(discordInvitationURL); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Development live stream").OnClick(func() {
				if err := browser.OpenURL(twitchURL); err != nil {
					log.Print(err)
				}
			}),
			g.MenuItem("Support us").OnClick(func() {
				if err := browser.OpenURL(supportURL); err != nil {
					log.Print(err)
				}
			}),
		}),
	}

	if a.focusedEditor != nil {
		a.focusedEditor.UpdateMainMenuLayout(&menuLayout)
	}

	menuBar := g.MainMenuBar().Layout(menuLayout)

	menuBar.Build()
}

func (a *App) buildViewMenu() g.Layout {
	result := make([]g.Widget, 0)

	result = append(result, g.Menu("Tool Windows").Layout(g.Layout{
		g.MenuItem("Project Explorer\tCtrl+Shift+P").
			Selected(a.projectExplorer.Visible).
			Enabled(true).
			OnClick(a.toggleProjectExplorer),

		g.MenuItem("MPQ Explorer\t\tCtrl+Shift+M").
			Selected(a.mpqExplorer.Visible).
			Enabled(a.project != nil).
			OnClick(a.toggleMPQExplorer),

		g.MenuItem("Console\t\t\t\t\tCtrl+Shift+C").
			Selected(a.console.Visible).
			OnClick(a.toggleConsole),
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
	if err != nil || file == "" {
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
	if err != nil || file == "" {
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
	if a.abyssWrapper.IsRunning() {
		if err := a.abyssWrapper.Kill(); err != nil {
			dialog.Message(err.Error()).Error()
		}

		return
	}

	a.console.Show()

	if err := a.abyssWrapper.Launch(a.config, a.console); err != nil {
		dialog.Message(err.Error()).Error()
	}
}

func (a *App) onProjectExportMPQClicked() {
}
