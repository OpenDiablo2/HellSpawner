package tiletypeimage

import (
	"image"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/AllenDang/giu"
)

const (
	half           = 2
	floorW, floorH = 60, 30
	wallW, wallH   = floorW / 2, floorH
	widthWallW     = 10
	doorW, doorH   = wallW / 2, wallH * 2 / 3
	cornerW        = 20
	cornerH        = int(float32(cornerW)/2*tg) * 2
	thickness      = 3

	// tg is a tangent of an angle between a floor's border and longer diagonal
	tg = (float32(floorH) / half) / (float32(floorW) / half)

	// ImageW - max width of an image
	ImageW = floorW + wallW
	// ImageH is a max height of image
	ImageH = floorH + wallH
)

// Builder allows to build a small tile preview depending on its type
type Builder struct {
	canvas *giu.Canvas
	pos    image.Point
	borderColor,
	fillingColor,
	wallColor color.RGBA
}

// TileTypeImage creates a new builder
func TileTypeImage(canvas *giu.Canvas, pos image.Point) *Builder {
	return &Builder{
		canvas:       canvas,
		pos:          pos,
		borderColor:  colornames.Green,
		fillingColor: colornames.Yellowgreen,
		wallColor:    colornames.Gray,
	}
}

// Floor adds a floor preview
func (b *Builder) Floor() *Builder {
	pos := b.pos.Add(image.Pt(floorW/half, wallH))
	p1 := pos.Add(image.Pt(0, 0))
	p2 := pos.Add(image.Pt(floorW/half, floorH/half))
	p3 := pos.Add(image.Pt(0, floorH))
	p4 := pos.Add(image.Pt(-floorW/half, floorH/half))

	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)
	b.canvas.AddQuadFilled(p1, p2, p3, p4, b.fillingColor)

	return b
}

// WestWall adds a west wall
func (b *Builder) WestWall(filling bool) *Builder {
	p3 := b.pos.Add(image.Pt(wallW, wallH))
	p4 := b.pos.Add(image.Pt(0, wallH+floorH/half))
	p1 := p4.Add(image.Pt(0, -wallH))
	p2 := p3.Add(image.Pt(0, -wallH))
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)

	if filling {
		b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)
	}

	return b
}

// NorthWall adds a north (right) wall
func (b *Builder) NorthWall(filling bool) *Builder {
	pos := b.pos.Add(image.Pt(wallW, 0))
	p3 := pos.Add(image.Pt(wallW, wallH+floorH/half))
	p4 := pos.Add(image.Pt(0, wallH))
	p1 := p4.Add(image.Pt(0, -wallH))
	p2 := p3.Add(image.Pt(0, -wallH))
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)

	if filling {
		b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)
	}

	return b
}

// EastWall adds an easter wall
func (b *Builder) EastWall() *Builder {
	pos := b.pos.Add(image.Pt(wallW, 0))
	my := float32(floorH/half) / float32(floorW/half) * float32(widthWallW)
	p3 := pos.Add(image.Pt(wallW, wallH+floorH/half))
	p4 := pos.Add(image.Pt(-widthWallW, wallH-int(my)))
	p1 := p4.Add(image.Pt(0, -wallH))
	p2 := p3.Add(image.Pt(0, -wallH))

	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)
	b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)

	return b
}

// SouthWall adds a wall on a south
func (b *Builder) SouthWall() *Builder {
	my := float32(floorH/half) / float32(floorW/half) * float32(widthWallW)
	p3 := b.pos.Add(image.Pt(wallW+widthWallW, wallH-int(my)))
	p4 := b.pos.Add(image.Pt(0, wallH+floorH/half))
	p1 := p4.Add(image.Pt(0, -wallH))
	p2 := p3.Add(image.Pt(0, -wallH))

	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)
	b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)

	return b
}

