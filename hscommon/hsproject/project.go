package hsproject

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
)

type Project struct {
	filePath       string
	ProjectName    string
	pathEntryCache *hscommon.PathEntry

	Description   string
	Author        string
	AuxiliaryMPQs []string
}

func CreateNew(fileName string) (*Project, error) {
	if strings.ToLower(filepath.Ext(fileName)) != ".hsp" {
		fileName += ".hsp"
	}

	result := &Project{
		filePath:       fileName,
		ProjectName:    "Untitled Project",
		pathEntryCache: nil,
	}

	if err := result.Save(); err != nil {
		return nil, err
	}

	if err := result.ensureProjectPaths(); err != nil {
		return nil, err
	}

	return result, nil
}

func (p *Project) Save() error {
	var err error
	var file []byte

	if file, err = json.MarshalIndent(p, "", "   "); err != nil {
		return err
	}
	if err = ioutil.WriteFile(p.filePath, file, os.FileMode(0644)); err != nil {
		return err
	}
	if err = p.ensureProjectPaths(); err != nil {
		return err
	}

	p.InvalidateFileStructure()

	return nil
}

func (p *Project) ValidateAuxiliaryMPQs(config *hsconfig.Config) bool {
	for idx := range p.AuxiliaryMPQs {
		realPath := filepath.Join(config.AuxiliaryMpqPath, p.AuxiliaryMPQs[idx])
		if _, err := os.Stat(realPath); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func LoadFromFile(fileName string) (*Project, error) {
	var err error
	var file []byte
	var result *Project

	if file, err = ioutil.ReadFile(fileName); err != nil {
		return nil, err
	}

	if err = json.Unmarshal(file, &result); err != nil {
		return nil, err
	}

	result.filePath = fileName

	if err = result.ensureProjectPaths(); err != nil {
		return nil, err
	}

	result.InvalidateFileStructure()

	return result, nil
}

func (p *Project) ensureProjectPaths() error {
	basePath := filepath.Dir(p.filePath)
	contentPath := filepath.Join(basePath, "content")

	if _, err := os.Stat(contentPath); os.IsNotExist(err) {
		if err := os.Mkdir(contentPath, os.FileMode(0755)); err != nil {
			return err
		}
	}

	return nil
}

func (p *Project) GetFileStructure() *hscommon.PathEntry {
	if p.pathEntryCache != nil {
		return p.pathEntryCache
	}

	if err := p.ensureProjectPaths(); err != nil {
		log.Fatal(err)
	}

	result := &hscommon.PathEntry{
		Name:     p.ProjectName,
		FullPath: p.filePath,
		Children: make([]*hscommon.PathEntry, 0),
	}

	p.getFileNodes(filepath.Join(filepath.Dir(p.filePath), "content"), result)

	p.pathEntryCache = result

	return result
}

func (p *Project) getFileNodes(path string, entry *hscommon.PathEntry) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for idx := range files {
		fileNode := &hscommon.PathEntry{
			Children: []*hscommon.PathEntry{},
			Name:     files[idx].Name(),
			FullPath: filepath.Join(path, files[idx].Name()),
		}

		if fileNode.Name[0] == '.' || fileNode.FullPath == p.filePath {
			continue
		}

		if files[idx].IsDir() {
			p.getFileNodes(fileNode.FullPath, fileNode)
		}

		entry.Children = append(entry.Children, fileNode)
	}
}

func (p *Project) InvalidateFileStructure() {
	p.pathEntryCache = nil
}
