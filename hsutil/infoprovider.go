package hsutil

import (
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsinput"
	"golang.org/x/image/font"
)

type InfoProvider interface {
	GetAppConfig() *hsconfig.AppConfig
	GetNormalFont() font.Face
	GetSymbolsFont() font.Face
	GetMonospaceFont() font.Face
	GetInputVector() *hsinput.InputVector
}
