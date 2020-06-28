package hsdc6window

import "github.com/gotk3/gotk3/gtk"

type DC6Window struct {
	*gtk.Window
}

const template = `
	<?xml version="1.0" encoding="UTF-8"?>
	<interface>
		<requires lib="gtk+" version="3.20"/>
		<object class="GtkWindow" id="dc6Window">
			<property name="default-width">600</property>
			<property name="default-height">500</property>
			<child>
				<object class="GtkScrolledWindow" id ="swContent">
				</object>
			</child>
		</object>
	</interface>
`
