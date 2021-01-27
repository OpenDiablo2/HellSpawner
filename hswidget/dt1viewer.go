package hswidget

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsenum"
)

const (
	gridMaxWidth    = 160
	gridMaxHeight   = 80
	gridDivisionsXY = 5
	subtileHeight   = gridMaxHeight / gridDivisionsXY
	subtileWidth    = gridMaxWidth / gridDivisionsXY
)

type dt1Controls struct {
	tileGroup   int32
	tileVariant int32
	// nolint:unused,structcheck // will be used
	tileType int32
	// nolint:unused,structcheck // will be used
	tileStyle int32
	// nolint:unused,structcheck // will be used
	tileSequence int32
	showGrid     bool
	showFloor    bool
	showWall     bool
	subtileFlag  int32
	scale        int32
}

// DT1ViewerState represents dt1 viewers state
type DT1ViewerState struct {
	*dt1Controls

	lastTileGroup int32

	tileGroups [][]*d2dt1.Tile
	textures   [][]map[string]*giu.Texture
}

// Dispose clears viewers state
func (is *DT1ViewerState) Dispose() {
	is.textures = nil
}

type tileIdentity string

func (tileIdentity) fromTile(tile *d2dt1.Tile) tileIdentity {
	str := fmt.Sprintf("%d:%d:%d", tile.Type, tile.Style, tile.Sequence)
	return tileIdentity(str)
}

// DT1ViewerWidget represents dt1 viewers widget
type DT1ViewerWidget struct {
	id  string
	dt1 *d2dt1.DT1
}

// DT1Viewer creates a new dt1 viewers widget
func DT1Viewer(id string, dt1 *d2dt1.DT1) *DT1ViewerWidget {
	result := &DT1ViewerWidget{
		id:  id,
		dt1: dt1,
	}

	result.registerKeyboardShortcuts()

	return result
}

func (p *DT1ViewerWidget) registerKeyboardShortcuts() {
	// noop
}

func (p *DT1ViewerWidget) getStateID() string {
	return fmt.Sprintf("DT1ViewerWidget_%s", p.id)
}

func (p *DT1ViewerWidget) getState() *DT1ViewerState {
	var state *DT1ViewerState

	s := giu.Context.GetState(p.getStateID())

	if s != nil {
		state = s.(*DT1ViewerState)
	} else {
		p.initState()
		p.makeTileTextures()
		state = p.getState()
	}

	return state
}

func (p *DT1ViewerWidget) setState(s giu.Disposable) {
	giu.Context.SetState(p.getStateID(), s)
}

func (p *DT1ViewerWidget) initState() {
	state := &DT1ViewerState{
		dt1Controls: &dt1Controls{
			showGrid:  true,
			showFloor: true,
			showWall:  true,
		},
		tileGroups: p.groupTilesByIdentity(),
	}

	p.setState(state)
}

