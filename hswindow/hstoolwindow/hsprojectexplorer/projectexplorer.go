package hsprojectexplorer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"

	"github.com/OpenDiablo2/dialog"

	"github.com/AllenDang/giu/imgui"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
)

const (
	refreshItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/reload.png"
)

type ProjectExplorerFileSelectedCallback func(path *hscommon.PathEntry)

type ProjectExplorer struct {
	hstoolwindow.ToolWindow

	fileSelectedCallback ProjectExplorerFileSelectedCallback
	nodeCache            map[string][]g.Widget
	refreshIconTexture   *g.Texture
}

func Create(fileSelectedCallback ProjectExplorerFileSelectedCallback) (*ProjectExplorer, error) {
	result := &ProjectExplorer{
		nodeCache:            make(map[string][]g.Widget),
		fileSelectedCallback: fileSelectedCallback,
	}
	result.Visible = false

	hscommon.CreateTextureFromFileAsync(refreshItemButtonPath, func(texture *g.Texture) {
		result.refreshIconTexture = texture
	})

	return result, nil
}

func (m *ProjectExplorer) Render(project *hsproject.Project) {
	if !m.Visible {
		return
	}

	g.Window("Project Explorer").IsOpen(&m.Visible).Pos(10, 30).Size(300, 400).Layout(g.Layout{
		g.Line(
			g.Custom(func() {
				imgui.PushStyleColor(imgui.StyleColorButton, imgui.Vec4{})
				imgui.PushStyleColor(imgui.StyleColorBorder, imgui.Vec4{})
				imgui.PushStyleVarVec2(imgui.StyleVarItemSpacing, imgui.Vec2{Y: 4})
				imgui.PushID("ProjectExplorerRefresh")
			}),
			g.ImageButton(m.refreshIconTexture).Size(16, 16).OnClick(func() { m.onRefreshProjectExplorerClicked(project) }),
			g.Tooltip("Refresh the view from the filesystem."),
			g.Custom(func() {
				imgui.PopID()
				imgui.PopStyleVar()
				imgui.PopStyleColorV(2)
			}),
		),
		g.Separator(),
		g.Child("ProjectExplorerProjectTreeContainer").Flags(g.WindowFlagsHorizontalScrollbar).Layout(m.getProjectTreeNodes(project)),
	})
}

func (m *ProjectExplorer) getProjectTreeNodes(project *hsproject.Project) g.Layout {

	if project == nil {
		return []g.Widget{g.Label("No project loaded...")}
	}

	fileStructure := project.GetFileStructure()

	if fileStructure == nil {
		return []g.Widget{g.Label("No file structure detected...")}
	}

	return []g.Widget{m.renderNodes(project.GetFileStructure(), project)}
}

func (m *ProjectExplorer) onRefreshProjectExplorerClicked(project *hsproject.Project) {
	project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onNewFontClicked(pathEntry *hscommon.PathEntry, project *hsproject.Project) {
	project.CreateNewFile(hsfiletypes.FileTypeFont, pathEntry)
}

func (m *ProjectExplorer) renderNodes(pathEntry *hscommon.PathEntry, project *hsproject.Project) g.Widget {

	if !pathEntry.IsDirectory {
		return m.createFileTreeItem(pathEntry, project)
	}

	// File items and empty dirs
	if len(pathEntry.Children) == 0 {
		return m.createDirectoryTreeItem(pathEntry, nil, project)
	}

	widgets := make([]g.Widget, len(pathEntry.Children))

	sortPaths(pathEntry)

	for idx := range pathEntry.Children {
		widgets[idx] = m.renderNodes(pathEntry.Children[idx], project)
	}

	return m.createDirectoryTreeItem(pathEntry, widgets, project)
}

func (m *ProjectExplorer) createFileTreeItem(pathEntry *hscommon.PathEntry, project *hsproject.Project) g.Widget {
	id := "##ProjectExplorerNode_" + pathEntry.FullPath
	var layout g.Layout = make([]g.Widget, 0)

	if pathEntry.IsRenaming {
		layout = g.Layout{
			g.Custom(func() {
				imgui.SetKeyboardFocusHere()
				if imgui.InputTextV("##RenameField_"+pathEntry.FullPath, &pathEntry.Name,
					int(g.InputTextFlagsAutoSelectAll|g.InputTextFlagsEnterReturnsTrue), nil) {
					pathEntry.IsRenaming = false
					m.onFileRenamed(pathEntry, project)
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
			g.MenuItem("Delete...").OnClick(func() { m.onDeleteFileClicked(pathEntry, project) }),
		}),
	)

	return layout
}

func (m *ProjectExplorer) createDirectoryTreeItem(pathEntry *hscommon.PathEntry, layout g.Layout, project *hsproject.Project) g.Widget {
	var id = pathEntry.Name + "##ProjectExplorerNode_" + pathEntry.FullPath

	if pathEntry.IsRenaming {
		return g.Layout{
			g.Custom(func() {
				imgui.SetKeyboardFocusHere()
				if imgui.InputTextV("##RenameField_"+pathEntry.FullPath, &pathEntry.Name,
					int(g.InputTextFlagsAutoSelectAll|g.InputTextFlagsEnterReturnsTrue), nil) {
					pathEntry.IsRenaming = false
					m.onFileRenamed(pathEntry, project)
				}
			}),
		}
	}

	contextMenuLayout := g.Layout{
		g.Menu("New").Layout(g.Layout{
			g.MenuItem("Folder").OnClick(func() { m.onNewFolderClicked(pathEntry, project) }),
			g.MenuItem("Font").OnClick(func() { m.onNewFontClicked(pathEntry, project) }),
		}),
	}

	if !pathEntry.IsRoot {
		contextMenuLayout = append(contextMenuLayout,
			g.Separator(),
			g.MenuItem("Rename").OnClick(func() { m.onRenameFileClicked(pathEntry) }),
			g.MenuItem("Delete Folder...").OnClick(func() { m.onDeleteFolderClicked(pathEntry, project) }),
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

func (m *ProjectExplorer) onDeleteFolderClicked(entry *hscommon.PathEntry, project *hsproject.Project) {
	if !dialog.Message("Are you sure you want to delete:\n%s", entry.FullPath).YesNo() {
		return
	}

	if err := os.RemoveAll(entry.FullPath); err != nil {
		dialog.Message("Could not delete:\n%s", entry.FullPath).Error()
		return
	}

	project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onDeleteFileClicked(entry *hscommon.PathEntry, project *hsproject.Project) {
	if !dialog.Message("Are you sure you want to delete:\n%s", entry.FullPath).YesNo() {
		return
	}
	if err := os.Remove(entry.FullPath); err != nil {
		dialog.Message("Could not delete:\n%s", entry.FullPath).Error()
		return
	}

	project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onRenameFileClicked(entry *hscommon.PathEntry) {
	entry.OldName = entry.Name
	entry.IsRenaming = true
}

func (m *ProjectExplorer) onFileRenamed(entry *hscommon.PathEntry, project *hsproject.Project) {
	if entry.Name == entry.OldName {
		entry.OldName = ""
		return
	}

	if len(entry.Name) == 0 {
		dialog.Message("Cannot rename file:\nFiles cannot have a blank name.").Error()
		entry.Name = entry.OldName
		entry.OldName = ""
		return
	}

	if len(filepath.Ext(entry.Name)) == 0 {
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

	project.InvalidateFileStructure()
}

func (m *ProjectExplorer) onNewFolderClicked(pathEntry *hscommon.PathEntry, project *hsproject.Project) {
	project.CreateNewFolder(pathEntry)
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
