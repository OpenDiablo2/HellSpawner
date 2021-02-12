package hswidget

import (
	"fmt"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
)

const (
	saveCancelButtonW, saveCancelButtonH = 80, 30
	bigListW                             = 200
	trueFalseListW                       = 60
)

// COFEditor contains data necessary do edit cof file
type COFEditor struct {
	cof *d2cof.COF
	id  string
}

// NewCofEditor creates a new cof editor
func NewCofEditor(textureLoader *hscommon.TextureLoader, id string) *COFEditor {
	result := &COFEditor{
		id: id,
	}

	return result
}

// nolint:funlen // can't reduce
func (p *COFEditor) makeAddLayerLayout() giu.Layout {
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

func (p *COFEditor) deleteCurrentLayer(index int32) {
	p.cof.NumberOfLayers--

	newLayers := make([]d2cof.CofLayer, 0)

	for n, i := range p.cof.CofLayers {
		if int32(n) != index {
			newLayers = append(newLayers, i)
		}
	}

	p.cof.CofLayers = newLayers
}

func (p *COFEditor) duplicateDirection() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	idx := state.COFViewerState.directionIndex

	p.cof.NumberOfDirections++

	p.cof.Priority = append(p.cof.Priority, p.cof.Priority[idx])

	// nolint:gomnd // directionIndex starts from 0, but len from 1
	state.directionIndex = int32(len(p.cof.Priority) - 1)
}

func (p *COFEditor) deleteCurrentDirection() {
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
func (p *COFEditor) CreateNewLayer() {
	stateID := fmt.Sprintf("COFWidget_%s", p.id)
	s := giu.Context.GetState(stateID)

	state := s.(*COFState)

	state.state = cofEditorStateAddLayer
}
