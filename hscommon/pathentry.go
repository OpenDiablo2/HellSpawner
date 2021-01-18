package hscommon

import "fmt"

// PathEntrySource represents the type of path entry.
type PathEntrySource int

const (
	// PathEntrySourceMPQ represents a PathEntry that is relative to a specific MPQ.
	PathEntrySourceMPQ PathEntrySource = iota

	// PathEntrySourceProject represents a PathEntry that is relative to the project.
	PathEntrySourceProject

	// PathEntryVirtual represents a PathEntry that is based on the composite view of
	// the project directory and all MPQs (Project first, then MPQs based on load order).
	PathEntryVirtual
)

// PathEntry defines a file/folder
type PathEntry struct {
	// Children represents child files/folders inside a folder.
	Children []*PathEntry `json:"children"`

	// Name is the visible name of the path entry.
	Name string `json:"name"`

	// FullPath is the actual path of the entry (filesystem, or mpq relative).
	FullPath string `json:"full_path"`

	// IsDirectory is true when this path represents a directory.
	IsDirectory bool `json:"is_directory"`

	// IsRoot is true When this path represents the root path (the project node).
	IsRoot bool `json:"is_root"`

	// IsRenaming is true when this path is currently being renamed in a tree view.
	IsRenaming bool `json:"is_renaming"`

	// OldName is the value of the path's Name before renaming started.
	// If renaming has not started, this value should be blank.
	OldName string `json:"old_name"`

	// PathEntrySource is the type of path entry this is (MPQ or Filesystem).
	Source PathEntrySource `json:"source"`

	// MPQFile represents the full path to the MPQ that contains this file (if this is an MPQ path).
	MPQFile string `json:"mpq_file"`
}

func (p *PathEntry) GetUniqueId() string {
	return fmt.Sprintf("%d_%s_%s", p.Source, p.MPQFile, p.FullPath)
}
