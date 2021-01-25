package hscommon

import (
	"image"
	"log"
	"sync"

	"github.com/enriquebris/goconcurrentqueue"

	g "github.com/ianling/giu"
)

// TextureLoadRequestItem represents texture request item
type TextureLoadRequestItem struct {
	rgb      *image.RGBA
	callback func(*g.Texture)
}

var canLoadTextures = false
var mutex = &sync.Mutex{}
var loadQueue = goconcurrentqueue.NewFIFO()

// StopLoadingTextures stops loading a texture
func StopLoadingTextures() {
	mutex.Lock()
	canLoadTextures = false
	mutex.Unlock()
}

// ResumeLoadingTextures resumes loading textures
func ResumeLoadingTextures() {
	mutex.Lock()
	canLoadTextures = true
	mutex.Unlock()
}

// ProcessTextureLoadRequests proceses texture loading request
func ProcessTextureLoadRequests() {
	go func() {
		for {
			item, err := loadQueue.DequeueOrWaitForNextElement()
			if err != nil {
				break
			}
			for {
				mutex.Lock()

				if !canLoadTextures {
					mutex.Unlock()
					continue
				}
				mutex.Unlock()
				break
			}

			loadRequest := item.(TextureLoadRequestItem)
			var texture *g.Texture

			if texture, err = g.NewTextureFromRgba(loadRequest.rgb); err != nil {
				log.Fatal(err)
			}

			loadRequest.callback(texture)
		}
	}()
}

// CreateTextureFromFileAsync creates an texture
func CreateTextureFromFileAsync(fileName string, callback func(*g.Texture)) {
	var imageData *image.RGBA
	var err error

	if imageData, err = g.LoadImage(fileName); err != nil {
		log.Fatal(err)
	}

	addTextureToLoadQueue(imageData, callback)
}

// CreateTextureFromARGB creates a texture fromo color given
func CreateTextureFromARGB(rgb *image.RGBA, callback func(*g.Texture)) {
	addTextureToLoadQueue(rgb, callback)
}

func addTextureToLoadQueue(rgb *image.RGBA, callback func(*g.Texture)) {
	err := loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      rgb,
		callback: callback,
	})
	if err != nil {
		log.Fatalf("failed to add texture load request to queue: %s", err)
	}
}
