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
	alignment    hscommon.HAlign
}

func CreateHBox() *VBox {
	result := &VBox{
		children:     []Widget{},
		dirty:        false,
		padding:      1,
		childSpacing: 4,
		expandChild:  false,
		alignment:    hscommon.HAlignTop,
	}

	return result
}

func (h *VBox) Render(screen *ebiten.Image, x, y, width, height int) {
	var childHeight int

	totalChildHeight := 0

	if h.alignment != hscommon.HAlignTop && !h.expandChild {
		for idx := range h.children {
			_, childHeight = h.children[idx].GetRequestedSize()
			totalChildHeight += childHeight
		}
		totalChildHeight += (len(h.children) - 1) * h.childSpacing
	}

	curY := 0
	curX := x + h.padding

	if h.expandChild {
		curY = y + h.padding
	} else {
		switch h.alignment {
		case hscommon.HAlignTop:
			curY = y + h.padding
		case hscommon.HAlignMiddle:
			curY = y + (height / 2) - (totalChildHeight / 2)
		case hscommon.HAlignBottom:
			curY = y + height - totalChildHeight
		default:
			log.Fatal("unknown HAlign type specified")
		}
	}

	if h.expandChild {
		childHeight = (height - (h.padding * 2) - ((len(h.children) - 1) * h.childSpacing)) / len(h.children)
	}

	for idx := range h.children {
		if !h.expandChild {
			_, childHeight = h.children[idx].GetRequestedSize()
		}
		h.children[idx].Render(screen, curX, curY, width, childHeight)
		curY += childHeight + h.childSpacing
	}
}

func (h *VBox) Update() {
	if h.dirty {
		h.Invalidate()
	}

	for idx := range h.children {
		h.children[idx].Update()
	}
}

func (h *VBox) GetRequestedSize() (int, int) {
	return 0, 0
}

func (h *VBox) Invalidate() {
	for idx := range h.children {
		h.children[idx].Invalidate()
	}
}

func (h *VBox) AddChild(widget Widget) {
	h.children = append(h.children, widget)
	h.dirty = true
}

func (h *VBox) SetAlignment(align hscommon.HAlign) {
	h.alignment = align
	h.dirty = true
}

func (h *VBox) GetAlignment() hscommon.HAlign {
	return h.alignment
}

func (h *VBox) SetChildSpacing(spacing int) {
	h.childSpacing = spacing
	h.dirty = true
}

func (h *VBox) GetChildSpacing() int {
	return h.childSpacing
}

func (h *VBox) SetPadding(padding int) {
	h.padding = padding
	h.dirty = true
}

func (h *VBox) GetPadding() int {
	return h.padding
}

func (h *VBox) SetExpandChild(expand bool) {
	h.expandChild = expand
	h.dirty = true
}

func (h *VBox) GetExpandChild() bool {
	return h.expandChild
}
