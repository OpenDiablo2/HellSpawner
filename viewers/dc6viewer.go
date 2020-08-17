package main

import (
	"flag"
    "errors"
    "log"
    "strings"
    "strconv"
    eb "github.com/hajimehoshi/ebiten"
    "github.com/OpenDiablo2/OpenDiablo2/d2common"
    "github.com/OpenDiablo2/OpenDiablo2/d2core/d2asset"
    "github.com/OpenDiablo2/OpenDiablo2/d2core/d2config"
    "github.com/OpenDiablo2/OpenDiablo2/d2core/d2input"
    "github.com/OpenDiablo2/OpenDiablo2/d2core/d2term"
    "github.com/OpenDiablo2/OpenDiablo2/d2core/d2render/ebiten"
    "github.com/OpenDiablo2/OpenDiablo2/d2common/d2resource"
    "github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dc6"
    "github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2mpq"
)

type playMode int

const (
	playModePause playMode = iota
	playModeForward
	playModeBackward
)

const defaultPalette string = d2resource.PaletteAct1
const debouncedKeyboardInputThreshold float64 = 0.15
const defaultPlayLoop bool = true
const defaultPlayMode playMode = playModeForward
const defaultTimeScale float64 = 1.0
const defaultWindowWidth uint = 800
const defaultWindowHeight uint = 600
const defaultZoomFactor float64 = 1.0
const zoomUnits float64 = 0.01
const defaultOffsetX int = 0
const defaultOffsetY int = 0
const defaultShowDebug bool = true 
var paletteList []string
var doingDrag bool = false
var lastCursorX = 0
var lastCursorY = 0

type DC6Window struct {
    renderer d2interface.Renderer
    animation d2interface.Animation
    filePath string
    dc6File *d2dc6.DC6
    terminal d2interface.Terminal
    inputManager d2interface.InputManager
    lastDebouncedKeyboardInput float64
    lastScreenAdvance float64
    lastFrameTime float64
    offsetX int
    offsetY int
    zoomScale float64
    showDebug bool
    lastTime float64
    timeScale float64
    frameEndIndex int
    playMode playMode
    playLoop bool
    paletteIndex int
    paletteFilePath string
}

