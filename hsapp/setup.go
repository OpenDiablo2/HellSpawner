package hsapp

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsds1editor"

	"github.com/go-gl/glfw/v3.3/glfw"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsdt1editor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsfonttableeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hspalettemapeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsstringtableeditor"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hspreferencesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hscofeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsdc6editor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsdcceditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsfonteditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hspaletteeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hssoundeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hstexteditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsprojectexplorer"
)

func (a *App) setup() error {
	var err error

	// Register the editors
	a.editorConstructors[hsfiletypes.FileTypeText] = hstexteditor.Create
	a.editorConstructors[hsfiletypes.FileTypeAudio] = hssoundeditor.Create
	a.editorConstructors[hsfiletypes.FileTypePalette] = hspaletteeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeDC6] = hsdc6editor.Create
	a.editorConstructors[hsfiletypes.FileTypeDCC] = hsdcceditor.Create
	a.editorConstructors[hsfiletypes.FileTypeCOF] = hscofeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeFont] = hsfonteditor.Create
	a.editorConstructors[hsfiletypes.FileTypeDT1] = hsdt1editor.Create
	a.editorConstructors[hsfiletypes.FileTypePL2] = hspalettemapeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeTBLStringTable] = hsstringtableeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeTBLFontTable] = hsfonttableeditor.Create
	a.editorConstructors[hsfiletypes.FileTypeDS1] = hsds1editor.Create

	// Register the tool windows
	if a.mpqExplorer, err = hsmpqexplorer.Create(a.openEditor, a.config); err != nil {
		return err
	}

	if a.projectExplorer, err = hsprojectexplorer.Create(a.openEditor); err != nil {
		return err
	}

	// Register the dialogs
	if a.aboutDialog, err = hsaboutdialog.Create(a.diabloRegularFont, a.diabloBoldFont, a.fontFixedSmall); err != nil {
		log.Fatal(err)
	}

	a.projectPropertiesDialog = hsprojectpropertiesdialog.Create(a.onProjectPropertiesChanged)
	a.preferencesDialog = hspreferencesdialog.Create(a.onPreferencesChanged)

	// Set up keyboard shortcuts
	glfw.GetCurrentContext().SetKeyCallback(hsinput.HandleInput)
	a.registerGlobalKeyboardShortcuts()

	return nil
}

func (a *App) registerGlobalKeyboardShortcuts() {
	hsinput.RegisterShortcut(a.onNewProjectClicked, g.KeyN, g.ModControl+g.ModShift, true)
	hsinput.RegisterShortcut(a.onOpenProjectClicked, g.KeyO, g.ModControl, true)
	hsinput.RegisterShortcut(a.onFilePreferencesClicked, g.KeyP, g.ModAlt, true)
	hsinput.RegisterShortcut(a.quit, g.KeyQ, g.ModAlt, true)
	hsinput.RegisterShortcut(a.onHelpAboutClicked, g.KeyF1, g.ModNone, true)

	hsinput.RegisterShortcut(a.closeActiveEditor, g.KeyW, g.ModControl, true)
	hsinput.RegisterShortcut(a.closePopups, g.KeyEscape, g.ModNone, true)

	hsinput.RegisterShortcut(a.toggleMPQExplorer, g.KeyM, g.ModControl+g.ModShift, true)
	hsinput.RegisterShortcut(a.toggleProjectExplorer, g.KeyP, g.ModControl+g.ModShift, true)
}
