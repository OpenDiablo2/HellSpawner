// Package hstexteditor contains text editor's data
package hstexteditor

import (
	"log"
	"strings"

	g "github.com/AllenDang/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 400, 300
	tableViewModW            = 80
	maxTableColumns          = 64
)

// static check, to ensure, if text editor implemented editoWindow
var _ hscommon.EditorWindow = &TextEditor{}

// TextEditor represents a text editor
type TextEditor struct {
	*hseditor.Editor

	text      string
	tableView bool
	tableRows []*g.TableRowWidget
	columns   int
}

// Create creates a new text editor
func Create(_ *hsconfig.Config,
	_ hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	_ []byte,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	result := &TextEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		text:   string(*data),
	}

	if w, h := result.CurrentSize(); w == 0 || h == 0 {
		result.Size(mainWindowW, mainWindowH)
	}

	lines := strings.Split(result.text, "\n")
	firstLine := lines[0]
	result.tableView = strings.Count(firstLine, "\t") > 0

	if !result.tableView {
		return result, nil
	}

	result.tableRows = make([]*g.TableRowWidget, len(lines))

	columns := strings.Split(firstLine, "\t")

	result.columns = len(columns)
	if result.columns > maxTableColumns {
		result.columns = maxTableColumns
		columns = columns[:maxTableColumns]

		log.Print("Waring: Table is too wide (more than 64 columns)! Only first 64 columns will be displayed" +
			"See: https://github.com/ocornut/imgui/issues/3572")
	}

	columnWidgets := make([]g.Widget, result.columns)

	for idx := range columns {
		columnWidgets[idx] = g.Label(columns[idx])
	}

	result.tableRows[0] = g.TableRow(columnWidgets...)

	for lineIdx := range lines[1:] {
		columns := strings.Split(lines[lineIdx+1], "\t")
		columnWidgets := make([]g.Widget, len(columns))

		for idx := range columns {
			columnWidgets[idx] = g.Label(columns[idx])
		}

		result.tableRows[lineIdx+1] = g.TableRow(columnWidgets...)
	}

	return result, nil
}

// Build builds an editor
func (e *TextEditor) Build() {
	if !e.tableView {
		e.IsOpen(&e.Visible).
			Layout(
				g.InputTextMultiline(&e.text).
					Flags(g.InputTextFlagsAllowTabInput),
			)
	} else {
		e.IsOpen(&e.Visible).
			Flags(g.WindowFlagsHorizontalScrollbar).
			Layout(
				g.Child().Border(false).Size(float32(e.columns*tableViewModW), 0).Layout(
					g.Table().FastMode(true).Freeze(0, 1).Rows(e.tableRows...),
				),
			)
	}
}

// UpdateMainMenuLayout updates mainMenu layout to it contains editor's options
func (e *TextEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Text Editor").Layout(g.Layout{
		g.MenuItem("Save\t\t\t\tCtrl+Shift+S").OnClick(e.Save),
		g.Separator(),
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Cleanup()
		}),
	})

	*l = append(*l, m)
}

// GenerateSaveData generates data to be saved
func (e *TextEditor) GenerateSaveData() []byte {
	data := []byte(e.text)

	return data
}

// Save saves an editor
func (e *TextEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *TextEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
