package hsmpqexplorer

import (
	"log"
	"path/filepath"
	"sync"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"

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

func Create(fileSelectedCallback MPQExplorerFileSelectedCallback, config *hsconfig.Config, x, y float32) (*MPQExplorer, error) {
	result := &MPQExplorer{
		ToolWindow:           hstoolwindow.New("MPQ Explorer", hsstate.ToolWindowTypeMPQExplorer, x, y),
		fileSelectedCallback: fileSelectedCallback,
		config:               config,
	}

	return result, nil
}

func (m *MPQExplorer) SetProject(project *hsproject.Project) {
	m.project = project
}

func (m *MPQExplorer) Build() {
	if m.project == nil {
		return
	}

	m.IsOpen(&m.Visible).
		Size(300, 400).
		Layout(g.Layout{
			g.Child("MpqExplorerContent").
				Border(false).
				Flags(g.WindowFlagsHorizontalScrollbar).
				Layout(m.getMpqTreeNodes()),
		})
}

func (m *MPQExplorer) getMpqTreeNodes() []g.Widget {
	if m.nodeCache != nil {
		return m.nodeCache
	}

	wg := sync.WaitGroup{}
	result := make([]g.Widget, len(m.project.AuxiliaryMPQs))
	wg.Add(len(m.project.AuxiliaryMPQs))

	for mpqIndex := range m.project.AuxiliaryMPQs {
		go func(idx int) {
			mpq, err := d2mpq.FromFile(filepath.Join(m.config.AuxiliaryMpqPath, m.project.AuxiliaryMPQs[idx]))
			if err != nil {
				log.Fatal("failed to load mpq: ", err)
			}
			nodes := m.project.GetMPQFileNodes(mpq, m.config)
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
