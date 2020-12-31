package hsaboutdialog

import (
	"image/color"
	"io/ioutil"
	"log"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/giu/imgui"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

type AboutDialog struct {
	hsdialog.Dialog
	titleFont   imgui.Font
	regularFont imgui.Font
	fixedFont   imgui.Font
	credits     string
	license     string
}

func Create(regularFont, titleFont, fixedFont imgui.Font) (*AboutDialog, error) {
	result := &AboutDialog{
		titleFont:   titleFont,
		regularFont: regularFont,
		fixedFont:   fixedFont,
	}
	var err error
	var data []byte

	if data, err = ioutil.ReadFile("LICENSE"); err != nil {
		log.Fatal(err)
	}
	result.license = string(data)

	if data, err = ioutil.ReadFile("CONTRIBUTORS"); err != nil {
		log.Fatal(err)
	}
	result.credits = string(data)

	return result, nil
}

func (a *AboutDialog) Render() {
	if !a.Visible {
		return
	}

	imgui.SetNextWindowFocus()
	g.WindowV("About HellSpawner", &a.Visible, g.WindowFlagsNoResize|g.WindowFlagsNoCollapse, 200, 100, 0, 0, g.Layout{
		g.Line(
			g.ImageWithFile("d2logo.png", 256, 256),
			g.Child("AboutHellSpawnerLayout", false, 500, 0, g.WindowFlagsNone, g.Layout{
				g.LabelV("HellSpawner", false, &color.RGBA{R: 255, G: 255, B: 255, A: 255}, &a.titleFont),
				g.LabelV("The OpenDiablo 2 Toolset", false, &color.RGBA{R: 255, G: 255, B: 255, A: 255}, &a.regularFont),
				g.LabelV("Local Build", false, &color.RGBA{R: 255, G: 255, B: 255, A: 255}, &a.fixedFont),
				g.Separator(),
				g.TabBarV("AboutHellSpawnerTabBar", g.TabBarFlagsNoCloseWithMiddleMouseButton, g.Layout{
					g.TabItem("Credits##AboutHellSpawner", g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("", &a.credits, 500, 150, g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll, nil, nil),
						g.Custom(func() { g.PopFont() }),
					}),
					g.TabItem("Licenses##AboutHellSpawner", g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("", &a.license, 500, 150, g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll, nil, nil),
						g.Custom(func() { g.PopFont() }),
					}),
				}),
			}),
		),
	})
}
