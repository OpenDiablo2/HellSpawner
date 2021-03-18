// Package hstexteditor contains text editor's data
package hstexteditor

import (
	"strings"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 400, 300
	tableViewModW            = 80
)

// static check, to ensure, if text editor implemented editoWindow
var _ hscommon.EditorWindow = &TextEditor{}

// TextEditor represents a text editor
type TextEditor struct {
	*hseditor.Editor

	text      string
	tableView bool
	tableRows g.Rows
	columns   int
}

// Create creates a new text editor
func Create(_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	result := &TextEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		text:   string(*data),
	}

	result.Path = pathEntry

	lines := strings.Split(result.text, "\n")
	firstLine := lines[0]
	result.tableView = strings.Count(firstLine, "\t") > 0

	if !result.tableView {
		return result, nil
	}

	result.tableRows = make([]*g.RowWidget, len(lines))

	columns := strings.Split(firstLine, "\t")
	result.columns = len(columns)
	columnWidgets := make([]g.Widget, len(columns))

	for idx := range columns {
		columnWidgets[idx] = g.Label(columns[idx])
	}

	result.tableRows[0] = g.Row(columnWidgets...)

	for lineIdx := range lines[1:] {
		columns := strings.Split(lines[lineIdx+1], "\t")
		columnWidgets := make([]g.Widget, len(columns))

		for idx := range columns {
			columnWidgets[idx] = g.Label(columns[idx])
		}

		result.tableRows[lineIdx+1] = g.Row(columnWidgets...)
	}

	return result, nil
}

// Build builds an editor
func (e *TextEditor) Build() {
	if !e.tableView {
		e.IsOpen(&e.Visible).Size(mainWindowW, mainWindowH).Layout(g.Layout{
			g.InputTextMultiline("", &e.text).Size(-1, -1).Flags(g.InputTextFlagsAllowTabInput),
		})
	} else {
		e.IsOpen(&e.Visible).Flags(g.WindowFlagsHorizontalScrollbar).Size(mainWindowW, mainWindowH).Layout(g.Layout{
			g.Child("").Border(false).Size(float32(e.columns*tableViewModW), 0).Layout(g.Layout{
				g.FastTable("").Border(true).Rows(e.tableRows),
			}),
		})
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

// RegisterKeyboardShortcuts adds a local shortcuts for this editor
func (e *TextEditor) RegisterKeyboardShortcuts(inputManager *hsinput.InputManager) {
	// Ctrl+Shift+S saves file
	inputManager.RegisterShortcut(func() {
		e.Save()
	}, g.KeyS, g.ModShift+g.ModControl, false)
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
