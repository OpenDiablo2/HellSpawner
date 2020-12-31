package hswidget

import (
	"fmt"
	image2 "image"
	"image/color"
	"log"

	"github.com/AllenDang/giu/imgui"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/AllenDang/giu"
)

type DC6ViewerState struct {
	controls struct {
		direction int32
		frame     int32
		scale     float32
	}
	textures []*giu.Texture
}

func (is *DC6ViewerState) Dispose() {
	is.textures = nil
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
		widget = giu.Image(nil, 32, 32)

		//Prevent multiple invocation to LoadImage.
		giu.Context.SetState(stateId, &DC6ViewerState{})

		sw := float32(p.dc6.Frames[0].Width)
		sh := float32(p.dc6.Frames[0].Height)
		widget = giu.Image(nil, sw, sh)

		rgb := make([]*image2.RGBA, p.dc6.Directions*p.dc6.FramesPerDirection)

		for frameIndex := 0; frameIndex < int(p.dc6.Directions*p.dc6.FramesPerDirection); frameIndex++ {
			rgb[frameIndex] = image2.NewRGBA(image2.Rect(0, 0, int(p.dc6.Frames[frameIndex].Width), int(p.dc6.Frames[frameIndex].Height)))
			decodedFrame := p.dc6.DecodeFrame(frameIndex)

			for y := 0; y < int(p.dc6.Frames[frameIndex].Height); y++ {
				for x := 0; x < int(p.dc6.Frames[frameIndex].Width); x++ {
					idx := x + (y * int(p.dc6.Frames[frameIndex].Width))
					val := decodedFrame[idx]
					rgb[frameIndex].Set(x, y, color.RGBA{R: val, G: val, B: val, A: 255})
				}
			}
		}

		go func() {
			textures := make([]*giu.Texture, p.dc6.Directions*p.dc6.FramesPerDirection)
			for frameIndex := 0; frameIndex < int(p.dc6.Directions*p.dc6.FramesPerDirection); frameIndex++ {
				var err error
				textures[frameIndex], err = giu.NewTextureFromRgba(rgb[frameIndex])
				if err != nil {
					log.Fatal(err)
				}
			}
			giu.Context.SetState(stateId, &DC6ViewerState{textures: textures})
		}()

		widget.Build()
	} else {
		viewerState := state.(*DC6ViewerState)

		imageScale := uint32(viewerState.controls.scale)
		curFrameIndex := int(viewerState.controls.frame) + (int(viewerState.controls.direction) * int(p.dc6.FramesPerDirection))

		if imageScale < 1 {
			imageScale = 1
		}

		giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
		var widget *giu.ImageWidget
		if viewerState.textures == nil || len(viewerState.textures) <= curFrameIndex || viewerState.textures[curFrameIndex] == nil {
			widget = giu.Image(nil, 32, 32)
		} else {
			widget = giu.Image(viewerState.textures[curFrameIndex],
				float32(p.dc6.Frames[curFrameIndex].Width*imageScale), float32(p.dc6.Frames[curFrameIndex].Height*imageScale))
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

				imgui.SliderFloat("Scale", &viewerState.controls.scale, 1, 8)
				imgui.EndGroup()
			}),
			giu.Separator(),
			widget,
		}.Build()
	}

}
