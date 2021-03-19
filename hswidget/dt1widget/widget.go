package dt1widget

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
)

const (
	inputIntW = 30
)

const (
	gridMaxWidth    = 160
	gridMaxHeight   = 80
	gridDivisionsXY = 5
	subtileHeight   = gridMaxHeight / gridDivisionsXY
	subtileWidth    = gridMaxWidth / gridDivisionsXY
	imageW, imageH  = 32, 32
)

type tileIdentity string

func (tileIdentity) fromTile(tile *d2dt1.Tile) tileIdentity {
	str := fmt.Sprintf("%d:%d:%d", tile.Type, tile.Style, tile.Sequence)
	return tileIdentity(str)
}

// widget represents dt1 viewers widget
type widget struct {
	id            string
	dt1           *d2dt1.DT1
	palette       *[256]d2interface.Color
	textureLoader *hscommon.TextureLoader
}

// Create creates a new dt1 viewers widget
func Create(palette *[256]d2interface.Color, textureLoader *hscommon.TextureLoader, id string, dt1 *d2dt1.DT1) giu.Widget {
	result := &widget{
		id:            id,
		dt1:           dt1,
		textureLoader: textureLoader,
		palette:       palette,
	}

	result.registerKeyboardShortcuts()

	return result
}

func (p *widget) registerKeyboardShortcuts() {
	// noop
}

// Build builds a viewer
func (p *widget) Build() {
	state := p.getState()

	if state.lastTileGroup != state.controls.tileGroup {
		state.lastTileGroup = state.controls.tileGroup
		state.controls.tileVariant = 0
	}

	tiles := state.tileGroups[int(state.controls.tileGroup)]
	tile := tiles[int(state.controls.tileVariant)]

	giu.Layout{
		p.makeTileSelector(),
		giu.Separator(),
		p.makeTileDisplay(state, tile),
		giu.Separator(),
		giu.TabBar("##TabBar_dt1_" + p.id).Layout(giu.Layout{
			giu.TabItem("Info").Layout(p.makeTileInfoTab(tile)),
			giu.TabItem("Material").Layout(p.makeMaterialTab(tile)),
			giu.TabItem("Subtile Flags").Layout(p.makeSubtileFlags(state, tile)),
		}),
	}.Build()
}

func (p *widget) groupTilesByIdentity() [][]*d2dt1.Tile {
	result := make([][]*d2dt1.Tile, 0)

	var tileID, groupID tileIdentity

OUTER:
	for tileIdx := range p.dt1.Tiles {
		tile := &p.dt1.Tiles[tileIdx]
		tileID = tileID.fromTile(tile)

		for groupIdx := range result {
			groupID = groupID.fromTile(result[groupIdx][0])

			if tileID == groupID {
				result[groupIdx] = append(result[groupIdx], tile)
				continue OUTER
			}
		}

		result = append(result, []*d2dt1.Tile{tile})
	}

	return result
}

func (p *widget) makeTileTextures() {
	state := p.getState()
	textureGroups := make([][]map[string]*giu.Texture, len(state.tileGroups))

	for groupIdx := range state.tileGroups {
		group := make([]map[string]*giu.Texture, len(state.tileGroups[groupIdx]))

		for variantIdx := range state.tileGroups[groupIdx] {
			variantIdx := variantIdx
			tile := state.tileGroups[groupIdx][variantIdx]

			floorPix, wallPix := p.makePixelBuffer(tile)

			tw, th := int(tile.Width), int(tile.Height)
			if th < 0 {
				th *= -1
			}

			rect := image.Rect(0, 0, tw, th)
			imgFloor, imgWall := image.NewRGBA(rect), image.NewRGBA(rect)
			imgFloor.Pix, imgWall.Pix = floorPix, wallPix

			p.textureLoader.CreateTextureFromARGB(imgFloor, func(tex *giu.Texture) {
				if group[variantIdx] == nil {
					group[variantIdx] = make(map[string]*giu.Texture)
				}

				group[variantIdx]["floor"] = tex
			})

			p.textureLoader.CreateTextureFromARGB(imgWall, func(tex *giu.Texture) {
				if group[variantIdx] == nil {
					group[variantIdx] = make(map[string]*giu.Texture)
				}

				group[variantIdx]["wall"] = tex
			})
		}

		textureGroups[groupIdx] = group
	}

	state.textures = textureGroups

	p.setState(state)
}

