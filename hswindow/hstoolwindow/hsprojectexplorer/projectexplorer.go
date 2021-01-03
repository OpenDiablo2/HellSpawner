package hsprojectexplorer

import (
	"os"
	"sort"
	"strings"

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

func (m *ProjectExplorer) onNewFontClicked() {

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
	return g.Layout{
		g.Selectable(pathEntry.Name + id).OnClick(func() {
			m.fileSelectedCallback(pathEntry)
		}),
		g.ContextMenu("Context" + id).Layout(g.Layout{
			g.MenuItem("Delete...").OnClick(func() { m.onDeleteFileClicked(pathEntry, project) }),
		}),
	}
}

func (m *ProjectExplorer) createDirectoryTreeItem(pathEntry *hscommon.PathEntry, layout g.Layout, project *hsproject.Project) g.Widget {
	var id = pathEntry.Name + "##ProjectExplorerNode_" + pathEntry.FullPath

	var menuLayout g.Layout

	if pathEntry.IsRoot {
		menuLayout = g.Layout{}
	} else {
		menuLayout = g.Layout{
			g.Custom(func() { imgui.PushID(id) }),
			g.ContextMenu("Context").Layout(g.Layout{
				g.Menu("New").Layout(g.Layout{
					g.MenuItem("Font").OnClick(m.onNewFontClicked),
				}),
				g.MenuItem("Delete Folder...").OnClick(func() { m.onDeleteFolderClicked(pathEntry, project) }),
			}),
			g.Custom(func() { imgui.PopID() }),
		}
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

func sortPaths(rootPath *hscommon.PathEntry) {
	sort.Slice(rootPath.Children, func(i, j int) bool {
		if ((len(rootPath.Children[i].Children) == 0) && (len(rootPath.Children[j].Children) == 0)) ||
			((len(rootPath.Children[i].Children) != 0) && (len(rootPath.Children[j].Children) != 0)) {
			return strings.ToLower(rootPath.Children[i].Name) < strings.ToLower(rootPath.Children[j].Name)
		}

		return len(rootPath.Children[i].Children) > len(rootPath.Children[j].Children)
	})
}
