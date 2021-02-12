package hswidget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
)

type cofEditorState int

const (
	cofEditorStateViewer cofEditorState = iota
	cofEditorStateAddLayer
	cofEditorStateConfirm
)

const (
	indicatorSize = 64
)

const (
	upDownArrowW, upDownArrowH           = 15, 15
	leftRightArrowW, leftRightArrowH     = 15, 15
	actionButtonW, actionButtonH         = 200, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	trueFalseListW                       = 60
	speedInputW                          = 40
)

const (
	maxSpeed = 100
)

// COFViewerState represents cof viewer's state
type COFViewerState struct {
	layerIndex     int32
	directionIndex int32
	frameIndex     int32
	layer          *d2cof.CofLayer
	confirmDialog  *PopUpConfirmDialog
}

// Dispose clears viewer's layers
func (s *COFViewerState) Dispose() {
	s.layer = nil
}

// COFEditorState represents state of cof editor
type COFEditorState struct {
	newLayerType        int32
	newLayerSelectable  int32
	newLayerTransparent int32
	newLayerDrawEffect  int32
	newLayerWeaponClass int32
}

// Dispose disposes editor's state
func (s *COFEditorState) Dispose() {
	// noop
}

// COFState represents cof editor's and viewer's state
type COFState struct {
	*COFViewerState
	*COFEditorState
	state cofEditorState
}

// Dispose clear widget's state
func (s *COFState) Dispose() {
	s.COFViewerState.Dispose()
	s.COFEditorState.Dispose()
}

// COFWidget represents cof viewer's widget
type COFWidget struct {
	id                string
	cof               *d2cof.COF
	upArrowTexture    *giu.Texture
	downArrowTexture  *giu.Texture
	leftArrowTexture  *giu.Texture
	rightArrowTexture *giu.Texture
}

// COFViewer creates a cof viewer widget
func COFViewer(textureLoader *hscommon.TextureLoader,
	upArrowTexture, downArrowTexture, rightArrowTexture, leftArrowTexture *giu.Texture,
	id string, cof *d2cof.COF) *COFWidget {
	result := &COFWidget{
		id:                id,
		cof:               cof,
		upArrowTexture:    upArrowTexture,
		downArrowTexture:  downArrowTexture,
		rightArrowTexture: rightArrowTexture,
		leftArrowTexture:  leftArrowTexture,
	}

	return result
}

