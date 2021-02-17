// Package hsfonttableeditor represents fontTableEditor's window
package hsfonttableeditor

import (
	"fmt"
	"sort"

	"github.com/OpenDiablo2/dialog"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 400, 300
)

// static check, to ensure, if font table editor implemented editoWindow
var _ hscommon.EditorWindow = &FontTableEditor{}

// FontTableEditor represents font table editor
type FontTableEditor struct {
	*hseditor.Editor
	fontTable *d2font.Font
	rows      g.Rows
}

// Create creates a new font table editor
func Create(_ *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	table, err := d2font.Load(*data)
	if err != nil {
		return nil, err
	}

	editor := &FontTableEditor{
		Editor:    hseditor.New(pathEntry, x, y, project),
		fontTable: table,
	}

	return editor, nil
}

func (e *FontTableEditor) init() {
	e.rows = make(g.Rows, 0)

	e.rows = append(e.rows, g.Row(
		g.Label("Index"),
		g.Label("Character"),
		g.Label("Width (px)"),
	))

	// we need to get keys from map[rune]*d2font.fontGlyph
	// and then sort them
	chars := make([]rune, len(e.fontTable.Glyphs))

	// reading runes from map
	idx := 0

	for r := range e.fontTable.Glyphs {
		chars[idx] = r
		idx++
	}

	// sorting runes
	sort.Slice(chars, func(i, j int) bool {
		return e.fontTable.Glyphs[chars[i]].FrameIndex() < e.fontTable.Glyphs[chars[j]].FrameIndex()
	})

	for _, idx := range chars {
		e.rows = append(e.rows, e.makeGlyphLayout(idx))
	}
}

// Build builds a font table editor's window
func (e *FontTableEditor) Build() {
	if e.rows == nil {
		e.init()
		return
	}

	tableLayout := g.Layout{g.Child("").
		Border(false).
		Layout(
			g.Layout{
				g.FastTable("").Border(true).Rows(e.rows),
			},
		)}

	e.IsOpen(&e.Visible).
		Flags(g.WindowFlagsHorizontalScrollbar).
		Size(mainWindowW, mainWindowH).
		Layout(tableLayout)
}

/*func (e *FontTableEditor) makeGlyphLayout(glyph interface {
	FrameIndex() int
	Size() (int, int)
}) *g.RowWidget {*/
func (e *FontTableEditor) makeGlyphLayout(r rune) *g.RowWidget {
	if e.fontTable.Glyphs[r] == nil {
		return &g.RowWidget{}
	}

	w, _ := e.fontTable.Glyphs[r].Size()

	row := g.Row(
		g.Label(fmt.Sprintf("%d", e.fontTable.Glyphs[r].FrameIndex())),
		g.Label(string(r)),
		g.Label(fmt.Sprintf("%d", w)),
	)

	return row
}

// UpdateMainMenuLayout updates mainMenu layout's to it contain FontTableEditor's options
func (e *FontTableEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Table Editor").Layout(g.Layout{
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
func (e *FontTableEditor) GenerateSaveData() []byte {
	// https://github.com/OpenDiablo2/HellSpawner/issues/181
	data, _ := e.Path.GetFileBytes()

	return data
}

// Save saves an editor
func (e *FontTableEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *FontTableEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
