package hsapp

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

const (
	bytesToMegabyte          = 1024 * 1024
	windowFrameHeight    int = 30
	windowFrameThickness int = 3
)

// App is an instance of the HellSpawner app.
type App struct {
	Config          AppConfig
	NormalFont      font.Face
	SymbolFont      font.Face
	MonospaceFont   font.Face
	ttNormal        *truetype.Font
	ttMono          *truetype.Font
	ttSymbols       *truetype.Font
	lastDeviceScale float64
	imgSquare       *ebiten.Image
	screenWidth     int
	screenHeight    int
	mouseX, mouseY  int

	isDraggingWindow  bool
	mouseDragStartX   int
	mouseDragStartY   int
	windowDragStartX  int
	windowDragStartY  int
	leftMouseReleased bool
}

// Create creates an instance of the HelSpawner app.
func Create() (*App, error) {
	result := &App{
		leftMouseReleased: true,
		isDraggingWindow:  false,
	}

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
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("HellSpawner")
	ebiten.SetWindowResizable(true)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowDecorated(false)

	// Configure the fonts for rendering
	if err = result.configureFonts(); err != nil {
		return nil, err
	}

	// Create a square for rendering colored rectangles
	if result.imgSquare, err = ebiten.NewImage(1, 1, ebiten.FilterNearest); err != nil {
		return nil, err
	}

	if err = result.imgSquare.Fill(color.Black); err != nil {
		return nil, err
	}

	// Store off the device scale factor so we can regenerate if we need to
	result.lastDeviceScale = ebiten.DeviceScaleFactor()

	return result, nil
}

func (a *App) Run() error {
	return ebiten.RunGame(a)
}

func (a *App) Update(*ebiten.Image) error {
	deviceScale := ebiten.DeviceScaleFactor()

	if !a.leftMouseReleased && !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		a.leftMouseReleased = true
	}

	a.mouseX, a.mouseY = ebiten.CursorPosition()

	// If the device scale has changed, we need to regenerate the fonts
	if deviceScale != a.lastDeviceScale {
		a.lastDeviceScale = deviceScale
		a.regenerateFonts()
	}

	if a.handleWindowDrag() {
		return nil
	}

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	windowBg := a.Config.Colors.WindowBackground
	frameColor := a.Config.Colors.WindowFrame
	buttonHighlightColor := a.Config.Colors.WindowButtonHighlight
	frameTitleColor := color.RGBA{
		R: a.Config.Colors.WindowFrameText[0],
		G: a.Config.Colors.WindowFrameText[1],
		B: a.Config.Colors.WindowFrameText[2],
		A: a.Config.Colors.WindowFrameText[3],
	}

	// Fill the window with the frame color
	_ = screen.Fill(color.RGBA{R: frameColor[0], G: frameColor[1], B: frameColor[2], A: 200})

	// Fill the desktop area with the background color
	a.drawColoredRect(screen, a.scaleToDevice(windowFrameThickness), a.scaleToDevice(windowFrameHeight+windowFrameThickness),
		a.screenWidth-a.scaleToDevice(windowFrameThickness*2), a.screenHeight-a.scaleToDevice(windowFrameHeight+(windowFrameThickness*2)),
		windowBg[0], windowBg[1], windowBg[2], windowBg[3])

	// Draw the caption
	text.Draw(screen, "HellSpawner", a.NormalFont, a.scaleToDevice(windowFrameThickness*2), a.scaleToDevice((windowFrameHeight/4)*3), frameTitleColor)

	// Draw the highlighter for the min/max/close buttons
	buttonLeftOffset := a.screenWidth - a.scaleToDevice((windowFrameThickness*2)+115)
	if a.mouseY >= 0 && a.mouseY < a.scaleToDevice(windowFrameHeight) && a.mouseX >= buttonLeftOffset && a.mouseX < a.screenWidth {
		posIndex := (a.mouseX - buttonLeftOffset) / a.scaleToDevice(40)
		a.drawColoredRect(screen,
			buttonLeftOffset+a.scaleToDevice(posIndex*40), 0,
			a.scaleToDevice(40), a.scaleToDevice(windowFrameHeight),
			buttonHighlightColor[0], buttonHighlightColor[1], buttonHighlightColor[2], buttonHighlightColor[3])
	}

	// Draw the minimize/maximize/close buttons
	text.Draw(screen, "⌵", a.SymbolFont, a.screenWidth-a.scaleToDevice((windowFrameThickness*2)+100), a.scaleToDevice(((windowFrameHeight/4)*3)-2), frameTitleColor)
	text.Draw(screen, "⛶", a.SymbolFont, a.screenWidth-a.scaleToDevice((windowFrameThickness*2)+60), a.scaleToDevice((windowFrameHeight/4)*3), frameTitleColor)
	text.Draw(screen, "ⅹ", a.SymbolFont, a.screenWidth-a.scaleToDevice((windowFrameThickness*2)+20), a.scaleToDevice(((windowFrameHeight/4)*3)-1), frameTitleColor)

	// Debug print stuff
	// m := &runtime.MemStats{}
	// runtime.ReadMemStats(m)
	// ebitenutil.DebugPrintAt(screen,
	// 	"Alloc:   "+strconv.FormatInt(int64(m.Alloc)/bytesToMegabyte, 10)+"\n"+
	// 		"Pause:   "+strconv.FormatInt(int64(m.PauseTotalNs/bytesToMegabyte), 10)+"\n"+
	//		"HeapSys: "+strconv.FormatInt(int64(m.HeapSys/bytesToMegabyte), 10)+"\n"+
	// 		"NumGC:   "+strconv.FormatInt(int64(m.NumGC), 10),
	// 	a.scaleToDevice(50), a.scaleToDevice(100))
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// Get the device scale factor so we can apply it
	deviceScale := ebiten.DeviceScaleFactor()

	// Store off the screen size for easy access
	a.screenWidth = int(float64(outsideWidth) * deviceScale)
	a.screenHeight = int(float64(outsideHeight) * deviceScale)

	// Return the actual resolution, determined by the virtual size and screen device scale
	return int(float64(outsideWidth) * deviceScale), int(float64(outsideHeight) * deviceScale)
}