// Build builds a cof viewer
func (p *COFWidget) Build() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	if s == nil {
		giu.Context.SetState(stateID, &COFState{
			state: cofEditorStateViewer,
			COFViewerState: &COFViewerState{
				layer:         &p.cof.CofLayers[0],
				confirmDialog: &PopUpConfirmDialog{},
			},
			COFEditorState: &COFEditorState{},
		})

		return
	}

	state := s.(*COFState)

	switch state.state {
	case cofEditorStateViewer:
		p.buildViewer()
	case cofEditorStateAddLayer:
		p.makeAddLayerLayout().Build()
	case cofEditorStateConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

// nolint:funlen // no need to reduce
func (p *COFWidget) buildViewer() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)
	state := s.(*COFState)

	var l1, l2, l3, l4 string

	numDirs := p.cof.NumberOfDirections
	numFrames := p.cof.FramesPerDirection

	l1 = fmt.Sprintf("Directions: %v", numDirs)

	if numDirs > 1 {
		l2 = fmt.Sprintf("Frames (x%v):", numDirs)
	} else {
		l2 = "Frames:"
	}

	// nolint:gomnd // constant
	fps := 25 * (float64(p.cof.Speed) / float64(256))
	if fps == 0 {
		fps = 25
	}

	l3 = fmt.Sprintf("FPS: %.1f", fps)
	// nolint:gomnd // miliseconds
	l4 = fmt.Sprintf("Duration: %.2fms", float64(numFrames)*(1/fps)*1000)

	layerStrings := make([]string, 0)
	for idx := range p.cof.CofLayers {
		layerStrings = append(layerStrings, strconv.Itoa(int(p.cof.CofLayers[idx].Type)))
	}

	layerList := giu.Combo("##"+p.id+"layer", layerStrings[state.COFViewerState.layerIndex], layerStrings, &state.layerIndex).
		Size(indicatorSize).OnChange(p.onUpdate)

	directionStrings := make([]string, 0)
	for idx := range p.cof.Priority {
		directionStrings = append(directionStrings, fmt.Sprintf("%d", idx))
	}

	directionList := giu.Combo("##"+p.id+"dir", directionStrings[state.COFViewerState.directionIndex],
		directionStrings, &state.directionIndex).
		Size(indicatorSize).OnChange(p.onUpdate)

	frameStrings := make([]string, 0)
	for idx := range p.cof.Priority[state.COFViewerState.directionIndex] {
		frameStrings = append(frameStrings, fmt.Sprintf("%d", idx))
	}

	frameList := giu.Combo("##"+p.id+"frame", frameStrings[state.COFViewerState.frameIndex], frameStrings, &state.frameIndex).
		Size(indicatorSize).OnChange(p.onUpdate)

	const vspace = 4 //nolint:unused // will be used

	speed := int32(p.cof.Speed)
	giu.TabBar("COFViewerTabs").Layout(giu.Layout{
		giu.TabItem("Animation").Layout(giu.Layout{
			giu.Label(l1),
			giu.Line(
				giu.Label(l2),
				giu.ImageButton(p.leftArrowTexture).Size(leftRightArrowW, leftRightArrowH).OnClick(func() {
					if p.cof.FramesPerDirection > 0 {
						p.cof.FramesPerDirection--
					}
				}),
				giu.Label(strconv.Itoa(numFrames)),
				giu.Custom(func() {
					imgui.PopID()
					imgui.PushID("##" + p.id + "IncreaseFramesPerDirection")
				}),
				giu.ImageButton(p.rightArrowTexture).Size(leftRightArrowW, leftRightArrowH).OnClick(func() {
					p.cof.FramesPerDirection++
				}),
				giu.Custom(func() {
					imgui.PopID()
					imgui.PushID("##" + p.id + "DecreaseFramesPerDirection")
				}),
			),
			giu.Line(
				giu.Label("Speed: "),
				giu.InputInt("##"+p.id+"CovViewerSpeedValue", &speed).Size(speedInputW).OnChange(func() {
					if speed <= maxSpeed {
						p.cof.Speed = int(speed)
					} else {
						p.cof.Speed = maxSpeed
					}
				}),
			),
			giu.Label(l3),
			giu.Label(l4),
		}),
		giu.TabItem("Layer").Layout(giu.Layout{
			giu.Layout{
				giu.Line(giu.Label("Selected Layer: "), layerList),
				giu.Separator(),
				p.makeLayerLayout(),
				giu.Separator(),
				giu.Button("Add a new layer...##"+p.id+"AddLayer").Size(actionButtonW, actionButtonH).OnClick(func() {
					p.CreateNewLayer()
				}),
				giu.Button("Delete current layer...##"+p.id+"DeleteLayer").Size(actionButtonW, actionButtonH).OnClick(func() {
					state.COFViewerState.confirmDialog = NewPopUpConfirmDialog(
						"##"+p.id+"DeleteLayerConfirm",
						"Do you raly want to remove this layer?",
						"If you'll click YES, all data from this layer will be lost. Continue?",
						func() {
							p.deleteCurrentLayer(state.COFViewerState.layerIndex)
							state.state = cofEditorStateViewer
						},
						func() {
							state.state = cofEditorStateViewer
						},
					)

					state.state = cofEditorStateConfirm
				}),
			},
		}),
		giu.TabItem("Priority").Layout(giu.Layout{
			giu.Line(
				giu.Label("Direction: "), directionList,
				giu.Label("Frame: "), frameList,
			),
			giu.Separator(),
			p.makeDirectionLayout(),
			giu.Button("Duplicate current direction...##"+p.id+"DuplicateDirection").Size(actionButtonW, actionButtonH).OnClick(func() {
				p.duplicateDirection()
			}),
			giu.Button("Delete current direction...##"+p.id+"DeleteDirection").Size(actionButtonW, actionButtonH).OnClick(func() {
				NewPopUpConfirmDialog("##"+p.id+"DeleteLayerConfirm",
					"Do you raly want to remove this direction?",
					"If you'll click YES, all data from this direction will be lost. Continue?",
					func() {
						p.deleteCurrentDirection()
						state.state = cofEditorStateViewer
					},
					func() {
						state.state = cofEditorStateViewer
					},
				)

				state.state = cofEditorStateConfirm
			}),
		}),
	}).Build()
}

