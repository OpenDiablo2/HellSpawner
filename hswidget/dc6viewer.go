package hswidget

import (
	"fmt"
	image2 "image"
	"image/color"
	"log"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

const (
	// nolint:gomnd // it was constant
	maxAlpha = uint8(255)
)

// DC6ViewerState represents dc6 viewer's state
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
	texture            *g.Texture
	rgb                []*image2.RGBA
}

// Dispose cleans state content
func (is *DC6ViewerState) Dispose() {
	is.texture = nil
}

// DC6ViewerWidget represents dc6viewer's widget
type DC6ViewerWidget struct {
	id            string
	dc6           *d2dc6.DC6
	textureLoader *hscommon.TextureLoader
}

// DC6Viewer creates new DC6ViewerWidget
func DC6Viewer(textureLoader *hscommon.TextureLoader, id string, dc6 *d2dc6.DC6) *DC6ViewerWidget {
	result := &DC6ViewerWidget{
		id:            id,
		dc6:           dc6,
		textureLoader: textureLoader,
	}

	return result
}

// Build builds a widget
func (p *DC6ViewerWidget) Build() {
	stateID := fmt.Sprintf("DC6ViewerWidget_%s", p.id)

	state := g.Context.GetState(stateID)
	if state == nil {
		p.buildNew(stateID)
	} else {
		viewerState := state.(*DC6ViewerState)

		vs := (viewerState.lastDirection != viewerState.controls.direction || viewerState.lastFrame != viewerState.controls.frame)
		if !viewerState.loadingTexture && vs {
			// Control values have changed, need to regenerate the texture
			viewerState.lastDirection = viewerState.controls.direction
			viewerState.lastFrame = viewerState.controls.frame
			viewerState.loadingTexture = true
			viewerState.texture = nil

			g.Context.SetState(stateID, viewerState)

			p.textureLoader.CreateTextureFromARGB(
				viewerState.rgb[viewerState.lastFrame+(viewerState.lastDirection*int32(viewerState.framesPerDirection))],
				func(tex *g.Texture,
				) {
					newState := g.Context.GetState(stateID).(*DC6ViewerState)

					newState.texture = tex
					newState.loadingTexture = false
					g.Context.SetState(stateID, newState)
				})
		}

		imageScale := uint32(viewerState.controls.scale)
		curFrameIndex := int(viewerState.controls.frame) + (int(viewerState.controls.direction) * int(p.dc6.FramesPerDirection))

		if imageScale < 1 {
			imageScale = 1
		}

		err := g.Context.GetRenderer().SetTextureMagFilter(g.TextureFilterNearest)
		if err != nil {
			log.Print(err)
		}

		var widget *g.ImageWidget
		w := float32(p.dc6.Frames[curFrameIndex].Width * imageScale)
		h := float32(p.dc6.Frames[curFrameIndex].Height * imageScale)
		if viewerState.texture == nil {
			widget = g.Image(nil).Size(w, h)
		} else {
			widget = g.Image(viewerState.texture).Size(w, h)
		}

		g.Layout{
			g.Label(fmt.Sprintf(
				"Version: %v\t Flags: %b\t Encoding: %v\t",
				p.dc6.Version,
				int64(p.dc6.Flags),
				p.dc6.Encoding,
			)),
			g.Label(fmt.Sprintf("Directions: %v\tFrames per Direction: %v", p.dc6.Directions, p.dc6.FramesPerDirection)),
			g.Custom(func() {
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
			g.Separator(),
			widget,
		}.Build()
	}
}

func (p *DC6ViewerWidget) buildNew(stateID string) {
	var widget *g.ImageWidget

	// Prevent multiple invocation to LoadImage.
	newState := &DC6ViewerState{
		lastFrame:          -1,
		lastDirection:      -1,
		framesPerDirection: p.dc6.FramesPerDirection,
	}

	sw := float32(p.dc6.Frames[0].Width)
	sh := float32(p.dc6.Frames[0].Height)
	widget = g.Image(nil).Size(sw, sh)

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

	g.Context.SetState(stateID, newState)

	widget.Build()
}
