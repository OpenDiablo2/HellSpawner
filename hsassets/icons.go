package hsassets

import (
	"bytes"
	_ "embed" // this is standard solution for embed
	"io"
)

// these variables are links to existing icones used in project
// nolint:gochecknoglobals // go:embed directive works only for globals
// https://github.com/golangci/golangci-lint/issues/1727
var (
	//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/reload.png
	ReloadIcon []byte

	//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_delete.png
	DeleteIcon []byte

	//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png
	DownArrowIcon []byte

	//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png
	UpArrowIcon []byte

	//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_left.png
	LeftArrowIcon []byte

	//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_right.png
	RightArrowIcon []byte
)

// MakeReader creates reader from variable
func MakeReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}
