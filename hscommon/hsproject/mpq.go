package hsproject

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

func (p *Project) GetMPQFileNodes(mpq d2interface.Archive, config *hsconfig.Config) *hscommon.PathEntry {
	result := &hscommon.PathEntry{
		Name:        filepath.Base(mpq.Path()),
		IsDirectory: true,
		Source:      hscommon.PathEntrySourceMPQ,
		MPQFile:     mpq.Path(),
	}

	files, err := mpq.Listfile()

	if err != nil {
		files, err = p.searchForMpqFiles(mpq, config)
		if err != nil {
			return result
		}
	}

	pathNodes := make(map[string]*hscommon.PathEntry)
	pathNodes[""] = result

	for idx := range files {
		elements := strings.FieldsFunc(files[idx], func(r rune) bool { return r == '\\' || r == '/' })

		path := ""

		for elemIdx := range elements {
			oldPath := path

			path += elements[elemIdx]
			if elemIdx < len(elements)-1 {
				path += `\`
			}

			if pathNodes[strings.ToLower(path)] == nil {
				pathNodes[strings.ToLower(path)] = &hscommon.PathEntry{
					Name:        elements[elemIdx],
					FullPath:    path,
					Source:      hscommon.PathEntrySourceMPQ,
					MPQFile:     mpq.Path(),
					IsDirectory: elemIdx < len(elements)-1,
				}

				pathNodes[strings.ToLower(oldPath)].Children =
					append(pathNodes[strings.ToLower(oldPath)].Children, pathNodes[strings.ToLower(path)])
			}
		}
	}

	hscommon.SortPaths(result)

	return result
}

// Search for files in MPQ's without listfiles using a list of known filenames
func (p *Project) searchForMpqFiles(mpq d2interface.Archive, config *hsconfig.Config) ([]string, error) {
	var files []string

	if config.ExternalListFile != "" {
		file, err := os.Open(config.ExternalListFile)
		if err != nil {
			return files, errors.New("Couldn't open listfile")
		}

		defer func() {
			err := file.Close()
			if err != nil {
				log.Print(err)
			}
		}()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fileName := scanner.Text()
			if mpq.Contains(fileName) {
				files = append(files, fileName)
			}
		}

		return files, scanner.Err()
	}

	return files, nil
}
