package hsmainwindow

import "C"
import (
	"path/filepath"
	"strings"

	"github.com/OpenDiablo2/D2Shared/d2data/d2mpq"
	"github.com/OpenDiablo2/HellSpawner/hsbuilder"
	"github.com/gotk3/gotk3/gtk"
)

// MainWindow represents the main window of HellSpawner
type MainWindow struct {
	*gtk.ApplicationWindow
	treeView  *gtk.TreeView
	treeStore *gtk.TreeStore
	mpqs      []*d2mpq.MPQ
	mpqPaths  map[string]*gtk.TreeIter
}

// Create creates a new instance of MainWindow
func Create(application *gtk.Application) (*MainWindow, error) {
	builder := hsbuilder.CreateBuilderFromTemplate(template)
	result := &MainWindow{
		ApplicationWindow: hsbuilder.ExtractApplicationWindow(builder, "mainApplicationWindow", application),
		mpqs:              make([]*d2mpq.MPQ, 0),
		mpqPaths:          make(map[string]*gtk.TreeIter),
	}

	result.treeStore = hsbuilder.ExtractWidget(builder, "mpqTreeStore").(*gtk.TreeStore)
	result.treeView = hsbuilder.ExtractWidget(builder, "mainTreeView").(*gtk.TreeView)

	result.wireUpMenuHandlers(builder)

	_, _ = result.Connect("destroy", func() { result.onWindowDestroyed() })

	return result, nil
}

// Append a row to the list store for the tree view
func (m *MainWindow) addRow(parent *gtk.TreeIter, file, path string) *gtk.TreeIter {
	// Get an iterator for a new row at the end of the list store
	iter := m.treeStore.Append(parent)

	// Set the contents of the list store row that the iterator represents
	_ = m.treeStore.SetValue(iter, 0, file)
	_ = m.treeStore.SetValue(iter, 1, path)

	return iter
}

func (m *MainWindow) onWindowDestroyed() {
	gtk.MainQuit()
}

func (m *MainWindow) wireUpMenuHandlers(builder *gtk.Builder) {
	miFileExit := hsbuilder.ExtractWidget(builder, "miFileExit").(*gtk.MenuItem)
	miFileAddExistingMPQ := hsbuilder.ExtractWidget(builder, "miFileAddExistingMPQ").(*gtk.MenuItem)

	_, _ = miFileExit.Connect("activate", func() { m.onFileExit() })
	_, _ = miFileAddExistingMPQ.Connect("activate", func() { m.onFileAddExistingMPQ() })
}

func (m *MainWindow) onFileExit() {
	gtk.MainQuit()
}

func (m *MainWindow) onFileAddExistingMPQ() {
	chooser, _ := gtk.FileChooserNativeDialogNew("Select MPQ(s)...", m, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Open", "Cancel")

	fileFilter, _ := gtk.FileFilterNew()
	fileFilter.AddPattern("*.mpq")

	chooser.AddFilter(fileFilter)
	chooser.SetModal(true)
	chooser.SetSelectMultiple(true)

	if chooser.Run() == int(gtk.RESPONSE_ACCEPT) {
		fileNames, _ := chooser.GetFilenames()
		for fileNameIdx := range fileNames {
			mpq, err := d2mpq.Load(fileNames[fileNameIdx])

			if err != nil {
				continue
			}

			m.mpqs = append(m.mpqs, mpq)
			mpqFileName := filepath.Base(fileNames[fileNameIdx])
			//mpqItem := m.addRow(nil, filepath.Base(fileNames[fileNameIdx]), fileNames[fileNameIdx])
			mpqFiles, _ := mpq.GetFileList()

			for idx := range mpqFiles {
				filePath := filepath.Clean(mpqFileName + "\\" + strings.ToLower(mpqFiles[idx]))
				parentNode := m.getFolderNode(filePath)
				fileParts := strings.Split(mpqFiles[idx], "\\")
				m.addRow(parentNode, fileParts[len(fileParts)-1], mpqFiles[idx])
			}
		}
	}
}

func (m *MainWindow) getFolderNode(path string) *gtk.TreeIter {
	pathParts := strings.Split(path, "\\")
	pathParts = pathParts[:len(pathParts)-1]
	fullPath := ""
	parentPath := ""

	// Ensure folder structure
	for idx := range pathParts {
		fullPath += pathParts[idx] + "\\"
		_, ok := m.mpqPaths[fullPath]

		if !ok {
			m.mpqPaths[fullPath] = m.addRow(m.mpqPaths[parentPath], pathParts[idx], parentPath)
		}

		parentPath += pathParts[idx] + "\\"
	}

	return m.mpqPaths[fullPath]
}

const template = `
	<?xml version="1.0" encoding="UTF-8"?>
	<interface>
		<requires lib="gtk+" version="3.20"/>
		<object class="GtkTreeStore" id="mpqTreeStore">
			<columns>
				<column type="gchararray"/>
				<column type="gchararray"/>
			</columns>
		</object>
		<object class="GtkApplicationWindow" id="mainApplicationWindow">
			<property name="title" translatable="yes">HellSpawner Toolset</property>
			<property name="default-width">300</property>
			<property name="default-height">500</property>
			<child>
				<object class="GtkBox">
					<property name="orientation">vertical</property>
					<child>
						<object class="GtkMenuBar" id="mainMenuBar">
							<child>
								<object class="GtkMenuItem">
									<property name="label">File</property>
									<child type="submenu">
										<object class="GtkMenu">
											<child>
												<object class="GtkMenuItem" id="miFileAddExistingMPQ">
													<property name="label">Add Existing MPQ...</property>
												</object>
											</child>
											<child>
												<object class="GtkMenuItem" id="miFileExit">
													<property name="label">Exit</property>
												</object>
											</child>
										</object>
									</child>
								</object>
							</child>
						</object>
					</child>
					<child>
						<object class="GtkScrolledWindow">
							<child>
								<object class="GtkTreeView" id="mainTreeView">
									<property name="model">mpqTreeStore</property>
									<property name="enable-search">True</property>
									<property name="enable-tree-lines">False</property>
									<property name="headers-visible">False</property>
									<child>
										<object class="GtkTreeViewColumn" id="test-column">
											<property name="title">Name</property>
											<child>
												<object class="GtkCellRendererText" id="test-renderer"/>
												<attributes>
												<attribute name="text">0</attribute>
												</attributes>
											</child>
										</object>
									</child>
								</object>
							</child>
						</object>
						<packing>
							<property name="expand">True</property>
							<property name="fill">True</property>
						</packing>
					</child>
					<child>
						<object class="GtkHSeparator"></object>
					</child>
					<child>
						<object class="GtkToolbar">
							<child>
								<object class="GtkToolButton">
									<property name="icon-name">list-add</property>
								</object>
							</child>
							<child>
								<object class="GtkToolButton">
									<property name="icon-name">list-remove</property>
								</object>
							</child>
							<child>
								<object class="GtkSeparatorToolItem"></object>
							</child>
							<child>
								<object class="GtkToolButton">
									<property name="icon-name">document-properties</property>
								</object>
							</child>
						</object>
					</child>
				</object>
			</child>
		</object>
	</interface>
`
