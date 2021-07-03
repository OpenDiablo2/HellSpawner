package hsfiletypes

import (
	"errors"
	"strings"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"
)

// FileType represents file type
type FileType int

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
	numFileTypes
)

// determinateTBLtype returns table type
func determineSubtypeTBL(data *[]byte) FileType {
	if _, err := d2tbl.LoadTextDictionary(*data); err == nil {
		return FileTypeTBLStringTable
	}

	d := *data
	if string(d[:4]) == "Woo!" {
		return FileTypeTBLFontTable
	}

	return FileTypeText
}

// String returns file type string
func (f FileType) String() string {
	table := map[FileType]string{
		FileTypeUnknown:        "unknown",
		FileTypeFont:           "Hellspawner font",
		FileTypePalette:        "palette",
		FileTypePL2:            "palette transform",
		FileTypeAudio:          "wav",
		FileTypeDCC:            "DCC image",
		FileTypeDC6:            "DC6 image",
		FileTypeCOF:            "COF animation data",
		FileTypeDT1:            "DT1 tileset",
		FileTypeTBLFontTable:   "Font Character Table",
		FileTypeTBLStringTable: "String Table",
		FileTypeText:           "text file",
		FileTypeDS1:            "DS1 Map Stamp",
		FileTypeAnimationData:  "Animation Dataset",
	}

	val, found := table[f]
	if !found {
		return table[FileTypeUnknown]
	}

	return val
}

// FileExtension returns file's extension
func (f FileType) FileExtension() string {
	table := map[FileType]string{
		FileTypeFont:           ".hsf",
		FileTypePalette:        ".dat",
		FileTypePL2:            ".pl2",
		FileTypeAudio:          ".wav",
		FileTypeDCC:            ".dcc",
		FileTypeDC6:            ".dc6",
		FileTypeCOF:            ".cof",
		FileTypeDT1:            ".dt1",
		FileTypeTBLFontTable:   ".tbl",
		FileTypeTBLStringTable: ".tbl",
		FileTypeText:           ".txt",
		FileTypeDS1:            ".ds1",
		FileTypeAnimationData:  ".d2",
	}

	return table[f]
}

type fileTypeCheckFn = func(data *[]byte) FileType

// subtypeCheckFn returns a function to check a file. This is important for
// distinguishing between files that share a common file extension.
func (f FileType) subtypeCheckFn() fileTypeCheckFn {
	table := map[FileType]fileTypeCheckFn{
		FileTypeTBL:            determineSubtypeTBL,
		FileTypeTBLFontTable:   determineSubtypeTBL,
		FileTypeTBLStringTable: determineSubtypeTBL,
	}

	return table[f]
}

// GetFileTypeFromExtension returns file type
func GetFileTypeFromExtension(extension string, data *[]byte) (FileType, error) {
	for fileType := FileType(0); fileType < numFileTypes; fileType++ {
		extensionsMatch := strings.EqualFold(fileType.FileExtension(), extension)
		if !extensionsMatch {
			continue
		}

		fnDetermine := fileType.subtypeCheckFn()
		if fnDetermine == nil {
			return fileType, nil
		}

		return fnDetermine(data), nil
	}

	return FileTypeUnknown, errors.New("filetype: no file type matches the extension provided")
}
