package hstextfilewindow

import (
	"strings"

	"github.com/gotk3/gotk3/glib"

	"github.com/OpenDiablo2/HellSpawner/hsbuilder"
	"github.com/gotk3/gotk3/gtk"
)

type TextFileWindow struct {
	*gtk.Window
	scrollWindow *gtk.ScrolledWindow
	textData     string
}

// Create creates a new instance of TextFileWindow
func Create(fileName, textData string) *TextFileWindow {
	builder := hsbuilder.CreateBuilderFromTemplate(template)
	result := &TextFileWindow{
		Window:       hsbuilder.ExtractWindow(builder, "textFileWindow"),
		scrollWindow: hsbuilder.ExtractWidget(builder, "swContent").(*gtk.ScrolledWindow),
		textData:     textData,
	}

	lines := strings.Split(textData, "\n")
	columns := strings.Split(lines[0], "\t")

	if len(columns) < 2 {
		result.createTextContent(textData)
	} else {
		result.createTableContent(lines[1:], columns)
	}

	result.Window.SetTitle(fileName)

	return result
}

func (t *TextFileWindow) createTextContent(textData string) {
	textControl, _ := gtk.TextViewNew()
	buffer, _ := textControl.GetBuffer()

	buffer.SetText(textData)
	t.scrollWindow.Add(textControl)
}

func (t *TextFileWindow) createTableContent(lines, columns []string) {
	treeView, _ := gtk.TreeViewNew()

	listTypes := make([]glib.Type, len(columns))
	colIndexes := make([]int, len(columns))

	for colIdx := range columns {
		treeView.AppendColumn(createColumn(strings.TrimSpace(columns[colIdx]), colIdx))

		listTypes[colIdx] = glib.TYPE_STRING
		colIndexes[colIdx] = colIdx
	}

	listStore, _ := gtk.ListStoreNew(listTypes...)
	treeView.SetModel(listStore)

	for lineIdx := range lines {
		cells := strings.Split(lines[lineIdx], "\t")
		iter := listStore.Append()
		items := make([]interface{}, len(cells))

		for cellIdx := range cells {
			items[cellIdx] = strings.TrimSpace(cells[cellIdx])
		}

		_ = listStore.Set(iter, colIndexes, items)
	}

	t.scrollWindow.Add(treeView)
}

func createColumn(title string, id int) *gtk.TreeViewColumn {
	cellRenderer, _ := gtk.CellRendererTextNew()
	column, _ := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)

	return column
}

const template = `
	<?xml version="1.0" encoding="UTF-8"?>
	<interface>
		<requires lib="gtk+" version="3.20"/>
		<object class="GtkWindow" id="textFileWindow">
			<property name="default-width">600</property>
			<property name="default-height">500</property>
			<child>
				<object class="GtkScrolledWindow" id ="swContent">
				</object>
			</child>
		</object>
	</interface>
`
