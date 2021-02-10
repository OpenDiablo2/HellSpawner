package hsenum

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

// GetWeaponClassString returns full name of weapon class given
// nolint:gocyclo // can't reduce
func GetWeaponClassString(cls d2enum.WeaponClass) string {
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
