package hscommon

import (
	"image"
	"log"
	"sync"

	g "github.com/AllenDang/giu"
)

var mutex sync.Mutex = sync.Mutex{}

func CreateTextureFromFileAsync(fileName string, callback func(*g.Texture)) {
	var texture *g.Texture
	var imageData *image.RGBA
	var err error

	if imageData, err = g.LoadImage(fileName); err != nil {
		log.Fatal(err)
	}

	go func() {
		mutex.Lock()
		if texture, err = g.NewTextureFromRgba(imageData); err != nil {
			log.Fatal(err)
		}
		mutex.Unlock()
		callback(texture)
	}()
}

func CreateTextureFromARGB(rgb *image.RGBA, callback func(*g.Texture)) {
	var texture *g.Texture
	var err error

	go func() {
		mutex.Lock()
		if texture, err = g.NewTextureFromRgba(rgb); err != nil {
			log.Fatal(err)
		}
		mutex.Unlock()
		callback(texture)
	}()
}
