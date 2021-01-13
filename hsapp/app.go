package hsapp

import (
	"errors"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/dialog"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hspreferencesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsprojectexplorer"
)

const baseWindowTitle = "HellSpawner"

type App struct {
	project *hsproject.Project
	config  *hsconfig.Config

	aboutDialog             *hsaboutdialog.AboutDialog
	preferencesDialog       *hspreferencesdialog.PreferencesDialog
	projectPropertiesDialog *hsprojectpropertiesdialog.ProjectPropertiesDialog

	projectExplorer *hsprojectexplorer.ProjectExplorer
	mpqExplorer     *hsmpqexplorer.MPQExplorer

	editors            []hscommon.EditorWindow
	editorConstructors map[hsfiletypes.FileType]func(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error)

	fontFixed         imgui.Font
	fontFixedSmall    imgui.Font
	diabloBoldFont    imgui.Font
	diabloRegularFont imgui.Font
}

func Create() (*App, error) {
	result := &App{
		editors:            make([]hscommon.EditorWindow, 0),
		editorConstructors: make(map[hsfiletypes.FileType]func(pathEntry *hscommon.PathEntry, data *[]byte) (hscommon.EditorWindow, error)),
		config:             hsconfig.Load(),
	}

	return result, nil
}

func (a *App) Run() {
	wnd := g.NewMasterWindow(baseWindowTitle, 1280, 720, 0, a.setupFonts)
	wnd.SetBgColor(color.RGBA{10, 10, 10, 255})

	sampleRate := beep.SampleRate(22050)
	if err := speaker.Init(sampleRate, sampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	if a.config.OpenMostRecentOnStartup && len(a.config.RecentProjects) > 0 {
		a.loadProjectFromFile(a.config.RecentProjects[0])
	}

	dialog.Init()
	hscommon.ProcessTextureLoadRequests()
	wnd.Run(a.render)
}

func (a *App) render() {
	hscommon.StopLoadingTextures()
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

	a.projectExplorer.Render(a.project)
	a.mpqExplorer.Render(a.project, a.config)
	a.preferencesDialog.Render()
	a.aboutDialog.Render()
	a.projectPropertiesDialog.Render()

	g.Update()
	hscommon.ResumeLoadingTextures()

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

	imgui.CurrentIO().Fonts().AddFontFromFileTTF("hsassets/fonts/NotoSans-Regular.ttf", 17)
	a.fontFixed = imgui.CurrentIO().Fonts().AddFontFromFileTTF("hsassets/fonts/CascadiaCode.ttf", 15)
	a.fontFixedSmall = imgui.CurrentIO().Fonts().AddFontFromFileTTF("hsassets/fonts/CascadiaCode.ttf", 12)
	a.diabloRegularFont = imgui.CurrentIO().Fonts().AddFontFromFileTTF("hsassets/fonts/DiabloRegular.ttf", 15)
	a.diabloBoldFont = imgui.CurrentIO().Fonts().AddFontFromFileTTF("hsassets/fonts/DiabloBold.ttf", 30)
	imgui.CurrentStyle().ScaleAllSizes(1.0)

	if err := a.setup(); err != nil {
		log.Fatal(err)
	}
}

func (a *App) GetFileBytes(pathEntry *hscommon.PathEntry) ([]byte, error) {
	if pathEntry.Source == hscommon.PathEntrySourceProject {
		if _, err := os.Stat(pathEntry.FullPath); os.IsNotExist(err) {
			return nil, err
		}

		return ioutil.ReadFile(pathEntry.FullPath)
	}

	mpq, err := d2mpq.FromFile(pathEntry.MPQFile)
	if err != nil {
		return nil, err
	}

	if mpq.Contains(pathEntry.FullPath) {
		return mpq.ReadFile(pathEntry.FullPath)
	}

	return nil, errors.New("could not locate file in mpq")
}

func (a *App) openEditor(path *hscommon.PathEntry) {
	uniqueId := path.GetUniqueId()
	for idx := range a.editors {
		if a.editors[idx].GetId() == uniqueId {
			a.editors[idx].BringToFront()
			return
		}
	}

	data, err := a.GetFileBytes(path)
	if err != nil {
		dialog.Message("Could not load file!").Error()
		return
	}

	fileType, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath))
	if err != nil {
		dialog.Message("No file type is defined for this extension!").Error()
		return
	}

	if a.editorConstructors[fileType] == nil {
		dialog.Message("No editor is defined for this file type!").Error()
	}

	go func() {
		editor, err := a.editorConstructors[fileType](path, &data)

		if err != nil {
			dialog.Message("Error creating editor!").Error()
			return
		}

		a.editors = append(a.editors, editor)
		editor.Show()
	}()
}

func (a *App) loadProjectFromFile(file string) {
	var project *hsproject.Project
	var err error

	if project, err = hsproject.LoadFromFile(file); err != nil {
		dialog.Message("Could not load project.").Title("Load HellSpawner Project Error").Error()
		return
	}

	if !project.ValidateAuxiliaryMPQs(a.config) {
		dialog.Message("Could not load project.\nCould not locate one or more auxiliary MPQs!").Title("Load HellSpawner Project Error").Error()
		return
	}

	a.project = project
	a.config.AddToRecentProjects(file)
	a.updateWindowTitle()
	a.reloadAuxiliaryMPQs()
	a.projectExplorer.Show()
}

func (a *App) updateWindowTitle() {
	if a.project == nil {
		glfw.GetCurrentContext().SetTitle(baseWindowTitle)
		return
	}
	glfw.GetCurrentContext().SetTitle(baseWindowTitle + " - " + a.project.ProjectName)
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
	a.reloadAuxiliaryMPQs()
}

func (a *App) onPreferencesChanged(config hsconfig.Config) {
	a.config = &config
	if err := a.config.Save(); err != nil {
		log.Fatal(err)
	}

	if a.project != nil {
		a.reloadAuxiliaryMPQs()
	}
}

func (a *App) reloadAuxiliaryMPQs() {
	a.project.ReloadAuxiliaryMPQs(a.config)
	a.mpqExplorer.Reset()
}

func (a *App) toggleProjectExplorer() {
	a.projectExplorer.ToggleVisibility()
}