func (a *DC6Window) render(target d2interface.Surface) error {
    var debugMessage strings.Builder
    cursorX, cursorY := eb.CursorPosition()
    animationWidth, animationHeight, animationOffsetX, animationOffsetY := Dimensions(a.dc6File, a.animation.GetCurrentFrame())
    debugMessage.WriteString("MPQ Filepath: ")
    debugMessage.WriteString(*flagMPQFilepath)
    debugMessage.WriteString("\nAsset Filepath: /")
    debugMessage.WriteString(a.filePath)
    debugMessage.WriteString("\nPalette Filepath: ")
    debugMessage.WriteString(a.paletteFilePath)
    debugMessage.WriteString("\nPalette Index: (")
    debugMessage.WriteString(strconv.Itoa(int(a.paletteIndex)))
    debugMessage.WriteString(" / ")
    debugMessage.WriteString(strconv.Itoa(int(len(paletteList) - 1)))
    debugMessage.WriteString(")")
    debugMessage.WriteString("\nAsset Frame Width: ")
    debugMessage.WriteString(strconv.Itoa(int(animationWidth)))
    debugMessage.WriteString("\nAsset Frame Height: ")
    debugMessage.WriteString(strconv.Itoa(int(animationHeight)))
    debugMessage.WriteString("\nAsset Frame Offset (x, y): (")
    debugMessage.WriteString(strconv.Itoa(int(animationOffsetX)))
    debugMessage.WriteString(", ")
    debugMessage.WriteString(strconv.Itoa(int(animationOffsetY)))
    debugMessage.WriteString(")\nCursor Position(x, y): (")
    debugMessage.WriteString(strconv.Itoa(cursorX))
    debugMessage.WriteString(", ")
    debugMessage.WriteString(strconv.Itoa(cursorY))
    debugMessage.WriteString(")\n")
    windowWidth, windowHeight := eb.WindowSize()
    debugMessage.WriteString("Window Dimensions(width, height): (")
    debugMessage.WriteString(strconv.Itoa(windowWidth))
    debugMessage.WriteString(", ")
    debugMessage.WriteString(strconv.Itoa(windowHeight))
    debugMessage.WriteString(")\nZoom Scale: ")
    debugMessage.WriteString(strconv.FormatFloat(a.zoomScale, 'f', -1, 64))
    debugMessage.WriteString("\nTranslation Offset X: ")
    debugMessage.WriteString(strconv.Itoa(a.offsetX))
    debugMessage.WriteString("\nTranslation Offset Y: ")
    debugMessage.WriteString(strconv.Itoa(a.offsetY))
    debugMessage.WriteString("\nAnimation Playing: ")

    if a.playMode == playModePause {
        debugMessage.WriteString("pause")
    } else if a.playMode == playModeForward {
        debugMessage.WriteString("forward")
    } else if a.playMode == playModeBackward {
        debugMessage.WriteString("backward")
    }

    debugMessage.WriteString("\nCurrent Frame Index: (")
    debugMessage.WriteString(strconv.Itoa(a.animation.GetCurrentFrame()))
    debugMessage.WriteString(" / ")
    debugMessage.WriteString(strconv.Itoa((a.animation.GetFrameCount()-1)))
    debugMessage.WriteString(")")
    debugMessage.WriteString("\n\nControls:")
    debugMessage.WriteString("\n[D] Toggle Debug Information")
    debugMessage.WriteString("\n[LEFT||RIGHT ARROWKEYS] Step Frame")
    debugMessage.WriteString("\n[CTRL + LEFT||RIGHT ARROWKEYS] Step Palette")
    debugMessage.WriteString("\n[MOUSEWHEEL] Modify Scale (+/-) 0.01")
    debugMessage.WriteString("\n[MOUSEWHEEL + CTRL] Modify Scale (+/-) 0.1")
    debugMessage.WriteString("\n[MOUSEWHEEL + SHIFT] Modify Scale (+/-) 0.1")
    debugMessage.WriteString("\n[MOUSEWHEEL + CTRL + SHIFT] Modify Scale (+/-) 1.0")
    debugMessage.WriteString("\n[MOUSEMOVE + MOUSEBUTTONLEFT] Translate (X, Y)")
    debugMessage.WriteString("\n[CTRL + 0] Center Asset")
    debugMessage.WriteString("\n[CTRL + SHIFT + 0] Reset Viewer")
    debugMessage.WriteString("\n[SPACEBAR] Play/Pause Animation")

    if a.showDebug == true {
        target.DrawTextf(debugMessage.String())
    }
    _, wheelYOff := eb.Wheel()

    if eb.IsKeyPressed(eb.KeyShift) == true && eb.IsKeyPressed(eb.KeyControl) == true {
        a.zoomScale += (100.0 * float64(zoomUnits * float64(wheelYOff)))
    } else if eb.IsKeyPressed(eb.KeyControl) == true && eb.IsKeyPressed(eb.KeyShift) == false {
        a.zoomScale += (10.0 * float64(zoomUnits * float64(wheelYOff)))
    } else if eb.IsKeyPressed(eb.KeyControl) == false && eb.IsKeyPressed(eb.KeyShift) == true {
        a.zoomScale += (10.0 * float64(zoomUnits * float64(wheelYOff)))
    } else {
        a.zoomScale += float64(zoomUnits * float64(wheelYOff))
    }

    newScale := float64(float64(int(defaultWindowWidth)) / (float64(windowWidth)))

    if eb.IsKeyPressed(eb.KeyControl) == true && eb.IsKeyPressed(eb.KeyShift) == true && eb.IsKeyPressed(eb.Key0) == true {
        a.zoomScale = 1.0
        a.offsetX = 0
        a.offsetY = 0
        eb.SetWindowSize(int(defaultWindowWidth), int(defaultWindowHeight))        
    } else if eb.IsKeyPressed(eb.KeyControl) == true && eb.IsKeyPressed(eb.Key0) == true {
        a.zoomScale = 1.0
        a.offsetX = 0
        a.offsetY = 0
        a.animation.SetCurrentFrame(0)
    } else {
        if doingDrag == true {
            a.offsetX -= (lastCursorX - cursorX)
            a.offsetY -= (lastCursorY - cursorY)
        }
    }

    translateX := int((float64((float64(windowWidth) / a.zoomScale) - float64(int(animationWidth))) / 2.0) * newScale * a.zoomScale)
    translateY := int((float64((float64(windowHeight) / a.zoomScale) - float64(int(animationHeight))) / 2.0) * newScale * a.zoomScale)
    target.PushTranslation((translateX + a.offsetX), (translateY + a.offsetY))
    target.PushScale((newScale * a.zoomScale), (newScale * a.zoomScale))
    defer target.PopN(2)
    a.animation.Render(target)

    if eb.IsMouseButtonPressed(eb.MouseButtonLeft) == true {
        if doingDrag == false {
            doingDrag = true
        }
    } else {
        if doingDrag == true {
            doingDrag = false
        }
    }

    target.PushTranslation((cursorX - (translateX + a.offsetX)), (cursorY - (translateY + a.offsetY)))
    defer target.PopN(1)

	if err := a.terminal.Render(target); err != nil {
		return err
	}

    lastCursorX = cursorX
    lastCursorY = cursorY
	return nil
}