func (p *COFWidget) onUpdate() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFState)

	clone := p.cof.CofLayers[state.COFViewerState.layerIndex]
	state.COFViewerState.layer = &clone

	giu.Context.SetState(p.id, state)
}

func (p *COFWidget) makeLayerLayout() giu.Layout {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFState)

	if state.COFViewerState.layer == nil {
		p.onUpdate()
	}

	layerName := hsenum.GetLayerName(state.COFViewerState.layer.Type)

	strType := fmt.Sprintf("Type: %s (%s)", state.COFViewerState.layer.Type, layerName)
	strShadow := fmt.Sprintf("Shadow: %t", state.COFViewerState.layer.Shadow > 0)
	strSelectable := fmt.Sprintf("Selectable: %t", state.COFViewerState.layer.Selectable)
	strTransparent := fmt.Sprintf("Transparent: %t", state.COFViewerState.layer.Transparent)

	effect := hsenum.GetDrawEffectName(state.COFViewerState.layer.DrawEffect)

	strEffect := fmt.Sprintf("Draw Effect: %s", effect)

	weapon := hsenum.GetWeaponClassString(state.COFViewerState.layer.WeaponClass)

	strWeaponClass := fmt.Sprintf("Weapon Class: (%s) %s", state.COFViewerState.layer.WeaponClass, weapon)

	return giu.Layout{
		giu.Label(strType),
		giu.Label(strShadow),
		giu.Label(strSelectable),
		giu.Label(strTransparent),
		giu.Label(strEffect),
		giu.Label(strWeaponClass),
	}
}

func (p *COFWidget) makeDirectionLayout() giu.Layout {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFState).COFViewerState

	frames := p.cof.Priority[state.directionIndex]
	layers := frames[int(state.frameIndex)%len(frames)]

	return giu.Layout{
		giu.Label("Render Order (first to last):"),
		giu.Custom(func() {
			for idx := range layers {
				currentIdx := idx
				giu.Line(
					giu.ImageButton(p.upArrowTexture).Size(upDownArrowW, upDownArrowH).OnClick(func() {
						if currentIdx > 0 {
							p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx-1],
								p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx] =
								p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx],
								p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx-1]
						}
					}),
					giu.Custom(func() {
						imgui.PopID()
						imgui.PushID(fmt.Sprintf("LayerPriorityUp_%d", currentIdx))
					}),
					giu.ImageButton(p.downArrowTexture).Size(upDownArrowW, upDownArrowH).OnClick(func() {
						if currentIdx < len(layers)-1 {
							p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx],
								p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx+1] =
								p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx+1],
								p.cof.Priority[state.directionIndex][state.frameIndex][currentIdx]
						}
					}),
					giu.Custom(func() {
						imgui.PopID()
						imgui.PushID(fmt.Sprintf("LayerPriorityDown_%d", currentIdx))
					}),
					giu.Label(fmt.Sprintf("%d: %s", idx, hsenum.GetLayerName(layers[idx]))),
				).Build()
			}
		}),
	}
}

