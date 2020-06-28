package main

import (
	"log"
	"os"

	"github.com/OpenDiablo2/HellSpawner/hswindows/hsmainwindow"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appID = "com.opendiablo2.hellspawner"

func main() {
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)

	if err != nil {
		log.Fatal("Error creating application", err)
	}

	_, err = application.Connect("activate", func() { onApplicationActivated(application) })

	if err != nil {
		log.Fatal("Error starting application", err)
	}

	os.Exit(application.Run(os.Args))
}

func onApplicationActivated(application *gtk.Application) {
	appWindow, err := hsmainwindow.Create(application)

	if err != nil {
		log.Fatal("Error creating window", err)
	}

	appWindow.ShowAll()

	gtk.Main()
}
