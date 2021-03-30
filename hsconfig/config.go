package hsconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsstate"

	"github.com/kirsle/configdir"
)

const (
	newFileMode = 0o644
)

const (
	maxRecentOpenedProjectsCount = 5
)

// Config represents HellSpawner's config
type Config struct {
	RecentProjects          []string
	AbyssEnginePath         string
	AuxiliaryMpqPath        string
	ExternalListFile        string
	OpenMostRecentOnStartup bool
	ProjectStates           map[string]hsstate.AppState
}

func getConfigPath() string {
	configPath := configdir.LocalConfig("hellspawner")
	if err := configdir.MakePath(configPath); err != nil {
		log.Fatal(err)
	}

	return filepath.Join(configPath, "environment.json")
}

func generateDefaultConfig() *Config {
	result := &Config{
		RecentProjects:          []string{},
		OpenMostRecentOnStartup: true,
		ProjectStates:           make(map[string]hsstate.AppState),
	}

	if err := result.Save(); err != nil {
		log.Fatalf("filed to save config: %s", err)
	}

	return result
}

// Load loads config
func Load() *Config {
	configFile := getConfigPath()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return generateDefaultConfig()
	}

	var err error

	var data []byte

	if data, err = ioutil.ReadFile(filepath.Clean(configFile)); err != nil {
		return generateDefaultConfig()
	}

	result := generateDefaultConfig()
	if err = json.Unmarshal(data, &result); err != nil {
		return generateDefaultConfig()
	}

	return result
}

// Save saves a new config
func (c *Config) Save() error {
	var err error

	var data []byte

	if data, err = json.MarshalIndent(c, "", "   "); err != nil {
		return fmt.Errorf("cannot marshal config: %w", err)
	}

	path := getConfigPath()
	if err := ioutil.WriteFile(path, data, os.FileMode(newFileMode)); err != nil {
		return fmt.Errorf("cannot write config at %s: %w", path, err)
	}

	return nil
}

// AddToRecentProjects adds a path to recent opened projects
func (c *Config) AddToRecentProjects(filePath string) {
	found := false

	for idx := range c.RecentProjects {
		if c.RecentProjects[idx] == filePath {
			found = true

			if idx != 0 {
				old := c.RecentProjects[0]
				c.RecentProjects[0] = filePath
				c.RecentProjects[idx] = old
			}
		}
	}

	if !found {
		recent := []string{filePath}

		for idx := range c.RecentProjects {
			if idx == maxRecentOpenedProjectsCount {
				break
			}

			recent = append(recent, c.RecentProjects[idx])
		}

		c.RecentProjects = recent
	}

	if err := c.Save(); err != nil {
		log.Fatalf("failed to save config: %s", err)
	}
}

// GetAuxMPQs returns paths to auxiliary mpq's
func (c *Config) GetAuxMPQs() []string {
	if c.AuxiliaryMpqPath == "" {
		return []string{}
	}

	result := make([]string, 0)

	err := filepath.Walk(c.AuxiliaryMpqPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".mpq" {
			result = append(result, path)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed to walk path for aux MPQs %s: %s", c.AuxiliaryMpqPath, err)
	}

	return result
}
