package hscommon

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
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
type TextureLoader interface {
	// controls
	StopLoadingTextures()
	ResumeLoadingTextures()
	ProcessTextureLoadRequests()

	// creators
	CreateTextureFromARGB(*image.RGBA, func(*g.Texture))
	CreateTextureFromFile([]byte, func(*g.Texture))
}

// textureLoader allows app to load textures and avoid unexpected panics
type textureLoader struct {
	canLoadTextures bool
	mutex           *sync.Mutex
	loadQueue       *goconcurrentqueue.FIFO
}

// NewTextureLoader creates a new texture loader
func NewTextureLoader() TextureLoader {
	result := &textureLoader{}
	result.canLoadTextures = false
	result.mutex = &sync.Mutex{}
	result.loadQueue = goconcurrentqueue.NewFIFO()

	return result
}

// StopLoadingTextures stops loading a texture
func (t *textureLoader) StopLoadingTextures() {
	t.mutex.Lock()
	t.canLoadTextures = false
	t.mutex.Unlock()
}

// ResumeLoadingTextures resumes loading textures
func (t *textureLoader) ResumeLoadingTextures() {
	t.mutex.Lock()
	t.canLoadTextures = true
	t.mutex.Unlock()
}

// ProcessTextureLoadRequests proceses texture loading request
func (t *textureLoader) ProcessTextureLoadRequests() {
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

// CreateTextureFromARGB creates a texture fromo color given
func (t *textureLoader) CreateTextureFromARGB(rgb *image.RGBA, callback func(*g.Texture)) {
	t.addTextureToLoadQueue(rgb, callback)
}

// CreateTextureFromFile creates a texture using io.Reader given
func (t *textureLoader) CreateTextureFromFile(fileData []byte, cb func(*g.Texture)) {
	fileReader := bytes.NewReader(fileData)

	rgba, err := convertToImage(fileReader)
	if err != nil {
		log.Fatal(err)
	}

	t.CreateTextureFromARGB(rgba, cb)
}

func (t *textureLoader) addTextureToLoadQueue(rgb *image.RGBA, callback func(*g.Texture)) {
	err := t.loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      rgb,
		callback: callback,
	})
	if err != nil {
		log.Fatalf("failed to add texture load request to queue: %s", err)
	}
}

func convertToImage(file io.Reader) (*image.RGBA, error) {
	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding png file: %w", err)
	}

	switch trueImg := img.(type) {
	case *image.RGBA:
		return trueImg, nil
	default:
		rgba := image.NewRGBA(trueImg.Bounds())
		draw.Draw(rgba, trueImg.Bounds(), trueImg, image.Pt(0, 0), draw.Src)

		return rgba, nil
	}
}
