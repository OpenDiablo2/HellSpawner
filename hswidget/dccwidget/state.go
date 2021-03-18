package dccwidget

import (
	"fmt"
	image2 "image"
	"image/color"
	"log"

	"github.com/ianling/giu"
)

type widgetState struct {
	controls struct {
		direction int32
		frame     int32
		scale     int32
	}

	textures []*giu.Texture
}

// Dispose cleans viewers state
func (s *widgetState) Dispose() {
	s.textures = nil
}

func (p *widget) getStateID() string {
	return fmt.Sprintf("widget_%s", p.id)
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
	p.setState(&widgetState{})

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

					var r, g, b uint8
					if p.palette != nil {
						col := p.palette[val]
						r, g, b = col.R(), col.G(), col.B()
					} else {
						r, g, b = val, val, val
					}

					RGBAcolor := color.RGBA{
						R: r,
						G: g,
						B: b,
						A: alpha,
					}

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
		p.setState(&widgetState{textures: textures})
	}()

	// display a temporary dummy image until the real one ready
	firstFrame := p.dcc.Directions[0].Frames[0]
	sw := float32(firstFrame.Width)
	sh := float32(firstFrame.Height)
	widget := giu.Image(nil).Size(sw, sh)
	widget.Build()
}

func (p *widget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}
