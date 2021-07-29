package animationwidget

import (
	"fmt"
	"image"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dcc"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/dialog"
	"github.com/ianling/giu"
)

type Widgeter interface {
	getDcImage() DcImage
	getID() string
	getTextureLoader() hscommon.TextureLoader
}

type State interface {
	getDirection() int32
	getImages() []*image.RGBA
	getTickTime() int32
	getTick() *int32
	getPlayMode() animationPlayMode
	setPlayMode(animationPlayMode)
	getRepeat() *bool
	getPlaying() *bool
	getTicker() *time.Ticker
}

type DcImage interface {}

func ExportGif(w Widgeter, s State) error {
	dc := w.getDcImage()
	var fpd int32

	switch dcImage := dc.(type) {
	case d2dc6.DC6:
		fpd = int32(dcImage.FramesPerDirection)
	case d2dcc.DCC:
		fpd = int32(dcImage.FramesPerDirection)
	default:
		return fmt.Errorf("DC File not supported")
	}

	firstFrame := s.getDirection() * fpd
	images := s.getImages()[firstFrame : firstFrame+fpd]

	err := hsutil.ExportToGif(images, s.getTickTime())

	if err != nil {
		return fmt.Errorf("error creating gif file: %w", err)
	}

	return nil
}

func makePlayerLayout(w Widgeter, s State) giu.Layout {
	playModeList := make([]string, 0)
	for i := playModeForward; i <= playModePingPong; i++ {
		playModeList = append(playModeList, i.String())
	}

	pm := int32(s.getPlayMode())
	id := w.getID()

	return giu.Layout{
		giu.Row(
			giu.Checkbox("Loop##"+id+"PlayRepeat", s.getRepeat()),
			giu.Combo("##"+id+"PlayModeList", playModeList[s.getPlayMode()], playModeList, &pm).OnChange(func() {
				s.setPlayMode(animationPlayMode(pm))
			}).Size(comboW),
			giu.InputInt("Tick time##"+id+"PlayTickTime", s.getTick()).Size(inputIntW).OnChange(func() {
				ticker := s.getTicker()
				ticker.Reset(time.Second * time.Duration(s.getTickTime() / miliseconds))
			}),
			hswidget.PlayPauseButton("##"+id+"PlayPauseAnimation", s.getPlaying(), w.getTextureLoader()).
				Size(playPauseButtonSize, playPauseButtonSize),
			giu.Button("Export GIF##"+id+"exportGif").OnClick(func() {
				err := ExportGif(w, s)
				if err != nil {
					dialog.Message(err.Error()).Error()
				}
			}),
		),
	}
}

func CreateAnimationWidget(tl hscommon.TextureLoader, state []byte, palette *[256]d2interface.Color, id string, dc DcImage) (giu.Widget, error) {
	widget := CreateWidget(palette, tl, id)

	switch dcImage := dc.(type) {
	case d2dc6.DC6:
		return createDc6Widget(state, widget, &dcImage), nil
	case d2dcc.DCC:
		return createDccWidget(state, widget, &dcImage), nil
	default:
		return nil, fmt.Errorf("DC File not supported")
	}
}