func (a *DC6Window) advance(elapsed float64, elapsedTimeUnscaled float64, current float64) error {
	a.lastScreenAdvance = current
    timeSinceLastDebouncedKeyboardInput := (current - a.lastDebouncedKeyboardInput)
    
    if timeSinceLastDebouncedKeyboardInput >= debouncedKeyboardInputThreshold {
        if eb.IsKeyPressed(eb.KeySpace) == true {
            if a.playMode == playModePause {
                a.playMode = playModeForward
                a.animation.PlayForward()
                a.lastDebouncedKeyboardInput = current
            } else if a.playMode == playModeForward {
                a.playMode = playModePause
                a.animation.Pause()
                a.lastDebouncedKeyboardInput = current
            }
        }

        if eb.IsKeyPressed(eb.KeyD) == true {
            if a.showDebug == false {
                a.showDebug = true
                a.lastDebouncedKeyboardInput = current
            } else {
                a.showDebug = false
                a.lastDebouncedKeyboardInput = current
            }
        }

        if eb.IsKeyPressed(eb.KeyControl) == true {
            if eb.IsKeyPressed(eb.KeyLeft) == true && eb.IsKeyPressed(eb.KeyRight) == false {
                newPaletteIndex := ((a.paletteIndex - 1) % len(paletteList))

                if newPaletteIndex == -1 {
                    newPaletteIndex = (len(paletteList)-1)
                }

                a.paletteIndex = newPaletteIndex
                a.paletteFilePath = paletteList[newPaletteIndex]
                newAnimation, newAnimationErr := d2asset.LoadAnimation(a.filePath, a.paletteFilePath)

                if newAnimationErr != nil {
                    log.Fatal(newAnimationErr)
                }
                
                a.animation = newAnimation
                
                if a.playMode == playModeForward {
                    a.animation.PlayForward()
                }

                a.animation.SetPlayLoop(a.playLoop)
                a.lastDebouncedKeyboardInput = current
            } else if eb.IsKeyPressed(eb.KeyLeft) == false && eb.IsKeyPressed(eb.KeyRight) == true {
                a.paletteIndex = ((a.paletteIndex + 1) % len(paletteList))
                a.paletteFilePath = paletteList[a.paletteIndex]
                newAnimation, newAnimationErr := d2asset.LoadAnimation(a.filePath, a.paletteFilePath)

                if newAnimationErr != nil {
                    log.Fatal(newAnimationErr)
                }
                
                a.animation = newAnimation

                if a.playMode == playModeForward {
                    a.animation.PlayForward()
                }

                a.lastDebouncedKeyboardInput = current
            }
        } else {
            if a.playMode == playModePause {
                if eb.IsKeyPressed(eb.KeyLeft) == true && eb.IsKeyPressed(eb.KeyRight) == false {
                    newFrameIndex := (((a.animation.GetCurrentFrame()) - 1) % a.animation.GetFrameCount())

                    if newFrameIndex == -1 {
                        newFrameIndex = (a.animation.GetFrameCount()-1)
                    }

                    a.animation.SetCurrentFrame(newFrameIndex)
                    a.lastDebouncedKeyboardInput = current
                } else if eb.IsKeyPressed(eb.KeyLeft) == false && eb.IsKeyPressed(eb.KeyRight) == true {
                    a.animation.SetCurrentFrame((((a.animation.GetCurrentFrame()) + 1) % a.animation.GetFrameCount()))
                    a.lastDebouncedKeyboardInput = current
                }
            }
        }
    }

	if a.playMode == playModePause {
		return nil
	} else {
        if err := a.animation.Advance(elapsed); err != nil {
            return err
        }
    }    

	if err := a.terminal.Advance(elapsed); err != nil {
		return err
	}

	return nil
}

