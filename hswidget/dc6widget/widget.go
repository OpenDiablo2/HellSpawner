package dc6widget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

const (
	// nolint:gomnd // it was constant
	maxAlpha = uint8(255)
)

// widget represents dc6viewer's widget
type widget struct {
	id            string
	dc6           *d2dc6.DC6
	textureLoader *hscommon.TextureLoader
}

// DC6Viewer creates new widget
func Create(textureLoader *hscommon.TextureLoader, id string, dc6 *d2dc6.DC6) giu.Widget {
	result := &widget{
		id:            id,
		dc6:           dc6,
		textureLoader: textureLoader,
	}

	return result
}

// Build builds a widget
func (p *widget) Build() {
	state := p.getState()

	switch state.mode {
	case dc6WidgetViewer:
		p.makeViewerLayout().Build()
	}
}

func (p *widget) makeViewerLayout() giu.Layout {
	viewerState := p.getState()

	vs := (viewerState.lastDirection != viewerState.controls.direction || viewerState.lastFrame != viewerState.controls.frame)
	if !viewerState.loadingTexture && vs {
		// Control values have changed, need to regenerate the texture
		viewerState.lastDirection = viewerState.controls.direction
		viewerState.lastFrame = viewerState.controls.frame
		viewerState.loadingTexture = true
		viewerState.texture = nil

		p.setState(viewerState)

		p.textureLoader.CreateTextureFromARGB(
			viewerState.rgb[viewerState.lastFrame+(viewerState.lastDirection*int32(viewerState.framesPerDirection))],
			func(tex *giu.Texture,
			) {
				newState := p.getState()

				newState.texture = tex
				newState.loadingTexture = false
				p.setState(newState)
			})
	}

	imageScale := uint32(viewerState.controls.scale)
	curFrameIndex := int(viewerState.controls.frame) + (int(viewerState.controls.direction) * int(p.dc6.FramesPerDirection))

	if imageScale < 1 {
		imageScale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Print(err)
	}

	var widget *giu.ImageWidget
	w := float32(p.dc6.Frames[curFrameIndex].Width * imageScale)
	h := float32(p.dc6.Frames[curFrameIndex].Height * imageScale)
	if viewerState.texture == nil {
		widget = giu.Image(nil).Size(w, h)
	} else {
		widget = giu.Image(viewerState.texture).Size(w, h)
	}

	return giu.Layout{
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
	}
}
