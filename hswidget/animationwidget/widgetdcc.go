package animationwidget

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"time"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	imageW, imageH = 32, 32
)

type dccWidget struct {
	*widget
	dcc *d2dcc.DCC
}

func (w *dccWidget) getDcImage() dcImage {
	return w.dcc
}

// Build build a widget
func (w *dccWidget) Build() {
	viewerState := w.getState()

	imageScale := uint32(viewerState.controls.scale)
	dirIdx := int(viewerState.controls.direction)
	frameIdx := viewerState.controls.frame

	textureIdx := dirIdx*len(w.dcc.Directions[dirIdx].Frames) + int(frameIdx)

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
		bw := w.dcc.Directions[dirIdx].Box.Width
		bh := w.dcc.Directions[dirIdx].Box.Height
		w := float32(uint32(bw) * imageScale)
		h := float32(uint32(bh) * imageScale)
		widget = giu.Image(viewerState.textures[textureIdx]).Size(w, h)
	}

	giu.Layout{
		giu.Row(
			giu.Label(fmt.Sprintf("Signature: %v", w.dcc.Signature)),
			giu.Label(fmt.Sprintf("Version: %v", w.dcc.Version)),
		),
		giu.Row(
			giu.Label(fmt.Sprintf("Directions: %v", w.dcc.NumberOfDirections)),
			giu.Label(fmt.Sprintf("Frames per Direction: %v", w.dcc.FramesPerDirection)),
		),
		giu.Custom(func() {
			imgui.BeginGroup()
			if w.dcc.NumberOfDirections > 1 {
				imgui.SliderInt("Direction", &viewerState.controls.direction, 0, int32(w.dcc.NumberOfDirections-1))
			}

			if w.dcc.FramesPerDirection > 1 {
				imgui.SliderInt("Frames", &viewerState.controls.frame, 0, int32(w.dcc.FramesPerDirection-1))
			}

			const minScale, maxScale = 1, 8

			imgui.SliderInt("Scale", &viewerState.controls.scale, minScale, maxScale)

			imgui.EndGroup()
		}),
		giu.Separator(),
		makePlayerLayout(w, viewerState),
		giu.Separator(),
		widget,
	}.Build()
}

func createDccWidget(state []byte, widget *widget, dcc *d2dcc.DCC) giu.Widget {
	dccWidget := &dccWidget{
		widget: widget,
		dcc:    dcc,
	}

	if giu.Context.GetState(dccWidget.widget.getStateID()) == nil && state != nil {
		s := dccWidget.getState()
		s.Decode(state)
		dccWidget.setState(s)
	}

	return dccWidget
}

func (w *dccWidget) getState() *dccWidgetState {
	var state *dccWidgetState

	s := giu.Context.GetState(w.widget.getStateID())

	if s != nil {
		state = s.(*dccWidgetState)
	} else {
		w.initState()
		state = w.getState()
	}

	return state
}

func (w *dccWidget) setState(s giu.Disposable) {
	giu.Context.SetState(w.widget.getStateID(), s)
}

func (w *dccWidget) runPlayer(state *dccWidgetState) {
	for range state.ticker.C {
		if !state.isPlaying {
			continue
		}

		numFrames := int32(w.dcc.FramesPerDirection - 1)
		isLastFrame := state.controls.frame == numFrames

		// update play direction
		switch state.playMode {
		case playModeForward:
			state.isForward = true
		case playModeBackward:
			state.isForward = false
		case playModePingPong:
			if isLastFrame || state.controls.frame == 0 {
				state.isForward = !state.isForward
			}
		}

		// now update the frame number
		if state.isForward {
			state.controls.frame++
		} else {
			state.controls.frame--
		}

		state.controls.frame = int32(hsutil.Wrap(int(state.controls.frame), w.dcc.FramesPerDirection))

		// next, check for stopping/repeat
		isStoppingFrame := (state.controls.frame == 0) || (state.controls.frame == numFrames)

		if isStoppingFrame && !state.repeat {
			state.isPlaying = false
		}
	}
}

func (w *dccWidget) makeImagePixel(val byte) color.RGBA {
	alpha := maxAlpha

	if val == 0 {
		alpha = 0
	}

	var r, g, b uint8

	if w.widget.palette != nil {
		col := w.widget.palette[val]
		r, g, b = col.R(), col.G(), col.B()
	} else {
		r, g, b = val, val, val
	}

	RGBAColor := color.RGBA{
		R: r,
		G: g,
		B: b,
		A: alpha,
	}

	return RGBAColor
}

func (w *dccWidget) initState() {
	// Prevent multiple invocation to LoadImage.
	state := &dccWidgetState{
		widgetState: widgetState{
			isPlaying: false,
			repeat:    false,
			tickTime:  defaultTickTime,
			playMode:  playModeForward,
		},
	}

	state.ticker = time.NewTicker(time.Second * time.Duration(state.tickTime) / miliseconds)

	w.setState(state)

	go w.runPlayer(state)

	totalFrames := w.dcc.NumberOfDirections * w.dcc.FramesPerDirection
	state.images = make([]*image.RGBA, totalFrames)

	for dirIdx := range w.dcc.Directions {
		fw := w.dcc.Directions[dirIdx].Box.Width
		fh := w.dcc.Directions[dirIdx].Box.Height

		for frameIdx := range w.dcc.Directions[dirIdx].Frames {
			absoluteFrameIdx := (dirIdx * w.dcc.FramesPerDirection) + frameIdx

			frame := w.dcc.Directions[dirIdx].Frames[frameIdx]
			pixels := frame.PixelData

			state.images[absoluteFrameIdx] = image.NewRGBA(image.Rect(0, 0, fw, fh))

			for y := 0; y < fh; y++ {
				for x := 0; x < fw; x++ {
					idx := x + (y * fw)
					if idx >= len(pixels) {
						continue
					}

					val := pixels[idx]

					RGBAColor := w.makeImagePixel(val)
					state.images[absoluteFrameIdx].Set(x, y, RGBAColor)
				}
			}
		}
	}

	go func() {
		textures := make([]*giu.Texture, totalFrames)

		for frameIndex := 0; frameIndex < totalFrames; frameIndex++ {
			frameIndex := frameIndex
			w.widget.textureLoader.CreateTextureFromARGB(state.images[frameIndex], func(t *giu.Texture) {
				textures[frameIndex] = t
			})
		}

		s := w.getState()
		s.textures = textures
		w.setState(s)
	}()
}
