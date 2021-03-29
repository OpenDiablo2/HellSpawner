//go:generate go install github.com/golangci/golangci-lint/cmd/golangci-lint
//go:generate go install github.com/client9/misspell/cmd/misspell
//go:generate go install golang.org/x/tools/cmd/goimports
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