func rangeByte(b byte, min, max float64) byte {
	// nolint:gomnd // constant
	return byte((float64(b)/255*(max-min) + min) * 255)
}

func (p *widget) makePixelBuffer(tile *d2dt1.Tile) (floorBuf, wallBuf []byte) {
	tw, th := int(tile.Width), int(tile.Height)
	if th < 0 {
		th *= -1
	}

	var tileYMinimum int32

	for _, block := range tile.Blocks {
		tileYMinimum = d2math.MinInt32(tileYMinimum, int32(block.Y))
	}

	tileYOffset := d2math.AbsInt32(tileYMinimum)

	floor := make([]byte, tw*th) // indices into palette
	wall := make([]byte, tw*th)  // indices into palette

	decodeTileGfxData(tile.Blocks, &floor, &wall, tileYOffset, tile.Width)

	// nolint:gomnd // constant
	floorBuf = make([]byte, tw*th*4) // rgba, fake palette values
	// nolint:gomnd // constant
	wallBuf = make([]byte, tw*th*4) // rgba, fake palette values

	for idx := range floor {
		var r, g, b, alpha byte

		floorVal := floor[idx]
		wallVal := wall[idx]

		// nolint:gomnd // constant
		rPos, gPos, bPos, aPos := idx*4+0, idx*4+1, idx*4+2, idx*4+3

		// the faux rgb color data here is just to make it look more interesting
		if p.palette != nil {
			col := p.palette[floorVal]
			r, g, b = col.R(), col.G(), col.B()
		} else {
			r = rangeByte(floorVal, 128, 256)
			g = 0
			b = rangeByte(rangeByte(floorVal, 0, 4), 128, 0)
		}

		floorBuf[rPos] = r
		floorBuf[gPos] = g
		floorBuf[bPos] = b

		if floorVal > 0 {
			alpha = 255
		} else {
			alpha = 0
		}

		floorBuf[aPos] = alpha

		if p.palette != nil {
			col := p.palette[wallVal]
			r, g, b = col.R(), col.G(), col.B()
		} else {
			r = 0
			g = rangeByte(wallVal, 64, 196)
			b = rangeByte(rangeByte(floorVal, 0, 4), 128, 0)
		}

		wallBuf[rPos] = r
		wallBuf[gPos] = g
		wallBuf[bPos] = b

		if wallVal > 0 {
			alpha = 255
		} else {
			alpha = 0
		}

		wallBuf[aPos] = alpha
	}

	return floorBuf, wallBuf
}

func (p *widget) makeTileSelector() giu.Layout {
	state := p.getState()

	if state.lastTileGroup != state.controls.tileGroup {
		state.lastTileGroup = state.controls.tileGroup
		state.controls.tileVariant = 0
	}

	numGroups := len(state.tileGroups) - 1
	numVariants := len(state.tileGroups[state.controls.tileGroup]) - 1

	// actual layout
	layout := giu.Layout{
		giu.SliderInt("Tile Group", &state.controls.tileGroup, 0, int32(numGroups)),
	}

	if numVariants > 1 {
		layout = append(layout, giu.SliderInt("Tile Variant", &state.controls.tileVariant, 0, int32(numVariants)))
	}

	p.setState(state)

	return layout
}

