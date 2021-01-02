package hscommon

type PathEntry struct {
	Children    []*PathEntry
	Name        string
	FullPath    string
	IsDirectory bool
	IsRoot      bool
}
