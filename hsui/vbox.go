package hsui

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/hajimehoshi/ebiten"
)

type VBox struct {
	children     []Widget
	dirty        bool
	padding      int
	childSpacing int
	expandChild  bool
	alignment    hscommon.VAlign
}

func CreateVBox() *VBox {
	result := &VBox{
		children:     []Widget{},
		dirty:        false,
		padding:      1,
		childSpacing: 4,
		expandChild:  false,
		alignment:    hscommon.VAlignTop,
	}

	return result
}

func (v *VBox) Render(screen *ebiten.Image, x, y, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

	visibleChildren := 0
	for idx := range v.children {
		cw, ch := v.children[idx].GetRequestedSize()
		if cw <= 0 || ch <= 0 {
			continue
		}
		visibleChildren++
	}

	var childHeight int

	totalChildHeight := 0

	if v.alignment != hscommon.VAlignTop && !v.expandChild {
		for idx := range v.children {
			_, childHeight = v.children[idx].GetRequestedSize()
			totalChildHeight += childHeight
		}
		totalChildHeight += (visibleChildren - 1) * v.childSpacing
	}

	curY := 0
	curX := x + v.padding

	if v.expandChild {
		curY = y + v.padding
	} else {
		switch v.alignment {
		case hscommon.VAlignTop:
			curY = y + v.padding
		case hscommon.VAlignMiddle:
			curY = y + (height / 2) - (totalChildHeight / 2)
		case hscommon.VAlignBottom:
			curY = y + height - totalChildHeight
		default:
			log.Fatal("unknown VAlign type specified")
		}
	}

	if v.expandChild {
		childHeight = (height - (v.padding * 2) - ((visibleChildren - 1) * v.childSpacing)) / visibleChildren
	}

	for idx := range v.children {
		if !v.expandChild {
			_, childHeight = v.children[idx].GetRequestedSize()
		} else {
			_, ch := v.children[idx].GetRequestedSize()
			if ch <= 0 {
				continue
			}
		}

		if childHeight <= 0 {
			continue
		}

		v.children[idx].Render(screen, curX, curY, width, childHeight)
		curY += childHeight + v.childSpacing
	}
}

func (v *VBox) Update() (dirty bool) {
	if v.dirty {
		v.Invalidate()
	}

	dirty = false
	for idx := range v.children {
		childDirty := v.children[idx].Update()

		if childDirty {
			dirty = true
		}
	}

	if dirty {
		v.dirty = true
	}

	return dirty
}

func (v *VBox) GetRequestedSize() (int, int) {
	w := 0
	h := 0

	for idx := range v.children {
		cw, ch := v.children[idx].GetRequestedSize()
		if cw < w {
			w = cw
		}
		h += ch
	}

	return w, h
}

func (v *VBox) Invalidate() {
	for idx := range v.children {
		v.children[idx].Invalidate()
	}
}

func (v *VBox) AddChild(widget Widget) {
	v.children = append(v.children, widget)
	v.dirty = true
}

func (v *VBox) SetAlignment(align hscommon.VAlign) {
	v.alignment = align
	v.dirty = true
}

func (v *VBox) GetAlignment() hscommon.VAlign {
	return v.alignment
}

func (v *VBox) SetChildSpacing(spacing int) {
	v.childSpacing = spacing
	v.dirty = true
}

func (v *VBox) GetChildSpacing() int {
	return v.childSpacing
}

func (v *VBox) SetPadding(padding int) {
	v.padding = padding
	v.dirty = true
}

func (v *VBox) GetPadding() int {
	return v.padding
}

func (v *VBox) SetExpandChild(expand bool) {
	v.expandChild = expand
	v.dirty = true
}

func (v *VBox) GetExpandChild() bool {
	return v.expandChild
}
