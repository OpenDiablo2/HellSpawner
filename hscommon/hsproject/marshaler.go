package hsproject

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2animdata"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dat"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2pl2"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2tbl"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsfiletypes"
)

type marshaler interface {
	Marshal() []byte
}

func getMarshallerByType(fileType hsfiletypes.FileType) marshaler {
	switch fileType {
	case hsfiletypes.FileTypeTBLFontTable:
		return &d2font.Font{}
	case hsfiletypes.FileTypeTBLStringTable:
		return &d2tbl.TextDictionary{}
	case hsfiletypes.FileTypeAnimationData:
		return &d2animdata.AnimationData{}
	case hsfiletypes.FileTypeCOF:
		return d2cof.New()
	case hsfiletypes.FileTypePalette:
		return d2dat.New()
	case hsfiletypes.FileTypePL2:
		return &d2pl2.PL2{}
	case hsfiletypes.FileTypeDS1:
		return &d2ds1.DS1{}
	case hsfiletypes.FileTypeDT1:
		return d2dt1.New()
	}

	return nil
}