// configureFonts loads the fonts the app needs.
func (a *App) configureFonts() error {
	var ttNormalBytes, ttMonoBytes, ttSymbolBytes []byte

	var err error

	// Load font data from the files into a byte array
	if ttNormalBytes, err = ioutil.ReadFile("NotoSans-Regular.ttf"); err != nil {
		return err
	}

	if ttSymbolBytes, err = ioutil.ReadFile("NotoSansSymbols-Medium.ttf"); err != nil {
		return err
	}

	if ttMonoBytes, err = ioutil.ReadFile("CascadiaCode.ttf"); err != nil {
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
		Size:    11,
		DPI:     96 * deviceScale,
		Hinting: font.HintingNone,
	})

	a.SymbolFont = truetype.NewFace(a.ttSymbols, &truetype.Options{
		Size:    11,
		DPI:     96 * deviceScale,
		Hinting: font.HintingNone,
	})

	a.MonospaceFont = truetype.NewFace(a.ttMono, &truetype.Options{
		Size:    11,
		DPI:     96 * deviceScale,
		Hinting: font.HintingNone,
	})
}

func (a *App) drawColoredRect(target *ebiten.Image, x, y, w, h int, r, g, b, alpha uint8) {
	drawOptions := &ebiten.DrawImageOptions{}

	drawOptions.GeoM.Translate(float64(x)*(1/float64(w)), float64(y)*(1/float64(h)))
	drawOptions.GeoM.Scale(float64(w), float64(h))
	drawOptions.ColorM.Translate(float64(r)/255, float64(g)/255, float64(b)/255, float64(alpha)/255)
	_ = target.DrawImage(a.imgSquare, drawOptions)
}

func (a *App) scaleToDevice(x int) int {
	return int(a.lastDeviceScale * float64(x))
}

func (a *App) handleWindowDrag() bool {
	// Handle actively dragging a window
	if a.isDraggingWindow {
		// If we let go of the mouse or loose window focus, then we need to stop dragging
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || !ebiten.IsFocused() {
			a.isDraggingWindow = false

			ebiten.SetCursorMode(ebiten.CursorModeVisible)
		}

		deltaX := a.mouseX - a.mouseDragStartX
		deltaY := a.mouseY - a.mouseDragStartY

		if deltaX != 0 && deltaY != 0 {
			ebiten.SetWindowPosition(a.windowDragStartX+deltaX, a.windowDragStartY+deltaY)
		}

		// Eat all inputs as long as we are dragging the window
		return true
	}

	// Only handle window dragging if the mouse is in the window frame area
	if a.mouseY < 0 || a.mouseY >= a.scaleToDevice(windowFrameHeight) {
		return false
	}

	if a.leftMouseReleased && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// Handle clicking on the buttons
		buttonLeftOffset := a.screenWidth - a.scaleToDevice((windowFrameThickness*2)+115)
		if a.mouseY >= 0 && a.mouseY < a.scaleToDevice(windowFrameHeight) && a.mouseX >= buttonLeftOffset && a.mouseX < a.screenWidth {
			posIndex := (a.mouseX - buttonLeftOffset) / a.scaleToDevice(40)
			switch posIndex {
			case 0: // Minimize
				ebiten.MinimizeWindow()
			case 1: // Maximize
				// TODO: This doesn't work on mac for some reason...
				if ebiten.IsWindowMaximized() {
					ebiten.RestoreWindow()
				} else {
					ebiten.MaximizeWindow()
				}
			case 2: // Close/Quit
				os.Exit(0)
			}

			return true
		}

		// If they mouse down on the frame header and we are not dragging, it's time to start dragging
		if !a.isDraggingWindow {
			ebiten.SetCursorMode(ebiten.CursorModeCaptured)
			a.windowDragStartX, a.windowDragStartY = ebiten.WindowPosition()
			a.mouseDragStartX = a.mouseX
			a.mouseDragStartY = a.mouseY
			a.isDraggingWindow = true

			return true
		}
	}

	return false
}