// nolint:funlen // can't reduce
func (p *COFWidget) makeAddLayerLayout() giu.Layout {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	trueFalse := []string{"false", "true"}

	compositeTypeList := make([]string, 0)
	for i := d2enum.CompositeTypeHead; i < d2enum.CompositeTypeMax; i++ {
		compositeTypeList = append(compositeTypeList, i.String()+" ("+hsenum.GetLayerName(i)+")")
	}

	drawEffectList := make([]string, d2enum.DrawEffectNone+1)
	for i := d2enum.DrawEffectPctTransparency25; i <= d2enum.DrawEffectNone; i++ {
		drawEffectList[int(i)] = strconv.Itoa(int(i)) + " (" + hsenum.GetDrawEffectName(i) + ")"
	}

	weaponClassList := make([]string, d2enum.WeaponClassTwoHandToHand+1)
	for i := d2enum.WeaponClassNone; i <= d2enum.WeaponClassTwoHandToHand; i++ {
		weaponClassList[int(i)] = i.String() + " (" + hsenum.GetWeaponClassString(i) + ")"
	}
	return giu.Layout{
		giu.Label("Select new COF's Layer parameters:"),
		giu.Separator(),
		giu.Line(
			giu.Label("Type: "),
			giu.Combo("##"+p.id+"AddLayerType", compositeTypeList[state.COFEditorState.newLayerType],
				compositeTypeList, &state.COFEditorState.newLayerType).Size(bigListW),
		),
		giu.Line(
			giu.Label("Selectable: "),
			giu.Combo("##"+p.id+"AddLayerSelectable", trueFalse[state.COFEditorState.newLayerSelectable],
				trueFalse, &state.COFEditorState.newLayerSelectable).Size(trueFalseListW),
		),
		giu.Line(
			giu.Label("Transparent: "),
			giu.Combo("##"+p.id+"AddLayerTransparent", trueFalse[state.COFEditorState.newLayerTransparent],
				trueFalse, &state.COFEditorState.newLayerTransparent).Size(trueFalseListW),
		),
		giu.Line(
			giu.Label("Draw effect: "),
			giu.Combo("##"+p.id+"AddLayerDrawEffect", drawEffectList[state.COFEditorState.newLayerDrawEffect],
				drawEffectList, &state.COFEditorState.newLayerDrawEffect).Size(bigListW),
		),
		giu.Line(
			giu.Label("Weapon class: "),
			giu.Combo("##"+p.id+"AddLayerWeaponClass", weaponClassList[state.COFEditorState.newLayerWeaponClass],
				weaponClassList, &state.COFEditorState.newLayerWeaponClass).Size(bigListW),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				newCofLayer := &d2cof.CofLayer{
					Type:        d2enum.CompositeType(state.COFEditorState.newLayerType),
					Selectable:  (state.COFEditorState.newLayerSelectable == 1),
					Transparent: (state.COFEditorState.newLayerTransparent == 1),
					DrawEffect:  d2enum.DrawEffect(state.COFEditorState.newLayerDrawEffect),
					WeaponClass: d2enum.WeaponClass(state.COFEditorState.newLayerWeaponClass),
				}

				p.cof.CofLayers = append(p.cof.CofLayers, *newCofLayer)

				p.cof.NumberOfLayers++

				for i := range p.cof.Priority {
					for j := range p.cof.Priority[i] {
						p.cof.Priority[i][j] = append(p.cof.Priority[i][j], newCofLayer.Type)
					}
				}

				state.state = cofEditorStateViewer
			}),
			giu.Button("Cancel##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
				state.state = cofEditorStateViewer
			}),
		),
	}
}

func (p *COFWidget) deleteCurrentLayer(index int32) {
	p.cof.NumberOfLayers--

	newLayers := make([]d2cof.CofLayer, 0)

	for n, i := range p.cof.CofLayers {
		if int32(n) != index {
			newLayers = append(newLayers, i)
		}
	}

	p.cof.CofLayers = newLayers
}

func (p *COFWidget) duplicateDirection() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	idx := state.COFViewerState.directionIndex

	p.cof.NumberOfDirections++

	p.cof.Priority = append(p.cof.Priority, p.cof.Priority[idx])

	// nolint:gomnd // directionIndex starts from 0, but len from 1
	state.directionIndex = int32(len(p.cof.Priority) - 1)
}

func (p *COFWidget) deleteCurrentDirection() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)
	index := state.COFViewerState.directionIndex

	p.cof.NumberOfDirections--

	newPriority := make([][][]d2enum.CompositeType, 0)

	for n, i := range p.cof.Priority {
		if int32(n) != index {
			newPriority = append(newPriority, i)
		}
	}

	p.cof.Priority = newPriority
}

// CreateNewLayer starts add-cof-layer dialog
func (p *COFWidget) CreateNewLayer() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	state.state = cofEditorStateAddLayer
}
