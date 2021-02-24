package dccwidget

import (
	"fmt"
	"log"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
)

const (
	// nolint:gomnd // constant = constant
	maxAlpha = uint8(255)
)

const (
	imageW, imageH = 32, 32
)

type widget struct {
	id  string
	dcc *d2dcc.DCC
}

// Create creates a new dcc widget
func Create(id string, dcc *d2dcc.DCC) giu.Widget {
	result := &widget{
		id:  id,
		dcc: dcc,
	}

	return result
}

// Build build a widget
// nolint:funlen // no need to change
func (p *widget) Build() {
	viewerState := p.getState()

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
		widget = giu.Image(nil).Size(imageW, imageH)
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