// nolint:funlen,gocognit,gocyclo // no need to change
func (p *widget) makeTileDisplay(state *widgetState, tile *d2dt1.Tile) *giu.Layout {
	layout := giu.Layout{}

	// nolint:gocritic // could be useful
	// curFrameIndex := int(state.controls.frame) + (int(state.controls.direction) * int(p.dt1.FramesPerDirection))

	if uint32(state.controls.scale) < 1 {
		state.controls.scale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Println(err)
	}

	w, h := float32(tile.Width), float32(tile.Height)
	if h < 0 {
		h *= -1
	}

	curGroup, curVariant := int(state.controls.tileGroup), int(state.controls.tileVariant)

	var floorTexture, wallTexture *giu.Texture

	if state.textures == nil ||
		len(state.textures) <= curGroup ||
		len(state.textures[curGroup]) <= curVariant ||
		state.textures[curGroup][curVariant] == nil {
		// do nothing
	} else {
		variant := state.textures[curGroup][curVariant]

		floorTexture = variant["floor"]
		wallTexture = variant["wall"]
	}

	imageControls := giu.Line(
		giu.Checkbox("Show Grid", &state.controls.showGrid),
		giu.Checkbox("Show Floor", &state.controls.showFloor),
		giu.Checkbox("Show Wall", &state.controls.showWall),
	)

	layout = append(layout, giu.Custom(func() {
		canvas := giu.GetCanvas()
		pos := giu.GetCursorScreenPos()

		gridOffsetY := int(h - gridMaxHeight + (subtileHeight >> 1))
		if tile.Type == 0 {
			// fucking weird special case...
			gridOffsetY -= subtileHeight
		}

		if state.controls.showGrid && (state.controls.showFloor || state.controls.showWall) {
			left := image.Point{X: 0 + pos.X, Y: pos.Y + gridOffsetY}

			halfTileW, halfTileH := subtileWidth>>1, subtileHeight>>1

			// make TL to BR lines
			// nolint:dupl // could be changed
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{
					X: left.X + (idx * halfTileW),
					Y: left.Y - (idx * halfTileH),
				}

				p2 := image.Point{
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y + (gridDivisionsXY * halfTileH),
				}

				// nolint:gomnd // const
				c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

				if idx == 0 || idx == gridDivisionsXY {
					c.R = 255
				}

				canvas.AddLine(p1, p2, c, 1)
			}

			// make TR to BL lines
			// nolint:dupl // is ok
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{
					X: left.X + (idx * halfTileW),
					Y: left.Y + (idx * halfTileH),
				}

				p2 := image.Point{
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y - (gridDivisionsXY * halfTileH),
				}

				// nolint:gomnd // const
				c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

				if idx == 0 || idx == gridDivisionsXY {
					c.R = 255
				}

				canvas.AddLine(p1, p2, c, 1)
			}
		}

		if state.controls.showFloor && floorTexture != nil {
			floorTL := image.Point{
				X: pos.X,
				Y: pos.Y,
			}

			floorBR := image.Point{
				X: floorTL.X + int(w),
				Y: floorTL.Y + int(h),
			}

			canvas.AddImage(floorTexture, floorTL, floorBR)
		}

		if state.controls.showWall && wallTexture != nil {
			wallTL := image.Point{
				X: pos.X,
				Y: pos.Y,
			}

			wallBR := image.Point{
				X: wallTL.X + int(w),
				Y: wallTL.Y + int(h),
			}

			canvas.AddImage(wallTexture, wallTL, wallBR)
		}
	}))

	if state.controls.showFloor || state.controls.showWall {
		layout = append(layout, giu.Dummy(w, h))
	}

	layout = append(layout, imageControls)

	return &layout
}

