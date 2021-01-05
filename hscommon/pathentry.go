package hscommon

import "fmt"

// PathEntrySource represents the type of path entry.
type PathEntrySource int

const (
	// PathEntrySourceMPQ represents a PathEntry that is relative to a specific MPQ.
	PathEntrySourceMPQ PathEntrySource = iota

	// PathEntrySourceProject represents a PathEntry that is relative to the project.
	PathEntrySourceProject

	// PathEntryCoalesced represents a PathEntry that is based on the coalesced view of
	// the project and all MPQs (Project first, then MPQs based on load order).
	PathEntryCoalesced
)

// PathEntry defines a file/folder
type PathEntry struct {
	// Children represents child files/folders inside a folder.
	Children []*PathEntry

	// Name is the visible name of the path entry.
	Name string

	// FullPath is the actual path of the entry (filesystem, or mpq relative).
	FullPath string

	// IsDirectory is true when this path represents a directory.
	IsDirectory bool

	// IsRoot is true When this path represents the root path (the project node).
	IsRoot bool

	// IsRenaming is true when this path is currently being renamed in a tree view.
	IsRenaming bool

	// OldName is the value of the path's Name before renaming started.
	// If renaming has not started, this value should be blank.
	OldName string

	// PathEntrySource is the type of path entry this is (MPQ or Filesystem).
	Source PathEntrySource

	// MPQFile represents the full path to the MPQ that contains this file (if this is an MPQ path).
	MPQFile string
}

func (p *PathEntry) GetUniqueId() string {
	return fmt.Sprintf("%d_%s_%s", p.Source, p.MPQFile, p.FullPath)
}
