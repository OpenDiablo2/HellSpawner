package hswidget

import (
	"fmt"
	image2 "image"
	"image/color"
	"log"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
)

// DCCViewerState represents dcc viewers state
type DCCViewerState struct {
	controls struct {
		direction int32
		frame     int32
		scale     int32
	}

	textures []*giu.Texture
}

// Dispose cleans viewers state
func (is *DCCViewerState) Dispose() {
	is.textures = nil
}

// DCCViewerWidget creates a new dcc widget
type DCCViewerWidget struct {
	id  string
	dcc *d2dcc.DCC
}

// DCCViewer creates a new dcc viewers widget
func DCCViewer(id string, dcc *d2dcc.DCC) *DCCViewerWidget {
	result := &DCCViewerWidget{
		id:  id,
		dcc: dcc,
	}

	return result
}

// Build build a widget
func (p *DCCViewerWidget) Build() {
	stateID := fmt.Sprintf("DCCViewerWidget_%s", p.id)
	state := giu.Context.GetState(stateID)

	if state == nil {
		// Prevent multiple invocation to LoadImage.
		giu.Context.SetState(stateID, &DCCViewerState{})

		totalFrames := p.dcc.NumberOfDirections * p.dcc.FramesPerDirection
		images := make([]*image2.RGBA, totalFrames)

		for dirIdx := range p.dcc.Directions {
			fw := p.dcc.Directions[dirIdx].Box.Width
			fh := p.dcc.Directions[dirIdx].Box.Height

			for frameIdx := range p.dcc.Directions[dirIdx].Frames {
				absoluteFrameIdx := (dirIdx * p.dcc.FramesPerDirection) + frameIdx

				frame := p.dcc.Directions[dirIdx].Frames[frameIdx]
				pixels := frame.PixelData

				images[absoluteFrameIdx] = image2.NewRGBA(image2.Rect(0, 0, fw, fh))

				for y := 0; y < fh; y++ {
					for x := 0; x < fw; x++ {
						idx := x + (y * fw)
						if idx >= len(pixels) {
							continue
						}

						val := pixels[idx]

						alpha := maxAlpha

						if val == 0 {
							alpha = 0
						}

						RGBAcolor := color.RGBA{R: val, G: val, B: val, A: alpha}

						images[absoluteFrameIdx].Set(x, y, RGBAcolor)
					}
				}
			}
		}

		go func() {
			textures := make([]*giu.Texture, totalFrames)

			for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
				var err error

				textures[frameIndex], err = giu.NewTextureFromRgba(images[frameIndex])
				if err != nil {
					log.Fatal(err)
				}
			}
			giu.Context.SetState(stateID, &DCCViewerState{textures: textures})
		}()

		// display a temporary dummy image until the real one ready
		firstFrame := p.dcc.Directions[0].Frames[0]
		sw := float32(firstFrame.Width)
		sh := float32(firstFrame.Height)
		widget := giu.Image(nil).Size(sw, sh)
		widget.Build()
	} else {
		viewerState := state.(*DCCViewerState)

		imageScale := uint32(viewerState.controls.scale)
		dirIdx := int(viewerState.controls.direction)
		frameIdx := viewerState.controls.frame

		textureIdx := dirIdx*len(p.dcc.Directions[dirIdx].Frames) + int(frameIdx)

		if imageScale < 1 {
			imageScale = 1
		}

		err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
		if err != nil {
			log.Print(err)
		}

		var widget *giu.ImageWidget
		if viewerState.textures == nil || len(viewerState.textures) <= int(frameIdx) || viewerState.textures[frameIdx] == nil {
			widget = giu.Image(nil).Size(32, 32)
		} else {
			bw := p.dcc.Directions[dirIdx].Box.Width
			bh := p.dcc.Directions[dirIdx].Box.Height
			w := float32(uint32(bw) * imageScale)
			h := float32(uint32(bh) * imageScale)
			widget = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
		}

		giu.Layout{
			giu.Line(
				giu.Label(fmt.Sprintf("Signature: %v", p.dcc.Signature)),
				giu.Label(fmt.Sprintf("Version: %v", p.dcc.Version)),
			),
			giu.Line(
				giu.Label(fmt.Sprintf("Directions: %v", p.dcc.NumberOfDirections)),
				giu.Label(fmt.Sprintf("Frames per Direction: %v", p.dcc.FramesPerDirection)),
			),
			giu.Custom(func() {
				imgui.BeginGroup()
				if p.dcc.NumberOfDirections > 1 {
					imgui.SliderInt("Direction", &viewerState.controls.direction, 0, int32(p.dcc.NumberOfDirections-1))
				}

				if p.dcc.FramesPerDirection > 1 {
					imgui.SliderInt("Frames", &viewerState.controls.frame, 0, int32(p.dcc.FramesPerDirection-1))
				}

				imgui.SliderInt("Scale", &viewerState.controls.scale, 1, 8)

				imgui.EndGroup()
			}),
			giu.Separator(),
			widget,
		}.Build()
	}
}
