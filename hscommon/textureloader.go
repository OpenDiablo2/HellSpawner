package hscommon

import (
	"image"
	"log"
	"sync"

	"github.com/enriquebris/goconcurrentqueue"

	g "github.com/AllenDang/giu"
)

type TextureLoadRequestItem struct {
	rgb      *image.RGBA
	callback func(*g.Texture)
}

var canLoadTextures = false
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

func CreateTextureFromFileAsync(fileName string, callback func(*g.Texture)) {
	var imageData *image.RGBA
	var err error

	if imageData, err = g.LoadImage(fileName); err != nil {
		log.Fatal(err)
	}

	err = loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      imageData,
		callback: callback,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTextureFromARGB(rgb *image.RGBA, callback func(*g.Texture)) {
	err := loadQueue.Enqueue(TextureLoadRequestItem{
		rgb:      rgb,
		callback: callback,
	})
	if err != nil {
		log.Fatal(err)
	}
}
