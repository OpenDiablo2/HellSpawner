package hsapp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/dialog"
	"github.com/go-gl/glfw/v3.3/glfw"

	"github.com/OpenDiablo2/HellSpawner/abysswrapper"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hspreferencesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsconsole"
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
	sampleDuration   = time.Second / 10

	autoSaveTimer = 120

	logFileSeparator = "-----%v-----\n"
	logFilePerms     = 0o644
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
	masterWindow *g.MasterWindow
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

	TextureLoader hscommon.TextureLoader

	showUsage bool
}

// Create creates new app instance
func Create() (*App, error) {
	result := &App{
		Flags:              &Flags{},
		editors:            make([]hscommon.EditorWindow, 0),
		editorConstructors: make(map[hsfiletypes.FileType]editorConstructor),
		TextureLoader:      hscommon.NewTextureLoader(),
		abyssWrapper:       abysswrapper.Create(),
	}

	if shouldTerminate := result.parseArgs(); shouldTerminate {
		return nil, nil
	}

	result.config = hsconfig.Load(*result.Flags.optionalConfigPath)

	return result, nil
}

// Run runs an app instance
func (a *App) Run() (err error) {
	defer a.Quit() // force-close and save everything (in case of crash)

	// setting up the logging here, as opposed to inside of app.setup(),
	// because of the deferred call to logfile.Close()
	if a.config.LoggingToFile || *a.Flags.logFile != "" {
		path := a.config.LogFilePath
		if *a.Flags.logFile != "" {
			path = *a.Flags.logFile
		}

		a.logFile, err = os.OpenFile(filepath.Clean(path), os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFilePerms)
		if err != nil {
			logErr("Error opening log file at %s: %v", a.config.LogFilePath, err)
		}

		defer func() {
			if logErr := a.logFile.Close(); logErr != nil {
				log.Fatal(logErr)
			}
		}()
	}

	err = a.setup()
	if err != nil {
		return err
	}

	if a.config.OpenMostRecentOnStartup && len(a.config.RecentProjects) > 0 {
		err = a.loadProjectFromFile(a.config.RecentProjects[0])
		if err != nil {
			return err
		}
	}

	a.masterWindow.Run(a.render)

	return nil
}

func (a *App) render() {
	a.TextureLoader.StopLoadingTextures()

	a.renderMainMenuBar()
	a.renderEditors()
	a.renderWindows()

	g.Update()

	a.TextureLoader.ResumeLoadingTextures()
}

func logErr(fmtErr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtErr, args...)
	log.Print(msg)
	dialog.Message(msg).Error()
}

func (a *App) createEditor(path *hscommon.PathEntry, state []byte, x, y, w, h float32) {
	data, err := path.GetFileBytes()
	if err != nil {
		const fmtErr = "Could not load file: %v"

		logErr(fmtErr, err)

		return
	}

	fileType, err := hsfiletypes.GetFileTypeFromExtension(filepath.Ext(path.FullPath), &data)
	if err != nil {
		const fmtErr = "Error reading file type: %v"

		logErr(fmtErr, err)

		return
	}

	if a.editorConstructors[fileType] == nil {
		const fmtErr = "Error opening editor: %v"

		logErr(fmtErr, err)

		return
	}

	editor, err := a.editorConstructors[fileType](a.config, a.TextureLoader, path, state, &data, x, y, a.project)
	if err != nil {
		const fmtErr = "Error creating editor: %v"

		logErr(fmtErr, err)

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

	// w, h = 0, because we're creating a new editor,
	// width and height aren't saved, so we give 0 and
	// editors without AutoResize flag sets w, h to default
	a.createEditor(path, nil, editorWindowDefaultX, editorWindowDefaultY, 0, 0)
}

func (a *App) loadProjectFromFile(file string) error {
	project, err := hsproject.LoadFromFile(file)
	if err != nil {
		return fmt.Errorf("could not load project from file %s, %w", file, err)
	}

	err = project.ValidateAuxiliaryMPQs(a.config)
	if err != nil {
		return fmt.Errorf("could not validate aux mpq's, %w", err)
	}

	a.project = project
	a.config.AddToRecentProjects(file)
	a.updateWindowTitle()

	err = a.reloadAuxiliaryMPQs()
	if err != nil {
		return err
	}

	a.projectExplorer.SetProject(a.project)
	a.mpqExplorer.SetProject(a.project)

	a.CloseAllOpenWindows()

	if state, ok := a.config.ProjectStates[a.project.GetProjectFilePath()]; ok {
		a.RestoreAppState(state)
	} else {
		// if we don't have a state saved for this project, just open the project explorer
		a.projectExplorer.Show()
	}

	return nil
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
		logErr("could not save project properties after changing, %s", err)
	}

	a.mpqExplorer.SetProject(a.project)
	a.updateWindowTitle()

	if err := a.reloadAuxiliaryMPQs(); err != nil {
		logErr("could not reload aux mpq's after changing project properties, %s", err)
	}
}

func (a *App) onPreferencesChanged(config *hsconfig.Config) {
	a.config = config
	if err := a.config.Save(); err != nil {
		logErr("after changing preferences, %s", err)
	}

	if a.project == nil {
		return
	}

	if err := a.reloadAuxiliaryMPQs(); err != nil {
		logErr("after changing preferences, %s", err)
	}
}

func (a *App) reloadAuxiliaryMPQs() error {
	if err := a.project.ReloadAuxiliaryMPQs(a.config); err != nil {
		return fmt.Errorf("could not reload aux mpq's in project, %w", err)
	}

	a.mpqExplorer.Reset()

	return nil
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
	a.focusedEditor = nil

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
		logErr("failed to save config: %s", err)
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
