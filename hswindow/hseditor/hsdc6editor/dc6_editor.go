package hsdc6editor

import (
	"fmt"
	"image"

	g "github.com/OpenDiablo2/giu"
	"github.com/OpenDiablo2/giu/imgui"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
)

func Create(path string, data []byte) (*DC6Editor, error) {
	dc6, err := d2dc6.Load(data)
	if err != nil {
		return nil, err
	}

	numFrames := dc6.Directions * dc6.FramesPerDirection

	result := &DC6Editor{
		path:          path,
		dc6:           dc6,
		decodedFrames: make([][]byte, numFrames),
		textures:      make([]*g.Texture, numFrames),
	}

	return result, nil
}

type DC6Editor struct {
	hseditor.Editor
	path     string
	dc6      *d2dc6.DC6
	controls struct {
		direction int32
		frame     int32
		scale     float32
	}
	decodedFrames [][]byte
	textures      []*g.Texture
}

func (e *DC6Editor) GetWindowTitle() string {
	return e.path + "##" + e.GetId()
}

func (e *DC6Editor) Cleanup() {
	e.Visible = false
}

func (e *DC6Editor) Render() {
	if !e.Visible {
		return
	}

	if e.ToFront {
		e.ToFront = false
		imgui.SetNextWindowFocus()
	}

	stateId := fmt.Sprintf("DC6Editor_%s", e.path)
	state := g.Context.GetState(stateId)

	currentFrameIndex := (int(e.controls.direction) * int(e.controls.frame)) % len(e.dc6.Frames)

	if e.decodedFrames[currentFrameIndex] == nil {
		e.decodedFrames[currentFrameIndex] = e.dc6.DecodeFrame(currentFrameIndex)
	}

	paletteIndices := e.decodedFrames[currentFrameIndex]
	currentFrame := e.dc6.Frames[currentFrameIndex]

	w, h := int(currentFrame.Width), int(currentFrame.Height)

	if state == nil {
		go func() {
			rect := image.Rect(0, 0, w, h)
			rgba := image.NewRGBA(rect)

			// fake because this should be coming from a palette
			// but we are creating a grayscale palette from the indices
			fakeRGBA := make([]byte, len(paletteIndices)*4)

			for idx, val := range paletteIndices {
				fakeRGBA[idx+0] = val
				fakeRGBA[idx+1] = val
				fakeRGBA[idx+2] = val
				fakeRGBA[idx+3] = 255
			}

			rgba.Pix = fakeRGBA

			tex, err := g.NewTextureFromRgba(rgba)
			if err == nil {
				g.Context.SetState(stateId, &g.ImageState{Texture: tex})
			}

			e.textures[currentFrameIndex] = tex
		}()

		return
	}

	imgState := state.(*g.ImageState)
	g.Image(imgState.Texture, float32(w), float32(h))

	g.WindowV(
		e.GetWindowTitle(),
		&e.Visible,
		g.WindowFlagsAlwaysAutoResize,
		0, 0,
		256, 256,
		g.Layout{
			g.Label(fmt.Sprintf("Version: %v", e.dc6.Version)),
			g.Label(fmt.Sprintf("Flags: %b", int64(e.dc6.Flags))),
			g.Label(fmt.Sprintf("Encoding: %v", e.dc6.Encoding)),
			g.Label(fmt.Sprintf("#Directions: %v", e.dc6.Directions)),
			g.Label(fmt.Sprintf("Frames/Direction: %v", e.dc6.FramesPerDirection)),
			g.Custom(func() {
				imgui.Text(fmt.Sprintf("%v", e.path))
				imgui.BeginGroup()
				imgui.SliderInt("Direction", &e.controls.direction, 0, int32(e.dc6.Directions-1))
				imgui.SliderInt("Frames", &e.controls.frame, 0, int32(e.dc6.FramesPerDirection-1))
				imgui.SliderFloat("Scale", &e.controls.scale, 1, 5)
				imgui.EndGroup()

			}),
		},
	)
}
