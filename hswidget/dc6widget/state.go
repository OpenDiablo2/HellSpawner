package dc6widget

import (
	"fmt"
	image2 "image"
	"image/color"

	"github.com/ianling/giu"
)

type dc6WidgetMode int

const (
	dc6WidgetViewer dc6WidgetMode = iota
)

// widgetState represents dc6 viewer's state
type widgetState struct {
	viewerState
	mode dc6WidgetMode
}

func (w *widgetState) Dispose() {
	w.viewerState.Dispose()
	w.mode = dc6WidgetViewer
}

// nolint:structcheck // :-/ linter bug?! thes values are deffinitly used
type viewerState struct {
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

// Dispose cleans state content
func (is *viewerState) Dispose() {
	is.texture = nil
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("DC6Widget_%s", p.id)
}

func (p *widget) getState() *widgetState {
	var state *widgetState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*widgetState)
	} else {
		p.initState()
		state = p.getState()
	}

	return state
}

func (p *widget) initState() {
	// Prevent multiple invocation to LoadImage.
	newState := &widgetState{
		mode: dc6WidgetViewer,
		viewerState: viewerState{
			lastFrame:          -1,
			lastDirection:      -1,
			framesPerDirection: p.dc6.FramesPerDirection,
		},
	}

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

				var r, g, b uint8

				if p.palette != nil {
					col := p.palette[val]
					r, g, b = col.R(), col.G(), col.B()
				} else {
					r, g, b = val, val, val
				}

				newState.rgb[frameIndex].Set(
					x, y,
					color.RGBA{
						R: r,
						G: g,
						B: b,
						A: alpha,
					},
				)
			}
		}
	}

	p.setState(newState)
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
