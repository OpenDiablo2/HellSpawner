package hsmpqexplorer

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow"
)

type MPQExplorerFileSelectedCallback func(path *hscommon.PathEntry)

type MPQExplorer struct {
	*hstoolwindow.ToolWindow
	config               *hsconfig.Config
	project              *hsproject.Project
	fileSelectedCallback MPQExplorerFileSelectedCallback
	nodeCache            []g.Widget
}

func Create(fileSelectedCallback MPQExplorerFileSelectedCallback, config *hsconfig.Config) (*MPQExplorer, error) {
	result := &MPQExplorer{
		ToolWindow:           hstoolwindow.New("MPQ Explorer"),
		fileSelectedCallback: fileSelectedCallback,
		config:               config,
	}

	return result, nil
}

func (m *MPQExplorer) SetProject(project *hsproject.Project) {
	m.project = project
}

func (m *MPQExplorer) Build() {
	m.IsOpen(&m.Visible).
		Pos(10, 30).
		Size(300, 400).
		Layout(g.Layout{
			g.Child("MpqExplorerContent").
				Border(false).
				Flags(g.WindowFlagsHorizontalScrollbar).
				Layout(m.getMpqTreeNodes(m.project, m.config)),
		})
}

func (m *MPQExplorer) getMpqTreeNodes(project *hsproject.Project, config *hsconfig.Config) []g.Widget {
	if m.nodeCache != nil {
		return m.nodeCache
	}

	wg := sync.WaitGroup{}
	result := make([]g.Widget, len(project.AuxiliaryMPQs))
	wg.Add(len(project.AuxiliaryMPQs))

	for mpqIndex := range project.AuxiliaryMPQs {
		go func(idx int) {
			mpq, err := d2mpq.FromFile(filepath.Join(m.config.AuxiliaryMpqPath, project.AuxiliaryMPQs[idx]))
			if err != nil {
				log.Fatal(err)
			}
			nodes := project.GetMPQFileNodes(mpq, config)
			result[idx] = m.renderNodes(nodes)

			wg.Done()
		}(mpqIndex)
	}

	wg.Wait()

	m.nodeCache = result
	return result
}

func (m *MPQExplorer) renderNodes(pathEntry *hscommon.PathEntry) g.Widget {
	if !pathEntry.IsDirectory {
		return g.Selectable(pathEntry.Name).OnClick(func() {
			go m.fileSelectedCallback(pathEntry)
		})
	}

	widgets := make([]g.Widget, len(pathEntry.Children))
	hscommon.SortPaths(pathEntry)

	wg := sync.WaitGroup{}
	wg.Add(len(pathEntry.Children))

	for childIdx := range pathEntry.Children {
		go func(idx int) {
			widgets[idx] = m.renderNodes(pathEntry.Children[idx])
			wg.Done()
		}(childIdx)
	}

	wg.Wait()

	return g.TreeNode(pathEntry.Name).Layout(widgets)
}

func (m *MPQExplorer) Reset() {
	m.nodeCache = nil
}
