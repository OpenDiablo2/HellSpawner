package hsapp

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/ioutil"
	"runtime"
	"strconv"

	"github.com/OpenDiablo2/HellSpawner/hscommon"

	"github.com/OpenDiablo2/HellSpawner/hsutil"
	"github.com/hajimehoshi/ebiten/ebitenutil"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsui"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"golang.org/x/image/font"
)

const (
	bytesToMegabyte          = 1024 * 1024
	windowFrameHeight    int = 30
	windowFrameThickness int = 3
	defaultWindowWidth       = 1280
	defaultWindowHeight      = 720
	defaultFontDPI           = 96
)

// App is an instance of the HellSpawner app.
type App struct {
	Config         hsconfig.AppConfig
	NormalFont     font.Face
	SymbolFont     font.Face
	MonospaceFont  font.Face
	ttNormal       *truetype.Font
	ttMono         *truetype.Font
	ttSymbols      *truetype.Font
	screenWidth    int
	screenHeight   int
	mouseX, mouseY int

	testbox     *hsui.VBox
	testpager   *hsui.Pager
	testTabView *hsui.TabView
}

func (a *App) GetAppConfig() *hsconfig.AppConfig {
	return &a.Config
}

func (a *App) GetNormalFont() font.Face {
	return a.NormalFont
}

func (a *App) GetSymbolsFont() font.Face {
	return a.SymbolFont
}

func (a *App) GetMonospaceFont() font.Face {
	return a.MonospaceFont
}

