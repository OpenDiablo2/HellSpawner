// Package hsfonteditor contains font editor's data
package hsfonteditor

import (
	"fmt"

	"github.com/OpenDiablo2/dialog"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes/hsfont"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 400, 300
	pathSize                 = 245
	browseW, browseH         = 30, 0
)

// FontEditor represents a font editor
type FontEditor struct {
	*hseditor.Editor
	*hsfont.Font
}

// Create creates a new font editor
func Create(_ *hscommon.TextureLoader, pathEntry *hscommon.PathEntry, data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	font, err := hsfont.LoadFromJSON(*data)
	if err != nil {
		return nil, err
	}

	result := &FontEditor{
		Editor: hseditor.New(pathEntry, x, y, project),
		Font:   font,
	}

	return result, nil
}

// Build builds an editor
func (e *FontEditor) Build() {
	e.IsOpen(&e.Visible).Size(mainWindowW, mainWindowH).Layout(g.Layout{
		g.Label("DC6 Path"),
		g.Line(
			g.InputText("##FontEditorDC6Path", &e.SpriteFile).Size(pathSize).Flags(g.InputTextFlagsReadOnly),
			g.Button("...##FontEditorDC6Browse").Size(browseW, browseH).OnClick(e.onBrowseDC6PathClicked),
		),
		g.Separator(),
		g.Label("TBL Path"),
		g.Line(
			g.InputText("##FontEditorTBLPath", &e.TableFile).Size(pathSize).Flags(g.InputTextFlagsReadOnly),
			g.Button("...##FontEditorTBLBrowse").Size(browseW, browseH).OnClick(e.onBrowseTBLPathClicked),
		),
		g.Separator(),
		g.Label("PL2 Path"),
		g.Line(
			g.InputText("##FontEditorPL2Path", &e.PaletteFile).Size(pathSize).Flags(g.InputTextFlagsReadOnly),
			g.Button("...##FontEditorPL2Browse").Size(browseW, browseH).OnClick(e.onBrowsePL2PathClicked),
		),
	})
}

func (e *FontEditor) onBrowseDC6PathClicked() {
	path := dialog.File().SetStartDir(e.Project.GetProjectFileContentPath())
	path.Filter("DC6 File", "dc6", "DC6")

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	e.SpriteFile = filePath
}

func (e *FontEditor) onBrowseTBLPathClicked() {
	path := dialog.File().SetStartDir(e.Project.GetProjectFileContentPath())
	path.Filter("TBL File", "tbl", "TBL")

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	e.TableFile = filePath
}

func (e *FontEditor) onBrowsePL2PathClicked() {
	path := dialog.File().SetStartDir(e.Project.GetProjectFileContentPath())
	path.Filter("PL2 File", "pl2", "PL2")

	filePath, err := path.Load()

	if err != nil || filePath == "" {
		return
	}

	e.PaletteFile = filePath
}

// UpdateMainMenuLayout updates main menu layout to it contains editors options
func (e *FontEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Editor").Layout(g.Layout{
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
func (e *FontEditor) GenerateSaveData() []byte {
	data, err := e.JSON()
	if err != nil {
		fmt.Println("failed to marshal font to JSON:, ", err)
		return nil
	}

	return data
}

// Save saves an editor
func (e *FontEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *FontEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
