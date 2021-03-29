// Package hsaboutdialog contains about dialog's data
package hsaboutdialog

import (
	"fmt"
	"io/ioutil"
	"strings"

	g "github.com/ianling/giu"
	"github.com/ianling/imgui-go"
	"github.com/jaytaylor/html2text"
	"github.com/russross/blackfriday"

	"github.com/OpenDiablo2/HellSpawner/hsassets"
	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hsdialog"
)

const (
	mainWindowW, mainWindowH = 256, 256
	mainLayoutW, mainLayoutH = 500, -1
)

const (
	white = 0xffffffff
)

// AboutDialog represents about dialog
type AboutDialog struct {
	*hsdialog.Dialog
	titleFont   imgui.Font
	regularFont imgui.Font
	fixedFont   imgui.Font
	credits     string
	license     string
	readme      string
	logo        *g.Texture
}

// Create creates a new AboutDialog
func Create(textureLoader *hscommon.TextureLoader, regularFont, titleFont, fixedFont imgui.Font) (*AboutDialog, error) {
	result := &AboutDialog{
		Dialog:      hsdialog.New("About HellSpawner"),
		titleFont:   titleFont,
		regularFont: regularFont,
		fixedFont:   fixedFont,
	}

	textureLoader.CreateTextureFromFile(hsassets.MakeReader(hsassets.HellSpawnerLogo), func(t *g.Texture) {
		result.logo = t
	})

	var err error

	var data []byte

	if data, err = ioutil.ReadFile("LICENSE"); err != nil {
		data = nil
	}

	result.license = string(data)

	if data, err = ioutil.ReadFile("CONTRIBUTORS"); err != nil {
		data = nil
	}

	result.credits = string(data)

	if data, err = ioutil.ReadFile("README.md"); err != nil {
		data = nil
	}

	// convert output md to html
	html := blackfriday.MarkdownBasic(data)
	// convert html to text
	text, err := html2text.FromString(string(html), html2text.Options{PrettyTables: true})
	if err != nil {
		return result, fmt.Errorf("error converting HTML to text: %w", err)
	}

	// set string's max length
	text = strings.Join(hsutil.SplitIntoLinesWithMaxWidth(text, 70), "\n")
	result.readme = text

	return result, nil
}

// Build build an about dialog
func (a *AboutDialog) Build() {
	colorWhite := hsutil.Color(white)
	a.IsOpen(&a.Visible).Layout(g.Layout{
		g.Line(
			g.Image(a.logo).Size(mainWindowW, mainWindowH),
			g.Child("AboutHellSpawnerLayout").Size(mainLayoutW, mainLayoutH).Layout(g.Layout{
				g.Label("HellSpawner").Color(&colorWhite).Font(&a.titleFont),
				g.Label("The OpenDiablo 2 Toolset").Color(&colorWhite).Font(&a.regularFont),
				g.Label("Local Build").Color(&colorWhite).Font(&a.fixedFont),
				g.Separator(),
				g.TabBar("AboutHellSpawnerTabBar").Flags(g.TabBarFlagsNoCloseWithMiddleMouseButton).Layout(g.Layout{
					g.TabItem("README##AboutHellSpawner").Layout(g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("##AboutHellSpawnerReadme", &a.readme).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { g.PopFont() }),
					}),
					g.TabItem("Credits##AboutHellSpawner").Layout(g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("##AboutHellSpawnerCredits", &a.credits).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { g.PopFont() }),
					}),
					g.TabItem("Licenses##AboutHellSpawner").Layout(g.Layout{
						g.Custom(func() { g.PushFont(a.fixedFont) }),
						g.InputTextMultiline("##AboutHellSpawnerLicense", &a.license).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly | g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { g.PopFont() }),
					}),
				}),
			}),
		),
	})
}
