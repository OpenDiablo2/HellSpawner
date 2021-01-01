package hsapp

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"

	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsdc6editor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hspaletteeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hssoundeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hstexteditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/OpenDiablo2/dialog"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
)

const baseWindowTitle = "HellSpawner"

type App struct {
	masterWindow *g.MasterWindow
	project      *hsproject.Project
	config       *hsconfig.Config

	aboutDialog             *hsaboutdialog.AboutDialog
	projectPropertiesDialog *hsprojectpropertiesdialog.ProjectPropertiesDialog

	mpqExplorer *hsmpqexplorer.MPQExplorer

	editors []hscommon.EditorWindow

	fontFixed         imgui.Font
	fontFixedSmall    imgui.Font
	diabloBoldFont    imgui.Font
	diabloRegularFont imgui.Font
}

func Create() (*App, error) {
	result := &App{}
	result.editors = make([]hscommon.EditorWindow, 0)
	result.config = hsconfig.Load()

	var err error

	if result.mpqExplorer, err = hsmpqexplorer.Create(result.openEditor); err != nil {
		return nil, err
	}

	result.projectPropertiesDialog = hsprojectpropertiesdialog.Create(result.onProjectPropertiesChanged)

	return result, nil
}

func (a *App) Run() {
	wnd := g.NewMasterWindow(baseWindowTitle, 1280, 720, 0, a.setupFonts)
	wnd.SetBgColor(color.RGBA{10, 10, 10, 255})

	sampleRate := beep.SampleRate(22050)
	if err := speaker.Init(sampleRate, sampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	dialog.Init()

	wnd.Main(a.render)
}

func (a *App) render() {
	a.renderMainMenuBar()

	idx := 0
	for idx < len(a.editors) {
		if !a.editors[idx].IsVisible() {
			a.editors[idx].Cleanup()
			a.editors = append(a.editors[:idx], a.editors[idx+1:]...)
			continue
		}

		a.editors[idx].Render()
		idx++
	}

	a.mpqExplorer.Render()
	a.aboutDialog.Render()
	a.projectPropertiesDialog.Render()

	g.Update()
}

func (a *App) loadMpq(fileName string) {
	a.mpqExplorer.AddMPQ(fileName)
	a.mpqExplorer.Show()
}

func (a *App) buildViewMenu() g.Layout {
	result := make([]g.Widget, 0)

	result = append(result, g.Menu("Tool Windows", g.Layout{
		g.MenuItemV("MPQ Explorer", a.mpqExplorer.Visible, true, a.toggleMPQExplorer),
	}))

	if len(a.editors) == 0 {
		return result
	}

	result = append(result, g.Separator())

	for idx := range a.editors {
		i := idx
		result = append(result, g.MenuItem(a.editors[idx].GetWindowTitle(), a.editors[i].BringToFront))
	}

	return result
}

func (a *App) renderMainMenuBar() {
	g.MainMenuBar(g.Layout{
		g.Menu("File##MainMenuFile", g.Layout{
			g.Menu("New##MainMenuFileNew", g.Layout{
				g.MenuItem("Project...##MainMenuFileNewProject", a.onNewProjectClicked),
			}),
			g.Menu("Open##MainMenuFileOpen", g.Layout{
				g.MenuItem("Project...##MainMenuFileOpenProject", a.onOpenProjectClicked),
				g.MenuItem("MPQ...##MainMenuFileOpenMPQ", a.onOpenMpqFileClicked),
			}),
			g.Menu("Open Recent##MainMenuOpenRecent", g.Layout{
				g.Custom(func() {
					if len(a.config.RecentProjects) == 0 {
						g.MenuItemV("No recent projects...##MainMenuOpenRecentItems", false, false, func() {}).Build()
						return
					}
					for idx := range a.config.RecentProjects {
						projectName := a.config.RecentProjects[idx]
						g.MenuItem(fmt.Sprintf("%s##MainMenuOpenRecent_%d", projectName, idx), func() {
							var err error
							var project *hsproject.Project
							if project, err = hsproject.LoadFromFile(projectName); err != nil {
								dialog.Message("Could not load project.").Error()
							}

							a.project = project
							a.config.AddToRecentProjects(projectName)
							a.updateWindowTitle()
						}).Build()
					}
				}),
			}),
			g.Separator(),
			g.MenuItemV("Preferences...##MainMenuFilePreferences", false, false, a.onFilePreferencesClicked),
			g.Separator(),
			g.MenuItem("Exit##MainMenuFileExit", func() { os.Exit(0) }),
		}),
		g.Menu("View##MainMenuView", a.buildViewMenu()),
		g.Menu("Project##MainMenuProject", g.Layout{
			g.MenuItemV("Run in OpenDiablo2##MainMenuProjectRun", false, a.project != nil, a.onProjectRunClicked),
			g.Separator(),
			g.MenuItemV("Properties...##MainMenuProjectProperties", false, a.project != nil, a.onProjectPropertiesClicked),
			g.Separator(),
			g.MenuItemV("Export MPQ...##MainMenuProjectExport", false, a.project != nil, a.onProjectExportMPQClicked),
		}),
		g.Menu("Help", g.Layout{
			g.MenuItem("About HellSpawner...##MainMenuHelpAbout", a.onHelpAboutClicked),
		}),
	}).Build()
}

func (a *App) setupFonts() {
	// Note: To support other languages we'll have to do something with glyph ranges here...
	// ranges := imgui.NewGlyphRanges()
	// rb := imgui.NewFontGlyphRangesBuilder()
	// rb.AddRanges(imgui.CurrentIO().Fonts().GlyphRangesJapanese())
	// rb.AddRanges(imgui.CurrentIO().Fonts().GlyphRangesChineseSimplifiedCommon())
	// rb.AddRanges(imgui.CurrentIO().Fonts().GlyphRangesCyrillic())
	// rb.AddRanges(imgui.CurrentIO().Fonts().GlyphRangesKorean())
	// rb.BuildRanges(ranges)
	// imgui.CurrentIO().Fonts().AddFontFromFileTTFV("NotoSans-Regular.ttf", 17, 0, imgui.CurrentIO().Fonts().GlyphRangesJapanese())

	imgui.CurrentIO().Fonts().AddFontFromFileTTF("NotoSans-Regular.ttf", 17)
	a.fontFixed = imgui.CurrentIO().Fonts().AddFontFromFileTTF("CascadiaCode.ttf", 15)
	a.fontFixedSmall = imgui.CurrentIO().Fonts().AddFontFromFileTTF("CascadiaCode.ttf", 12)
	a.diabloRegularFont = imgui.CurrentIO().Fonts().AddFontFromFileTTF("DiabloRegular.ttf", 15)
	a.diabloBoldFont = imgui.CurrentIO().Fonts().AddFontFromFileTTF("DiabloBold.ttf", 30)
	imgui.CurrentStyle().ScaleAllSizes(1.0)

	var err error
	if a.aboutDialog, err = hsaboutdialog.Create(a.diabloRegularFont, a.diabloBoldFont, a.fontFixedSmall); err != nil {
		log.Fatal(err)
	}
}

func (a *App) openEditor(path *hsmpqexplorer.PathEntry) {
	for idx := range a.editors {
		if a.editors[idx].GetId() == path.FullPath {
			a.editors[idx].BringToFront()
			return
		}
	}

	ext := strings.ToLower(path.FullPath[len(path.FullPath)-4:])
	parts := strings.Split(path.FullPath, "|")
	mpqFile := parts[0]
	filePath := cleanMpqPathName(parts[1])
	mpq, err := d2mpq.Load(mpqFile)

	if err != nil {
		log.Fatal(err)
	}

	switch ext {
	case ".txt":
		text, err := mpq.ReadTextFile(filePath)

		if err != nil {
			log.Fatal(err)
		}

		editor, err := hstexteditor.Create(path.Name, text, a.fontFixed)

		if err != nil {
			log.Fatal(err)
		}

		a.editors = append(a.editors, editor)
		editor.SetId(path.FullPath)
		editor.Show()
	case ".wav":
		audioStream, err := mpq.ReadFileStream(filePath)

		editor, err := hssoundeditor.Create(path.Name, audioStream)

		if err != nil {
			log.Fatal(err)
		}

		a.editors = append(a.editors, editor)
		editor.SetId(path.FullPath)
		editor.Show()
	case ".dat":
		data, err := mpq.ReadFile(filePath)
		if err != nil {
			return
		}

		editor, err := hspaletteeditor.Create(path.Name, path.FullPath, data)

		if err != nil {
			log.Fatal(err)
		}

		a.editors = append(a.editors, editor)
		editor.SetId(path.FullPath)

		editor.Show()
	case ".dc6":
		data, err := mpq.ReadFile(filePath)
		if err != nil {
			return
		}

		editor, err := hsdc6editor.Create(path.Name, path.FullPath, data)

		if err != nil {
			log.Fatal(err)
		}

		a.editors = append(a.editors, editor)
		editor.SetId(path.FullPath)

		editor.Show()
	}

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
	a.project = project
	a.config.AddToRecentProjects(file)
	a.updateWindowTitle()
}

func (a *App) onOpenMpqFileClicked() {
	file, err := dialog.File().Filter("MPQ Archive", "mpq").Load()
	if err != nil || len(file) == 0 {
		return
	}
	a.loadMpq(file)
}

func (a *App) onOpenProjectClicked() {
	file, err := dialog.File().Filter("HellSpawner Project", "hsp").Load()
	if err != nil || len(file) == 0 {
		return
	}
	var project *hsproject.Project
	if project, err = hsproject.LoadFromFile(file); err != nil {
		dialog.Message("Could not load project.").Error()
	}

	a.project = project
	a.config.AddToRecentProjects(file)
	a.updateWindowTitle()
}

func (a *App) onProjectPropertiesClicked() {
	a.projectPropertiesDialog.Show(a.project)
}

func (a *App) updateWindowTitle() {
	if a.project == nil {
		glfw.GetCurrentContext().SetTitle(baseWindowTitle)
		return
	}
	glfw.GetCurrentContext().SetTitle(baseWindowTitle + " - " + a.project.ProjectName)
}

func (a *App) onFilePreferencesClicked() {

}

func (a *App) onHelpAboutClicked() {
	a.aboutDialog.Show()
}

func (a *App) onProjectRunClicked() {

}

func (a *App) onProjectExportMPQClicked() {

}

func (a *App) toggleMPQExplorer() {
	a.mpqExplorer.ToggleVisibility()
}

func (a *App) onProjectPropertiesChanged(project hsproject.Project) {
	a.project = &project
	if err := a.project.Save(); err != nil {
		log.Fatal(err)
	}
	a.updateWindowTitle()
}

func cleanMpqPathName(name string) string {
	name = strings.ReplaceAll(name, "/", "\\")

	if string(name[0]) == "\\" {
		name = name[1:]
	}

	return name
}
