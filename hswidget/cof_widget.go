package hswidget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"
	"github.com/ianling/imgui-go"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
)

const (
	indicatorSize                        = 64
	upDownArrowW, upDownArrowH           = 15, 15
	leftRightArrowW, leftRightArrowH     = 15, 15
	actionButtonW, actionButtonH         = 200, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	trueFalseListW                       = 60
	speedInputW                          = 40
)

type cofTextures struct {
	up    *giu.Texture
	down  *giu.Texture
	left  *giu.Texture
	right *giu.Texture
}

// COFWidget represents cof viewer's widget
type COFWidget struct {
	id       string
	cof      *d2cof.COF
	textures cofTextures
}

// COFViewer creates a cof viewer widget
func COFViewer(
	up, down, right, left *giu.Texture,
	id string, cof *d2cof.COF,
) *COFWidget {
	result := &COFWidget{
		id:  id,
		cof: cof,
	}

	result.textures.up = up
	result.textures.down = down
	result.textures.left = left
	result.textures.right = right

	return result
}

// Build builds a cof viewer
func (p *COFWidget) Build() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	if s == nil {
		p.setDefaultState(stateID)
		return
	}

	state := s.(*COFState)

	// builds appropriate menu (depends on state)
	switch state.mode {
	case cofEditorModeViewer:
		p.makeViewerLayout().Build()
	case cofEditorModeAddLayer:
		p.makeAddLayerLayout().Build()
	case cofEditorModeConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

func (p *COFWidget) setDefaultState(id string) {
	defaultState := &COFState{
		mode: cofEditorModeViewer,
		viewerState: &viewerState{
			layer:         &p.cof.CofLayers[0],
			confirmDialog: &PopUpConfirmDialog{},
		},
		newLayerFields: &newLayerFields{
			selectable: 1,
			drawEffect: int32(d2enum.DrawEffectNone),
		},
	}

	giu.Context.SetState(id, defaultState)
}

// this likely needs to be a method of d2cof.COF
func speedToFPS(speed int) float64 {
	const (
		baseFPS      = 25
		speedDivisor = 256
	)

	fps := baseFPS * (float64(speed) / speedDivisor)
	if fps == 0 {
		fps = baseFPS
	}

	return fps
}

func (p *COFWidget) makeViewerLayout() giu.Layout {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)
	state := s.(*COFState)

	layerStrings := make([]string, 0)
	for idx := range p.cof.CofLayers {
		layerStrings = append(layerStrings, strconv.Itoa(int(p.cof.CofLayers[idx].Type)))
	}

	currentLayerName := layerStrings[state.viewerState.layerIndex]
	layerList := giu.Combo("##"+p.id+"layer", currentLayerName, layerStrings, &state.layerIndex)
	layerList.Size(indicatorSize).OnChange(p.onUpdate)

	directionStrings := make([]string, 0)
	for idx := range p.cof.Priority {
		directionStrings = append(directionStrings, fmt.Sprintf("%d", idx))
	}

	directionString := directionStrings[state.viewerState.directionIndex]
	directionList := giu.Combo("##"+p.id+"dir", directionString, directionStrings, &state.directionIndex)
	directionList.Size(indicatorSize).OnChange(p.onUpdate)

	frameStrings := make([]string, 0)
	for idx := range p.cof.Priority[state.viewerState.directionIndex] {
		frameStrings = append(frameStrings, fmt.Sprintf("%d", idx))
	}

	frameString := frameStrings[state.viewerState.frameIndex]
	frameList := giu.Combo("##"+p.id+"frame", frameString, frameStrings, &state.frameIndex)
	frameList.Size(indicatorSize).OnChange(p.onUpdate)

	return giu.Layout{
		giu.TabBar("COFViewerTabs").Layout(giu.Layout{
			giu.TabItem("Animation").Layout(p.makeAnimationTab()),
			giu.TabItem("Layer").Layout(p.makeLayerTab(state, layerList)),
			giu.TabItem("Priority").Layout(p.makePriorityTab(state, directionList, frameList)),
		}),
	}
}

