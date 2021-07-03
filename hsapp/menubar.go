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

	"github.com/OpenDiablo2/HellSpawner/hscommon"
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

	mCloseProject := menuItem("MainMenuCloseProject", "Close Project", "")
	mCloseProject.OnClick(a.onCloseProjectClicked).Enabled(a.project != nil)

	mPreferences := menuItem("MainMenuFilePreferences", "Preferences...", "Alt+P")
	mPreferences.OnClick(a.onFilePreferencesClicked)

	mExit := menuItem("MainMenuFile", "Exit", "Alt+Q")
	fnExit := func() {
		a.Quit()
		os.Exit(0)
	}

	m.Layout(
		mNew.Layout(mNewProject),
		mOpen.Layout(mOpenProject),
		a.openRecentProjectMenu(),
		mSaveProject,
		g.Separator(),
		mCloseProject,
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
				if err := a.loadProjectFromFile(projectName); err != nil {
					logErr("could not open recent file %s", err)
				}
			}).Build()
		}
	}

	m.Layout(g.Custom(fnRecent))

	return m
}

func (a *App) renderMainMenuBar() {
	menuLayout := g.Layout{
		a.fileMenu(),
		a.viewMenu(),
		a.projectMenu(),
		a.helpMenu(),
	}

	if a.focusedEditor != nil {
		a.focusedEditor.UpdateMainMenuLayout(&menuLayout)
	}

	menuBar := g.MainMenuBar().Layout(menuLayout)

	menuBar.Build()
}

func (a *App) viewMenu() *g.MenuWidget {
	viewMenu := menu("MainMenu", "View")
	hasProject := a.project != nil

	toolWindows := g.Menu("Tool Windows").Layout(g.Layout{
		g.MenuItem("Project Explorer\tCtrl+Shift+P").
			Selected(a.projectExplorer.Visible && hasProject).
			Enabled(hasProject).
			OnClick(a.toggleProjectExplorer),

		g.MenuItem("MPQ Explorer\t\tCtrl+Shift+M").
			Selected(a.mpqExplorer.Visible && hasProject).
			Enabled(hasProject).
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

func openURL(url string) {
	if err := browser.OpenURL(url); err != nil {
		log.Print(err)
	}
}

func (a *App) onOpenURL(url string) func() {
	return func() {
		openURL(url)
	}
}

func (a *App) helpMenu() *g.MenuWidget {
	menuHelp := menu("MainMenu", "Help")
	menuHelpAbout := menuItem("MainMenuHelp", "About HellSpawner...", "F1").
		OnClick(a.onHelpAboutClicked)
	menuHelpGithub := menuItem("MainMenuHelp", "GitHub repository", "").
		OnClick(a.onOpenURL(githubURL))
	menuHelpDiscord := menuItem("MainMenuHelp", "Join Discord server", "").
		OnClick(a.onOpenURL(discordInvitationURL))
	menuHelpTwitch := menuItem("MainMenuHelp", "Development live stream", "").
		OnClick(a.onOpenURL(twitchURL))
	menuHelpSupport := menuItem("MainMenuHelp", "Support us", "").
		OnClick(a.onOpenURL(supportURL))
	menuHelpBug := menuItem("MainMenuHelp", "Report Bug on GitHub", "").
		OnClick(a.onReportBugClicked)

	return menuHelp.Layout(
		menuHelpAbout,
		g.Separator(),
		menuHelpGithub,
		menuHelpDiscord,
		menuHelpTwitch,
		menuHelpSupport,
		g.Separator(),
		menuHelpBug,
	)
}

func (a *App) onNewProjectClicked() {
	file, err := dialog.File().Filter("HellSpawner Project", "hsp").Save()
	if err != nil || file == "" {
		logErr("could not create new project, %s", err)
	}

	project, err := hsproject.CreateNew(file)
	if err != nil {
		logErr("could not create new project file, %s", err)
	}

	ppath := project.GetProjectFilePath()
	if err := a.loadProjectFromFile(ppath); err != nil {
		logErr("could not load new project from file %s, %s", ppath, err)
	}
}

func (a *App) onOpenProjectClicked() {
	file, err := dialog.File().Filter("HellSpawner Project", "hsp").Load()
	if err != nil || file == "" {
		return
	}

	if err := a.loadProjectFromFile(file); err != nil {
		logErr("could not open project file %s, %s", file, err)
	}
}

func (a *App) onProjectPropertiesClicked() {
	a.projectPropertiesDialog.Show(a.project, a.config)
}

func (a *App) onFilePreferencesClicked() {
	a.preferencesDialog.Show(a.config)
}

func (a *App) onCloseProjectClicked() {
	if save := dialog.Message("Do you want to save current project?").YesNo(); save {
		if err := a.project.Save(); err != nil {
			a.console.BringToFront()
			a.console.SetVisible(true)

			const errSave = "could not save project"
			_, _ = a.console.Write([]byte(errSave))
		}
	}
	}

	a.project = nil

	a.projectExplorer.SetProject(nil)
	a.mpqExplorer.SetProject(nil)
	a.CloseAllOpenWindows()
	a.updateWindowTitle()
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
		log.Println(err)
	}
}

func (a *App) renderEditors() {
	idx := 0
	for idx < len(a.editors) {
		editor := a.editors[idx]
		if !editor.IsVisible() {
			editor.Cleanup()

			if editor.HasFocus() {
				a.focusedEditor = nil
			}

			a.editors = append(a.editors[:idx], a.editors[idx+1:]...)

			continue
		}

		hadFocus := editor.HasFocus()

		editor.Build()

		// if this window didn't have focus before, but it does now,
		// unregister any other window's shortcuts, and register this window's keyboard shortcuts instead
		if !hadFocus && editor.HasFocus() {
			a.InputManager.UnregisterWindowShortcuts()

			editor.RegisterKeyboardShortcuts(a.InputManager)

			a.focusedEditor = editor
		}

		idx++
	}
}

func (a *App) renderWindows() {
	windows := []hscommon.Renderable{
		a.projectExplorer,
		a.mpqExplorer,
		a.console,
		a.preferencesDialog,
		a.aboutDialog,
		a.projectPropertiesDialog,
	}

	for _, tw := range windows {
		if tw.IsVisible() {
			tw.Build()
		}
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
