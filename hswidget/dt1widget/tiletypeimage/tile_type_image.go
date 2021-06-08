package tiletypeimage

import (
	"image"
	"image/color"

	"golang.org/x/image/colornames"

	"github.com/ianling/giu"
)

const (
	floorW, floorH = 60, 30
	wallW, wallH   = floorW / 2, 50
	doorW, doorH   = wallW / 2, wallH / 2

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
	pos := b.pos.Add(image.Pt(0, wallH))
	p1 := pos.Add(image.Pt(0, 0))
	p2 := pos.Add(image.Pt(floorW/2, -1*floorH/2))
	p3 := pos.Add(image.Pt(floorW, 0))
	p4 := pos.Add(image.Pt(floorW/2, floorH/2))

	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 5)
	b.canvas.AddQuadFilled(p1, p2, p3, p4, b.fillingColor)

	return b
}

// WestWall adds a west wall
func (b *Builder) WestWall(filling bool) *Builder {
	p1 := b.pos.Add(image.Pt(0, wallH/3))
	p2 := b.pos.Add(image.Pt(wallW, 0))
	p3 := b.pos.Add(image.Pt(wallW, wallH-floorH/2))
	p4 := b.pos.Add(image.Pt(0, wallH))
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 3)

	if filling {
		b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)
	}

	return b
}

// NorthWall adds a north (right) wall
func (b *Builder) NorthWall(filling bool) *Builder {
	pos := b.pos.Add(image.Pt(wallW, 0))
	p1 := pos.Add(image.Pt(0, 0))
	p2 := pos.Add(image.Pt(wallW, wallH/3))
	p3 := pos.Add(image.Pt(wallW, wallH))
	p4 := pos.Add(image.Pt(0, wallH-floorH/2))
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 3)

	if filling {
		b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)
	}

	return b
}

// EastWall adds an easter wall
func (b *Builder) EastWall() *Builder {
	pos := b.pos.Add(image.Pt(wallW, floorH/2))
	p1 := pos.Add(image.Pt(0, wallH/5))
	p2 := pos.Add(image.Pt(wallW, 0))
	p3 := pos.Add(image.Pt(wallW, wallH-floorH/2))
	p4 := pos.Add(image.Pt(0, wallH))
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 3)
	b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)

	return b
}

// SoathWall adds a wall on a soath
func (b *Builder) SoathWall() *Builder {
	pos := b.pos.Add(image.Pt(0, floorH/2))
	p1 := pos.Add(image.Pt(0, 0))
	p2 := pos.Add(image.Pt(wallW, wallH/5))
	p3 := pos.Add(image.Pt(wallW, wallH))
	p4 := pos.Add(image.Pt(0, wallH-floorH/2))
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 3)
	b.canvas.AddQuadFilled(p1, p2, p3, p4, b.wallColor)

	return b
}

// WestDoor builds wall with a doors on a west edge
func (b *Builder) WestDoor() *Builder {
	p1 := b.pos.Add(image.Pt(0, wallH/3))
	p2 := b.pos.Add(image.Pt(wallW, 0))
	p3 := b.pos.Add(image.Pt(wallW, wallH-floorH/2))
	p4 := b.pos.Add(image.Pt(0, wallH))

	// bottom of the doors
	tg := float32(floorH/2) / float32(floorW/2)
	w := (wallW-doorW)/2 + doorW
	mod := float32((wallW-doorW)/2) * tg
	h := wallH - (floorH / 2) + mod
	d3 := b.pos.Add(image.Pt(w, int(h)))

	w = (wallW - doorW) / 2
	mod = float32((wallW-doorW)/2+doorW) * tg
	h = wallH - (floorH / 2) + mod
	d4 := b.pos.Add(image.Pt(w, int(h)))

	d1 := d4.Add(image.Pt(0, -doorH))
	d2 := d3.Add(image.Pt(0, -doorH))

	// wall border
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 3)
	// door border
	b.canvas.AddQuad(d1, d2, d3, d4, b.borderColor, 3)
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
	p1 := pos.Add(image.Pt(0, 0))
	p2 := pos.Add(image.Pt(wallW, wallH/3))
	p3 := pos.Add(image.Pt(wallW, wallH))
	p4 := pos.Add(image.Pt(0, wallH-floorH/2))

	// bottom of the doors
	tg := float32(floorH/2) / float32(floorW/2)
	w := (wallW-doorW)/2 + doorW
	mod := float32((wallW+doorW)/2) * tg
	h := wallH - (floorH / 2) + mod
	d3 := pos.Add(image.Pt(w, int(h)))

	w = (wallW - doorW) / 2
	mod = float32((wallW-doorW)/2) * tg
	h = wallH - (floorH / 2) + mod
	d4 := pos.Add(image.Pt(w, int(h)))

	d1 := d4.Add(image.Pt(0, -doorH))
	d2 := d3.Add(image.Pt(0, -doorH))

	// wall border
	b.canvas.AddQuad(p1, p2, p3, p4, b.borderColor, 3)
	// door border
	b.canvas.AddQuad(d1, d2, d3, d4, b.borderColor, 3)
	// wall filling
	b.canvas.AddQuadFilled(p1, p2, d2, d1, b.wallColor)
	b.canvas.AddQuadFilled(p2, p3, d3, d2, b.wallColor)
	b.canvas.AddQuadFilled(p3, p4, d4, d3, b.wallColor)
	b.canvas.AddQuadFilled(p4, p1, d1, d4, b.wallColor)

	return b
}
