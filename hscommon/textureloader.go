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

// TextureLoader represents a texture loader
type TextureLoader struct {
	canLoadTextures bool
	mutex           *sync.Mutex
	loadQueue       *goconcurrentqueue.FIFO
}

// NewTextureLoader creates a new texture loader
func NewTextureLoader() *TextureLoader {
	result := &TextureLoader{}
	result.canLoadTextures = false
	result.mutex = &sync.Mutex{}
	result.loadQueue = goconcurrentqueue.NewFIFO()

	return result
}

// StopLoadingTextures stops loading a texture
func (t *TextureLoader) StopLoadingTextures() {
	t.mutex.Lock()
	t.canLoadTextures = false
	t.mutex.Unlock()
}

// ResumeLoadingTextures resumes loading textures
func (t *TextureLoader) ResumeLoadingTextures() {
	t.mutex.Lock()
	t.canLoadTextures = true
	t.mutex.Unlock()
}

// ProcessTextureLoadRequests proceses texture loading request
func (t *TextureLoader) ProcessTextureLoadRequests() {
	go func() {
		for {
			item, err := t.loadQueue.DequeueOrWaitForNextElement()
			if err != nil {
				break
			}

			for {
				t.mutex.Lock()

				if !t.canLoadTextures {
					t.mutex.Unlock()
					continue
				}
				t.mutex.Unlock()

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
func (t *TextureLoader) CreateTextureFromFileAsync(fileName string, callback func(*g.Texture)) {
	var imageData *image.RGBA

	var err error

	if imageData, err = g.LoadImage(fileName); err != nil {
		log.Fatal(err)
	}

	t.addTextureToLoadQueue(imageData, callback)
}

// CreateTextureFromARGB creates a texture fromo color given
func (t *TextureLoader) CreateTextureFromARGB(rgb *image.RGBA, callback func(*g.Texture)) {
	t.addTextureToLoadQueue(rgb, callback)
}

func (t *TextureLoader) addTextureToLoadQueue(rgb *image.RGBA, callback func(*g.Texture)) {
	err := t.loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      rgb,
		callback: callback,
	})
	if err != nil {
		log.Fatalf("failed to add texture load request to queue: %s", err)
	}
}
