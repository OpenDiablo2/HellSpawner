package hsassets

import (
	"bytes"
	_ "embed"
	"io"
)

//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/reload.png
var ReloadIcon []byte

//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_delete.png
var DeleteIcon []byte

//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png
var DownArrowIcon []byte

//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png
var UpArrowIcon []byte

//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_left.png
var LeftArrowIcon []byte

//go:embed 3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_right.png
var RightArrowIcon []byte

func MakeReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}
