package cofwidget

import (
	"fmt"
	"strconv"

	"github.com/AllenDang/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
	"github.com/OpenDiablo2/HellSpawner/hswidget"
)

const (
	layerListW                           = 64
	buttonWidthHeight                    = 15
	actionButtonW, actionButtonH         = 200, 30
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	speedInputW                          = 40
)

type widget struct {
	id            string
	cof           *d2cof.COF
	textureLoader hscommon.TextureLoader
}

// Create a new COF widget
func Create(
	state []byte,
	textureLoader hscommon.TextureLoader,
	id string, cof *d2cof.COF,
) giu.Widget {
	result := &widget{
		id:            id,
		cof:           cof,
		textureLoader: textureLoader,
	}

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
		result.setState(s)
	}

	return result
}

// Build builds a cof viewer
func (p *widget) Build() {
	state := p.getState()

	// builds appropriate menu (depends on state)
	switch state.mode {
	case modeViewer:
		p.makeViewerLayout().Build()
	case modeAddLayer:
		p.makeAddLayerLayout().Build()
	case modeConfirm:
		giu.Layout{
			giu.Label("Please confirm your decision"),
			state.confirmDialog,
		}.Build()
	}
}

func (p *widget) makeViewerLayout() giu.Layout {
	state := p.getState()

	return giu.Layout{
		giu.TabBar("COFViewerTabs").Layout(giu.Layout{
			giu.TabItem("Animation").Layout(p.makeAnimationTab(state)),
			giu.TabItem("Layer").Layout(p.makeLayerTab(state)),
			giu.TabItem("Priority").Layout(p.makePriorityTab(state)),
		}),
	}
}

func (p *widget) makeAnimationTab(state *widgetState) giu.Layout {
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

	setSpeed := func() {
		if p.cof.Speed >= maxSpeed {
			p.cof.Speed = maxSpeed
		}
	}

	speedLabel := giu.Label(strSpeed)
	speedInput := hswidget.MakeInputInt(
		"##"+p.id+"CovViewerSpeedValue",
		speedInputW,
		&p.cof.Speed,
		setSpeed,
	)

	return giu.Layout{
		giu.Label(strLabelDirections),
		p.layoutAnimFrames(state),
		giu.Line(speedLabel, speedInput),
		giu.Label(strLabelFPS),
		giu.Label(strLabelDuration),
	}
}

