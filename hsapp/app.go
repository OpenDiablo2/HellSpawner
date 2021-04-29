package hsapp

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsconsole"

	"github.com/OpenDiablo2/HellSpawner/abysswrapper"

	"github.com/OpenDiablo2/HellSpawner/hsinput"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/dialog"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hspreferencesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsprojectexplorer"
)

const (
	baseWindowTitle          = "HellSpawner"
	baseWindowW, baseWindowH = 1280, 720
	editorWindowDefaultX     = 320
	editorWindowDefaultY     = 30
	projectExplorerDefaultX  = 0
	projectExplorerDefaultY  = 25
	mpqExplorerDefaultX      = 30
	mpqExplorerDefaultY      = 30
	consoleDefaultX          = 10
	consoleDefaultY          = 500

	samplesPerSecond = 22050

	bgColor = 0x0a0a0aff

	autoSaveTimer = 120

	logFileSeparator = "-----%v-----\n"
	logFilePerms     = 0o600
)

const (
	baseFontSize          = 17
	fixedFontSize         = 15
	fixedSmallFontSize    = 12
	diabloRegularFontSize = 15
	diabloBoldFontSize    = 30
)

type editorConstructor func(
	config *hsconfig.Config,
	textureLoader hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	state []byte,
	data *[]byte,
	x, y float32,
	project *hsproject.Project,
) (hscommon.EditorWindow, error)

// App represents an app
type App struct {
	*Flags
	project      *hsproject.Project
	config       *hsconfig.Config
	abyssWrapper *abysswrapper.AbyssWrapper
	logFile      *os.File

	aboutDialog             *hsaboutdialog.AboutDialog
	preferencesDialog       *hspreferencesdialog.PreferencesDialog
	projectPropertiesDialog *hsprojectpropertiesdialog.ProjectPropertiesDialog

	projectExplorer *hsprojectexplorer.ProjectExplorer
	mpqExplorer     *hsmpqexplorer.MPQExplorer
	console         *hsconsole.Console

	editors            []hscommon.EditorWindow
	editorConstructors map[hsfiletypes.FileType]editorConstructor

	editorManagerMutex sync.RWMutex
	focusedEditor      hscommon.EditorWindow

	fontFixed         imgui.Font
	fontFixedSmall    imgui.Font
	diabloBoldFont    imgui.Font
	diabloRegularFont imgui.Font

	InputManager  *hsinput.InputManager
	TextureLoader hscommon.TextureLoader
}

// Create creates new app instance
func Create() (*App, error) {
	tl := hscommon.NewTextureLoader()
	result := &App{
		Flags:              &Flags{},
		editors:            make([]hscommon.EditorWindow, 0),
		editorConstructors: make(map[hsfiletypes.FileType]editorConstructor),
		TextureLoader:      tl,
	}

	im := hsinput.NewInputManager()
	result.InputManager = im

	result.abyssWrapper = abysswrapper.Create()

	result.parseArgs()

	result.config = hsconfig.Load(*result.Flags.optionalConfigPath)

	return result, nil
}

// Run runs an app instance
func (a *App) Run() {
	var err error

	wnd := g.NewMasterWindow(baseWindowTitle, baseWindowW, baseWindowH, 0, a.setupFonts)
	wnd.SetBgColor(hsutil.Color(bgColor))

	sampleRate := beep.SampleRate(samplesPerSecond)

	// nolint:gomnd // this is 0.1 of second
	if err = speaker.Init(sampleRate, sampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	dialog.Init()

	// initialize auto-save timer
	go func() {
		time.Sleep(autoSaveTimer * time.Second)
		a.Save()
	}()

	a.TextureLoader.ProcessTextureLoadRequests()

	defer a.Quit() // force-close and save everything (in case of crash)

	a.logFile, err = os.OpenFile(a.config.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFilePerms)
	if err != nil {
		log.Printf("Error opening log file at %s: %v", a.config.LogFilePath, err)
	}

	defer func() {
		err := a.logFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := a.setup(); err != nil {
		log.Panic(err)
	}

	if a.config.OpenMostRecentOnStartup && len(a.config.RecentProjects) > 0 {
		a.loadProjectFromFile(a.config.RecentProjects[0])
	}

	a.TextureLoader.ProcessTextureLoadRequests()

	wnd.SetInputCallback(a.InputManager.HandleInput)
	wnd.Run(a.render)
}

func (a *App) render() {
	a.TextureLoader.StopLoadingTextures()
	a.renderMainMenuBar()

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

	g.Update()
	a.TextureLoader.ResumeLoadingTextures()
}

func (a *App) createEditor(path *hscommon.PathEntry, state []byte, x, y, w, h float32) {
	data, err := path.GetFileBytes()
	if err != nil {
		log.Printf("Could not load file: %v", err)
		dialog.Message("Could not load file!").Error()

		return
	}

	fileType, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath), &data)
	if err != nil {
		log.Printf("Error reading file type: %v", err)
		dialog.Message("No file type is defined for this extension!").Error()

		return
	}

	if a.editorConstructors[fileType] == nil {
		log.Printf("Error loading editor: %v", err)
		dialog.Message("No editor is defined for this file type!").Error()

		return
	}

	editor, err := a.editorConstructors[fileType](a.config, a.TextureLoader, path, state, &data, x, y, a.project)
	if err != nil {
		log.Printf("Error creating editor: %v", err)
		dialog.Message("Error creating editor: %s", err).Error()

		return
	}

	editor.Size(w, h)

	a.editorManagerMutex.Lock()
	a.editors = append(a.editors, editor)
	a.editorManagerMutex.Unlock()
	editor.Show()
	editor.BringToFront()
}

