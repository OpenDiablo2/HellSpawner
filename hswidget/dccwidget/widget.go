package dccwidget

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/OpenDiablo2/dialog"
	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

const (
	inputIntW           = 30
	playPauseButtonSize = 15
	comboW              = 125
)

const (
	maxAlpha = uint8(255)
)

const (
	imageW, imageH = 32, 32
)

type widget struct {
	id            string
	dcc           *d2dcc.DCC
	palette       *[256]d2interface.Color
	textureLoader hscommon.TextureLoader
}

// Create creates a new dcc widget
func Create(tl hscommon.TextureLoader, state []byte, palette *[256]d2interface.Color, id string, dcc *d2dcc.DCC) giu.Widget {
	result := &widget{
		id:            id,
		dcc:           dcc,
		palette:       palette,
		textureLoader: tl,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		if err := json.Unmarshal(state, s); err != nil {
			log.Printf("error decoding dcc widget state: %v", err)
		}

		// update ticker
		s.ticker.Reset(time.Second * time.Duration(s.TickTime) / miliseconds)
		result.setState(s)
	}

	return result
}

// Build build a widget
func (p *widget) Build() {
	viewerState := p.getState()

	imageScale := uint32(viewerState.Controls.Scale)
	dirIdx := int(viewerState.Controls.Direction)
	frameIdx := viewerState.Controls.Frame

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
		giu.Row(
			giu.Label(fmt.Sprintf("Signature: %v", p.dcc.Signature)),
			giu.Label(fmt.Sprintf("Version: %v", p.dcc.Version)),
		),
		giu.Row(
			giu.Label(fmt.Sprintf("Directions: %v", p.dcc.NumberOfDirections)),
			giu.Label(fmt.Sprintf("Frames per Direction: %v", p.dcc.FramesPerDirection)),
		),
		giu.Custom(func() {
			imgui.BeginGroup()
			if p.dcc.NumberOfDirections > 1 {
				imgui.SliderInt("Direction", &viewerState.Controls.Direction, 0, int32(p.dcc.NumberOfDirections-1))
			}

			if p.dcc.FramesPerDirection > 1 {
				imgui.SliderInt("Frames", &viewerState.Controls.Frame, 0, int32(p.dcc.FramesPerDirection-1))
			}

			const minScale, maxScale = 1, 8

			imgui.SliderInt("Scale", &viewerState.Controls.Scale, minScale, maxScale)

			imgui.EndGroup()
		}),
		giu.Separator(),
		p.makePlayerLayout(viewerState),
		giu.Separator(),
		widget,
	}.Build()
}

func (p *widget) makePlayerLayout(state *widgetState) giu.Layout {
	playModeList := make([]string, 0)
	for i := playModeForward; i <= playModePingPong; i++ {
		playModeList = append(playModeList, i.String())
	}

	pm := int32(state.PlayMode)

	return giu.Layout{
		giu.Row(
			giu.Checkbox("Loop##"+p.id+"PlayRepeat", &state.Repeat),
			giu.Combo("##"+p.id+"PlayModeList", playModeList[state.PlayMode], playModeList, &pm).OnChange(func() {
				state.PlayMode = animationPlayMode(pm)
			}).Size(comboW),
			giu.InputInt("Tick time##"+p.id+"PlayTickTime", &state.TickTime).Size(inputIntW).OnChange(func() {
				state.ticker.Reset(time.Second * time.Duration(state.TickTime) / miliseconds)
			}),
			hswidget.PlayPauseButton("##"+p.id+"PlayPauseAnimation", &state.IsPlaying, p.textureLoader).
				Size(playPauseButtonSize, playPauseButtonSize),
			giu.Button("Export GIF##"+p.id+"exportGif").OnClick(func() {
				err := p.exportGif(state)
				if err != nil {
					dialog.Message(err.Error()).Error()
				}
			}),
		),
	}
}

func (p *widget) exportGif(state *widgetState) error {
	fpd := int32(p.dcc.FramesPerDirection)
	firstFrame := state.Controls.Direction * fpd
	images := state.images[firstFrame : firstFrame+fpd]

	err := hsutil.ExportToGif(images, state.TickTime)
	if err != nil {
		return fmt.Errorf("error creating gif file: %w", err)
	}

	return nil
}
