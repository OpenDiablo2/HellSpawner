package hswidget

import (
	"fmt"
	image2 "image"
	"image/color"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/AllenDang/giu"
)

const (
	maxAlpha = uint8(255)
)

type DC6ViewerState struct {
	controls struct {
		direction int32
		frame     int32
		scale     int32
	}
	loadingTexture     bool
	lastFrame          int32
	lastDirection      int32
	framesPerDirection uint32
	texture            *giu.Texture
	rgb                []*image2.RGBA
}

func (is *DC6ViewerState) Dispose() {
	is.texture = nil
}

type DC6ViewerWidget struct {
	id  string
	dc6 *d2dc6.DC6
}

func DC6Viewer(id string, dc6 *d2dc6.DC6) *DC6ViewerWidget {
	result := &DC6ViewerWidget{
		id:  id,
		dc6: dc6,
	}

	return result
}

func (p *DC6ViewerWidget) Build() {
	stateId := fmt.Sprintf("DC6ViewerWidget_%s", p.id)
	state := giu.Context.GetState(stateId)
	var widget *giu.ImageWidget

	if state == nil {
		//Prevent multiple invocation to LoadImage.
		newState := &DC6ViewerState{
			lastFrame:          -1,
			lastDirection:      -1,
			framesPerDirection: p.dc6.FramesPerDirection,
		}

		sw := float32(p.dc6.Frames[0].Width)
		sh := float32(p.dc6.Frames[0].Height)
		widget = giu.Image(nil).Size(sw, sh)

		newState.rgb = make([]*image2.RGBA, p.dc6.Directions*p.dc6.FramesPerDirection)

		for frameIndex := 0; frameIndex < int(p.dc6.Directions*p.dc6.FramesPerDirection); frameIndex++ {
			newState.rgb[frameIndex] = image2.NewRGBA(image2.Rect(0, 0, int(p.dc6.Frames[frameIndex].Width), int(p.dc6.Frames[frameIndex].Height)))
			decodedFrame := p.dc6.DecodeFrame(frameIndex)

			for y := 0; y < int(p.dc6.Frames[frameIndex].Height); y++ {
				for x := 0; x < int(p.dc6.Frames[frameIndex].Width); x++ {
					idx := x + (y * int(p.dc6.Frames[frameIndex].Width))
					val := decodedFrame[idx]

					alpha := maxAlpha

					if val == 0 {
						alpha = 0
					}

					newState.rgb[frameIndex].Set(x, y, color.RGBA{R: val, G: val, B: val, A: alpha})
				}
			}
		}

		giu.Context.SetState(stateId, newState)

		widget.Build()
	} else {
		viewerState := state.(*DC6ViewerState)

		if !viewerState.loadingTexture && (viewerState.lastDirection != viewerState.controls.direction || viewerState.lastFrame != viewerState.controls.frame) {
			// Control values have changed, need to regenerate the texture
			viewerState.lastDirection = viewerState.controls.direction
			viewerState.lastFrame = viewerState.controls.frame
			viewerState.loadingTexture = true
			viewerState.texture = nil

			giu.Context.SetState(stateId, viewerState)

			hscommon.CreateTextureFromARGB(viewerState.rgb[viewerState.lastFrame+(viewerState.lastDirection*int32(viewerState.framesPerDirection))], func(tex *g.Texture) {
				newState := giu.Context.GetState(stateId).(*DC6ViewerState)

				newState.texture = tex
				newState.loadingTexture = false
				giu.Context.SetState(stateId, newState)
			})
		}

		imageScale := uint32(viewerState.controls.scale)
		curFrameIndex := int(viewerState.controls.frame) + (int(viewerState.controls.direction) * int(p.dc6.FramesPerDirection))

		if imageScale < 1 {
			imageScale = 1
		}

		_ = giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
		var widget *giu.ImageWidget
		w := float32(p.dc6.Frames[curFrameIndex].Width * imageScale)
		h := float32(p.dc6.Frames[curFrameIndex].Height * imageScale)
		if viewerState.texture == nil {
			widget = giu.Image(nil).Size(w, h)
		} else {

			widget = giu.Image(viewerState.texture).Size(w, h)
		}

		giu.Layout{
			giu.Label(fmt.Sprintf(
				"Version: %v\t Flags: %b\t Encoding: %v\t",
				p.dc6.Version,
				int64(p.dc6.Flags),
				p.dc6.Encoding,
			)),
			giu.Label(fmt.Sprintf("Directions: %v\tFrames per Direction: %v", p.dc6.Directions, p.dc6.FramesPerDirection)),
			giu.Custom(func() {
				imgui.BeginGroup()
				if p.dc6.Directions > 1 {
					imgui.SliderInt("Direction", &viewerState.controls.direction, 0, int32(p.dc6.Directions-1))
				}

				if p.dc6.FramesPerDirection > 1 {
					imgui.SliderInt("Frames", &viewerState.controls.frame, 0, int32(p.dc6.FramesPerDirection-1))
				}

				imgui.SliderInt("Scale", &viewerState.controls.scale, 1, 8)

				imgui.EndGroup()
			}),
			giu.Separator(),
			widget,
		}.Build()
	}

}
