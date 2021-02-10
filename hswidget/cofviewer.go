package hswidget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"
)

const (
	indicatorSize = 64
)

type COFEditorState int

const (
	COFEditorStateViewer COFEditorState = iota
	COFEditorStateAddLayer
	COFEditorStateConfirm
)

// COFViewerState represents cof viewer's state
type COFViewerState struct {
	layerIndex     int32
	directionIndex int32
	frameIndex     int32
	state          COFEditorState
	layer          *d2cof.CofLayer
	newCofLayer    *d2cof.CofLayer
	confirmDialog  confirmDialog
}

type confirmDialog struct {
	confirmHeader  string
	confirmMessage string
	confirmed      bool
	cb             func()
}

func newCofLayer() *d2cof.CofLayer {
	return &d2cof.CofLayer{
		Type:        d2enum.CompositeTypeHead,
		Shadow:      1,
		Selectable:  true,
		Transparent: false,
		DrawEffect:  d2enum.DrawEffectNone,
		WeaponClass: d2enum.WeaponClassNone,
	}
}

// Dispose clears viewer's layers
func (s *COFViewerState) Dispose() {
	s.layer = nil
}

// COFViewerWidget represents cof viewer's widget
type COFViewerWidget struct {
	id  string
	cof *d2cof.COF
}

// COFViewer creates a cof viewer widget
func COFViewer(id string, cof *d2cof.COF) *COFViewerWidget {
	result := &COFViewerWidget{
		id:  id,
		cof: cof}

	return result
}

// Build builds a cof viewer
func (p *COFViewerWidget) Build() {
	stateID := fmt.Sprintf("COFViewerWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	if s == nil {
		giu.Context.SetState(stateID, &COFViewerState{
			layer:         &p.cof.CofLayers[0],
			state:         COFEditorStateViewer,
			confirmDialog: confirmDialog{},
		})

		return
	}

	state := s.(*COFViewerState)

	switch state.state {
	case COFEditorStateViewer:
		p.buildViewer(stateID, state)
	case COFEditorStateAddLayer:
		p.buildAddLayer(state)
	case COFEditorStateConfirm:
		p.buildPopUpConfirm(state)
	}
}

func (p *COFViewerWidget) buildPopUpConfirm(state *COFViewerState) {
	open := true
	giu.Layout{
		giu.Label("Please confirm"),
		giu.PopupModal(state.confirmDialog.confirmHeader).IsOpen(&open).Layout(giu.Layout{
			giu.Label(state.confirmDialog.confirmMessage),
			giu.Separator(),
			giu.Line(
				giu.Button("YES##"+p.id+"confirmDialog").Size(40, 25).OnClick(func() {
					state.confirmDialog.cb()
					state.state = COFEditorStateViewer
				}),
				giu.Button("NO##"+p.id+"confirmDialog").Size(40, 25).OnClick(func() {
					state.state = COFEditorStateViewer
				}),
			),
		}),
	}.Build()
}

