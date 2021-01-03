package hsfont

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Font struct {
	filePath    string
	TableFile   string
	SpriteFile  string
	PaletteFile string
}

func NewFile(filePath string) *Font {
	result := &Font{
		filePath: filePath,
	}

	return result
}

func (f *Font) SaveToFile() error {
	var data []byte
	var err error

	if data, err = json.MarshalIndent(f, "", "   "); err != nil {
		log.Fatal(err)
	}

	if err = ioutil.WriteFile(f.filePath, data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}
