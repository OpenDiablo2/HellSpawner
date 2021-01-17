package hsconfig

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kirsle/configdir"
)

type Config struct {
	RecentProjects          []string
	AbyssEnginePath         string
	AuxiliaryMpqPath        string
	ExternalListFile        string
	OpenMostRecentOnStartup bool
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
	}

	err := result.Save()
	if err != nil {
		log.Fatalf("filed to save config: %s", err)
	}

	return result
}

func Load() *Config {
	configFile := getConfigPath()

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return generateDefaultConfig()
	}

	var err error
	var data []byte
	if data, err = ioutil.ReadFile(configFile); err != nil {
		return generateDefaultConfig()
	}

	var result *Config
	if err = json.Unmarshal(data, &result); err != nil {
		return generateDefaultConfig()
	}

	return result
}

func (c *Config) Save() error {
	var err error
	var data []byte

	if data, err = json.MarshalIndent(c, "", "   "); err != nil {
		return err
	}

	if err = ioutil.WriteFile(getConfigPath(), data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}

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
			if idx == 5 {
				break
			}

			recent = append(recent, c.RecentProjects[idx])
		}

		c.RecentProjects = recent
	}

	err := c.Save()
	if err != nil {
		log.Fatalf("failed to save config: %s", err)
	}
}

func (c *Config) GetAuxMPQs() []string {
	if len(c.AuxiliaryMpqPath) == 0 {
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