func (p *widget) makeTileInfoTab(tile *d2dt1.Tile) giu.Layout {
	var tileTypeImage *giu.ImageWithFileWidget

	// we're creating list of tile names
	tileTypeList := make([]string, d2enum.TileRightWallWithDoor+1)
	for i := d2enum.TileFloor; i <= d2enum.TileRightWallWithDoor; i++ {
		tileTypeList[int(i)] = i.String()
	}

	// tileTypeIdx is current index on tile types' list
	var tileTypeIdx int32
	// if tileTypeIdx is in range of known names (hsenum.GetTileTypeString)
	// then this index is set to tile.Type
	// else, we're adding Unknown+#tile.Type to list
	// and setting tileTypeIdx to this index
	if tile.Type <= int32(d2enum.TileRightWallWithDoor) {
		tileTypeIdx = tile.Type
	} else {
		// nolint:makezero // this is OK
		tileTypeList = append(tileTypeList, "Unknown (#"+strconv.Itoa(int(tile.Type))+")")
		tileTypeIdx = int32(len(tileTypeList) - 1)
	}

	tileImageFile := getTileTypeImage(d2enum.TileType(tile.Type))

	tileTypeImage = giu.ImageWithFile("./hsassets/images/" + tileImageFile)

	tileTypeInfo := giu.Layout{
		giu.Line(
			giu.Label("Type: "),
			giu.InputInt("##"+p.id+"tileTypeInt", &tile.Type).Size(inputIntW),
			giu.Combo("##"+p.id+"tileTypeList", tileTypeList[tileTypeIdx], tileTypeList, &tile.Type),
		),
	}

	if tileTypeImage != nil {
		tileTypeInfo = append(tileTypeInfo,
			tileTypeImage.Size(imageW, imageH),
		)
	}

	w, h := tile.Width, tile.Height
	if h < 0 {
		h *= -1
	}

	roofHeight := int32(tile.RoofHeight)

	return giu.Layout{
		giu.Line(
			giu.InputInt("##"+p.id+"inputWidth", &w).Size(inputIntW).OnChange(func() {
				tile.Width = w
			}),
			giu.Label(" x "),
			giu.InputInt("##"+p.id+"winputHeight", &h).Size(inputIntW).OnChange(func() {
				tile.Height = h
			}),
			giu.Label("pixels"),
		),
		giu.Dummy(1, 4),

		giu.Line(
			giu.Label("Direction: "),
			giu.InputInt("##"+p.id+"tileDirection", &tile.Direction).Size(inputIntW),
		),
		giu.Dummy(1, 4),

		giu.Line(
			giu.Label("RoofHeight:"),
			giu.InputInt("##"+p.id+"roofHeight", &roofHeight).Size(inputIntW).OnChange(func() {
				tile.RoofHeight = int16(roofHeight)
			}),
		),
		giu.Dummy(1, 4),

		tileTypeInfo,
		giu.Dummy(1, 4),

		giu.Line(
			giu.Label("Style:"),
			giu.InputInt("##"+p.id+"style", &tile.Style).Size(inputIntW),
		),
		giu.Dummy(1, 4),

		giu.Line(
			giu.Label("Sequence:"),
			giu.InputInt("##"+p.id+"sequence", &tile.Sequence).Size(inputIntW),
		),
		giu.Dummy(1, 4),

		giu.Line(
			giu.Label("RarityFrameIndex:"),
			giu.InputInt("##"+p.id+"rarityFrameIndex", &tile.RarityFrameIndex).Size(inputIntW),
		),
		// giu.Line(
		//	giu.Label(fmt.Sprintf("SubTileFlags: %v", tile.SubTileFlags)),
		// ),
		// giu.Line(
		//	giu.Label(fmt.Sprintf("Blocks: %v", tile.Blocks)),
		// ),
	}
}

func (p *widget) makeMaterialTab(tile *d2dt1.Tile) giu.Layout {
	return giu.Layout{
		giu.Label("Material Flags"),
		giu.FastTable("##" + p.id + "materialFlags").
			Border(false).
			Rows(giu.Rows{
				giu.Row(
					giu.Checkbox("Other", &tile.MaterialFlags.Other),
					giu.Checkbox("Water", &tile.MaterialFlags.Water),
				),
				giu.Row(
					giu.Checkbox("WoodObject", &tile.MaterialFlags.WoodObject),
					giu.Checkbox("InsideStone", &tile.MaterialFlags.InsideStone),
				),
				giu.Row(
					giu.Checkbox("OutsideStone", &tile.MaterialFlags.OutsideStone),
					giu.Checkbox("Dirt", &tile.MaterialFlags.Dirt),
				),
				giu.Row(
					giu.Checkbox("Sand", &tile.MaterialFlags.Sand),
					giu.Checkbox("Wood", &tile.MaterialFlags.Wood),
				),
				giu.Row(
					giu.Checkbox("Lava", &tile.MaterialFlags.Lava),
					giu.Checkbox("Snow", &tile.MaterialFlags.Snow),
				),
			}),
	}
}

