package hsfiletypes

import (
	"errors"
	"strings"
)

type FileType int

type fileTypeInfoStruct struct {
	Name      string
	Extension string
}

const (
	FileTypeUnknown FileType = iota
	FileTypeText
	FileTypeFont
	FileTypePalette
	FileTypeAudio
	FileTypeDCC
	FileTypeDC6
	FileTypeCOF
	FileTypeDT1
	FileTypePL2
)

func fileExtensionInfo() map[FileType]fileTypeInfoStruct {
	return map[FileType]fileTypeInfoStruct{
		FileTypeUnknown: {},
		FileTypeText:    {Name: "Text", Extension: ".txt"},
		FileTypeFont:    {Name: "Font", Extension: ".hsf"},
		FileTypePalette: {Name: "Palette", Extension: ".dat"},
		FileTypePL2:     {Name: "Palette Map", Extension: ".pl2"},
		FileTypeAudio:   {Name: "Audio", Extension: ".wav"},
		FileTypeDCC:     {Name: "DCC", Extension: ".dcc"},
		FileTypeDC6:     {Name: "DC6", Extension: ".dc6"},
		FileTypeCOF:     {Name: "COF", Extension: ".cof"},
		FileTypeDT1:     {Name: "DT1", Extension: ".dt1"},
	}
}

func (f FileType) String() string {
	return fileExtensionInfo()[f].Name
}

func (f FileType) FileExtension() string {
	return fileExtensionInfo()[f].Extension
}

func GetFileTypeFromExtension(extension string) (FileType, error) {
	info := fileExtensionInfo()
	for idx := range info {
		if strings.EqualFold(info[idx].Extension, extension) {
			return idx, nil
		}
	}

	return FileTypeUnknown, errors.New("filetype: no file type matches the extension provided")
}
