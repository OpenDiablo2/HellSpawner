package hswidget

import (
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2cof"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsutil"
)

const (
	upItemButtonPath   = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png"
	downItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png"
)

type COFEditor struct {
	newCofLayer      *d2cof.CofLayer
	cof              *d2cof.COF
	id               string
	upArrowTexture   *giu.Texture
	downArrowTexture *giu.Texture
}

func NewCofEditor(textureLoader *hscommon.TextureLoader, id string) *COFEditor {
	result := &COFEditor{
		id:          id,
		newCofLayer: newCofLayer(),
	}

	textureLoader.CreateTextureFromFileAsync(upItemButtonPath, func(texture *giu.Texture) {
		result.upArrowTexture = texture
	})

	textureLoader.CreateTextureFromFileAsync(downItemButtonPath, func(texture *giu.Texture) {
		result.downArrowTexture = texture
	})

	return result
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

func (p *COFEditor) makeAddLayerLayout(state *COFViewerState) giu.Layout {
	if p.newCofLayer == nil {
		p.newCofLayer = newCofLayer()

		return nil
	}

	var selectable int32 = hsutil.BoolToInt(p.newCofLayer.Selectable)
	var transparent int32 = hsutil.BoolToInt(p.newCofLayer.Transparent)
	var drawEffect int32 = int32(p.newCofLayer.DrawEffect)
	var weaponClass int32 = int32(p.newCofLayer.WeaponClass)

	trueFalse := []string{"false", "true"}

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

	/*p.newCofLayer.Type = d2enum.CompositeType(first)

	var compositeType int32 = int32(p.newCofLayer.Type)*/
	var compositeType int32

	drawEffectList := make([]string, int(d2enum.DrawEffectNone)+1)
	for i := d2enum.DrawEffectPctTransparency25; d2enum.DrawEffect(i) <= d2enum.DrawEffectNone; i++ {
		drawEffectList[int(i)] = strconv.Itoa(int(i)) + " (" + hsenum.GetDrawEffectName(i) + ")"
	}

	weaponClassList := make([]string, int(d2enum.WeaponClassTwoHandToHand)+1)
	for i := d2enum.WeaponClassNone; d2enum.WeaponClass(i) <= d2enum.WeaponClassTwoHandToHand; i++ {
		weaponClassList[int(i)] = i.String() + " (" + hsenum.GetWeaponClassString(i) + ")"
	}

	return giu.Layout{
		giu.Label("Select new COF's Layer parameters:"),
		giu.Separator(),
		giu.Line(
			giu.Label("Type: "),
			giu.Combo("##"+p.id+"AddLayerType", compositeTypeList[compositeType], compositeTypeList, &compositeType).Size(200).OnChange(func() {
				p.newCofLayer.Type = d2enum.CompositeType(compositeType)
			}),
		),
		giu.Line(
			giu.Label("Selectable: "),
			giu.Combo("##"+p.id+"AddLayerSelectable", trueFalse[selectable], trueFalse, &selectable).Size(60).OnChange(func() {
				p.newCofLayer.Selectable = hsutil.IntToBool(selectable)
			}),
		),
		giu.Line(
			giu.Label("Transparent: "),
			giu.Combo("##"+p.id+"AddLayerTransparent", trueFalse[transparent], trueFalse, &transparent).Size(60).OnChange(func() {
				p.newCofLayer.Transparent = hsutil.IntToBool(transparent)
			}),
		),
		giu.Line(
			giu.Label("Draw effect: "),
			giu.Combo("##"+p.id+"AddLayerDrawEffect", drawEffectList[drawEffect], drawEffectList, &drawEffect).Size(200).OnChange(func() {
				p.newCofLayer.DrawEffect = d2enum.DrawEffect(drawEffect)
			}),
		),
		giu.Line(
			giu.Label("Weapon class: "),
			giu.Combo("##"+p.id+"AddLayerWeaponClass", weaponClassList[weaponClass], weaponClassList, &weaponClass).Size(200).OnChange(func() {
				p.newCofLayer.WeaponClass = d2enum.WeaponClass(weaponClass)
			}),
		),
		giu.Separator(),
		giu.Line(
			giu.Button("Save##AddLayer").Size(80, 30).OnClick(func() {
				p.cof.CofLayers = append(p.cof.CofLayers, *p.newCofLayer)
				p.cof.NumberOfLayers++

				for i := range p.cof.Priority {
					for j := range p.cof.Priority[i] {
						p.cof.Priority[i][j] = append(p.cof.Priority[i][j], p.newCofLayer.Type)
					}
				}

				state.state = COFEditorStateViewer
			}),
			giu.Button("Close##AddLayer").Size(80, 30).OnClick(func() { state.state = COFEditorStateViewer }),
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

func (p *COFEditor) deleteCurrentDirection(index int32) {
	p.cof.NumberOfDirections--

	newPriority := make([][][]d2enum.CompositeType, 0)
	for n, i := range p.cof.Priority {
		if int32(n) != index {
			newPriority = append(newPriority, i)
		}
	}

	p.cof.Priority = newPriority
}
