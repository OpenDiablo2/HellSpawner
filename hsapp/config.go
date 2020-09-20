package hsapp

type AppConfig struct {
	Colors AppColorConfig `json:"colors"`
}

type AppColorConfig struct {
	WindowBackground      []uint8 `json:"windowBackground"`
	WindowFrame           []uint8 `json:"windowFrame"`
	WindowButtonHighlight []uint8 `json:"windowButtonHighlight"`
	WindowFrameText       []uint8 `json:"windowFrameText"`
}
