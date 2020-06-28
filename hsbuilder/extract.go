package hsbuilder

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"

	"github.com/gotk3/gotk3/gtk"
)

// ExtractWindow extracts a window from a builder
func ExtractWindow(builder *gtk.Builder, id string) *gtk.Window {
	obj, err := builder.GetObject(id)

	if err != nil {
		fmt.Printf("failed to extract object with id %s\n", id)
		return nil
	}

	window, ok := obj.(*gtk.Window)

	if !ok {
		fmt.Printf("object %s is not a *gtk.Window\n", id)
		return nil
	}

	return window
}

// ExtractApplicationWindow extracts an application window from a builder
func ExtractApplicationWindow(builder *gtk.Builder, id string, application *gtk.Application) *gtk.ApplicationWindow {
	obj, err := builder.GetObject(id)

	if err != nil {
		fmt.Printf("failed to extract object with id %s\n", id)
		return nil
	}

	window, ok := obj.(*gtk.ApplicationWindow)

	if !ok {
		fmt.Printf("object %s is not a *gtk.ApplicationWindow\n", id)
		return nil
	}

	return window
}

// ExtractWidget returns the widget based on the id
func ExtractWidget(builder *gtk.Builder, id string) glib.IObject {
	obj, err := builder.GetObject(id)

	if err != nil {
		fmt.Printf("failed to extract object with id %s\n", id)
		return nil
	}

	return obj
}
