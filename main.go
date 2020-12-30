package main

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hsapp"
)

func main() {
	app, err := hsapp.Create()

	if err != nil {
		log.Fatal(err)
	}

	app.Run()
}
