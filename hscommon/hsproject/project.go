package hsproject

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/dialog"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes/hsfont"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
)

const (
	projectExtension = ".hsp"
)

// Project represents HellSpawner's project
type Project struct {
	ProjectName   string
	Description   string
	Author        string
	AuxiliaryMPQs []string

	filePath       string
	pathEntryCache *hscommon.PathEntry
	mpqs           []d2interface.Archive
}

// CreateNew creates new project
func CreateNew(fileName string) (*Project, error) {
	defaultProjectName := filepath.Base(fileName)

	if strings.EqualFold(filepath.Ext(fileName), projectExtension) {
		fileName += projectExtension
	}

	result := &Project{
		filePath:       fileName,
		ProjectName:    defaultProjectName,
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

// GetProjectFileContentPath returns path to project's content
func (p *Project) GetProjectFileContentPath() string {
	return filepath.Join(filepath.Dir(p.filePath), "content")
}

// GetProjectFilePath returns project's file path
func (p *Project) GetProjectFilePath() string {
	return p.filePath
}

// Save saves project
func (p *Project) Save() error {
	var err error

	var file []byte

	if file, err = json.MarshalIndent(p, "", "   "); err != nil {
		return err
	}

	if err := ioutil.WriteFile(p.filePath, file, os.FileMode(0644)); err != nil {
		return err
	}

	if err := p.ensureProjectPaths(); err != nil {
		return err
	}

	p.InvalidateFileStructure()

	return nil
}

// ValidateAuxiliaryMPQs creates auxiliary mpq's list
func (p *Project) ValidateAuxiliaryMPQs(config *hsconfig.Config) bool {
	for idx := range p.AuxiliaryMPQs {
		realPath := filepath.Join(config.AuxiliaryMpqPath, p.AuxiliaryMPQs[idx])
		if _, err := os.Stat(realPath); os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// LoadFromFile loads projects file
func LoadFromFile(fileName string) (*Project, error) {
	var err error

	var file []byte

	var result *Project

	if file, err = ioutil.ReadFile(fileName); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(file, &result); err != nil {
		return nil, err
	}

	result.filePath = fileName

	if err := result.ensureProjectPaths(); err != nil {
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

// GetFileStructure returns project's file structure
func (p *Project) GetFileStructure() *hscommon.PathEntry {
	if p.pathEntryCache != nil {
		return p.pathEntryCache
	}

	if err := p.ensureProjectPaths(); err != nil {
		log.Fatal(err)
	}

	result := &hscommon.PathEntry{
		Name:        p.ProjectName,
		Children:    make([]*hscommon.PathEntry, 0),
		IsDirectory: true,
		IsRoot:      true,
		Source:      hscommon.PathEntrySourceProject,
	}

	result.FullPath = filepath.Join(filepath.Dir(p.filePath), "content")
	p.getFileNodes(result.FullPath, result)

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
			Source:   hscommon.PathEntrySourceProject,
		}

		if fileNode.Name[0] == '.' || fileNode.FullPath == p.filePath {
			continue
		}

		if files[idx].IsDir() {
			fileNode.IsDirectory = true
			p.getFileNodes(fileNode.FullPath, fileNode)
		}

		entry.Children = append(entry.Children, fileNode)
	}
}

// InvalidateFileStructure cleans project's files structure
func (p *Project) InvalidateFileStructure() {
	p.pathEntryCache = nil
}

// RenameFile renames project's file
func (p *Project) RenameFile(path string) {
	pathEntry := p.FindPathEntry(path)
	if pathEntry == nil {
		return
	}

	pathEntry.OldName = pathEntry.Name

	pathEntry.IsRenaming = true
}

// FindPathEntry search for path entry in project's cahe
func (p *Project) FindPathEntry(path string) *hscommon.PathEntry {
	if p.pathEntryCache == nil {
		return nil
	}

	return p.searchPathEntries(p.pathEntryCache, path)
}

func (p *Project) searchPathEntries(pathEntry *hscommon.PathEntry, path string) *hscommon.PathEntry {
	if pathEntry.FullPath == path {
		return p.pathEntryCache
	}

	for child := range pathEntry.Children {
		if pathEntry.Children[child].FullPath == path {
			return pathEntry.Children[child]
		}

		if found := p.searchPathEntries(pathEntry.Children[child], path); found != nil {
			return found
		}
	}

	return nil
}

// CreateNewFolder creates a new directory
func (p *Project) CreateNewFolder(path *hscommon.PathEntry) {
	basePath := path.FullPath

	filePathFormat := filepath.Join(basePath, "untitled%d")

	var fileName string

	for i := 0; ; i++ {
		possibleFileName := fmt.Sprintf(filePathFormat, i)
		if _, err := os.Stat(possibleFileName); os.IsNotExist(err) {
			fileName = possibleFileName

			break
		}

		if i > 100 {
			dialog.Message("Could not create a new project folder!").Error()

			return
		}
	}

	if err := os.Mkdir(fileName, 0775); err != nil {
		dialog.Message("Could not create a new project folder!").Error()

		return
	}

	p.InvalidateFileStructure()
	p.GetFileStructure()
	p.RenameFile(fileName)
}

// CreateNewFile creates a new file
func (p *Project) CreateNewFile(fileType hsfiletypes.FileType, path *hscommon.PathEntry) {
	basePath := path.FullPath

	filePathFormat := filepath.Join(basePath, "untitled%d"+fileType.FileExtension())

	var fileName string

	for i := 0; ; i++ {
		possibleFileName := fmt.Sprintf(filePathFormat, i)
		if _, err := os.Stat(possibleFileName); os.IsNotExist(err) {
			fileName = possibleFileName

			break
		}

		if i > 100 {
			dialog.Message("Could not create a new project file!").Error()

			return
		}
	}

	if fileType == hsfiletypes.FileTypeFont {
		_, err := hsfont.NewFile(fileName)
		if err != nil {
			log.Fatalf("failed to save font: %s", err)
		}
	}

	p.InvalidateFileStructure()

	// Force regeneration of file structure so that rename can find the file
	p.GetFileStructure()
	p.RenameFile(fileName)
}

// ReloadAuxiliaryMPQs reloads auxiliary MPQs
func (p *Project) ReloadAuxiliaryMPQs(config *hsconfig.Config) {
	p.mpqs = make([]d2interface.Archive, len(p.AuxiliaryMPQs))

	wg := sync.WaitGroup{}
	wg.Add(len(p.AuxiliaryMPQs))

	for mpqIdx := range p.AuxiliaryMPQs {
		go func(idx int) {
			fileName := filepath.Join(config.AuxiliaryMpqPath, p.AuxiliaryMPQs[idx])
			data, err := d2mpq.FromFile(fileName)

			if err != nil {
				log.Fatal(err)
			}

			p.mpqs[idx] = data

			wg.Done()
		}(mpqIdx)
	}

	wg.Wait()
}
