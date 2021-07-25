// Package hsaboutdialog provides the "About" window implementation, which shows information about hellspawner.
package hsaboutdialog

import (
	"fmt"
	"io/ioutil"
	"strings"

	g "github.com/AllenDang/giu"
	"github.com/AllenDang/imgui-go"
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
func Create(textureLoader hscommon.TextureLoader, regularFont, titleFont, fixedFont imgui.Font) (*AboutDialog, error) {
	result := &AboutDialog{
		Dialog:      hsdialog.New("About HellSpawner"),
		titleFont:   titleFont,
		regularFont: regularFont,
		fixedFont:   fixedFont,
	}

	textureLoader.CreateTextureFromFile(hsassets.HellSpawnerLogo, func(t *g.Texture) {
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
	const maxColumns = 70
	text = strings.Join(hsutil.SplitIntoLinesWithMaxWidth(text, maxColumns), "\n")
	result.readme = text

	return result, nil
}

// Build build an about dialog
func (a *AboutDialog) Build() {
	colorWhite := hsutil.Color(white)
	a.IsOpen(&a.Visible).Layout(
		g.Row(
			g.Image(a.logo).Size(mainWindowW, mainWindowH),
			g.Child().Size(mainLayoutW, mainLayoutH).Layout(
				g.Style().SetColor(g.StyleColorText, colorWhite).To(
					g.Custom(func() { imgui.PushFont(a.titleFont) }),
					g.Label("HellSpawner"),
					g.Custom(func() { imgui.PopFont() }),
					g.Custom(func() { imgui.PushFont(a.regularFont) }),
					g.Label("The OpenDiablo 2 Toolset"),
					g.Custom(func() { imgui.PopFont() }),
					g.Custom(func() { imgui.PushFont(a.fixedFont) }),
					g.Label("Local Build"),
					g.Custom(func() { imgui.PopFont() }),
				),
				g.Separator(),
				g.TabBar().Flags(g.TabBarFlagsNoCloseWithMiddleMouseButton).TabItems(
					g.TabItem("README##AboutHellSpawner").Layout(
						g.Custom(func() { imgui.PushFont(a.fixedFont) }),
						g.InputTextMultiline(&a.readme).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { imgui.PopFont() }),
					),
					g.TabItem("Credits##AboutHellSpawner").Layout(
						g.Custom(func() { imgui.PushFont(a.fixedFont) }),
						g.InputTextMultiline(&a.credits).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { imgui.PopFont() }),
					),
					g.TabItem("Licenses##AboutHellSpawner").Layout(
						g.Custom(func() { imgui.PushFont(a.fixedFont) }),
						g.InputTextMultiline(&a.license).
							Size(-1, -1).Flags(g.InputTextFlagsReadOnly|g.InputTextFlagsNoHorizontalScroll),
						g.Custom(func() { imgui.PopFont() }),
					),
				),
			),
		),
	).Build()
}
