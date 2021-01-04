package hsmpqexplorer

import (
	"log"
	"path/filepath"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
)

type MPQExplorerFileSelectedCallback func(path *hscommon.PathEntry)

type MPQExplorer struct {
	hstoolwindow.ToolWindow
	config               *hsconfig.Config
	fileSelectedCallback MPQExplorerFileSelectedCallback
	nodeCache            []g.Widget
}

//result := make([]g.Widget, len(pathNodes[""].Children))
//
//for idx := range rootNode.Children {
//result[idx] = renderNodes(rootNode.Children[idx], m)
//}

func Create(fileSelectedCallback MPQExplorerFileSelectedCallback, config *hsconfig.Config) (*MPQExplorer, error) {
	result := &MPQExplorer{
		fileSelectedCallback: fileSelectedCallback,
		config:               config,
	}
	result.Visible = false

	return result, nil
}

func (m *MPQExplorer) Render(project *hsproject.Project, config *hsconfig.Config) {
	if !m.Visible {
		return
	}

	g.Window("MPQ Explorer").IsOpen(&m.Visible).Pos(10, 30).Size(300, 400).Layout(g.Layout{
		g.Child("MpqExplorerContent").Border(false).Flags(g.WindowFlagsHorizontalScrollbar).Layout(m.getMpqTreeNodes(project, config)),
	})
}

func (m *MPQExplorer) getMpqTreeNodes(project *hsproject.Project, config *hsconfig.Config) []g.Widget {
	if m.nodeCache != nil {
		return m.nodeCache
	}

	result := make([]g.Widget, len(project.AuxiliaryMPQs))
	for idx := range project.AuxiliaryMPQs {
		mpq, err := d2mpq.Load(filepath.Join(m.config.AuxiliaryMpqPath, project.AuxiliaryMPQs[idx]))
		if err != nil {
			log.Fatal(err)
		}
		nodes := project.GetMPQFileNodes(mpq, config)
		result[idx] = m.renderNodes(nodes)
	}

	m.nodeCache = result
	return result
}

func (m *MPQExplorer) renderNodes(pathEntry *hscommon.PathEntry) g.Widget {
	if !pathEntry.IsDirectory {
		return g.Selectable(pathEntry.Name).OnClick(func() {
			m.fileSelectedCallback(pathEntry)
		})
	}

	widgets := make([]g.Widget, len(pathEntry.Children))

	hscommon.SortPaths(pathEntry)

	for idx := range pathEntry.Children {
		widgets[idx] = m.renderNodes(pathEntry.Children[idx])
	}

	return g.TreeNode(pathEntry.Name).Layout(widgets)
}

func (m *MPQExplorer) Reset() {
	m.nodeCache = nil
}
