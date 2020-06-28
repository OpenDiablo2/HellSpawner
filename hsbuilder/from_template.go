package hsbuilder

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

// CreateBuilderFromTemplate creates a builder based on the specified template
func CreateBuilderFromTemplate(template string) *gtk.Builder {
	builder, err := gtk.BuilderNew()

	if err != nil {
		log.Panic(err)
	}

	if err := builder.AddFromString(template); err != nil {
		log.Panic(err)
	}

	return builder
}