// Build builds a viewer
func (p *DT1ViewerWidget) Build() {
	state := p.getState()

	if state.lastTileGroup != state.dt1Controls.tileGroup {
		state.lastTileGroup = state.dt1Controls.tileGroup
		state.dt1Controls.tileVariant = 0
	}

	tiles := state.tileGroups[int(state.dt1Controls.tileGroup)]
	tile := tiles[int(state.dt1Controls.tileVariant)]

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

func (p *DT1ViewerWidget) groupTilesByIdentity() [][]*d2dt1.Tile {
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

func (p *DT1ViewerWidget) makeTileTextures() {
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

			hscommon.CreateTextureFromARGB(imgFloor, func(tex *giu.Texture) {
				if group[variantIdx] == nil {
					group[variantIdx] = make(map[string]*giu.Texture)
				}

				group[variantIdx]["floor"] = tex
			})

			hscommon.CreateTextureFromARGB(imgWall, func(tex *giu.Texture) {
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

func (p *DT1ViewerWidget) makePixelBuffer(tile *d2dt1.Tile) (floorBuf, wallBuf []byte) {
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
		var alpha byte

		floorVal := floor[idx]
		wallVal := wall[idx]

		// nolint:gomnd // constant
		r, g, b, a := idx*4+0, idx*4+1, idx*4+2, idx*4+3

		// the faux rgb color data here is just to make it look more interesting
		floorBuf[r] = rangeByte(floorVal, 128, 256)
		floorBuf[g] = 0
		floorBuf[b] = rangeByte(rangeByte(floorVal, 0, 4), 128, 0)

		if floorVal > 0 {
			alpha = 255
		} else {
			alpha = 0
		}

		floorBuf[a] = alpha

		wallBuf[r] = 0
		wallBuf[g] = rangeByte(wallVal, 64, 196)
		wallBuf[b] = rangeByte(rangeByte(floorVal, 0, 4), 128, 0)

		if wallVal > 0 {
			alpha = 255
		} else {
			alpha = 0
		}

		wallBuf[a] = alpha
	}

	return floorBuf, wallBuf
}

func (p *DT1ViewerWidget) makeTileSelector() giu.Layout {
	state := p.getState()

	if state.lastTileGroup != state.dt1Controls.tileGroup {
		state.lastTileGroup = state.dt1Controls.tileGroup
		state.dt1Controls.tileVariant = 0
	}

	numGroups := len(state.tileGroups) - 1
	numVariants := len(state.tileGroups[state.dt1Controls.tileGroup]) - 1

	// actual layout
	layout := giu.Layout{
		giu.SliderInt("Tile Group", &state.dt1Controls.tileGroup, 0, int32(numGroups)),
	}

	if numVariants > 1 {
		layout = append(layout, giu.SliderInt("Tile Variant", &state.dt1Controls.tileVariant, 0, int32(numVariants)))
	}

	p.setState(state)

	return layout
}

// nolint:funlen,gocognit,gocyclo // no need to change
func (p *DT1ViewerWidget) makeTileDisplay(state *DT1ViewerState, tile *d2dt1.Tile) *giu.Layout {
	layout := giu.Layout{}

	// nolint:gocritic // could be useful
	// curFrameIndex := int(state.dt1Controls.frame) + (int(state.dt1Controls.direction) * int(p.dt1.FramesPerDirection))

	if uint32(state.dt1Controls.scale) < 1 {
		state.dt1Controls.scale = 1
	}

	err := giu.Context.GetRenderer().SetTextureMagFilter(giu.TextureFilterNearest)
	if err != nil {
		log.Println(err)
	}

	w, h := float32(tile.Width), float32(tile.Height)
	if h < 0 {
		h *= -1
	}

	curGroup, curVariant := int(state.dt1Controls.tileGroup), int(state.dt1Controls.tileVariant)

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
		giu.Checkbox("Show Grid", &state.dt1Controls.showGrid),
		giu.Checkbox("Show Floor", &state.dt1Controls.showFloor),
		giu.Checkbox("Show Wall", &state.dt1Controls.showWall),
	)

	layout = append(layout, giu.Custom(func() {
		canvas := giu.GetCanvas()
		pos := giu.GetCursorScreenPos()

		gridOffsetY := int(h - gridMaxHeight + (subtileHeight >> 1))
		if tile.Type == 0 {
			// fucking weird special case...
			gridOffsetY -= subtileHeight
		}

		if state.dt1Controls.showGrid && (state.dt1Controls.showFloor || state.dt1Controls.showWall) {
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

				// nolint:gomnd // const
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

		if state.dt1Controls.showFloor && floorTexture != nil {
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

		if state.dt1Controls.showWall && wallTexture != nil {
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

	if state.dt1Controls.showFloor || state.dt1Controls.showWall {
		layout = append(layout, giu.Dummy(w, h))
	}

	layout = append(layout, imageControls)

	return &layout
}

// nolint:gocyclo // can't reduce
func getTileTypeString(t int32) string {
	switch t {
	case hsenum.TileFloor:
		return "floor"
	case hsenum.TileSpecialTile1, hsenum.TileSpecialTile2:
		return "special"
	case hsenum.TileShadow:
		return "shadow"
	case hsenum.TileTree:
		return "wall/object"
	case hsenum.TileRoof:
		return "roof"
	case hsenum.TileLeftWall:
		return "Left Wall"
	case hsenum.TileRightWall:
		return "Upper Wall"
	case hsenum.TileRightPartOfNorthCornerWall:
		return "Upper part of an Upper-Left corner"
	case hsenum.TileLeftPartOfNorthCornerWall:
		return "Left part of an Upper-Left corner"
	case hsenum.TileLeftEndWall:
		return "Upper-Right corner"
	case hsenum.TileRightEndWall:
		return "Lower-Left corner"
	case hsenum.TileSouthCornerWall:
		return "Lower-Right corner"
	case hsenum.TileLeftWallWithDoor:
		return "Left Wall with Door object, but not always"
	case hsenum.TileRightWallWithDoor:
		return "Upper Wall with Door object, but not always"
	default:
		return "lower wall ?"
	}
}

func getTileTypeImage(t int32) string {
	switch t {
	case hsenum.TileFloor:
		return "floor.png"
	case hsenum.TileLeftWall:
		return "wall_west.png"
	case hsenum.TileRightWall:
		return "wall_north.png"
	case hsenum.TileRightPartOfNorthCornerWall:
		return "corner_upper_north.png"
	case hsenum.TileLeftPartOfNorthCornerWall:
		return "corner_upper_west.png"
	case hsenum.TileLeftEndWall:
		return "corner_upper_east.png"
	case hsenum.TileRightEndWall:
		return "corner_lower_south.png"
	case hsenum.TileSouthCornerWall:
		return "corner_lower_east.png"
	case hsenum.TileLeftWallWithDoor:
		return "door_west.png"
	case hsenum.TileRightWallWithDoor:
		return "door_north.png"
	default:
		return ""
	}
}

func (p *DT1ViewerWidget) makeTileInfoTab(tile *d2dt1.Tile) giu.Layout {
	var tileTypeImage *giu.ImageWithFileWidget

	strType := getTileTypeString(tile.Type)

	tileImageFile := getTileTypeImage(tile.Type)

	tileTypeImage = giu.ImageWithFile("./hsassets/images/" + tileImageFile)

	tileTypeInfo := giu.Layout{
		giu.Label(fmt.Sprintf("Type: %d (%s)", int(tile.Type), strType)),
	}

	if tileTypeImage != nil {
		tileTypeInfo = giu.Layout{
			giu.Label(fmt.Sprintf("Type: %d (%s)", int(tile.Type), strType)),
			tileTypeImage.Size(imageW, imageH),
		}
	}

	w, h := float32(tile.Width), float32(tile.Height)
	if h < 0 {
		h *= -1
	}

	return giu.Layout{
		giu.Label(fmt.Sprintf("%d x %d pixels", int(w), int(h))),
		giu.Dummy(1, 4),

		giu.Label(fmt.Sprintf("Direction: %d", int(tile.Direction))),
		giu.Dummy(1, 4),

		giu.Label(fmt.Sprintf("RoofHeight: %d", int(tile.RoofHeight))),
		giu.Dummy(1, 4),

		tileTypeInfo,
		giu.Dummy(1, 4),

		giu.Label(fmt.Sprintf("Style: %d", int(tile.Style))),
		giu.Dummy(1, 4),

		giu.Label(fmt.Sprintf("Sequence: %d", int(tile.Sequence))),
		giu.Dummy(1, 4),

		giu.Label(fmt.Sprintf("RarityFrameIndex: %d", int(tile.RarityFrameIndex))),
		// giu.Line(
		//	giu.Label(fmt.Sprintf("SubTileFlags: %v", tile.SubTileFlags)),
		// ),
		// giu.Line(
		//	giu.Label(fmt.Sprintf("Blocks: %v", tile.Blocks)),
		// ),
	}
}

func (p *DT1ViewerWidget) makeMaterialTab(tile *d2dt1.Tile) giu.Layout {
	isOther := tile.MaterialFlags.Other
	isWater := tile.MaterialFlags.Water
	isWoodObject := tile.MaterialFlags.WoodObject
	isInsideStone := tile.MaterialFlags.InsideStone
	isOutsideStone := tile.MaterialFlags.OutsideStone
	isDirt := tile.MaterialFlags.Dirt
	isSand := tile.MaterialFlags.Sand
	isWood := tile.MaterialFlags.Wood
	isLava := tile.MaterialFlags.Lava
	isSnow := tile.MaterialFlags.Snow

	return giu.Layout{
		giu.Label("Material Flags"),
		giu.Line(giu.Checkbox("Other", &isOther), giu.Checkbox("Water", &isWater)),
		giu.Line(giu.Checkbox("WoodObject", &isWoodObject), giu.Checkbox("InsideStone", &isInsideStone)),
		giu.Line(giu.Checkbox("OutsideStone", &isOutsideStone), giu.Checkbox("Dirt", &isDirt)),
		giu.Line(giu.Checkbox("Sand", &isSand), giu.Checkbox("Wood", &isWood)),
		giu.Line(giu.Checkbox("Lava", &isLava), giu.Checkbox("Snow", &isSnow)),
	}
}

// TileGroup returns current tile group
func (p *DT1ViewerWidget) TileGroup() int32 {
	state := p.getState()
	return state.tileGroup
}

// SetTileGroup sets current tile group
func (p *DT1ViewerWidget) SetTileGroup(tileGroup int32) {
	state := p.getState()
	if int(tileGroup) > len(state.tileGroups) {
		tileGroup = int32(len(state.tileGroups))
	} else if tileGroup < 0 {
		tileGroup = 0
	}

	state.tileGroup = tileGroup
}

type subtileFlag byte

func (f subtileFlag) from(flags d2dt1.SubTileFlags) subtileFlag {
	if flags.BlockWalk {
		f |= 1 << 0
	}

	if flags.BlockLOS {
		f |= 1 << 1
	}

	if flags.BlockJump {
		// nolint:gomnd // const
		f |= 1 << 2
	}

	if flags.BlockPlayerWalk {
		// nolint:gomnd // const
		f |= 1 << 3
	}

	if flags.Unknown1 {
		// nolint:gomnd // const
		f |= 1 << 4
	}

	if flags.BlockLight {
		// nolint:gomnd // const
		f |= 1 << 5
	}

	if flags.Unknown2 {
		// nolint:gomnd // const
		f |= 1 << 6
	}

	if flags.Unknown3 {
		// nolint:gomnd // const
		f |= 1 << 7
	}

	return f
}

// String returns current subtiles name
func (f subtileFlag) String() string {
	lookup := map[subtileFlag]string{
		1 << 0: "block walk",
		1 << 1: "block light and line of sight",
		1 << 2: "block jump/teleport",
		1 << 3: "block player walk, allow merc walk",
		1 << 4: "unknown #4",
		1 << 5: "block light only",
		1 << 6: "unknown #6",
		1 << 7: "unknown #7",
	}

	str, found := lookup[f]
	if !found {
		return "undefined"
	}

	return str
}

// nolint:unused // will be used
func (f subtileFlag) blockWalk() bool {
	// nolint:gomnd // const
	return ((f >> 0) & 0b1) > 0
}

// nolint:unused // will be used
func (f subtileFlag) blockLightAndLOS() bool {
	// nolint:gomnd // const
	return ((f >> 1) & 0b1) > 0
}

// nolint:unused // will be used
func (f subtileFlag) blockJumpAndTeleport() bool {
	// nolint:gomnd // const
	return ((f >> 2) & 0b1) > 0
}

// nolint:unused // will be used
func (f subtileFlag) blockPlayerAllowMercWalk() bool {
	// nolint:gomnd // const
	return ((f >> 3) & 0b1) > 0
}

// nolint:unused // I suppose, it will be used
func (f subtileFlag) unknown4() bool {
	// nolint:gomnd // const
	return ((f >> 4) & 0b1) > 0
}

// nolint:unused // will be used
func (f subtileFlag) blockLightOnly() bool {
	// nolint:gomnd // const
	return ((f >> 5) & 0b1) > 0
}

// nolint:unused // will be used
func (f subtileFlag) unknown6() bool {
	// nolint:gomnd // const
	return ((f >> 6) & 0b1) > 0
}

// nolint:unused // will be used
func (f subtileFlag) unknown7() bool {
	// nolint:gomnd // const
	return ((f >> 7) & 0b1) > 0
}

func getFlagFromPos(x, y int) int {
	var subtileLookup = [5][5]int{
		{20, 21, 22, 23, 24},
		{15, 16, 17, 18, 19},
		{10, 11, 12, 13, 14},
		{5, 6, 7, 8, 9},
		{0, 1, 2, 3, 4},
	}

	return subtileLookup[y][x]
}

func (p *DT1ViewerWidget) makeSubtileFlags(state *DT1ViewerState, tile *d2dt1.Tile) giu.Layout {
	if tile.Height < 0 {
		tile.Height *= -1
	}

	return giu.Layout{
		giu.SliderInt("Subtile Type", &state.dt1Controls.subtileFlag, 0, 7),
		giu.Label(subtileFlag(1 << state.dt1Controls.subtileFlag).String()),
		giu.Dummy(0, 4),
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
					// nolint:gomnd // const
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
					flag := subtileFlag(0).from(tile.SubTileFlags[getFlagFromPos(flagOffsetIdx, 4-idx)])

					hasFlag := (flag & (1 << state.dt1Controls.subtileFlag)) > 0

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
					// nolint:gomnd // const
					c.R = 255
				}

				canvas.AddLine(p1, p2, c, 1)
			}
		}),

		giu.Dummy(gridMaxWidth, gridMaxHeight),
	}
}

// this is copied from `OpenDiablo2/d2common/d2fileformats/d2dt1`,
// we want to render the isometric (floor) and rle (wall) pixel buffers separately
func decodeTileGfxData(blocks []d2dt1.Block, floorPixBuf, wallPixBuf *[]byte, tileYOffset, tileWidth int32) {
	for _, block := range blocks {
		switch block.Format {
		case d2dt1.BlockFormatIsometric:
			decodeFloorBlock(&block, floorPixBuf, tileYOffset, tileWidth)
		case d2dt1.BlockFormatRLE:
			decodeWallBlock(&block, wallPixBuf, tileYOffset, tileWidth)
		}
	}
}

// nolint:gomnd // 3D isometric decoding
func decodeFloorBlock(block *d2dt1.Block, floorPixBuf *[]byte, tileYOffset, tileWidth int32) {
	xjump := []int32{14, 12, 10, 8, 6, 4, 2, 0, 2, 4, 6, 8, 10, 12, 14}
	nbpix := []int32{4, 8, 12, 16, 20, 24, 28, 32, 28, 24, 20, 16, 12, 8, 4}
	blockX := int32(block.X)
	blockY := int32(block.Y)
	length := int32(256)
	x := int32(0)
	y := int32(0)
	idx := 0

	for length > 0 {
		x = xjump[y]
		n := nbpix[y]
		length -= n

		for n > 0 {
			offset := ((blockY + y + tileYOffset) * tileWidth) + (blockX + x)
			(*floorPixBuf)[offset] = block.EncodedData[idx]
			x++
			n--
			idx++
		}
		y++
	}
}

func decodeWallBlock(block *d2dt1.Block, wallPixBuf *[]byte, tileYOffset, tileWidth int32) {
	// RLE Encoding
	blockX := int32(block.X)
	blockY := int32(block.Y)
	x := int32(0)
	y := int32(0)
	idx := 0
	length := block.Length

	for length > 0 {
		b1 := block.EncodedData[idx]
		b2 := block.EncodedData[idx+1]
		idx += 2
		length -= 2

		if (b1 | b2) == 0 {
			x = 0
			y++

			continue
		}

		x += int32(b1)
		length -= int32(b2)

		for b2 > 0 {
			offset := ((blockY + y + tileYOffset) * tileWidth) + (blockX + x)
			(*wallPixBuf)[offset] = block.EncodedData[idx]
			idx++
			x++
			b2--
		}
	}
}
