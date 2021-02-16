package ds1widget

import (
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
)

// Warning: this is 1:1 copy from
// github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2ds1.(*DS1).setupLayerStreamType()
// but this method is unexported for now, so...
// see https://github.com/OpenDiablo2/OpenDiablo2/pull/1059
func (p *DS1Widget) recreateLayerStreamTypes() {
	var layerStream []d2enum.LayerStreamType

	// nolint:gomnd // this is constant version
	// see in OpenDiablo2
	if p.ds1.Version < 4 {
		layerStream = []d2enum.LayerStreamType{
			d2enum.LayerStreamWall1,
			d2enum.LayerStreamFloor1,
			d2enum.LayerStreamOrientation1,
			d2enum.LayerStreamSubstitute,
			d2enum.LayerStreamShadow,
		}
	} else {
		// nolint:gomnd // constant (each wall layer has d2enum.LayerStreamWall and Orientation)
		layerStream = make([]d2enum.LayerStreamType,
			(p.ds1.NumberOfWalls*2)+p.ds1.NumberOfFloors+p.ds1.NumberOfShadowLayers+p.ds1.NumberOfSubstitutionLayers)

		layerIdx := 0
		for i := 0; i < int(p.ds1.NumberOfWalls); i++ {
			layerStream[layerIdx] = d2enum.LayerStreamType(int(d2enum.LayerStreamWall1) + i)
			layerStream[layerIdx+1] = d2enum.LayerStreamType(int(d2enum.LayerStreamOrientation1) + i)
			layerIdx += 2
		}
		for i := 0; i < int(p.ds1.NumberOfFloors); i++ {
			layerStream[layerIdx] = d2enum.LayerStreamType(int(d2enum.LayerStreamFloor1) + i)
			layerIdx++
		}
		if p.ds1.NumberOfShadowLayers > 0 {
			layerStream[layerIdx] = d2enum.LayerStreamShadow
			layerIdx++
		}
		if p.ds1.NumberOfSubstitutionLayers > 0 {
			layerStream[layerIdx] = d2enum.LayerStreamSubstitute
		}
	}

	p.ds1.LayerStreamTypes = layerStream
}