// WestDoor builds wall with a doors on a west edge
func (b *Builder) WestDoor() *Builder {
	p3 := b.pos.Add(image.Pt(wallW, wallH))
	p4 := b.pos.Add(image.Pt(0, wallH+floorH/half))
	p1 := p4.Add(image.Pt(0, -wallH))
	p2 := p3.Add(image.Pt(0, -wallH))

	// bottom of the doors
	w := (wallW-doorW)/half + doorW
	mod := float32((wallW-doorW)/half) * tg
	h := wallH - (floorH / half) + mod
	d3 := b.pos.Add(image.Pt(w, int(h)+floorH/half))

	w = (wallW - doorW) / half
	mod = float32((wallW-doorW)/half+doorW) * tg
	h = wallH - (floorH / half) + mod
	d4 := b.pos.Add(image.Pt(w, int(h)+floorH/half))

	d1 := d4.Add(image.Pt(0, -doorH))
	d2 := d3.Add(image.Pt(0, -doorH))

	// wall border
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)
	// door border
	b.canvas.AddQuad(d1, d2, d3, d4, b.borderColor, thickness)
	// wall filling
	b.canvas.AddQuadFilled(p1, p2, d2, d1, b.wallColor)
	b.canvas.AddQuadFilled(p2, p3, d3, d2, b.wallColor)
	b.canvas.AddQuadFilled(p3, p4, d4, d3, b.wallColor)
	b.canvas.AddQuadFilled(p4, p1, d1, d4, b.wallColor)

	return b
}

// NorthDoor builds a wall with a doors on north edge
func (b *Builder) NorthDoor() *Builder {
	pos := b.pos.Add(image.Pt(wallW, 0))
	p3 := pos.Add(image.Pt(wallW, wallH+floorH/half))
	p4 := pos.Add(image.Pt(0, wallH))
	p1 := p4.Add(image.Pt(0, -wallH))
	p2 := p3.Add(image.Pt(0, -wallH))

	// bottom of the doors
	w := (wallW-doorW)/half + doorW
	mod := float32((wallW+doorW)/half) * tg
	h := wallH - (floorH / half) + mod
	d3 := pos.Add(image.Pt(w, int(h)+floorH/half))

	w = (wallW - doorW) / half
	mod = float32((wallW-doorW)/half) * tg
	h = wallH - (floorH / half) + mod
	d4 := pos.Add(image.Pt(w, int(h)+floorH/half))

	d1 := d4.Add(image.Pt(0, -doorH))
	d2 := d3.Add(image.Pt(0, -doorH))

	// wall border
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, thickness)
	// door border
	b.canvas.AddQuad(d1, d2, d3, d4, b.borderColor, thickness)
	// wall filling
	b.canvas.AddQuadFilled(p1, p2, d2, d1, b.wallColor)
	b.canvas.AddQuadFilled(p2, p3, d3, d2, b.wallColor)
	b.canvas.AddQuadFilled(p3, p4, d4, d3, b.wallColor)
	b.canvas.AddQuadFilled(p4, p1, d1, d4, b.wallColor)

	return b
}

// Corner draws a "lower-right" corner
func (b *Builder) Corner() *Builder {
	padding := (floorW - cornerW) / half
	pos := b.pos

	b1 := pos.Add(image.Pt(padding, wallH+cornerH/2))
	b2 := pos.Add(image.Pt(floorW/half, wallH+cornerH))
	b3 := pos.Add(image.Pt(padding+cornerW, wallH+cornerH/2))

	u1 := b1.Add(image.Pt(0, -wallH))
	u2 := b2.Add(image.Pt(0, -wallH))
	u3 := b3.Add(image.Pt(0, -wallH))
	u4 := u2.Add(image.Pt(0, -cornerH))

	// borders
	b.canvas.AddQuad(u1, u2, b2, b1, b.borderColor, thickness)
	b.canvas.AddQuad(u2, u3, b3, b2, b.borderColor, thickness)
	b.canvas.AddQuad(u1, u2, u3, u4, b.borderColor, thickness)

	// filling
	b.canvas.AddQuadFilled(u1, u2.Add(image.Pt(-1, 0)), b2.Add(image.Pt(-1, 0)), b1, b.wallColor)
	b.canvas.AddQuadFilled(u2.Add(image.Pt(1, 0)), u3, b3, b2.Add(image.Pt(1, 0)), b.wallColor)
	b.canvas.AddQuadFilled(u1.Add(image.Pt(1, 0)), u2.Add(image.Pt(0, -1)), u3.Add(image.Pt(-1, 0)), u4, b.wallColor)

	return b
}
