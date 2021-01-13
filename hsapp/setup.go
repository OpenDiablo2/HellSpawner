package hsapp

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsdt1editor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hspalettemapeditor"

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

	return nil
}