// this should also probably be a method of COF
func calculateDuration(cof *d2cof.COF) float64 {
	const (
		milliseconds = 1000
	)

	frameDelay := milliseconds / speedToFPS(cof.Speed)

	return float64(cof.FramesPerDirection) * frameDelay
}

func (p *COFWidget) makeAnimationTab() giu.Layout {
	const (
		fmtFPS        = "FPS: %.1f"
		fmtDuration   = "Duration: %.2fms"
		fmtDirections = "Directions: %v"
		strSpeed      = "Speed: "
		maxSpeed      = 256
	)

	numDirs := p.cof.NumberOfDirections
	fps := speedToFPS(p.cof.Speed)
	duration := calculateDuration(p.cof)

	strLabelDirections := fmt.Sprintf(fmtDirections, numDirs)
	strLabelFPS := fmt.Sprintf(fmtFPS, fps)
	strLabelDuration := fmt.Sprintf(fmtDuration, duration)

	speed32 := int32(p.cof.Speed)
	setSpeed := func() {
		p.cof.Speed = int(speed32)

		if speed32 >= maxSpeed {
			p.cof.Speed = maxSpeed
		}
	}

	speedLabel := giu.Label(strSpeed)
	speedInput := giu.InputInt("##"+p.id+"CovViewerSpeedValue", &speed32).Size(speedInputW).OnChange(setSpeed)

	return giu.Layout{
		giu.Label(strLabelDirections),
		p.layoutAnimFrames(),
		giu.Line(speedLabel, speedInput),
		giu.Label(strLabelFPS),
		giu.Label(strLabelDuration),
	}
}

func (p *COFWidget) makeLayerTab(state *COFState, layerList giu.Widget) giu.Layout {
	addLayerButtonID := fmt.Sprintf("Add a new layer...##%sAddLayer", p.id)
	addLayerButton := giu.Button(addLayerButtonID).Size(actionButtonW, actionButtonH)
	addLayerButton.OnClick(func() {
		p.CreateNewLayer()
	})

	deleteLayerButtonID := fmt.Sprintf("Delete current layer...##%sDeleteLayer", p.id)
	deleteLayerButton := giu.Button(deleteLayerButtonID).Size(actionButtonW, actionButtonH)
	deleteLayerButton.OnClick(func() {
		const (
			strPrompt  = "Do you really want to remove this layer?"
			strMessage = "If you'll click YES, all data from this layer will be lost. Continue?"
		)

		fnYes := func() {
			p.deleteCurrentLayer(state.viewerState.layerIndex)
			state.mode = cofEditorModeViewer
		}

		fnNo := func() {
			state.mode = cofEditorModeViewer
		}

		id := fmt.Sprintf("##%sDeleteLayerConfirm", p.id)
		state.viewerState.confirmDialog = NewPopUpConfirmDialog(id, strPrompt, strMessage, fnYes, fnNo)

		state.mode = cofEditorModeConfirm
	})

	layout := giu.Layout{
		giu.Line(giu.Label("Selected Layer: "), layerList),
		giu.Separator(),
		p.makeLayerLayout(),
		giu.Separator(),
		addLayerButton,
		deleteLayerButton,
	}

	return layout
}