// Create creates an instance of the HelSpawner app.
func Create() (*App, error) {
	result := &App{}

	var err error

	var configBytes []byte

	// Load the configuration file
	// TODO: Check for a custom config file
	if configBytes, err = ioutil.ReadFile("config.json"); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(configBytes, &result.Config); err != nil {
		return nil, err
	}

	// Set up the initial window layout and properties
	ebiten.SetWindowSize(defaultWindowWidth, defaultWindowHeight)
	ebiten.SetWindowTitle("HellSpawner")
	ebiten.SetWindowResizable(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowDecorated(true)

	// Configure the fonts for rendering
	if err = result.configureFonts(); err != nil {
		return nil, err
	}

	// Store off the device scale factor so we can regenerate if we need to
	hsutil.SetDeviceScale(ebiten.DeviceScaleFactor())

	result.createTestBox()
	result.createTestTabView()

	return result, nil
}

func (a *App) createTestBox() {
	a.testbox = hsui.CreateVBox()

	// button captions
	const (
		alignVertTop    = "Align Top"
		alignVertMiddle = "Align Middle"
		alignVertBottom = "Align Bottom"
		visToggle       = "Vis Toggle"
		expandToggle    = "Toggle Expand Child"
		spaceInc        = "Child Spacing +"
		spaceDec        = "Child Spacing -"
	)

	buttons := make(map[string]*hsui.Button)

	buttonConfigs := []struct {
		caption  string
		callback func()
	}{
		{
			alignVertTop,
			func() { a.testbox.SetAlignment(hscommon.VAlignTop) },
		},
		{
			alignVertMiddle,
			func() { a.testbox.SetAlignment(hscommon.VAlignMiddle) },
		},
		{
			alignVertBottom,
			func() { a.testbox.SetAlignment(hscommon.VAlignBottom) },
		},
		{
			visToggle,
			func() { buttons[alignVertMiddle].ToggleVisible() },
		},
		{
			expandToggle,
			func() { a.testbox.ToggleExpandChild() },
		},
		{
			spaceInc,
			func() { a.testbox.SetChildSpacing(a.testbox.GetChildSpacing() + 1) },
		},
		{
			spaceDec,
			func() { a.testbox.SetChildSpacing(a.testbox.GetChildSpacing() - 1) },
		},
	}

	for idx := range buttonConfigs {
		cfg := &buttonConfigs[idx]
		button := hsui.CreateButton(a, cfg.caption, cfg.callback)
		buttons[cfg.caption] = button

		a.testbox.AddChild(button)
	}

	hbox := hsui.CreateHBox()
	hbox.SetExpandChild(true)
	hbox.AddChild(hsui.CreateButton(a, "Left", func() {}))
	hbox.AddChild(hsui.CreateButton(a, "Center", func() {}))
	hbox.AddChild(hsui.CreateButton(a, "Right", func() {}))

	a.testbox.AddChild(hbox)
}

func (a *App) createTestPager() {
	numPages := 10

	// each page is a grid of buttons
	minGridOrder := 1
	maxGridOrder := 10

	pager := hsui.CreatePager(300, 300, nil)

	fn := map[bool]func(){
		false: pager.SelectPreviousChild,
		true:  pager.SelectNextChild,
	}

	// just making a grid of buttons as a test
	for pageIdx, order := 0, minGridOrder; pageIdx < numPages && order <= maxGridOrder; pageIdx++ {
		outerVbox := hsui.CreateVBox()
		outerVbox.SetExpandChild(true)
		rows, columns := order, order

		for rowIdx := 0; rowIdx < rows; rowIdx++ {
			row := hsui.CreateHBox()
			row.SetExpandChild(true)

			for colIdx := 0; colIdx < columns; colIdx++ {
				caption := "next"

				button := hsui.CreateButton(a, caption, fn[true])
				fn[true]()

				row.AddChild(button)
			}

			outerVbox.AddChild(row)
		}

		pager.AddChild(outerVbox)

		order++
	}

	pager.SetSelectedChild(0)

	a.testpager = pager
}

func NOOP() {}

func (a *App) createTestTabView() {

	// each page is a grid of buttons
	minGridOrder := 5
	maxGridOrder := 10

	numPages := maxGridOrder - minGridOrder

	tabView := hsui.CreateTabView(a, 300, 300)

	// just making a grid of buttons as a test
	for pageIdx, order := 0, minGridOrder; pageIdx < numPages && order <= maxGridOrder; pageIdx++ {
		tabTitle := fmt.Sprintf("%dx%d", pageIdx+1, pageIdx+1)

		outerVbox := hsui.CreateVBox()
		outerVbox.SetExpandChild(true)
		rows, columns := order, order

		tabView.AddTab(tabTitle, outerVbox, true)

		for rowIdx := 0; rowIdx < rows; rowIdx++ {
			row := hsui.CreateHBox()
			row.SetExpandChild(true)

			for colIdx := 0; colIdx < columns; colIdx++ {
				caption := "test"

				button := hsui.CreateButton(a, caption, NOOP)

				row.AddChild(button)
			}

			outerVbox.AddChild(row)
		}

		order++
	}

	a.testTabView = tabView
}

func (a *App) Run() error {
	return ebiten.RunGame(a)
}

func (a *App) Update(*ebiten.Image) error {
	deviceScale := ebiten.DeviceScaleFactor()
	a.mouseX, a.mouseY = ebiten.CursorPosition()

	// If the device scale has changed, we need to regenerate the fonts
	if deviceScale != hsutil.GetLastDeviceScale() {
		hsutil.SetDeviceScale(deviceScale)
		a.regenerateFonts()
		a.testbox.Invalidate()
		a.testTabView.Invalidate()
	}

	a.testbox.Update()
	a.testTabView.Update()

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	frameColor := a.Config.Colors.WindowBackground

	// Fill the window with the frame color
	_ = screen.Fill(color.RGBA{R: frameColor[0], G: frameColor[1], B: frameColor[2], A: frameColor[3]})

	const testSplitPoint = 300
	a.testbox.Render(screen, 0, 0, testSplitPoint, a.screenHeight)
	a.testTabView.Render(screen, testSplitPoint, 0, a.screenWidth-testSplitPoint, a.screenHeight)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Store off the screen size for easy access
	a.screenWidth = outsideWidth
	a.screenHeight = outsideHeight

	// Return the actual resolution, determined by the virtual size and screen device scale
	return hsutil.ScaleToDevice(outsideWidth), hsutil.ScaleToDevice(outsideHeight)
}

// configureFonts loads the fonts the app needs.
func (a *App) configureFonts() error {
	var ttNormalBytes, ttMonoBytes, ttSymbolBytes []byte

	var err error

	// Load font data from the files into a byte array
	if ttNormalBytes, err = ioutil.ReadFile(a.Config.Fonts.Normal.Face); err != nil {
		return err
	}

	if ttSymbolBytes, err = ioutil.ReadFile(a.Config.Fonts.Symbols.Face); err != nil {
		return err
	}

	if ttMonoBytes, err = ioutil.ReadFile(a.Config.Fonts.Monospaced.Face); err != nil {
		return err
	}

	// Parse the TTF font files
	if a.ttNormal, err = truetype.Parse(ttNormalBytes); err != nil {
		return err
	}

	if a.ttSymbols, err = truetype.Parse(ttSymbolBytes); err != nil {
		return err
	}

	if a.ttMono, err = truetype.Parse(ttMonoBytes); err != nil {
		return err
	}

	// Generate the fonts
	a.regenerateFonts()

	return nil
}

// regenerateFonts regenerates the fonts (at startup and when dragging to a different DPI display).
func (a *App) regenerateFonts() {
	deviceScale := ebiten.DeviceScaleFactor()

	a.NormalFont = truetype.NewFace(a.ttNormal, &truetype.Options{
		Size:    float64(a.Config.Fonts.Normal.Size),
		DPI:     defaultFontDPI * deviceScale,
		Hinting: font.HintingNone,
	})

	a.SymbolFont = truetype.NewFace(a.ttSymbols, &truetype.Options{
		Size:    float64(a.Config.Fonts.Symbols.Size),
		DPI:     defaultFontDPI * deviceScale,
		Hinting: font.HintingNone,
	})

	a.MonospaceFont = truetype.NewFace(a.ttMono, &truetype.Options{
		Size:    float64(a.Config.Fonts.Monospaced.Size),
		DPI:     defaultFontDPI * deviceScale,
		Hinting: font.HintingNone,
	})
}

func (a *App) printDebugInfo(screen *ebiten.Image) {
	// Debug print stuff
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)
	ebitenutil.DebugPrintAt(screen,
		"Alloc:   "+strconv.FormatInt(int64(m.Alloc)/bytesToMegabyte, 10)+"\n"+
			"Pause:   "+strconv.FormatInt(int64(m.PauseTotalNs/bytesToMegabyte), 10)+"\n"+
			"HeapSys: "+strconv.FormatInt(int64(m.HeapSys/bytesToMegabyte), 10)+"\n"+
			"NumGC:   "+strconv.FormatInt(int64(m.NumGC), 10),
		hsutil.ScaleToDevice(550), hsutil.ScaleToDevice(100))
}