func (p *COFViewerWidget) buildAddLayer(state *COFViewerState) {
	var selectable int32 = boolToInt(state.newCofLayer.Selectable)
	var transparent int32 = boolToInt(state.newCofLayer.Transparent)
	var weaponClass int32 = int32(state.newCofLayer.WeaponClass)

	trueFalse := []string{"false", "true"}

	weaponClassList := make([]string, int(d2enum.WeaponClassTwoHandToHand)+1)
	for i := d2enum.WeaponClassNone; d2enum.WeaponClass(i) <= d2enum.WeaponClassTwoHandToHand; i++ {
		weaponClassList[int(i)] = i.String() + " (" + p.getWeaponClass(i) + ")"
	}

	//compositeTypeList := make([]string, int(d2enum.CompositeTypeMax))
	compositeTypeList := make([]string, 0)
	first := d2enum.CompositeTypeHead
	for i := d2enum.CompositeTypeHead; i < d2enum.CompositeTypeMax; i++ {
		contains := false
		for _, j := range p.cof.CofLayers {
			if j.Type == i {
				contains = true
				if first == j.Type {
					first++
				}

				break
			}
		}

		if !contains {
			compositeTypeList = append(compositeTypeList, i.String()+" ("+getLayerName(i)+")")
		}
	}

	state.newCofLayer.Type = d2enum.CompositeType(first)

	var compositeType int32 = int32(state.newCofLayer.Type)

	giu.Layout{
		giu.Label("Select new COF's Layer parameters:"),
		giu.Separator(),
		giu.Line(
			giu.Label("Type: "),
			giu.Combo("##"+p.id+"AddLayerType", compositeTypeList[compositeType], compositeTypeList, &compositeType).Size(200).OnChange(func() {
				state.newCofLayer.Type = d2enum.CompositeType(compositeType)
			}),
		),
		giu.Line(
			giu.Label("Selectable: "),
			giu.Combo("##"+p.id+"AddLayerSelectable", trueFalse[selectable], trueFalse, &selectable).Size(60).OnChange(func() {
				state.newCofLayer.Selectable = intToBool(selectable)
			}),
		),
		giu.Line(
			giu.Label("Transparent: "),
			giu.Combo("##"+p.id+"AddLayerTransparent", trueFalse[transparent], trueFalse, &transparent).Size(60).OnChange(func() {
				state.newCofLayer.Selectable = intToBool(selectable)
			}),
		),
		giu.Line(
			giu.Label("WeaponClass: "),
			giu.Combo("##"+p.id+"AddLayerWeaponClass", weaponClassList[weaponClass], weaponClassList, &weaponClass).Size(200).OnChange(func() {
				state.newCofLayer.WeaponClass = d2enum.WeaponClass(weaponClass)
			}),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##AddLayer").Size(80, 30).OnClick(func() {
				p.cof.CofLayers = append(p.cof.CofLayers, *state.newCofLayer)
				p.cof.NumberOfLayers++

				for i := range p.cof.Priority {
					for j := range p.cof.Priority[i] {
						p.cof.Priority[i][j] = append(p.cof.Priority[i][j], state.newCofLayer.Type)
					}
				}

				state.state = COFEditorStateViewer
			}),
			giu.Button("Close##AddLayer").Size(80, 30).OnClick(func() { state.state = COFEditorStateViewer }),
		),
	}.Build()
}

func intToBool(i int32) bool {
	if i == 1 {
		return true
	} else {
		return false
	}

	return false
}

func boolToInt(b bool) int32 {
	if b {
		return 1
	}

	return 0
}

