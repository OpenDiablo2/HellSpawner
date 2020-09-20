package hsui

import (
	"log"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/hajimehoshi/ebiten"
)

type HBox struct {
	children     []Widget
	dirty        bool
	padding      int
	childSpacing int
	expandChild  bool
	alignment    hscommon.HAlign
}

func CreateHBox() *HBox {
	result := &HBox{
		children:     []Widget{},
		dirty:        false,
		padding:      1,
		childSpacing: 4,
		expandChild:  false,
		alignment:    hscommon.HAlignLeft,
	}

	return result
}

func (h *HBox) Render(screen *ebiten.Image, x, y, width, height int) {
	if width <= 0 || height <= 0 {
		return
	}

	var childWidth int

	totalChildWidth := 0

	if h.alignment != hscommon.HAlignLeft && !h.expandChild {
		for idx := range h.children {
			childWidth, _ = h.children[idx].GetRequestedSize()
			totalChildWidth += childWidth
		}
		totalChildWidth += (len(h.children) - 1) * h.childSpacing
	}

	curY := y + h.padding
	curX := 0

	if h.expandChild {
		curX = x + h.padding
	} else {
		switch h.alignment {
		case hscommon.HAlignLeft:
			curX = x + h.padding
		case hscommon.HAlignCenter:
			curX = x + (width / 2) - (totalChildWidth / 2)
		case hscommon.HAlignRight:
			curY = y + height - totalChildWidth
		default:
			log.Fatal("unknown HAlign type specified")
		}
	}

	if h.expandChild {
		childWidth = (width - (h.padding * 2) - ((len(h.children) - 1) * h.childSpacing)) / len(h.children)
	}

	for idx := range h.children {
		if !h.expandChild {
			childWidth, _ = h.children[idx].GetRequestedSize()
		}
		h.children[idx].Render(screen, curX, curY, childWidth, height)
		curX += childWidth + h.childSpacing
	}
}

func (h *HBox) Update() {
	if h.dirty {
		h.Invalidate()
	}

	for idx := range h.children {
		h.children[idx].Update()
	}
}

func (h *HBox) GetRequestedSize() (int, int) {
	tw := 0
	th := 0

	for idx := range h.children {
		cw, ch := h.children[idx].GetRequestedSize()
		if th < ch {
			th = ch
		}
		tw += cw
	}

	return tw, th
}

func (h *HBox) Invalidate() {
	for idx := range h.children {
		h.children[idx].Invalidate()
	}
}

func (h *HBox) AddChild(widget Widget) {
	h.children = append(h.children, widget)
	h.dirty = true
}

func (h *HBox) SetAlignment(align hscommon.HAlign) {
	h.alignment = align
	h.dirty = true
}

func (h *HBox) GetAlignment() hscommon.HAlign {
	return h.alignment
}

func (h *HBox) SetChildSpacing(spacing int) {
	h.childSpacing = spacing
	h.dirty = true
}

func (h *HBox) GetChildSpacing() int {
	return h.childSpacing
}

func (h *HBox) SetPadding(padding int) {
	h.padding = padding
	h.dirty = true
}

func (h *HBox) GetPadding() int {
	return h.padding
}

func (h *HBox) SetExpandChild(expand bool) {
	h.expandChild = expand
	h.dirty = true
}

func (h *HBox) GetExpandChild() bool {
	return h.expandChild
}
