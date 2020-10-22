package main

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hsapp"
)

func main() {
	var app *hsapp.App = nil

	var err error

	// Create the HellSpawner app instance
	if app, err = hsapp.Create(); err != nil {
		log.Fatal(err)
	}

	// Run the HellSpawner app
	if err = app.Run(); err != nil {
		log.Fatal(err)
	}
}
