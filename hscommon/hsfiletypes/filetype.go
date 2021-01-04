package hsfiletypes

type FileType int

const (
	FileTypeFont FileType = iota
)

func (f FileType) String() string {
	return [...]string{
		"Font",
	}[f]
}

func (f FileType) FileExtension() string {
	return [...]string{
		".hsf",
	}[f]
}
