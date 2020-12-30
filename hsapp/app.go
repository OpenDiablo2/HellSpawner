package hsapp

import (
	"image/color"
	"log"
	"strings"
	"time"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hssoundeditor"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"

	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor/hstexteditor"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/AllenDang/giu/imgui"

	g "github.com/AllenDang/giu"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hstoolwindow/hsmpqexplorer"
	"github.com/sqweek/dialog"
)

type App struct {
	mpqExplorer *hsmpqexplorer.MPQExplorer
	editors     []hscommon.EditorWindow

	fontFixed imgui.Font
}

func Create() (*App, error) {
	result := &App{}
	result.editors = make([]hscommon.EditorWindow, 0)
	var err error

	result.mpqExplorer, err = hsmpqexplorer.Create(result.openEditor)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *App) Run() {
	wnd := g.NewMasterWindow("HellSpawner", 1280, 720, 0, a.setupFonts)
	wnd.SetBgColor(color.RGBA{10, 10, 10, 255})

	sampleRate := beep.SampleRate(22050)
	if err := speaker.Init(sampleRate, sampleRate.N(time.Second/10)); err != nil {
		log.Fatal(err)
	}

	wnd.Main(a.render)
}

func (a *App) render() {
	a.renderMainMenuBar()

	for idx := range a.editors {
		a.editors[idx].Render()
	}

	a.mpqExplorer.Render()

	g.Update()
}

func (a *App) loadMpq(fileName string) {
	a.mpqExplorer.AddMPQ(fileName)
}

func (a *App) buildViewMenu() g.Layout {
	result := make([]g.Widget, 0)

	result = append(result, g.Menu("Tool Windows", g.Layout{
		g.MenuItemV("MPQ Explorer", a.mpqExplorer.Visible, true, func() { a.mpqExplorer.ToggleVisibility() }),
	}))

	if len(a.editors) == 0 {
		return result
	}

	result = append(result, g.Separator())

	for idx := range a.editors {
		result = append(result, g.MenuItem(a.editors[idx].GetWindowTitle(), func() {
			a.editors[idx].Show()
		}))
	}

	return result
}

func (a *App) renderMainMenuBar() {
	g.MainMenuBar(g.Layout{
		g.Menu("File", g.Layout{
			g.MenuItem("Open MPQ...", func() {
				file, err := dialog.File().Filter("MPQ Archive", "mpq").Load()
				if err != nil || len(file) == 0 {
					return
				}
				a.loadMpq(file)
			}),
			g.MenuItem("Exit", nil),
		}),
		g.Menu("View", a.buildViewMenu()),
		g.Menu("Help", g.Layout{
			g.MenuItem("About HellSpawner...", func() {
				dialog.Message("%s", "HellSpawner IDE\nPre-Alpha").Title("About HellSpawner").Info()
			}),
		}),
	}).Build()
}

func (a *App) setupFonts() {
	imgui.CurrentIO().Fonts().AddFontFromFileTTF("NotoSans-Regular.ttf", 17)
	a.fontFixed = imgui.CurrentIO().Fonts().AddFontFromFileTTF("CascadiaCode.ttf", 17)
	imgui.CurrentStyle().ScaleAllSizes(1.0)
}

func (a *App) openEditor(path *hsmpqexplorer.PathEntry) {
	ext := strings.ToLower(path.FullPath[len(path.FullPath)-4:])
	parts := strings.Split(path.FullPath, ":")
	mpqFile := parts[0]
	filePath := cleanMpqPathName(parts[1])
	mpq, err := d2mpq.Load(mpqFile)

	if err != nil {
		log.Fatal(err)
	}

	switch ext {
	case ".txt":
		text, err := mpq.ReadTextFile(filePath)

		if err != nil {
			log.Fatal(err)
		}

		editor, err := hstexteditor.Create(path.Name, text, a.fontFixed)

		if err != nil {
			log.Fatal(err)
		}

		a.editors = append(a.editors, editor)
		editor.Show()
	case ".wav":
		audioStream, err := mpq.ReadFileStream(filePath)

		editor, err := hssoundeditor.Create(path.Name, audioStream)

		if err != nil {
			log.Fatal(err)
		}

		a.editors = append(a.editors, editor)

		editor.Show()
	}
}

func cleanMpqPathName(name string) string {
	name = strings.ReplaceAll(name, "/", "\\")

	if string(name[0]) == "\\" {
		name = name[1:]
	}

	return name
}
