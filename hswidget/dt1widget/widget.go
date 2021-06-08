package dt1widget

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"strconv"

	"golang.org/x/image/colornames"

	"github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2dt1"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2interface"
	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2math"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hswidget/dt1widget/tiletypeimage"
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
	halfTileW       = subtileWidth >> 1
	halfTileH       = subtileHeight >> 1
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
	textureLoader hscommon.TextureLoader
}

// Create creates a new dt1 viewers widget
func Create(state []byte, palette *[256]d2interface.Color, textureLoader hscommon.TextureLoader, id string, dt1 *d2dt1.DT1) giu.Widget {
	result := &widget{
		id:            id,
		dt1:           dt1,
		textureLoader: textureLoader,
		palette:       palette,
	}

	result.registerKeyboardShortcuts()

	if giu.Context.GetState(result.getStateID()) == nil && state != nil {
		s := result.getState()
		s.Decode(state)
	}

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

	if len(state.tileGroups) == 0 {
		giu.Layout{
			giu.Label("Nothing to display"),
		}.Build()

		return
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
			if len(floorPix) == 0 || len(wallPix) == 0 {
				continue
			}

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

func (p *widget) makePixelBuffer(tile *d2dt1.Tile) (floorBuf, wallBuf []byte) {
	const (
		rOff = iota // rg,b offsets
		gOff
		bOff
		aOff
		bpp // bytes per pixel
	)

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

	floorBuf = make([]byte, tw*th*bpp)
	wallBuf = make([]byte, tw*th*bpp)

	for idx := range floor {
		var r, g, b, alpha byte

		floorVal := floor[idx]
		wallVal := wall[idx]

		rPos, gPos, bPos, aPos := idx*bpp+rOff, idx*bpp+gOff, idx*bpp+bOff, idx*bpp+aOff

		// the faux rgb color data here is just to make it look more interesting
		if p.palette != nil {
			col := p.palette[floorVal]
			r, g, b = col.R(), col.G(), col.B()
		} else {
			r = floorVal
			g = floorVal
			b = floorVal
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
			r = wallVal
			g = wallVal
			b = wallVal
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

	imageControls := giu.Row(
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
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{
					X: left.X + (idx * halfTileW),
					Y: left.Y - (idx * halfTileH),
				}

				p2 := image.Point{
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y + (gridDivisionsXY * halfTileH),
				}

				c := colornames.Green

				if idx == 0 || idx == gridDivisionsXY {
					c = colornames.Yellowgreen
				}

				canvas.AddLine(p1, p2, c, 1)
			}

			// make TR to BL lines
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{
					X: left.X + (idx * halfTileW),
					Y: left.Y + (idx * halfTileH),
				}

				p2 := image.Point{
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y - (gridDivisionsXY * halfTileH),
				}

				c := colornames.Green

				if idx == 0 || idx == gridDivisionsXY {
					c = colornames.Yellowgreen
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

	tileTypeInfo := giu.Layout{
		giu.Row(
			giu.Label("Type: "),
			giu.InputInt("##"+p.id+"tileTypeInt", &tile.Type).Size(inputIntW),
			giu.Combo("##"+p.id+"tileTypeList", tileTypeList[tileTypeIdx], tileTypeList, &tile.Type),
		),
	}

	w, h := tile.Width, tile.Height
	if h < 0 {
		h *= -1
	}

	roofHeight := int32(tile.RoofHeight)

	const (
		vspaceHeight = 4 // px
	)

	spacer := giu.Dummy(1, vspaceHeight)

	return giu.Layout{
		giu.Row(
			giu.InputInt("##"+p.id+"inputWidth", &w).Size(inputIntW).OnChange(func() {
				tile.Width = w
			}),
			giu.Label(" x "),
			giu.InputInt("##"+p.id+"winputHeight", &h).Size(inputIntW).OnChange(func() {
				tile.Height = h
			}),
			giu.Label("pixels"),
		),
		spacer,

		giu.Row(
			giu.Label("Direction: "),
			giu.InputInt("##"+p.id+"tileDirection", &tile.Direction).Size(inputIntW),
		),
		spacer,

		giu.Row(
			giu.Label("RoofHeight:"),
			giu.InputInt("##"+p.id+"roofHeight", &roofHeight).Size(inputIntW).OnChange(func() {
				tile.RoofHeight = int16(roofHeight)
			}),
		),
		spacer,

		tileTypeInfo,
		drawTileTypeImage(d2enum.TileType(tile.Type)),
		giu.Dummy(1, tiletypeimage.ImageH),

		giu.Row(
			giu.Label("Style:"),
			giu.InputInt("##"+p.id+"style", &tile.Style).Size(inputIntW),
		),
		spacer,

		giu.Row(
			giu.Label("Sequence:"),
			giu.InputInt("##"+p.id+"sequence", &tile.Sequence).Size(inputIntW),
		),
		spacer,

		giu.Row(
			giu.Label("RarityFrameIndex:"),
			giu.InputInt("##"+p.id+"rarityFrameIndex", &tile.RarityFrameIndex).Size(inputIntW),
		),
	}
}

func (p *widget) makeMaterialTab(tile *d2dt1.Tile) giu.Layout {
	return giu.Layout{
		giu.Label("Material Flags"),
		giu.Table("##"+p.id+"materialFlags").FastMode(true).
			Rows(giu.TableRow(
				giu.Checkbox("Other", &tile.MaterialFlags.Other),
				giu.Checkbox("Water", &tile.MaterialFlags.Water),
			),
				giu.TableRow(
					giu.Checkbox("WoodObject", &tile.MaterialFlags.WoodObject),
					giu.Checkbox("InsideStone", &tile.MaterialFlags.InsideStone),
				),
				giu.TableRow(
					giu.Checkbox("OutsideStone", &tile.MaterialFlags.OutsideStone),
					giu.Checkbox("Dirt", &tile.MaterialFlags.Dirt),
				),
				giu.TableRow(
					giu.Checkbox("Sand", &tile.MaterialFlags.Sand),
					giu.Checkbox("Wood", &tile.MaterialFlags.Wood),
				),
				giu.TableRow(
					giu.Checkbox("Lava", &tile.MaterialFlags.Lava),
					giu.Checkbox("Snow", &tile.MaterialFlags.Snow),
				),
			),
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

	const (
		maxSubtileIndex = 7
		spacerHeight    = 4 // px
	)

	return giu.Layout{
		giu.SliderInt("Subtile Type", &state.controls.subtileFlag, 0, maxSubtileIndex),
		giu.Label(subTileString(state.controls.subtileFlag)),
		p.makeSubTilePreview(tile, state),
		giu.Dummy(gridMaxWidth, gridMaxHeight),
		giu.Label("Click to Add/Remove flags"),
	}
}

func (p *widget) makeSubTilePreview(tile *d2dt1.Tile, state *widgetState) giu.Layout {
	return giu.Layout{
		giu.Custom(func() {
			canvas := giu.GetCanvas()
			pos := giu.GetCursorScreenPos()

			left := image.Point{X: 0 + pos.X, Y: (gridMaxHeight >> 1) + pos.Y}

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

				c := colornames.Green

				if idx == 0 || idx == gridDivisionsXY {
					c = colornames.Yellowgreen
				}

				for flagOffsetIdx := 0; flagOffsetIdx < gridDivisionsXY; flagOffsetIdx++ {
					if idx == gridDivisionsXY {
						continue
					}

					ox := (flagOffsetIdx + 1) * halfTileW
					oy := flagOffsetIdx * halfTileH

					flagPoint := image.Point{X: p1.X + ox, Y: p1.Y + oy}

					col := colornames.Yellow

					subtileIdx := getFlagFromPos(flagOffsetIdx, idx%gridDivisionsXY)
					flag := tile.SubTileFlags[subtileIdx].Encode()

					hasFlag := (flag & (1 << state.controls.subtileFlag)) > 0

					p.handleSubtileHoverAndClick(subtileIdx, flagPoint, canvas)

					if hasFlag {
						const circleRadius = 3 // px

						canvas.AddCircle(flagPoint, circleRadius, col, 1)
					}
				}

				canvas.AddLine(p1, p2, c, 1)
			}

			// make TR to BL lines
			for idx := 0; idx <= gridDivisionsXY; idx++ {
				p1 := image.Point{ // bottom left point
					X: left.X + (idx * halfTileW),
					Y: left.Y + (idx * halfTileH),
				}

				p2 := image.Point{ // top-right point
					X: p1.X + (gridDivisionsXY * halfTileW),
					Y: p1.Y - (gridDivisionsXY * halfTileH),
				}

				c := colornames.Green

				if idx == 0 || idx == gridDivisionsXY {
					c = colornames.Yellowgreen
				}

				canvas.AddLine(p1, p2, c, 1)
			}
		}),
	}
}

func (p *widget) handleSubtileHoverAndClick(subtileIdx int, flagPoint image.Point, canvas *giu.Canvas) {
	mousePos := giu.GetMousePos()
	delta := mousePos.Sub(flagPoint)
	dx, dy := int(math.Abs(float64(delta.X))), int(math.Abs(float64(delta.Y)))
	closeEnough := (dx < halfTileH) && (dy < halfTileH)

	// draw a crosshair on the point if hovered
	if closeEnough {
		highlight := color.RGBA{255, 255, 255, 64}

		p1, p2 := flagPoint.Sub(image.Point{X: -halfTileW}), flagPoint.Sub(image.Point{X: halfTileW})
		canvas.AddLine(p1, p2, highlight, 1)

		p3, p4 := flagPoint.Sub(image.Point{Y: -halfTileH}), flagPoint.Sub(image.Point{Y: halfTileH})
		canvas.AddLine(p3, p4, highlight, 1)
	}

	// on mouse release, toggle the flag
	if closeEnough && giu.IsMouseReleased(giu.MouseButtonLeft) {
		bit := p.getSubTileFieldToEdit(subtileIdx)
		*bit = !(*bit)
	}
}
