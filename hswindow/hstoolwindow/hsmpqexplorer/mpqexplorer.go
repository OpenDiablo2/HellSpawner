package hsmpqexplorer

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
)

type MPQExplorerFileSelectedCallback func(path *hscommon.PathEntry)

type MPQExplorer struct {
	hstoolwindow.ToolWindow
	config               hsconfig.Config
	fileSelectedCallback MPQExplorerFileSelectedCallback
	mpqs                 []d2interface.Archive
	nodeCache            map[string][]g.Widget
}

func Create(fileSelectedCallback MPQExplorerFileSelectedCallback, config hsconfig.Config) (*MPQExplorer, error) {
	result := &MPQExplorer{
		nodeCache:            make(map[string][]g.Widget),
		fileSelectedCallback: fileSelectedCallback,
		config:               config,
	}
	result.Visible = false

	return result, nil
}

func (m *MPQExplorer) Render() {
	if !m.Visible {
		return
	}

	g.Window("MPQ Explorer").IsOpen(&m.Visible).Pos(10, 30).Size(300, 400).Layout(g.Layout{
		g.Child("MpqExplorerContent").Border(false).Flags(g.WindowFlagsHorizontalScrollbar).Layout(m.getMpqTreeNodes()),
	})
}

func (m *MPQExplorer) getMpqTreeNodes() []g.Widget {
	result := make([]g.Widget, len(m.mpqs))
	for idx := range m.mpqs {
		result[idx] = g.TreeNode(filepath.Base(m.mpqs[idx].Path())).Layout(m.getMpqFileNodes(m.mpqs[idx]))
	}
	return result
}

// Search for files in MPQ's without listfiles using a list of known filenames
func (m *MPQExplorer) searchForMpqFiles(mpq d2interface.Archive) ([]string, error) {
	var files []string

	if m.config.ExternalListfile != "" {
		file, err := os.Open(m.config.ExternalListfile)
		if err != nil {
			return files, errors.New("Couldn't open listfile")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileName := scanner.Text()
			if mpq.FileExists(fileName) {
				files = append(files, fileName)
			}
		}
		return files, scanner.Err()
	}

	return files, nil
}

func (m *MPQExplorer) getMpqFileNodes(mpq d2interface.Archive) []g.Widget {
	if m.nodeCache[mpq.Path()] != nil {
		return m.nodeCache[mpq.Path()]
	}

	pathNodes := make(map[string]*hscommon.PathEntry)
	rootNode := &hscommon.PathEntry{Name: "/"}
	pathNodes[""] = rootNode

	files, err := mpq.GetFileList()

	if err != nil {
		files, err = m.searchForMpqFiles(mpq)
		if err != nil {
			return []g.Widget{}
		}
	}

	for idx := range files {
		elements := strings.FieldsFunc(files[idx], func(r rune) bool { return r == '\\' || r == '/' })
		path := ""
		for elemIdx := range elements {
			oldPath := path

			path += "/" + elements[elemIdx]
			if pathNodes[strings.ToLower(path)] == nil {
				pathNodes[strings.ToLower(path)] = &hscommon.PathEntry{
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

func sortPaths(rootPath *hscommon.PathEntry) {
	sort.Slice(rootPath.Children, func(i, j int) bool {
		if ((len(rootPath.Children[i].Children) == 0) && (len(rootPath.Children[j].Children) == 0)) ||
			((len(rootPath.Children[i].Children) != 0) && (len(rootPath.Children[j].Children) != 0)) {
			return strings.ToLower(rootPath.Children[i].Name) < strings.ToLower(rootPath.Children[j].Name)
		}

		return len(rootPath.Children[i].Children) > len(rootPath.Children[j].Children)
	})
}

func renderNodes(pathEntry *hscommon.PathEntry, m *MPQExplorer) g.Widget {
	if len(pathEntry.Children) == 0 {
		return g.Selectable(pathEntry.Name).OnClick(func() {
			m.fileSelectedCallback(pathEntry)
		})
	}

	widgets := make([]g.Widget, len(pathEntry.Children))

	sortPaths(pathEntry)

	for idx := range pathEntry.Children {
		widgets[idx] = renderNodes(pathEntry.Children[idx], m)
	}

	return g.TreeNode(pathEntry.Name).Layout(widgets)
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

func (m *MPQExplorer) Reset() {
	m.mpqs = make([]d2interface.Archive, 0)
}
