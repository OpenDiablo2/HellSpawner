package hsconfig

import (
	"os"
	"path"
)

const (
	// ConfigDirName is the name of the config directory in the user config dir
	ConfigDirName = "Hellspawner"

	// ConfigFileName is the actual config file name
	ConfigFileName = "config.json"
)

// DefaultConfigPath returns the absolute path for the default config file location
func DefaultConfigPath() string {
	if configDir, err := os.UserConfigDir(); err == nil {
		return path.Join(configDir, ConfigDirName, ConfigFileName)
	}

	return LocalConfigPath()
}

// LocalConfigPath returns the absolute path to the directory of the OpenDiablo2 executable
func LocalConfigPath() string {
	return path.Join(path.Dir(os.Args[0]), ConfigFileName)
}

// DefaultConfig creates and returns a default configuration
func DefaultConfig() *AppConfig {
	config := &AppConfig{
		Colors: AppColorConfig{
			WindowBackground:     []uint8{5, 5, 5, 255},
			Primary:              []uint8{60, 15, 15, 255},
			PrimaryHighlight:     []uint8{80, 15, 15, 255},
			Text:                 []uint8{255, 255, 255, 255},
			Disabled:             []uint8{24, 24, 24, 255},
			DisabledText:         []uint8{128, 128, 128, 255},
			Tab:                  []uint8{100, 30, 30, 255},
			TabSelected:          []uint8{255, 128, 128, 48},
			TabHighlight:         []uint8{128, 128, 128, 255},
			TabSelectedHighlight: []uint8{128, 128, 128, 255},
		},
		Fonts: FontConfig{
			Normal: FontItemConfig{
				Face: "NotoSans-Regular.ttf",
				Size: 10,
			},
			Symbols: FontItemConfig{
				Face: "NotoSansSymbols-Medium.ttf",
				Size: 10,
			},
			Monospaced: FontItemConfig{
				Face: "CascadiaCode.ttf",
				Size: 10,
			},
			Info: FontItemConfig{
				Face: "NotoSans-Regular.ttf",
				Size: 9,
			},
		},
	}

	return config
}
