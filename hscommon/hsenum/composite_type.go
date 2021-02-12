package hsenum

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

// GetLayerName returns full name of layer
// nolint:gocyclo // can't reduce
func GetLayerName(i interface{}) string {
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