func (p *COFViewerWidget) buildViewer(stateID string, state *COFViewerState) {
	var l1, l2, l3, l4 string

	numDirs := p.cof.NumberOfDirections
	numFrames := p.cof.FramesPerDirection

	l1 = fmt.Sprintf("Directions: %v", numDirs)

	if numDirs > 1 {
		l2 = fmt.Sprintf("Frames (x%v): %v", numDirs, numFrames)
	} else {
		l2 = fmt.Sprintf("Frames: %v", numFrames)
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

	layerList := giu.Combo("##"+p.id+"layer", layerStrings[state.layerIndex], layerStrings, &state.layerIndex).
		Size(indicatorSize).OnChange(p.onUpdate)

	directionStrings := make([]string, 0)
	for idx := range p.cof.Priority {
		directionStrings = append(directionStrings, fmt.Sprintf("%d", idx))
	}

	directionList := giu.Combo("##"+p.id+"dir", directionStrings[state.directionIndex], directionStrings, &state.directionIndex).
		Size(indicatorSize).OnChange(p.onUpdate)

	frameStrings := make([]string, 0)
	for idx := range p.cof.Priority[state.directionIndex] {
		frameStrings = append(frameStrings, fmt.Sprintf("%d", idx))
	}

	frameList := giu.Combo("##"+p.id+"frame", frameStrings[state.frameIndex], frameStrings, &state.frameIndex).
		Size(indicatorSize).OnChange(p.onUpdate)

	const vspace = 4 //nolint:unused // will be used

	giu.TabBar("COFViewerTabs").Layout(giu.Layout{
		giu.TabItem("Animation").Layout(giu.Layout{
			giu.Label(l1),
			giu.Label(l2),
			giu.Label(l3),
			giu.Label(l4),
		}),
		giu.TabItem("Layer").Layout(giu.Layout{
			giu.Layout{
				giu.Line(giu.Label("Selected Layer: "), layerList),
				giu.Separator(),
				p.makeLayerLayout(),
				giu.Button("Add a new layer...##"+p.id+"AddLayer").Size(200, 30).OnClick(func() { state.newCofLayer = newCofLayer(); state.state = COFEditorStateAddLayer }),
				giu.Button("Delete current layer...##"+p.id+"DeleteLayer").Size(200, 30).OnClick(func() {
					state.confirmDialog = confirmDialog{
						confirmHeader:  "Do you raly want to remove this layer?",
						confirmMessage: "If you'll click YES, all data from this layer will be lost. Continue?",
						cb:             func() { p.deleteCurrentLayer(state.layerIndex) },
					}

					state.state = COFEditorStateConfirm
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
		}),
	}).Build()
}

func (p *COFViewerWidget) deleteCurrentLayer(index int32) {
	p.cof.NumberOfLayers--

	newLayers := make([]d2cof.CofLayer, 0)
	for n, i := range p.cof.CofLayers {
		if int32(n) != index {
			newLayers = append(newLayers, i)
		}
	}

	p.cof.CofLayers = newLayers
}

func (p *COFViewerWidget) onUpdate() {
	stateID := fmt.Sprintf("COFViewerWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFViewerState)

	clone := p.cof.CofLayers[state.layerIndex]
	state.layer = &clone

	giu.Context.SetState(p.id, state)
}

func (p *COFViewerWidget) makeLayerLayout() giu.Layout {
	stateID := fmt.Sprintf("COFViewerWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFViewerState)

	if state.layer == nil {
		p.onUpdate()
	}

	layerName := getLayerName(state.layer.Type)

	strType := fmt.Sprintf("Type: %s (%s)", state.layer.Type, layerName)
	strShadow := fmt.Sprintf("Shadow: %t", state.layer.Shadow > 0)
	strSelectable := fmt.Sprintf("Selectable: %t", state.layer.Selectable)
	strTransparent := fmt.Sprintf("Transparent: %t", state.layer.Transparent)

	effect := p.getDrawEffect(state.layer.DrawEffect)

	strEffect := fmt.Sprintf("Draw Effect: %s", effect)

	weapon := p.getWeaponClass(state.layer.WeaponClass)

	strWeaponClass := fmt.Sprintf("Weapon Class: (%s) %s", state.layer.WeaponClass, weapon)

	return giu.Layout{
		giu.Label(strType),
		giu.Label(strShadow),
		giu.Label(strSelectable),
		giu.Label(strTransparent),
		giu.Label(strEffect),
		giu.Label(strWeaponClass),
	}
}

func (p *COFViewerWidget) getDrawEffect(eff d2enum.DrawEffect) string {
	var effect string

	switch eff {
	case d2enum.DrawEffectPctTransparency25:
		effect = "25% alpha"
	case d2enum.DrawEffectPctTransparency50:
		effect = "50% alpha"
	case d2enum.DrawEffectPctTransparency75:
		effect = "75% alpha"
	case d2enum.DrawEffectModulate:
		effect = "Modulate"
	case d2enum.DrawEffectBurn:
		effect = "Burn"
	case d2enum.DrawEffectNormal:
		effect = "Normal"
	case d2enum.DrawEffectMod2XTrans:
		effect = "Mod2XTrans"
	case d2enum.DrawEffectMod2X:
		effect = "Mod2X"
	case d2enum.DrawEffectNone:
		// nolint:goconst // that's not a constant
		effect = "None"
	}

	return effect
}

// nolint:gocyclo // can't reduce
func (p *COFViewerWidget) getWeaponClass(cls d2enum.WeaponClass) string {
	var weapon string

	switch cls {
	case d2enum.WeaponClassNone:
		// nolint:goconst // that's not a constant
		weapon = "None"
	case d2enum.WeaponClassHandToHand:
		weapon = "Hand To Hand"
	case d2enum.WeaponClassBow:
		weapon = "Bow"
	case d2enum.WeaponClassOneHandSwing:
		weapon = "One Hand Swing"
	case d2enum.WeaponClassOneHandThrust:
		weapon = "One Hand Thrust"
	case d2enum.WeaponClassStaff:
		weapon = "Staff"
	case d2enum.WeaponClassTwoHandSwing:
		weapon = "Two Hand Swing"
	case d2enum.WeaponClassTwoHandThrust:
		weapon = "Two Hand Thrust"
	case d2enum.WeaponClassCrossbow:
		weapon = "Crossbow"
	case d2enum.WeaponClassLeftJabRightSwing:
		weapon = "Left Jab Right Swing"
	case d2enum.WeaponClassLeftJabRightThrust:
		weapon = "Left Jab Right Thrust"
	case d2enum.WeaponClassLeftSwingRightSwing:
		weapon = "Left Swing Right Swing"
	case d2enum.WeaponClassLeftSwingRightThrust:
		weapon = "Left Swing Right Thrust"
	case d2enum.WeaponClassOneHandToHand:
		weapon = "One Hand To Hand"
	case d2enum.WeaponClassTwoHandToHand:
		weapon = "Two Hand To Hand"
	}

	return weapon
}

// nolint:gocyclo // can't reduce
func getLayerName(i interface{}) string {
	var t d2enum.CompositeType

	switch j := i.(type) {
	case int:
		t = d2enum.CompositeType(j)
	case d2enum.CompositeType:
		t = j
	}

	var layerName string

	switch t {
	case d2enum.CompositeTypeHead:
		layerName = "Head"
	case d2enum.CompositeTypeTorso:
		layerName = "Torso"
	case d2enum.CompositeTypeLegs:
		layerName = "Legs"
	case d2enum.CompositeTypeRightArm:
		layerName = "Right Arm"
	case d2enum.CompositeTypeLeftArm:
		layerName = "Left Arm"
	case d2enum.CompositeTypeRightHand:
		layerName = "Right Hand"
	case d2enum.CompositeTypeLeftHand:
		layerName = "Left Hand"
	case d2enum.CompositeTypeShield:
		layerName = "Shield"
	case d2enum.CompositeTypeSpecial1:
		layerName = "Special 1"
	case d2enum.CompositeTypeSpecial2:
		layerName = "Special 2"
	case d2enum.CompositeTypeSpecial3:
		layerName = "Special 3"
	case d2enum.CompositeTypeSpecial4:
		layerName = "Special 4"
	case d2enum.CompositeTypeSpecial5:
		layerName = "Special 5"
	case d2enum.CompositeTypeSpecial6:
		layerName = "Special 6"
	case d2enum.CompositeTypeSpecial7:
		layerName = "Special 7"
	case d2enum.CompositeTypeSpecial8:
		layerName = "Special 8"
	}

	return layerName
}

func (p *COFViewerWidget) makeDirectionLayout() giu.Layout {
	stateID := fmt.Sprintf("COFViewerWidget_%s", p.id)
	state := giu.Context.GetState(stateID).(*COFViewerState)

	frames := p.cof.Priority[state.directionIndex]
	layers := frames[int(state.frameIndex)%len(frames)]

	return giu.Layout{
		giu.Label("Render Order (first to last):"),
		giu.Custom(func() {
			for idx := range layers {
				giu.Label(fmt.Sprintf("\t%d: %s", idx, getLayerName(layers[idx]))).Build()
			}
		}),
	}
}