func (a *DC6Window) update(target d2interface.Surface) error {
	currentTime := d2common.Now()
	elapsedTimeUnscaled := currentTime - a.lastTime
	elapsedTime := elapsedTimeUnscaled * a.timeScale
	a.lastTime = currentTime

	if err := a.advance(elapsedTime, elapsedTimeUnscaled, currentTime); err != nil {
		return err
	}

	if err := a.render(target); err != nil {
		return err
	}

	if target.GetDepth() > 0 {
		return errors.New("detected surface stack leak")
	}

	return nil
}

func loadDC6(dc6Path string) (*d2dc6.DC6, error) {
	dc6Data, err := d2asset.LoadFile(dc6Path)
	if err != nil {
		return nil, err
	}

	dc6, err := d2dc6.Load(dc6Data)

	if err != nil {
		return nil, err
	}

	return dc6, nil
}

func Dimensions(inputDC6 *d2dc6.DC6, inputFrameIndex int) (uint32, uint32, int32, int32) {
    var resultWidth uint32 = 0
    var resultHeight uint32 = 0
    var resultOffsetX int32 = int32(defaultWindowWidth)
    var resultOffsetY int32 = int32(defaultWindowHeight)

    if len(inputDC6.Frames) > 0 {
        resultWidth = inputDC6.Frames[inputFrameIndex].Width
        resultHeight = inputDC6.Frames[inputFrameIndex].Height
    }

    return resultWidth, resultHeight, resultOffsetX, resultOffsetY
}

func initalizePaletteList() {
    paletteList = append(paletteList, d2resource.PaletteAct1)
    paletteList = append(paletteList, d2resource.PaletteAct2)
    paletteList = append(paletteList, d2resource.PaletteAct3)
    paletteList = append(paletteList, d2resource.PaletteAct4)
    paletteList = append(paletteList, d2resource.PaletteAct5)
    paletteList = append(paletteList, d2resource.PaletteEndGame)
    paletteList = append(paletteList, d2resource.PaletteEndGame2)
    paletteList = append(paletteList, d2resource.PaletteFechar)
    paletteList = append(paletteList, d2resource.PaletteLoading)
    paletteList = append(paletteList, d2resource.PaletteMenu0)
    paletteList = append(paletteList, d2resource.PaletteMenu1)
    paletteList = append(paletteList, d2resource.PaletteMenu2)
    paletteList = append(paletteList, d2resource.PaletteMenu3)
    paletteList = append(paletteList, d2resource.PaletteMenu4)
    paletteList = append(paletteList, d2resource.PaletteSky)
    paletteList = append(paletteList, d2resource.PaletteStatic)
    paletteList = append(paletteList, d2resource.PaletteTrademark)
    paletteList = append(paletteList, d2resource.PaletteUnits)
}

