// Package hsfont contains data for font file types
package hsfont

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	newFilePerms = 0644
)

// Font represents font
type Font struct {
	filePath    string
	TableFile   string
	SpriteFile  string
	PaletteFile string
}

// NewFile creates a new font
func NewFile(filePath string) (*Font, error) {
	result := &Font{
		filePath: filePath,
	}

	if err := result.SaveToFile(); err != nil {
		return nil, err
	}

	return result, nil
}

// LoadFromJSON loads a new font from json
func LoadFromJSON(data []byte) (*Font, error) {
	var font *Font = &Font{}

	err := json.Unmarshal(data, font)

	return font, fmt.Errorf("cannot unmarshal font: %w", err)
}

// JSON exports font to json
func (f *Font) JSON() ([]byte, error) {
	data, err := json.MarshalIndent(f, "", "   ")

	return data, fmt.Errorf("cannot marshal Font: %w", err)
}

// SaveToFile saves font
func (f *Font) SaveToFile() error {
	var data []byte

	var err error

	data, err = f.JSON()
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(f.filePath, data, os.FileMode(newFilePerms)); err != nil {
		return fmt.Errorf("cannot write to file %s: %w", f.filePath, err)
	}

	return nil
}
