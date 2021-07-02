package hsutil

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"log"
	"os"
	"path/filepath"

	"github.com/OpenDiablo2/dialog"
)

const milliseconds = 1000

// BoolToInt converts bool into 32-bit intager
// if b is true, then returns 1, else 0
func BoolToInt(b bool) int32 {
	if b {
		return 1
	}

	return 0
}

// Wrap integer to max: wrap(450, 360) == 90
func Wrap(x, max int) int {
	wrapped := x % max

	if wrapped < 0 {
		return max + wrapped
	}

	return wrapped
}

// ExportToGif converts images area to GIF format and saves it under the path selected by user
// tutorial: http://tech.nitoyon.com/en/blog/2016/01/07/go-animated-gif-gen/
func ExportToGif(images []*image.RGBA, delay int32) error {
	filePath, err := dialog.File().Title("Save").Filter("gif images", "gif").Save()
	if err != nil {
		return fmt.Errorf("error reading filepath: %w", err)
	}

	outGif := &gif.GIF{}

	// reload static image and construct outGif
	for _, img := range images {
		// FROM TUTORIAL:
		// Read each frame GIF image with gif.Decode. If we read JPEG images, we have to convert them programatically
		// (goanigiffy does this by calling gif.Encode and gif.Decode).
		g := bytes.NewBuffer([]byte{})

		err := gif.Encode(g, img, nil) // nolint:govet // I want to reuse this ;-)
		if err != nil {
			return fmt.Errorf("error encoding gif: %w", err)
		}

		inGif, err := gif.Decode(g)
		if err != nil {
			return fmt.Errorf("error decoding gif image: %w", err)
		}

		outGif.Image = append(outGif.Image, inGif.(*image.Paletted))
		outGif.Delay = append(outGif.Delay, int(delay/milliseconds))
	}

	// save gif image
	file, err := os.OpenFile(filepath.Clean(filePath), os.O_WRONLY|os.O_CREATE, 0o600)
	if err != nil {
		return fmt.Errorf("error creating a new file: %w", err)
	}

	defer func() {
		err := file.Close() // nolint:govet // I want to re-use err
		if err != nil {
			log.Printf("Error closing file %s: %v", filePath, err)
		}
	}()

	err = gif.EncodeAll(file, outGif)
	if err != nil {
		return fmt.Errorf("error saving to output gif: %w", err)
	}

	return nil
}
