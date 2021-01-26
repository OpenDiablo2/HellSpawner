// Package hsprojectexplorer contains project explorer's data
package hsprojectexplorer

import (
	"image/color"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"

	"github.com/OpenDiablo2/dialog"

	"github.com/ianling/imgui-go"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
)

const (
	refreshItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/reload.png"
)

// ProjectExplorerFileSelectedCallback represents callback on project file selected
type ProjectExplorerFileSelectedCallback func(path *hscommon.PathEntry)

// ProjectExplorer represents a project explorer
type ProjectExplorer struct {
	*hstoolwindow.ToolWindow

	project              *hsproject.Project
	fileSelectedCallback ProjectExplorerFileSelectedCallback
	nodeCache            map[string][]g.Widget
	refreshIconTexture   *g.Texture
}

// Create creates a new project explorer
func Create(fileSelectedCallback ProjectExplorerFileSelectedCallback, x, y float32) (*ProjectExplorer, error) {
	result := &ProjectExplorer{
		ToolWindow:           hstoolwindow.New("Project Explorer", hsstate.ToolWindowTypeProjectExplorer, x, y),
		nodeCache:            make(map[string][]g.Widget),
		fileSelectedCallback: fileSelectedCallback,
	}
	result.Visible = false

	hscommon.CreateTextureFromFileAsync(refreshItemButtonPath, func(texture *g.Texture) {
		result.refreshIconTexture = texture
	})

	return result, nil
}

// SetProject sets explored project
func (m *ProjectExplorer) SetProject(project *hsproject.Project) {
	m.project = project
}

// Build builds explorer
func (m *ProjectExplorer) Build() {
	if m.project == nil {
		return
	}

	header := g.Line(
		m.makeRefreshButtonLayout(),
	)

	tree := g.Child("ProjectExplorerProjectTreeContainer").
		Flags(g.WindowFlagsHorizontalScrollbar).
		Layout(m.getProjectTreeNodes())

	m.IsOpen(&m.Visible).
		Size(300, 400).
		Layout(g.Layout{
			header,
			g.Separator(),
			tree,
		})
}

func (m *ProjectExplorer) makeRefreshButtonLayout() g.Layout {
	button := g.ImageButton(m.refreshIconTexture).
		Size(16, 16).
		OnClick(func() {
			m.onRefreshProjectExplorerClicked()
		})

	const tooltipText = "Refresh the view from the filesystem."

	if m.project == nil {
		button.TintColor(color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x20})
	}

	return g.Layout{
		g.Custom(func() {
			imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{})
			imgui.PushStyleColor(imgui.StyleColorBorder, imgui.Vec4{})
			imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{Y: 4})
			imgui.PushID("ProjectExplorerRefresh")
		}),

		button,

		g.Tooltip(tooltipText),

		g.Custom(func() {
			imgui.PopID()
			imgui.PopStyleVar()
			imgui.PopStyleColorV(2)
		}),
	}
}

func (m *ProjectExplorer) getProjectTreeNodes() g.Layout {
	if m.project == nil {
		return []g.Widget{g.Label("No project loaded...")}
	}

	fileStructure := m.project.GetFileStructure()

	if fileStructure == nil {
		return []g.Widget{g.Label("No file structure detected...")}
	}

	return []g.Widget{m.renderNodes(m.project.GetFileStructure())}
}

