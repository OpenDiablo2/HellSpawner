package hsapp

import (
	"fmt"
	"log"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsconsole"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsds1editor"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsdt1editor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsfonttableeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hspalettemapeditor"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsstringtableeditor"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsaboutdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hspreferencesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog/hsprojectpropertiesdialog"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hsanimdataeditor"
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
	a.editorConstructors[hsfiletypes.FileTypeAnimationData] = hsanimdataeditor.Create
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
	if a.mpqExplorer, err = hsmpqexplorer.Create(a.openEditor, a.config, mpqExplorerDefaultX, mpqExplorerDefaultY); err != nil {
		return fmt.Errorf("error creating a MPQ explorer: %w", err)
	}

	if a.projectExplorer, err = hsprojectexplorer.Create(a.TextureLoader,
		a.openEditor, projectExplorerDefaultX,
		projectExplorerDefaultY); err != nil {
		return fmt.Errorf("error creating a project explorer: %w", err)
	}

	a.console = hsconsole.Create(a.fontFixed, consoleDefaultX, consoleDefaultY, a.logFile)

	log.SetFlags(log.Lshortfile)
	log.SetOutput(a.console)

	t := time.Now()
	y, m, d := t.Date()
	log.Printf(logFileSeparator, fmt.Sprintf("%d-%d-%d, %d:%d:%d", y, m, d, t.Hour(), t.Minute(), t.Second()))

	// Register the dialogs
	if a.aboutDialog, err = hsaboutdialog.Create(a.TextureLoader, a.diabloRegularFont, a.diabloBoldFont, a.fontFixedSmall); err != nil {
		return fmt.Errorf("error creating an about dialog: %w", err)
	}

	a.projectPropertiesDialog = hsprojectpropertiesdialog.Create(a.TextureLoader, a.onProjectPropertiesChanged)
	a.preferencesDialog = hspreferencesdialog.Create(a.onPreferencesChanged, a.masterWindow.SetBgColor)

	// Set up keyboard shortcuts
	a.registerGlobalKeyboardShortcuts()

	return nil
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
	imgui.CurrentIO().Fonts().AddFontFromMemoryTTF(hsassets.FontNotoSansRegular, baseFontSize)
	a.fontFixed = imgui.CurrentIO().Fonts().AddFontFromMemoryTTF(hsassets.FontCascadiaCode, fixedFontSize)
	a.fontFixedSmall = imgui.CurrentIO().Fonts().AddFontFromMemoryTTF(hsassets.FontCascadiaCode, fixedSmallFontSize)
	a.diabloRegularFont = imgui.CurrentIO().Fonts().AddFontFromMemoryTTF(hsassets.FontDiabloRegular, diabloRegularFontSize)
	a.diabloBoldFont = imgui.CurrentIO().Fonts().AddFontFromMemoryTTF(hsassets.FontDiabloBold, diabloBoldFontSize)
	imgui.CurrentStyle().ScaleAllSizes(1)
}

func (a *App) registerGlobalKeyboardShortcuts() {
	a.InputManager.RegisterShortcut(a.onNewProjectClicked, g.KeyN, g.ModControl+g.ModShift, true)
	a.InputManager.RegisterShortcut(a.onOpenProjectClicked, g.KeyO, g.ModControl, true)
	a.InputManager.RegisterShortcut(a.Save, g.KeyS, g.ModControl, true)
	a.InputManager.RegisterShortcut(a.onFilePreferencesClicked, g.KeyP, g.ModAlt, true)
	a.InputManager.RegisterShortcut(a.Quit, g.KeyQ, g.ModAlt, true)
	a.InputManager.RegisterShortcut(a.onHelpAboutClicked, g.KeyF1, g.ModNone, true)

	a.InputManager.RegisterShortcut(a.closeActiveEditor, g.KeyW, g.ModControl, true)
	a.InputManager.RegisterShortcut(func() { a.closePopups(); a.closeActiveEditor() }, g.KeyEscape, g.ModNone, true)

	a.InputManager.RegisterShortcut(a.toggleMPQExplorer, g.KeyM, g.ModControl+g.ModShift, true)
	a.InputManager.RegisterShortcut(a.toggleProjectExplorer, g.KeyP, g.ModControl+g.ModShift, true)
	a.InputManager.RegisterShortcut(a.toggleConsole, g.KeyC, g.ModControl+g.ModShift, true)
}
