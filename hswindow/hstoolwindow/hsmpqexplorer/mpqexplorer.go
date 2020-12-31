package hsmpqexplorer

import (
	"log"
	"path/filepath"
	"sort"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type MPQExplorerFileSelectedCallback func(path *PathEntry)

type PathEntry struct {
	Children []*PathEntry
	Name     string
	FullPath string
}

type MPQExplorer struct {
	hstoolwindow.ToolWindow
	fileSelectedCallback MPQExplorerFileSelectedCallback
	mpqs                 []d2interface.Archive
	nodeCache            map[string][]g.Widget
}

func Create(fileSelectedCallback MPQExplorerFileSelectedCallback) (*MPQExplorer, error) {
	result := &MPQExplorer{
		nodeCache:            make(map[string][]g.Widget),
		fileSelectedCallback: fileSelectedCallback,
	}
	result.Visible = false

	return result, nil
}

func (m *MPQExplorer) Render() {
	if !m.Visible {
		return
	}

	g.WindowV("MPQ Explorer", &m.Visible, g.WindowFlagsNone, 10, 30, 300, 400, g.Layout{
		g.Child("", false, 0, 0, g.WindowFlagsHorizontalScrollbar, m.getMpqTreeNodes()),
	})
}

func (m *MPQExplorer) getMpqTreeNodes() []g.Widget {
	result := make([]g.Widget, len(m.mpqs))
	for idx := range m.mpqs {
		result[idx] = g.TreeNode(filepath.Base(m.mpqs[idx].Path()), g.TreeNodeFlagsNone, m.getMpqFileNodes(m.mpqs[idx]))
	}
	return result
}

func (m *MPQExplorer) getMpqFileNodes(mpq d2interface.Archive) []g.Widget {
	if m.nodeCache[mpq.Path()] != nil {
		return m.nodeCache[mpq.Path()]
	}

	pathNodes := make(map[string]*PathEntry)
	rootNode := &PathEntry{Name: "/"}
	pathNodes[""] = rootNode

	files, err := mpq.GetFileList()

	if err != nil {
		log.Fatal(err)
	}

	for idx := range files {
		elements := strings.FieldsFunc(files[idx], func(r rune) bool { return r == '\\' || r == '/' })
		path := ""
		for elemIdx := range elements {
			oldPath := path

			path += "/" + elements[elemIdx]
			if pathNodes[strings.ToLower(path)] == nil {
				pathNodes[strings.ToLower(path)] = &PathEntry{
					Name:     elements[elemIdx],
					FullPath: mpq.Path() + "|" + path,
				}
				pathNodes[strings.ToLower(oldPath)].Children =
					append(pathNodes[strings.ToLower(oldPath)].Children, pathNodes[strings.ToLower(path)])
			}
		}
	}

	sortPaths(rootNode)

	result := make([]g.Widget, len(pathNodes[""].Children))

	for idx := range rootNode.Children {
		result[idx] = renderNodes(rootNode.Children[idx], m)
	}

	m.nodeCache[mpq.Path()] = result
	return result
}

func sortPaths(rootPath *PathEntry) {
	sort.Slice(rootPath.Children, func(i, j int) bool {
		if ((len(rootPath.Children[i].Children) == 0) && (len(rootPath.Children[j].Children) == 0)) ||
			((len(rootPath.Children[i].Children) != 0) && (len(rootPath.Children[j].Children) != 0)) {
			return strings.ToLower(rootPath.Children[i].Name) < strings.ToLower(rootPath.Children[j].Name)
		}

		return len(rootPath.Children[i].Children) > len(rootPath.Children[j].Children)
	})
}

func renderNodes(pathEntry *PathEntry, m *MPQExplorer) g.Widget {
	if len(pathEntry.Children) == 0 {
		return g.Selectable(pathEntry.Name, func() {
			m.fileSelectedCallback(pathEntry)
		})
	}

	widgets := make([]g.Widget, len(pathEntry.Children))

	sortPaths(pathEntry)

	for idx := range pathEntry.Children {
		widgets[idx] = renderNodes(pathEntry.Children[idx], m)
	}

	return g.TreeNode(pathEntry.Name, g.TreeNodeFlagsNone, widgets)
}

func (m *MPQExplorer) AddMPQ(fileName string) {
	for idx := range m.mpqs {
		if m.mpqs[idx].Path() == fileName {
			return
		}
	}

	data, err := d2mpq.Load(fileName)
	if err != nil {
		log.Fatal(err)
	}
	m.mpqs = append(m.mpqs, data)
}