func (p *COFWidget) makePriorityTab(state *COFState, directionList, frameList giu.Widget) giu.Layout {
	const (
		strPrompt  = "Do you really want to remove this direction?"
		strMessage = "If you'll click YES, all data from this direction will be lost. Continue?"
	)

	duplicateButtonID := fmt.Sprintf("Duplicate current direction...##%sDuplicateDirection", p.id)
	duplicateButton := giu.Button(duplicateButtonID).Size(actionButtonW, actionButtonH)
	duplicateButton.OnClick(func() {
		p.duplicateDirection()
	})

	deleteButtonID := fmt.Sprintf("Delete current direction...##%sDeleteDirection", p.id)
	deleteButton := giu.Button(deleteButtonID).Size(actionButtonW, actionButtonH)
	deleteButton.OnClick(func() {
		fnYes := func() {
			p.deleteCurrentDirection()
			state.mode = cofEditorModeViewer
		}

		fnNo := func() {
			state.mode = cofEditorModeViewer
		}

		popupID := fmt.Sprintf("##%sDeleteLayerConfirm", p.id)

		NewPopUpConfirmDialog(popupID, strPrompt, strMessage, fnYes, fnNo)
		state.mode = cofEditorModeConfirm
	})

	return giu.Layout{
		giu.Line(
			giu.Label("Direction: "), directionList,
			giu.Label("Frame: "), frameList,
		),
		giu.Separator(),
		p.makeDirectionLayout(),
		duplicateButton,
		deleteButton,
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

//nolint:unparam // width and height are always 15 at the time of writing, but may change
func makeImageButton(id string, w, h int, t *giu.Texture, fn func()) giu.Layout {
	return giu.Layout{
		giu.ImageButton(t).Size(float32(w), float32(h)).OnClick(fn),
		giu.Custom(func() {
			// make this button unique across all editor instances
			// at the time of writing, ImageButton uses the texture ID as the button ID
			// so it wont be unique across multiple instances if we use the same texture...
			// we need to step over giu and manually tell imgui to pop the last ID and
			// push the desired one onto the stack
			imgui.PopID()
			imgui.PushID(id)
		}),
	}
}

// the layout ends up looking like this:
// Frames (x6):  <- 10 ->
// you use the arrows to set the number of frames per direction
func (p *COFWidget) layoutAnimFrames() *giu.LineWidget {
	numFrames := p.cof.FramesPerDirection
	numDirs := p.cof.NumberOfDirections

	strLabel := "Frames:"
	if numDirs > 1 {
		strLabel = fmt.Sprintf("Frames (x%v):", numDirs)
	}

	fnDecrease := func() {
		p.cof.FramesPerDirection = max(p.cof.FramesPerDirection-1, 0)
	}

	fnIncrease := func() {
		p.cof.FramesPerDirection++
	}

	label := giu.Label(strLabel)

	leftButtonID := fmt.Sprintf("##%sDecreaseFramesPerDirection", p.id)
	rightButtonID := fmt.Sprintf("##%sIncreaseFramesPerDirection", p.id)

	left := makeImageButton(leftButtonID, leftRightArrowW, leftRightArrowH, p.textures.left, fnDecrease)
	frameCount := giu.Label(fmt.Sprintf("%d", numFrames))
	right := makeImageButton(rightButtonID, leftRightArrowW, leftRightArrowH, p.textures.right, fnIncrease)

	return giu.Line(label, left, frameCount, right)
}

func (p *COFWidget) onUpdate() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFState)

	clone := p.cof.CofLayers[state.viewerState.layerIndex]
	state.viewerState.layer = &clone

	giu.Context.SetState(p.id, state)
}

