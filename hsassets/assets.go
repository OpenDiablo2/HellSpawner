package hsassets

import (
	_ "embed" // this is standard solution for embed
)

// these variables are links to existing icons used in project
// nolint:gochecknoglobals // go:embed directive works only for globals
// https://github.com/golangci/golangci-lint/issues/1727
var (
	//go:embed icons/reload.png
	ReloadIcon []byte

	//go:embed icons/stock_delete.png
	DeleteIcon []byte

	//go:embed icons/stock_down.png
	DownArrowIcon []byte

	//go:embed icons/stock_up.png
	UpArrowIcon []byte

	//go:embed icons/stock_left.png
	LeftArrowIcon []byte

	//go:embed icons/stock_right.png
	RightArrowIcon []byte

	//go:embed icons/player_play.png
	PlayButtonIcon []byte

	//go:embed icons/player_pause.png
	PauseButtonIcon []byte
)

// these variables are links to existing fonts used in project
// nolint:gochecknoglobals // like above: this is ok for go:embed
var (
	//go:embed fonts/NotoSans-Regular.ttf
	FontNotoSansRegular []byte
	//go:embed fonts/CascadiaCode.ttf
	FontCascadiaCode []byte
	//go:embed fonts/DiabloRegular.ttf
	FontDiabloRegular []byte
	//go:embed fonts/DiabloBold.ttf
	FontDiabloBold []byte
	//go:embed fonts/source-han-serif/SourceHanSerifTC-Regular.otf
	FontSourceHanSerif []byte
)

// HellSpawnerLogo is a logo image from about dialog
//go:embed images/d2logo.png
// nolint:gochecknoglobals // go-embed
var HellSpawnerLogo []byte

// ImageShrug is an image, which is displayed in ds1 editor, when no objects are there
//go:embed images/shrug.png
// nolint:gochecknoglobals // the same reson as above
var ImageShrug []byte
