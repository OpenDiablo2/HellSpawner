package hstextfilewindow

import (
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

	result.createTextContent(textData)

	result.Window.SetTitle(fileName)

	return result
}

func (t *TextFileWindow) createTextContent(textData string) {
	textControl, _ := gtk.TextViewNew()
	buffer, _ := textControl.GetBuffer()

	buffer.SetText(textData)
	t.scrollWindow.Add(textControl)
}

const template = `
	<?xml version="1.0" encoding="UTF-8"?>
	<interface>
		<requires lib="gtk+" version="3.20"/>
		<object class="GtkWindow" id="textFileWindow">
			<child>
				<object class="GtkScrolledWindow" id ="swContent">
				</object>
			</child>
		</object>
	</interface>
`