func (p *COFWidget) makeLayerLayout() giu.Layout {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFState)

	if state.viewerState.layer == nil {
		p.onUpdate()
	}

	layerName := hsenum.GetLayerName(state.viewerState.layer.Type)

	strType := fmt.Sprintf("Type: %s (%s)", state.viewerState.layer.Type, layerName)
	strShadow := fmt.Sprintf("Shadow: %t", state.viewerState.layer.Shadow > 0)
	strSelectable := fmt.Sprintf("Selectable: %t", state.viewerState.layer.Selectable)
	strTransparent := fmt.Sprintf("Transparent: %t", state.viewerState.layer.Transparent)

	effect := hsenum.GetDrawEffectName(state.viewerState.layer.DrawEffect)

	strEffect := fmt.Sprintf("Draw Effect: %s", effect)

	weapon := hsenum.GetWeaponClassString(state.viewerState.layer.WeaponClass)

	strWeaponClass := fmt.Sprintf("Weapon Class: (%s) %s", state.viewerState.layer.WeaponClass, weapon)

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
	const (
		strRenderOrderLabel = "Render Order (first to last):"
		fmtIncreasePriority = "LayerPriorityUp_%d"
		fmtDecreasePriority = "LayerPriorityDown_%d"
		fmtLayerLabel       = "%d: %s"
	)

	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFState).viewerState

	frames := p.cof.Priority[state.directionIndex]
	layers := frames[int(state.frameIndex)%len(frames)]

	// increase / decrease callback function providers, based on layer index
	makeIncPriorityFn := func(idx int) func() {
		return func() {
			if idx <= 0 {
				return
			}

			list := &p.cof.Priority[state.directionIndex][state.frameIndex]
			(*list)[idx-1], (*list)[idx] = (*list)[idx], (*list)[idx-1]
		}
	}

	makeDecPriorityFn := func(idx int) func() {
		return func() {
			list := &p.cof.Priority[state.directionIndex][state.frameIndex]

			if idx >= len(*list)-1 {
				return
			}

			(*list)[idx], (*list)[idx+1] = (*list)[idx+1], (*list)[idx]
		}
	}

	// each layer line looks like:
	// <- -> 0: Name
	// the left/right buttons use the callbacks created by the previous funcs for index=0
	buildLayerPriorityLine := func(idx int) {
		currentIdx := idx

		strIncPri := fmt.Sprintf(fmtIncreasePriority, currentIdx)
		strDecPri := fmt.Sprintf(fmtDecreasePriority, currentIdx)

		fnIncPriority := makeIncPriorityFn(currentIdx)
		fnDecPriority := makeDecPriorityFn(currentIdx)

		increasePriority := makeImageButton(strIncPri, upDownArrowW, upDownArrowH, p.textures.up, fnIncPriority)
		decreasePriority := makeImageButton(strDecPri, upDownArrowW, upDownArrowH, p.textures.down, fnDecPriority)

		strLayerName := hsenum.GetLayerName(layers[idx])
		strLayerLabel := fmt.Sprintf(fmtLayerLabel, idx, strLayerName)

		layerNameLabel := giu.Label(strLayerLabel)

		giu.Line(increasePriority, decreasePriority, layerNameLabel).Build()
	}

	// finally, a func that we can pass to giu.Custom
	buildLayerLines := func() {
		for idx := range layers {
			buildLayerPriorityLine(idx)
		}
	}

	return giu.Layout{
		giu.Label(strRenderOrderLabel),
		giu.Custom(buildLayerLines),
	}
}

func (p *COFWidget) makeAddLayerLayout() giu.Layout {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	trueFalse := []string{"false", "true"}

	// available is a list of available (not currently used) composite types
	available := make([]d2enum.CompositeType, 0)

	for i := d2enum.CompositeTypeHead; i < d2enum.CompositeTypeMax; i++ {
		contains := false

		for _, j := range p.cof.CofLayers {
			if i == j.Type {
				contains = true

				break
			}
		}

		if !contains {
			available = append(available, i)
		}
	}

	compositeTypeList := make([]string, len(available))
	for n, i := range available {
		compositeTypeList[n] = i.String() + " (" + hsenum.GetLayerName(i) + ")"
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
			giu.Combo("##"+p.id+"AddLayerType", compositeTypeList[state.newLayerFields.layerType],
				compositeTypeList, &state.newLayerFields.layerType).Size(bigListW),
		),
		giu.Line(
			giu.Label("Shadow: "),
			giu.Combo("##"+p.id+"AddLayerShadow", trueFalse[state.newLayerFields.shadow],
				trueFalse, &state.newLayerFields.shadow).Size(trueFalseListW),
		),
		giu.Line(
			giu.Label("Selectable: "),
			giu.Combo("##"+p.id+"AddLayerSelectable", trueFalse[state.newLayerFields.selectable],
				trueFalse, &state.newLayerFields.selectable).Size(trueFalseListW),
		),
		giu.Line(
			giu.Label("Transparent: "),
			giu.Combo("##"+p.id+"AddLayerTransparent", trueFalse[state.newLayerFields.transparent],
				trueFalse, &state.newLayerFields.transparent).Size(trueFalseListW),
		),
		giu.Line(
			giu.Label("Draw effect: "),
			giu.Combo("##"+p.id+"AddLayerDrawEffect", drawEffectList[state.newLayerFields.drawEffect],
				drawEffectList, &state.newLayerFields.drawEffect).Size(bigListW),
		),
		giu.Line(
			giu.Label("Weapon class: "),
			giu.Combo("##"+p.id+"AddLayerWeaponClass", weaponClassList[state.newLayerFields.weaponClass],
				weaponClassList, &state.newLayerFields.weaponClass).Size(bigListW),
		),
		giu.Separator(),
		p.makeSaveCancelButtonLine(available, state),
	}
}