func getPaletteIndex(inputPaletteFilePath string) (int) {
    result := -1

    for i := 0; i < len(paletteList); i++ {
        if inputPaletteFilePath == paletteList[i] {
            result = i
            break
        }
    }

    return result
}

func Create(fileName string, paletteFilePath string) (*DC6Window, error) {
    initalizePaletteList()
    result := &DC6Window{}
	configLoadErr := d2config.Load()
 
    if configLoadErr != nil {
		return result, configLoadErr
	}

    inputManager := d2input.NewInputManager()
	term, termErr := d2term.New(inputManager)

	if termErr != nil {
		return result, termErr
	}

	renderer, rendererErr := ebiten.CreateRenderer()

	if rendererErr != nil {
		return result, rendererErr
	}

    result.renderer = renderer
    result.terminal = term
    result.inputManager = inputManager
    d2AssetErr := d2asset.Initialize(renderer, term)

    if d2AssetErr != nil {
        return result, d2AssetErr
    }

    result.filePath = fileName
    dc6File, dc6FileErr := loadDC6(result.filePath)
        
    if dc6FileErr != nil {
        return result, dc6FileErr
    }

    result.dc6File = dc6File
    result.playLoop = defaultPlayLoop
    result.playMode = defaultPlayMode
    result.timeScale = defaultTimeScale
    result.lastTime = d2common.Now()
    result.lastDebouncedKeyboardInput = d2common.Now()
    result.showDebug = defaultShowDebug
    result.offsetX = defaultOffsetX
    result.offsetY = defaultOffsetY
    result.zoomScale = defaultZoomFactor
    animation, animationErr := d2asset.LoadAnimation(result.filePath, paletteFilePath)

    if animationErr != nil {
        return result, animationErr
    }

    result.paletteFilePath = paletteFilePath
    result.paletteIndex = getPaletteIndex(paletteFilePath)
    result.animation = animation
    result.animation.PlayForward()
    result.animation.SetPlayLoop(result.playLoop)
    result.frameEndIndex = (result.animation.GetFrameCount()-1)
    eb.SetCursorVisible(true)
    runErr := renderer.Run(result.update, int(defaultWindowWidth), int(defaultWindowHeight), "HellSpawner DC6 Viewer")
    
    if runErr != nil {
		return result, runErr
	}

	return result, nil
}

func OpenDC6FileWindow(mpqPath string, filePath string, paletteFilePath string) (error) {
    _, windowErr := Create(filePath, paletteFilePath)

    if windowErr != nil {
        log.Fatal(windowErr)
    }

    return nil
}

var flagMPQFilepath = flag.String("mpq", "", "mpq filepath")
var flagAssetFilepath = flag.String("asset", "", "asset filepath")
var flagPaletteFilepath = flag.String("palette", "", "palette filepath")

func main() {
    flag.Parse()
    initialPalette := *flagPaletteFilepath
    
    if *flagMPQFilepath == "" {
        log.Fatal(errors.New("main :: mpq flag must be set to the mpq filepath."))
    }

    if *flagAssetFilepath == "" {
        log.Fatal(errors.New("main :: asset flag must be set to the asset filepath."))
    }

    if *flagPaletteFilepath == "" {
        initialPalette = defaultPalette
    }

    _, mpqErr := d2mpq.Load(*flagMPQFilepath)

	if mpqErr != nil {
        log.Fatal(mpqErr)
    }
    
    OpenDC6FileWindow(*flagMPQFilepath, *flagAssetFilepath, initialPalette)
}
