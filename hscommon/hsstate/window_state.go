package hsstate

// WindowState holds information about windows.
type WindowState struct {
	Visible bool    `json:"visible"`
	PosX    float32 `json:"x"`
	PosY    float32 `json:"y"`
}