func (p *COFWidget) makeSaveCancelButtonLine(available []d2enum.CompositeType, state *COFState) *giu.LineWidget {
	return giu.Line(
		giu.Button("Save##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
			newCofLayer := &d2cof.CofLayer{
				Type:        available[state.newLayerFields.layerType],
				Shadow:      byte(state.newLayerFields.selectable),
				Selectable:  state.newLayerFields.selectable == 1,
				Transparent: state.newLayerFields.transparent == 1,
				DrawEffect:  d2enum.DrawEffect(state.newLayerFields.drawEffect),
				WeaponClass: d2enum.WeaponClass(state.newLayerFields.weaponClass),
			}

			p.cof.CofLayers = append(p.cof.CofLayers, *newCofLayer)

			p.cof.NumberOfLayers++

			for i := range p.cof.Priority {
				for j := range p.cof.Priority[i] {
					p.cof.Priority[i][j] = append(p.cof.Priority[i][j], newCofLayer.Type)
				}
			}

			// this sets layer index to just added layer
			state.viewerState.layerIndex = int32(p.cof.NumberOfLayers - 1)

			state.mode = cofEditorModeViewer
		}),
		giu.Button("Cancel##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(func() {
			state.mode = cofEditorModeViewer
		}),
	)
}

func (p *COFWidget) deleteCurrentLayer(index int32) {
	p.cof.NumberOfLayers--

	newPriority := make([][][]d2enum.CompositeType, p.cof.NumberOfDirections)

	for dn := range p.cof.Priority {
		newPriority[dn] = make([][]d2enum.CompositeType, p.cof.FramesPerDirection)
		for fn := range p.cof.Priority[dn] {
			newPriority[dn][fn] = make([]d2enum.CompositeType, p.cof.NumberOfLayers)

			for ln := range p.cof.Priority[dn][fn] {
				if p.cof.CofLayers[index].Type != p.cof.Priority[dn][fn][ln] {
					newPriority[dn][fn] = append(newPriority[dn][fn], p.cof.Priority[dn][fn][ln])
				}
			}
		}
	}

	p.cof.Priority = newPriority

	newLayers := make([]d2cof.CofLayer, 0)

	for n, i := range p.cof.CofLayers {
		if int32(n) != index {
			newLayers = append(newLayers, i)
		}
	}

	p.cof.CofLayers = newLayers

	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	if state.viewerState.layerIndex != 0 {
		state.viewerState.layerIndex--
	}
}

func (p *COFWidget) duplicateDirection() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	idx := state.viewerState.directionIndex

	p.cof.NumberOfDirections++

	p.cof.Priority = append(p.cof.Priority, p.cof.Priority[idx])

	// nolint:gomnd // directionIndex starts from 0, but len from 1
	state.directionIndex = int32(len(p.cof.Priority) - 1)
}

func (p *COFWidget) deleteCurrentDirection() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)
	index := state.viewerState.directionIndex

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

	state.mode = cofEditorModeAddLayer
}
