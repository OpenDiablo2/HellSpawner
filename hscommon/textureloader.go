package hscommon

import (
	"image"
	"log"
	"sync"
	"time"

	"github.com/enriquebris/goconcurrentqueue"

	g "github.com/AllenDang/giu"
)

type TextureLoadRequestItem struct {
	rgb      *image.RGBA
	callback func(*g.Texture)
}

var canLoadTextures = true
var mutex = &sync.Mutex{}
var loadQueue = goconcurrentqueue.NewFIFO()

func StopLoadingTextures() {
	mutex.Lock()
	canLoadTextures = false
	mutex.Unlock()
}

func ResumeLoadingTextures() {
	mutex.Lock()
	canLoadTextures = true
	mutex.Unlock()
}

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
					time.Sleep(50 * time.Millisecond)
					continue
				}
				mutex.Unlock()
				break
			}

			loadRequest := item.(TextureLoadRequestItem)
			var texture *g.Texture
			log.Printf("Loading texture...")

			if texture, err = g.NewTextureFromRgba(loadRequest.rgb); err != nil {
				log.Fatal(err)
			}

			log.Printf("Done loading texture")
			loadRequest.callback(texture)
		}
	}()
}

func CreateTextureFromFileAsync(fileName string, callback func(*g.Texture)) {
	var imageData *image.RGBA
	var err error

	if imageData, err = g.LoadImage(fileName); err != nil {
		log.Fatal(err)
	}

	loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      imageData,
		callback: callback,
	})
}

func CreateTextureFromARGB(rgb *image.RGBA, callback func(*g.Texture)) {
	loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      rgb,
		callback: callback,
	})
}
