package hsconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type AppConfig struct {
	Colors AppColorConfig `json:"colors"`
	Fonts  FontConfig     `json:"fonts"`
}

// Save saves the configuration object to disk
func (c *AppConfig) Save(path string) error {
	configDir := filepath.Dir(path)
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return err
	}

	configFile, err := os.Create(path)
	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	if _, err := configFile.Write(buf); err != nil {
		return err
	}

	return configFile.Close()
}

type AppColorConfig struct {
	WindowBackground     []uint8 `json:"windowBackground"`
	Text                 []uint8 `json:"text"`
	Primary              []uint8 `json:"primary"`
	PrimaryHighlight     []uint8 `json:"primaryHighlight"`
	Disabled             []uint8 `json:"disabled"`
	DisabledText         []uint8 `json:"disabledText"`
	Tab                  []uint8 `json:"tab"`
	TabSelected          []uint8 `json:"tabSelected"`
	TabHighlight         []uint8 `json:"tabHighlight"`
	TabSelectedHighlight []uint8 `json:"tabSelectedHighlight"`
}

type FontConfig struct {
	Normal     FontItemConfig `json:"normal"`
	Symbols    FontItemConfig `json:"symbols"`
	Monospaced FontItemConfig `json:"monospaced"`
	Info       FontItemConfig `json:"info"`
}

type FontItemConfig struct {
	Face string `json:"face"`
	Size int    `json:"size"`
}
