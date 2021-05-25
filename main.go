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
	}

	app.Run()
}
