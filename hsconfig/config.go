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
	WindowBackground     []int `json:"windowBackground"`
	Text                 []int `json:"text"`
	Primary              []int `json:"primary"`
	PrimaryHighlight     []int `json:"primaryHighlight"`
	Disabled             []int `json:"disabled"`
	DisabledText         []int `json:"disabledText"`
	Tab                  []int `json:"tab"`
	TabSelected          []int `json:"tabSelected"`
	TabHighlight         []int `json:"tabHighlight"`
	TabSelectedHighlight []int `json:"tabSelectedHighlight"`
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
