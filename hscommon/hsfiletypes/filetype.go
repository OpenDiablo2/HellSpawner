package hsfiletypes

import (
	"errors"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"
)

// FileType represents file type
type FileType int

type fileTypeInfoStruct struct {
	FileType
	Name         string
	Extension    string
	subTypeCheck func(*[]byte) (FileType, error)
}

// enumerate known file types
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
	FileTypeTBL
	FileTypeTBLStringTable
	FileTypeTBLFontTable
	FileTypeDS1
	FileTypeAnimationData
)

// determinateTBLtype returns table type
func determineTBLtype(data *[]byte) (FileType, error) {
	_, err := d2tbl.LoadTextDictionary(*data)
	if err == nil {
		return FileTypeTBLStringTable, nil
	}

	if string((*data)[:4]) == "Woo!" {
		return FileTypeTBLFontTable, nil
	}

	return FileTypeText, nil
}

func fileExtensionInfo() []fileTypeInfoStruct {
	return []fileTypeInfoStruct{
		{FileType: FileTypeUnknown},
		{FileType: FileTypeFont, Name: "Font", Extension: ".hsf"},
		{FileType: FileTypePalette, Name: "Palette", Extension: ".dat"},
		{FileType: FileTypePL2, Name: "Palette Map", Extension: ".pl2"},
		{FileType: FileTypeAudio, Name: "Audio", Extension: ".wav"},
		{FileType: FileTypeDCC, Name: "DCC", Extension: ".dcc"},
		{FileType: FileTypeDC6, Name: "DC6", Extension: ".dc6"},
		{FileType: FileTypeCOF, Name: "COF", Extension: ".cof"},
		{FileType: FileTypeDT1, Name: "DT1", Extension: ".dt1"},
		{FileType: FileTypeTBL, Name: "TBL", Extension: ".tbl", subTypeCheck: determineTBLtype},
		{FileType: FileTypeText, Name: "Text", Extension: ".txt"},
		{FileType: FileTypeDS1, Name: "DS1", Extension: ".ds1"},
		{FileType: FileTypeAnimationData, Name: "AnimationData", Extension: ".d2"},
	}
}

// String returns file type string
func (f FileType) String() string {
	return fileExtensionInfo()[f].Name
}

// FileExtension returns file's extension
func (f FileType) FileExtension() string {
	return fileExtensionInfo()[f].Extension
}

// GetFileTypeFromExtension returns file type
func GetFileTypeFromExtension(extension string, data *[]byte) (FileType, error) {
	for _, info := range fileExtensionInfo() {
		if strings.EqualFold(info.Extension, extension) {
			if info.subTypeCheck == nil {
				return info.FileType, nil
			}

			return info.subTypeCheck(data)
		}
	}

	return FileTypeUnknown, errors.New("filetype: no file type matches the extension provided")
}