// TileGroup returns current tile group
func (p *widget) TileGroup() int32 {
	state := p.getState()
	return state.tileGroup
}

// SetTileGroup sets current tile group
func (p *widget) SetTileGroup(tileGroup int32) {
	state := p.getState()
	if int(tileGroup) > len(state.tileGroups) {
		tileGroup = int32(len(state.tileGroups))
	} else if tileGroup < 0 {
		tileGroup = 0
	}

	state.tileGroup = tileGroup
}

func (p *widget) makeSubtileFlags(state *widgetState, tile *d2dt1.Tile) giu.Layout {
	if tile.Height < 0 {
		tile.Height *= -1
	}

	return giu.Layout{
		giu.SliderInt("Subtile Type", &state.controls.subtileFlag, 0, 7),
		giu.Label(subTileString(state.controls.subtileFlag)),
		giu.Label("Edit:"),
		giu.Custom(func() {
			for y := 0; y < gridDivisionsXY; y++ {
				layout := giu.Layout{}
				for x := 0; x < gridDivisionsXY; x++ {
					layout = append(layout,
						giu.Checkbox("##"+strconv.Itoa(y*gridDivisionsXY+x),
							p.getSubTileFieldToEdit(y+x*gridDivisionsXY),
						),
					)
				}

				giu.Line(layout...).Build()
			}
		}),
		giu.Dummy(0, 4),
		giu.Label("Preview:"),
		p.makeSubTilePreview(tile, state),

		giu.Dummy(gridMaxWidth, gridMaxHeight),
	}
}

func (p *widget) makeSubTilePreview(tile *d2dt1.Tile, state *widgetState) giu.Layout {
	return giu.Layout{
		giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()

			left := image.Point{X: 0 + pos.X, Y: (gridMaxHeight >> 1) + pos.Y}

			halfTileW, halfTileH := subtileWidth>>1, subtileHeight>>1

			// make TL to BR lines
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{ // top-left point
					X: left.X + (idx * halfTileW),
					Y: left.Y - (idx * halfTileH),
				}

				p2 := image.Point{ // bottom-right point
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y + (gridDivisionsXY * halfTileH),
				}

				// nolint:gomnd // const
				c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

				if idx == 0 || idx == gridDivisionsXY {
					c.R = 255
				}

				for flagOffsetIdx := 0; flagOffsetIdx < gridDivisionsXY; flagOffsetIdx++ {
					if idx == gridDivisionsXY {
						continue
					}

					ox := (flagOffsetIdx + 1) * halfTileW
					oy := flagOffsetIdx * halfTileH

					flagPoint := image.Point{
						X: p1.X + ox,
						Y: p1.Y + oy,
					}

					// nolint:gomnd // const
					col := color.RGBA{
						R: 0,
						G: 255,
						B: 255,
						A: 255,
					}

					// nolint:gomnd // constant
					flag := tile.SubTileFlags[getFlagFromPos(flagOffsetIdx, 4-idx)].Encode()

					hasFlag := (flag & (1 << state.controls.subtileFlag)) > 0

					if hasFlag {
						canvas.AddCircle(flagPoint, 3, col, 1, 0)
					}
				}

				canvas.AddLine(p1, p2, c, 1)
			}

			// make TR to BL lines
			// nolint:dupl // also ok
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{ // bottom left point
					X: left.X + (idx * halfTileW),
					Y: left.Y + (idx * halfTileH),
				}

				p2 := image.Point{ // top-right point
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y - (gridDivisionsXY * halfTileH),
				}

				// nolint:gomnd // const
				c := color.RGBA{R: 0, G: 255, B: 0, A: 255}

				if idx == 0 || idx == gridDivisionsXY {
					c.R = 255
				}

				canvas.AddLine(p1, p2, c, 1)
			}
		}),
	}
}