func (p *widget) makeLayerTab(state *widgetState) giu.Layout {
	addLayerButtonID := fmt.Sprintf("Add a new layer...##%sAddLayer", p.id)
	addLayerButton := giu.Button(addLayerButtonID).Size(actionButtonW, actionButtonH)
	addLayerButton.OnClick(func() {
		p.createNewLayer()
	})

	if state.viewerState.layer == nil {
		return giu.Layout{addLayerButton}
	}

	layerStrings := make([]string, 0)
	for idx := range p.cof.CofLayers {
		layerStrings = append(layerStrings, strconv.Itoa(int(p.cof.CofLayers[idx].Type)))
	}

	currentLayerName := layerStrings[state.viewerState.layerIndex]
	layerList := giu.Combo("##"+p.id+"layer", currentLayerName, layerStrings, &state.layerIndex)
	layerList.Size(layerListW).OnChange(p.onUpdate)

	deleteLayerButtonID := fmt.Sprintf("Delete current layer...##%sDeleteLayer", p.id)
	deleteLayerButton := giu.Button(deleteLayerButtonID).Size(actionButtonW, actionButtonH)
	deleteLayerButton.OnClick(func() {
		const (
			strPrompt  = "Do you really want to remove this layer?"
			strMessage = "If you'll click YES, all data from this layer will be lost. Continue?"
		)

		fnYes := func() {
			p.deleteCurrentLayer(state.viewerState.layerIndex)
			state.mode = modeViewer
		}

		fnNo := func() {
			state.mode = modeViewer
		}

		id := fmt.Sprintf("##%sDeleteLayerConfirm", p.id)
		state.viewerState.confirmDialog = hswidget.NewPopUpConfirmDialog(id, strPrompt, strMessage, fnYes, fnNo)

		state.mode = modeConfirm
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

func (p *widget) createNewLayer() {
	state := p.getState()

	state.mode = modeAddLayer
}

func (p *widget) makePriorityTab(state *widgetState) giu.Layout {
	if len(p.cof.Priority) == 0 {
		return giu.Layout{
			giu.Label("Nothing here"),
		}
	}

	directionStrings := make([]string, 0)
	for idx := range p.cof.Priority {
		directionStrings = append(directionStrings, fmt.Sprintf("%d", idx))
	}

	directionString := directionStrings[state.viewerState.directionIndex]
	directionList := giu.Combo("##"+p.id+"dir", directionString, directionStrings, &state.directionIndex)
	directionList.Size(layerListW).OnChange(p.onUpdate)

	frameStrings := make([]string, 0)
	for idx := range p.cof.Priority[state.viewerState.directionIndex] {
		frameStrings = append(frameStrings, fmt.Sprintf("%d", idx))
	}

	frameString := frameStrings[state.viewerState.frameIndex]
	frameList := giu.Combo("##"+p.id+"frame", frameString, frameStrings, &state.frameIndex)
	frameList.Size(layerListW).OnChange(p.onUpdate)

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
			state.mode = modeViewer
		}

		fnNo := func() {
			state.mode = modeViewer
		}

		popupID := fmt.Sprintf("%sDeleteLayerConfirm", p.id)

		state.confirmDialog = hswidget.NewPopUpConfirmDialog(popupID, strPrompt, strMessage, fnYes, fnNo)
		state.mode = modeConfirm
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

// the layout ends up looking like this:
// Frames (x6):  <- 10 ->
// you use the arrows to set the number of frames per direction
func (p *widget) layoutAnimFrames(state *widgetState) *giu.LineWidget {
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

	left := hswidget.MakeImageButton(leftButtonID, buttonWidthHeight, buttonWidthHeight, state.textures.left, fnDecrease)
	frameCount := giu.Label(fmt.Sprintf("%d", numFrames))
	right := hswidget.MakeImageButton(rightButtonID, buttonWidthHeight, buttonWidthHeight, state.textures.right, fnIncrease)

	return giu.Line(label, left, frameCount, right)
}

func (p *widget) onUpdate() {
	state := p.getState()

	clone := p.cof.CofLayers[state.viewerState.layerIndex]
	state.viewerState.layer = &clone

	giu.Context.SetState(p.id, state)
}

func (p *widget) makeLayerLayout() giu.Layout {
	state := p.getState()

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

func (p *widget) makeDirectionLayout() giu.Layout {
	const (
		strRenderOrderLabel = "Render Order (first to last):"
		fmtIncreasePriority = "LayerPriorityUp_%d"
		fmtDecreasePriority = "LayerPriorityDown_%d"
		fmtLayerLabel       = "%d: %s"
	)

	state := p.getState()

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

		increasePriority := hswidget.MakeImageButton(strIncPri, buttonWidthHeight, buttonWidthHeight, state.textures.up, fnIncPriority)
		decreasePriority := hswidget.MakeImageButton(strDecPri, buttonWidthHeight, buttonWidthHeight, state.textures.down, fnDecPriority)

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

func (p *widget) makeAddLayerLayout() giu.Layout {
	state := p.getState()

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
			hswidget.MakeCheckboxFromByte("##"+p.id+"AddLayerShadow", &state.newLayerFields.shadow),
		),
		giu.Line(
			giu.Label("Selectable: "),
			giu.Checkbox("##"+p.id+"AddLayerSelectable", &state.newLayerFields.selectable),
		),
		giu.Line(
			giu.Label("Transparent: "),
			giu.Checkbox("##"+p.id+"AddLayerTransparent", &state.newLayerFields.transparent),
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

func (p *widget) makeSaveCancelButtonLine(available []d2enum.CompositeType, state *widgetState) *giu.LineWidget {
	fnSave := func() {
		newCofLayer := &d2cof.CofLayer{
			Type:        available[state.newLayerFields.layerType],
			Shadow:      state.newLayerFields.shadow,
			Selectable:  state.newLayerFields.selectable,
			Transparent: state.newLayerFields.transparent,
			DrawEffect:  d2enum.DrawEffect(state.newLayerFields.drawEffect),
			WeaponClass: d2enum.WeaponClass(state.newLayerFields.weaponClass),
		}

		p.cof.CofLayers = append(p.cof.CofLayers, *newCofLayer)

		p.cof.NumberOfLayers++

		for dirIdx := range p.cof.Priority {
			for frameIdx := range p.cof.Priority[dirIdx] {
				p.cof.Priority[dirIdx][frameIdx] = append(p.cof.Priority[dirIdx][frameIdx], newCofLayer.Type)
			}
		}

		// this sets layer index to just added layer
		state.viewerState.layerIndex = int32(p.cof.NumberOfLayers - 1)
		state.viewerState.layer = newCofLayer
		state.mode = modeViewer
	}

	fnCancel := func() {
		state.mode = modeViewer
	}

	buttonSave := giu.Button("Save##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(fnSave)
	buttonCancel := giu.Button("Cancel##AddLayer").Size(saveCancelButtonW, saveCancelButtonH).OnClick(fnCancel)

	return giu.Line(buttonSave, buttonCancel)
}

func (p *widget) deleteCurrentLayer(index int32) {
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

	state := p.getState()

	if state.viewerState.layerIndex != 0 {
		state.viewerState.layerIndex--
	}
}

func (p *widget) duplicateDirection() {
	state := p.getState()

	idx := state.viewerState.directionIndex

	p.cof.NumberOfDirections++

	p.cof.Priority = append(p.cof.Priority, p.cof.Priority[idx])

	state.directionIndex = int32(len(p.cof.Priority) - 1)
}

func (p *widget) deleteCurrentDirection() {
	state := p.getState()

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
