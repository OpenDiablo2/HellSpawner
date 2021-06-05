package hsapp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/OpenDiablo2/dialog"
	"github.com/gravestench/osinfo"
	g "github.com/ianling/giu"
	"github.com/pkg/browser"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
)

const (
	githubURL            = "https://github.com/OpenDiablo2/HellSpawner"
	discordInvitationURL = "https://discord.gg/pRy8tdc"
	twitchURL            = "https://www.twitch.tv/essial/"
	supportURL           = "https://www.patreon.com/bePatron?u=37261055"
)

func (a *App) fileMenu() *g.MenuWidget {
	m := menu("MainMenu", "File")

	mNew := menu("MainMenuFile", "New")
	mNewProject := menuItem("MainMenuFileNew", "Project...", "Ctrl+Shift+N")
	mNewProject.OnClick(a.onNewProjectClicked)

	mOpen := menu("MainMenuFile", "Open")
	mOpenProject := menuItem("MainMenuFileOpen", "Project...", "Ctrl+Shift+O")
	mOpenProject.OnClick(a.onOpenProjectClicked)

	mSaveProject := menuItem("MainMenuFileSaveProject", "Save Project", "Ctrl+S")
	mSaveProject.OnClick(a.Save)

	mPreferences := menuItem("MainMenuFilePreferences", "Preferences...", "Alt+P")
	mPreferences.OnClick(a.onFilePreferencesClicked)

	mExit := menuItem("MainMenuFile", "Exit", "Alt+Q")
	fnExit := func() {
		a.Quit()
		os.Exit(0)
	}

	m.Layout(
		mNew.Layout(
			mNewProject,
		),
		mOpen.Layout(
			mOpenProject,
		),
		a.openRecentProjectMenu(),
		mSaveProject,
		g.Separator(),
		mPreferences,
		g.Separator(),
		mExit.OnClick(fnExit),
	)

	return m
}

func (a *App) openRecentProjectMenu() *g.MenuWidget {
	m := menu("MainMenuFileOpenRecent", "Recent Project...")

	fnRecent := func() {
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
	}

	m.Layout(g.Custom(fnRecent))

	return m
}

func (a *App) renderMainMenuBar() {
	openURL := func(url string) func() {
		return func() {
			if err := browser.OpenURL(url); err != nil {
				log.Print(err)
			}
		}
	}

	menuLayout := g.Layout{
		a.fileMenu(),
		a.viewMenu(),
		a.projectMenu(),
		g.Menu("Help").Layout(g.Layout{
			g.MenuItem("About HellSpawner...\tF1##MainMenuHelpAbout").OnClick(a.onHelpAboutClicked),
			g.Separator(),
			g.MenuItem("GitHub repository").OnClick(openURL(githubURL)),
			g.MenuItem("Join Discord server").OnClick(openURL(discordInvitationURL)),
			g.MenuItem("Development live stream").OnClick(openURL(twitchURL)),
			g.MenuItem("Support us").OnClick(openURL(supportURL)),
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

func (a *App) viewMenu() *g.MenuWidget {
	viewMenu := menu("MainMenu", "View")

	toolWindows := g.Menu("Tool Windows").Layout(g.Layout{
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
	})

	items := []g.Widget{
		toolWindows,
	}

	if len(a.editors) > 0 {
		items = append(items, g.Separator())

		for i := range a.editors {
			editorItem := g.MenuItem(a.editors[i].GetWindowTitle()).OnClick(a.editors[i].BringToFront)
			items = append(items, editorItem)
		}
	}

	return viewMenu.Layout(items...)
}

func (a *App) projectMenu() *g.MenuWidget {
	const (
		runAbyssEngine  = "Run in Abyss Engine"
		stopAbyssEngine = "Stop Abyss Engine"
	)

	projectOpened := a.project != nil
	enginePathSet := len(a.config.AbyssEnginePath) > 0

	label := runAbyssEngine
	if a.abyssWrapper.IsRunning() {
		label = stopAbyssEngine
	}

	projectMenu := menu("MainMenu", "Project")

	projectMenuRun := menuItem("MainMenuProject", label, "").
		Enabled(projectOpened && enginePathSet).
		OnClick(a.onProjectRunClicked)

	projectMenuProperties := menuItem("MainMenuProject", "Properties...", "").
		Enabled(projectOpened).
		OnClick(a.onProjectPropertiesClicked)

	projectMenuExportMPQ := menuItem("MainMenuProject", "Export MPQ...", "").
		Enabled(projectOpened).
		OnClick(a.onProjectExportMPQClicked)

	return projectMenu.Layout(
		projectMenuRun,
		g.Separator(),
		projectMenuProperties,
		g.Separator(),
		projectMenuExportMPQ,
	)
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

// NOTE: some characters in URLs cannot be dirrectly written, because they have
// another meaning (e.g. #). Instead we need to use ASCII code (for # %23).
// for ascii codes see https://www.w3schools.com/tags/ref_urlencode.ASP
func (a *App) onReportBugClicked() {
	osInfo := osinfo.NewOS()

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
		"%23%23 Desktop:",
		"- OS: " + osInfo.Name,
		"- Version: " + osInfo.Version,
		"- Arch: " + osInfo.Arch,
		"- Go version: " + runtime.Version(),
		"",
		"%23%23 Additional context",
		"Add any other context about the problem here.",
		"",
		"%23%23 Config (your config file):",
		"<details><summary>config file</summary><br><pre>",
		strings.ReplaceAll(string(config), "\n", "%0D"),
		"</pre></details>",
	}

	err = browser.OpenURL("https://github.com/OpenDiablo2/HellSpawner/issues/new?body=" + strings.Join(body, "%0D"))
	if err != nil {
		log.Fatal(err)
	}
}

func makeMenuID(name, group, shortcut string) string {
	const (
		sep  = "##"
		fmt2 = "%v%v%v%v"
		fmt3 = "%v\t(%v)%v%v%v"
	)

	if len(shortcut) > 0 {
		return fmt.Sprintf(fmt3, name, shortcut, sep, group, name)
	}

	return fmt.Sprintf(fmt2, name, sep, group, name)
}

func menuID(group, name string) string {
	return makeMenuID(name, group, "")
}

func itemID(group, name, shortcut string) string {
	return makeMenuID(name, group, shortcut)
}

func menu(group, name string) *g.MenuWidget {
	return g.Menu(menuID(group, name))
}

func menuItem(group, name, shortcut string) *g.MenuItemWidget {
	return g.MenuItem(itemID(group, name, shortcut))
}
