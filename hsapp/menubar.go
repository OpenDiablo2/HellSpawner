package hsapp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"
	"github.com/pkg/browser"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
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
			g.MenuItem("Report Bug on GitHub##MainMenuHelpBug").OnClick(a.onReportBugClicked),
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

func (a *App) onReportBugClicked() {
	// NOTE: some characters in URLs cannot be dirrectly writen, because they have
	// another meaning (e.g. #). Instead we need to use ASCII code (for # %23).
	// for ascii codes see https://www.w3schools.com/tags/ref_urlencode.ASP

	osInfo := hscommon.NewOS()

	config, err := json.MarshalIndent(a.config, " ", "   ")
	if err != nil {
		log.Printf("Unable to Marshal config file: %v", err)
	}

	// issue's body (from bug_report.md)
	body := []string{
		"%23%23 Describe the bug",
		"A clear and concise description of what the bug is.",
		"",
		"%23%23 To Reproduce",
		"Steps to reproduce the behavior:",
		"1. Go to '...'",
		"2. Click on '....'",
		"3. Scroll down to '....'",
		"4. See error",
		"",
		"%23%23 Expected behavior",
		"A clear and concise description of what you expected to happen.",
		"",
		"%23%23 Screenshots",
		"If applicable, add screenshots to help explain your problem.",
		"",
		"%23%23 Desktop (please complete the following information):",
		"- OS: " + osInfo.Name,
		"- Version: " + osInfo.Version,
		"- Arch: " + osInfo.Arch,
		"- Go version: " + runtime.Version(),
		"",
		"%23%23 Additional context",
		"Add any other context about the problem here.",
		"",
		"%23%23 Config (your config file)",
		"<details><summary>config file</summary><br><pre>",
		strings.ReplaceAll(string(config), "\n", "%0D"),
		"</pre></details>",
	}

	err = browser.OpenURL("https://github.com/OpenDiablo2/HellSpawner/issues/new?body=" + strings.Join(body, "%0D"))
	if err != nil {
		log.Fatal(err)
	}
}
