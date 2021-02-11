package hsenum

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

const (
	stringNone = "None"
)

// GetDrawEffectName returns name of draw effect given
func GetDrawEffectName(eff d2enum.DrawEffect) string {
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
		effect = stringNone
	}

	return effect
}