func (m *ProjectExplorer) onRefreshProjectExplorerClicked() {
	if m.project == nil {
		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onNewFontClicked(pathEntry *hscommon.PathEntry) {
	m.project.CreateNewFile(hsfiletypes.FileTypeFont, pathEntry)
}

func (m *ProjectExplorer) renderNodes(pathEntry *hscommon.PathEntry) g.Widget {
	if !pathEntry.IsDirectory {
		return m.createFileTreeItem(pathEntry)
	}

	// File items and empty dirs
	if len(pathEntry.Children) == 0 {
		return m.createDirectoryTreeItem(pathEntry, nil)
	}

	widgets := make([]g.Widget, len(pathEntry.Children))

	sortPaths(pathEntry)

	for idx := range pathEntry.Children {
		widgets[idx] = m.renderNodes(pathEntry.Children[idx])
	}

	return m.createDirectoryTreeItem(pathEntry, widgets)
}

func (m *ProjectExplorer) createFileTreeItem(pathEntry *hscommon.PathEntry) g.Widget {
	id := "##ProjectExplorerNode_" + pathEntry.FullPath

	var layout g.Layout = make([]g.Widget, 0)

	if pathEntry.IsRenaming {
		layout = g.Layout{
			g.Custom(func() {
				imgui.SetKeyboardFocusHere()
				if imgui.InputTextV("##RenameField_"+pathEntry.FullPath, &pathEntry.Name,
					int(g.InputTextFlagsAutoSelectAll|g.InputTextFlagsEnterReturnsTrue), nil) {
					pathEntry.IsRenaming = false
					m.onFileRenamed(pathEntry)
				}
			}),
		}
	} else {
		layout = append(layout, g.Selectable(pathEntry.Name+id).OnClick(func() {
			m.fileSelectedCallback(pathEntry)
		}))
	}

	layout = append(layout,
		g.ContextMenu("Context"+id).Layout(g.Layout{
			g.MenuItem("Rename").OnClick(func() { m.onRenameFileClicked(pathEntry) }),
			g.MenuItem("Delete...").OnClick(func() { m.onDeleteFileClicked(pathEntry) }),
		}),
	)

	return layout
}

func (m *ProjectExplorer) createDirectoryTreeItem(pathEntry *hscommon.PathEntry, layout g.Layout) g.Widget {
	var id = pathEntry.Name + "##ProjectExplorerNode_" + pathEntry.FullPath

	if pathEntry.IsRenaming {
		return g.Layout{
			g.Custom(func() {
				imgui.SetKeyboardFocusHere()
				if imgui.InputTextV("##RenameField_"+pathEntry.FullPath, &pathEntry.Name,
					int(g.InputTextFlagsAutoSelectAll|g.InputTextFlagsEnterReturnsTrue), nil) {
					pathEntry.IsRenaming = false
					m.onFileRenamed(pathEntry)
				}
			}),
		}
	}

	contextMenuLayout := g.Layout{
		g.Menu("New").Layout(g.Layout{
			g.MenuItem("Folder").OnClick(func() { m.onNewFolderClicked(pathEntry) }),
			g.MenuItem("Font").OnClick(func() { m.onNewFontClicked(pathEntry) }),
		}),
	}

	if !pathEntry.IsRoot {
		contextMenuLayout = append(contextMenuLayout,
			g.Separator(),
			g.MenuItem("Rename").OnClick(func() { m.onRenameFileClicked(pathEntry) }),
			g.MenuItem("Delete Folder...").OnClick(func() { m.onDeleteFolderClicked(pathEntry) }),
		)
	}

	menuLayout := g.Layout{
		g.Custom(func() { imgui.PushID(id) }),
		g.ContextMenu("Context").Layout(contextMenuLayout),
		g.Custom(func() { imgui.PopID() }),
	}

	if layout == nil {
		return g.TreeNode(id).Layout(menuLayout)
	}

	return g.TreeNode(id).Layout(append(menuLayout, layout...))
}

func (m *ProjectExplorer) onDeleteFolderClicked(entry *hscommon.PathEntry) {
	if !dialog.Message("Are you sure you want to delete:\n%s", entry.FullPath).YesNo() {
		return
	}

	if err := os.RemoveAll(entry.FullPath); err != nil {
		dialog.Message("Could not delete:\n%s", entry.FullPath).Error()

		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onDeleteFileClicked(entry *hscommon.PathEntry) {
	if !dialog.Message("Are you sure you want to delete:\n%s", entry.FullPath).YesNo() {
		return
	}

	if err := os.Remove(entry.FullPath); err != nil {
		dialog.Message("Could not delete:\n%s", entry.FullPath).Error()

		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onRenameFileClicked(entry *hscommon.PathEntry) {
	entry.OldName = entry.Name
	entry.IsRenaming = true
}

func (m *ProjectExplorer) onFileRenamed(entry *hscommon.PathEntry) {
	if entry.Name == entry.OldName {
		entry.OldName = ""

		return
	}

	if entry.Name == "" {
		dialog.Message("Cannot rename file:\nFiles cannot have a blank name.").Error()

		entry.Name = entry.OldName

		entry.OldName = ""

		return
	}

	if filepath.Ext(entry.Name) == "" {
		entry.Name += filepath.Ext(entry.OldName)
	}

	if !strings.EqualFold(filepath.Ext(entry.OldName), filepath.Ext(entry.Name)) {
		dialog.Message("Cannot rename file:\nFile extension cannot be changed.").Error()

		entry.Name = entry.OldName

		entry.OldName = ""

		return
	}

	basePath := filepath.Dir(entry.FullPath)

	oldPath := filepath.Join(basePath, entry.OldName)

	newPath := filepath.Join(basePath, entry.Name)

	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		dialog.Message("Cannot rename file:\nAlready exists.").Error()

		entry.Name = entry.OldName
		entry.OldName = ""

		return
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		dialog.Message("Could not rename file:\n" + err.Error()).Error()

		entry.Name = entry.OldName
		entry.OldName = ""

		return
	}

	m.project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onNewFolderClicked(pathEntry *hscommon.PathEntry) {
	m.project.CreateNewFolder(pathEntry)
}

func sortPaths(rootPath *hscommon.PathEntry) {
	sort.Slice(rootPath.Children, func(i, j int) bool {
		if rootPath.Children[i].IsDirectory == rootPath.Children[j].IsDirectory {
			var nameI, nameJ string

			if len(rootPath.Children[i].OldName) > 0 {
				nameI = rootPath.Children[i].OldName
			} else {
				nameI = rootPath.Children[i].Name
			}

			if len(rootPath.Children[j].OldName) > 0 {
				nameJ = rootPath.Children[j].OldName
			} else {
				nameJ = rootPath.Children[j].Name
			}

			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		}

		return rootPath.Children[i].IsDirectory && !rootPath.Children[j].IsDirectory
	})
}