func (a *App) openEditor(path *hscommon.PathEntry) {
	a.editorManagerMutex.RLock()

	uniqueID := path.GetUniqueID()
	for idx := range a.editors {
		if a.editors[idx].GetID() == uniqueID {
			a.editors[idx].BringToFront()
			a.editorManagerMutex.RUnlock()

			return
		}
	}

	a.editorManagerMutex.RUnlock()

	// w, h = 0, because we're createing a new editor,
	// width and height aren't saved, so we give 0 and
	// editors without AutoResize flag sets w, h to default
	a.createEditor(path, nil, editorWindowDefaultX, editorWindowDefaultY, 0, 0)
}

func (a *App) loadProjectFromFile(file string) {
	var project *hsproject.Project

	var err error

	if project, err = hsproject.LoadFromFile(file); err != nil {
		log.Printf("Error loading project: %v", err)
		dialog.Message("Could not load project.").Title("Load HellSpawner Project Error").Error()

		return
	}

	if !project.ValidateAuxiliaryMPQs(a.config) {
		log.Printf("Error loading mpqs: %v", err)
		dialog.Message("Could not load project.\nCould not locate one or more auxiliary MPQs!").Title("Load HellSpawner Project Error").Error()

		return
	}

	a.project = project
	a.config.AddToRecentProjects(file)
	a.updateWindowTitle()
	a.reloadAuxiliaryMPQs()
	a.projectExplorer.SetProject(a.project)
	a.mpqExplorer.SetProject(a.project)

	a.CloseAllOpenWindows()

	if state, ok := a.config.ProjectStates[a.project.GetProjectFilePath()]; ok {
		a.RestoreAppState(state)
	} else {
		// if we don't have a state saved for this project, just open the project explorer
		a.projectExplorer.Show()
	}
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

func (a *App) onProjectPropertiesChanged(project *hsproject.Project) {
	a.project = project
	if err := a.project.Save(); err != nil {
		log.Fatal(err)
	}

	a.mpqExplorer.SetProject(a.project)
	a.updateWindowTitle()
	a.reloadAuxiliaryMPQs()
}

func (a *App) onPreferencesChanged(config *hsconfig.Config) {
	a.config = config
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

func (a *App) closeActiveEditor() {
	for _, editor := range a.editors {
		if editor.HasFocus() {
			// don't call Cleanup here. the Render loop will call Cleanup when it notices that this editor isn't visible
			editor.SetVisible(false)
			return
		}
	}
}

func (a *App) closePopups() {
	a.projectPropertiesDialog.Cleanup()
	a.aboutDialog.Cleanup()
	a.preferencesDialog.Cleanup()
}

func (a *App) toggleConsole() {
	a.console.ToggleVisibility()
}

// CloseAllOpenWindows closes all opened windows
func (a *App) CloseAllOpenWindows() {
	a.closePopups()
	a.projectExplorer.Cleanup()
	a.mpqExplorer.Cleanup()

	for _, editor := range a.editors {
		editor.Cleanup()
	}
}

// Save saves app state
func (a *App) Save() {
	if a.project != nil {
		a.config.ProjectStates[a.project.GetProjectFilePath()] = a.State()
	}

	if err := a.config.Save(); err != nil {
		log.Print("failed to save config: ", err)
		return
	}

	if a.focusedEditor != nil {
		a.focusedEditor.Save()
	}
}

// Quit quits the app
func (a *App) Quit() {
	if a.abyssWrapper.IsRunning() {
		_ = a.abyssWrapper.Kill()
	}

	a.Save()

	a.CloseAllOpenWindows()
}
