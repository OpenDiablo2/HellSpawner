package hsfiletypes

import (
	"errors"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"
)

type FileType int

type fileTypeInfoStruct struct {
	Name         string
	Extension    string
	subTypeCheck func(*[]byte) (FileType, error)
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
	FileTypeTBL
	FileTypeTBLStringTable
)

func determineTBLtype(data *[]byte) (FileType, error) {
	_, err := d2tbl.LoadTextDictionary(*data)
	if err == nil {
		return FileTypeTBLStringTable, err
	}

	return FileTypeUnknown, errors.New("unknown file type")
}

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
		FileTypeTBL:     {Name: "TBL", Extension: ".tbl", subTypeCheck: determineTBLtype},
	}
}

func (f FileType) String() string {
	return fileExtensionInfo()[f].Name
}

func (f FileType) FileExtension() string {
	return fileExtensionInfo()[f].Extension
}

func GetFileTypeFromExtension(extension string, data *[]byte) (FileType, error) {
	info := fileExtensionInfo()
	for idx := range info {
		if strings.EqualFold(info[idx].Extension, extension) {
			if info[idx].subTypeCheck == nil {
				return idx, nil
			}

			return info[idx].subTypeCheck(data)
		}
	}

	return FileTypeUnknown, errors.New("filetype: no file type matches the extension provided")
}
