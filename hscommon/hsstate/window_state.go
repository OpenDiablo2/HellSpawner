package hsstate

// WindowState holds information about windows.
type WindowState struct {
	Visible bool    `json:"visible"`
	PosX    float32 `json:"x"`
	PosY    float32 `json:"y"`
	Width   float32 `json:"w"`
	Height  float32 `json:"h"`
}
