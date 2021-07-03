package main

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hsapp"
)

func main() {
	log.SetFlags(log.Lshortfile)

	app, err := hsapp.Create()
	if err != nil {
		log.Fatal(err)
	} else if app == nil {
		return // we've terminated early
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
