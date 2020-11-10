package hsapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hsconfig"
	"github.com/OpenDiablo2/HellSpawner/hsui"
	"github.com/OpenDiablo2/HellSpawner/hsutil"
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
	asset          *d2asset.AssetManager
	ttNormal       *truetype.Font
	ttMono         *truetype.Font
	ttSymbols      *truetype.Font
	screenWidth    int
	screenHeight   int
	mouseX, mouseY int

	rootWidget  *hsui.Modal
	mainTabView *hsui.TabView
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

// Create creates an instance of the HellSpawner app.
func Create() (*App, error) {
	var err error

	result := &App{
		rootWidget: hsui.CreateModal(),
	}

	err = result.initAssetManager()
	if err != nil {
		return nil, err
	}

	result.mainTabView = hsui.CreateTabView(result, defaultWindowWidth, defaultWindowHeight)
	result.rootWidget.Push(result.mainTabView)

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

	result.initTests()

	return result, nil
}

func (a *App) initAssetManager() error {
	assetManager, err := d2asset.NewAssetManager()
	if err != nil {
		return err
	}

	a.asset = assetManager

	localDir := filepath.Dir(hsconfig.LocalConfigPath())
	configDir := filepath.Dir(hsconfig.DefaultConfigPath())

	// ensure the config dir exists
	if err := os.MkdirAll(configDir, 0750); err != nil {
		return err
	}

	// bootstrap the two config dir locations we will check
	if _, err = a.asset.Loader.AddSource(localDir); err != nil {
		return err
	}

	if _, err = a.asset.Loader.AddSource(configDir); err != nil {
		return err
	}

	// try to load the config file
	if configBytes, err := a.asset.LoadFile(hsconfig.ConfigFileName); err == nil {
		// unmarshal the loaded data if we found the config file
		return json.Unmarshal(configBytes, &a.Config)
	}

	// create a default and save to disk if we didnt find one
	a.Config = *hsconfig.DefaultConfig()

	return a.Config.Save(hsconfig.DefaultConfigPath())
}

func (a *App) initTests() {
	a.createTestBox()
	a.createTestPager()
	a.createTestTabView()
}

func (a *App) createTestBox() {
	testbox := hsui.CreateVBox()

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
			func() { testbox.SetAlignment(hscommon.VAlignTop) },
		},
		{
			alignVertMiddle,
			func() { testbox.SetAlignment(hscommon.VAlignMiddle) },
		},
		{
			alignVertBottom,
			func() { testbox.SetAlignment(hscommon.VAlignBottom) },
		},
		{
			visToggle,
			func() { buttons[alignVertMiddle].ToggleVisible() },
		},
		{
			expandToggle,
			func() { testbox.ToggleExpandChild() },
		},
		{
			spaceInc,
			func() { testbox.SetChildSpacing(testbox.GetChildSpacing() + 1) },
		},
		{
			spaceDec,
			func() { testbox.SetChildSpacing(testbox.GetChildSpacing() - 1) },
		},
	}

	for idx := range buttonConfigs {
		cfg := &buttonConfigs[idx]
		button := hsui.CreateButton(a, cfg.caption, cfg.callback)
		buttons[cfg.caption] = button

		testbox.AddChild(button)
	}

	hbox := hsui.CreateHBox()
	hbox.SetExpandChild(true)
	hbox.AddChild(hsui.CreateButton(a, "Left", func() {}))
	hbox.AddChild(hsui.CreateButton(a, "Center", func() {}))
	hbox.AddChild(hsui.CreateButton(a, "Right", func() {}))

	testbox.AddChild(hbox)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	tabIconPath, err := filepath.Abs(filepath.Join(dir, "assets", "images", "star.png"))
	if err != nil {
		fmt.Println(err)
		tabIconPath = "" // should be okay
	}

	a.mainTabView.AddTab("testbox", tabIconPath, testbox, false)
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

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	tabIconPath, err := filepath.Abs(filepath.Join(dir, "assets", "images", "star.png"))
	if err != nil {
		fmt.Println(err)
		tabIconPath = "" // should be okay
	}

	pager.SetSelectedChild(0)

	a.mainTabView.AddTab("pager test", tabIconPath, pager, false)
}

func (a *App) createTestTabView() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	tabIconPath, err := filepath.Abs(filepath.Join(dir, "assets", "images", "star.png"))
	if err != nil {
		fmt.Println(err)
		tabIconPath = "" // should be okay
	}

	// each page is a grid of buttons
	testTabView := hsui.CreateTabView(a, 300, 300)

	rand.Seed(time.Now().UnixNano())

	outerVbox := hsui.CreateVBox()
	outerVbox.SetExpandChild(true)
	rows, columns := rand.Intn(10)+1, rand.Intn(10)+1
	tabTitle := fmt.Sprintf("test %dx%d", rows, columns)

	testTabView.AddTab(tabTitle, tabIconPath, outerVbox, true)

	for rowIdx := 0; rowIdx < rows; rowIdx++ {
		row := hsui.CreateHBox()
		row.SetExpandChild(true)

		for colIdx := 0; colIdx < columns; colIdx++ {
			img, _ := hsui.CreateImage(tabIconPath)
			if err != nil {
				continue
			}

			img.SetFit(true)

			row.AddChild(img)
		}

		outerVbox.AddChild(row)
	}

	a.mainTabView.AddTab("test TabView", tabIconPath, testTabView, false)
}

func (a *App) Run() error {
	return ebiten.RunGame(a)
}

func (a *App) Update() error {
	deviceScale := ebiten.DeviceScaleFactor()
	a.mouseX, a.mouseY = ebiten.CursorPosition()

	// If the device scale has changed, we need to regenerate the fonts
	if deviceScale != hsutil.GetLastDeviceScale() {
		hsutil.SetDeviceScale(deviceScale)
		a.regenerateFonts()
		a.rootWidget.Invalidate()
	}

	a.rootWidget.Update()

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	frameColor := hsutil.ArrayToRGBA(a.Config.Colors.WindowBackground)

	// Fill the window with the frame color
	screen.Fill(frameColor)

	const testSplitPoint = 300

	a.rootWidget.Render(screen, 0, 0, a.screenWidth, a.screenHeight)
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
